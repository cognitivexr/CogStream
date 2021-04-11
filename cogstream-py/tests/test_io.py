import time
from unittest import TestCase

import numpy as np

from cogstream.engine.io import SocketResultWriter, SocketResultScanner, ResultPacket, SocketFrameWriter, \
    SocketFrameScanner, FramePacket


class FakeSocketPipe:
    buf: bytes

    def __init__(self):
        self.buf = bytes()

    def sendall(self, data: bytes, flags: int = ...) -> None:
        self.buf += data

    def recv(self, bufsize: int, flags: int = ...) -> bytes:
        b = self.buf[0:bufsize]
        self.buf = self.buf[bufsize:]
        return b


class TestSocketResultWriterScanner(TestCase):
    def test_write_scan(self):
        pipe = FakeSocketPipe()

        writer = SocketResultWriter(pipe)
        scanner = SocketResultScanner(pipe)

        then = time.time()

        e1 = ResultPacket(0, 1, then, b'foo')
        e2 = ResultPacket(0, 2, then, b'bar')

        writer.write(e1)
        writer.write(e2)

        a1 = scanner.next()
        a2 = scanner.next()

        self.assertEqual(e1, a1)
        self.assertEqual(e2, a2)

    def test_write_scan_empty_payload(self):
        pipe = FakeSocketPipe()

        writer = SocketResultWriter(pipe)
        scanner = SocketResultScanner(pipe)

        then = time.time()

        e1 = ResultPacket(0, 1, then, b'')

        writer.write(e1)

        a1 = scanner.next()

        self.assertEqual(e1, a1)


class TestSocketFrameWriterScanner(TestCase):
    def test_write_scan(self):
        pipe = FakeSocketPipe()

        writer = SocketFrameWriter(pipe)
        scanner = SocketFrameScanner(pipe)

        e1 = FramePacket(1, 42, time.time(), b'metadata', np.asarray([1, 2, 3]).tobytes())
        e2 = FramePacket(1, 43, time.time(), b'metadata', np.asarray([2, 4, 6]).tobytes())

        writer.write(e1)
        writer.write(e2)

        a1 = scanner.next()
        a2 = scanner.next()

        self.assertEqual(e1, a1)
        self.assertEqual(e2, a2)

    def test_write_scan_empty_metadata(self):
        pipe = FakeSocketPipe()

        writer = SocketFrameWriter(pipe)
        scanner = SocketFrameScanner(pipe)

        e1 = FramePacket(1, 42, time.time(), None, np.asarray([1, 2, 3]).tobytes())
        e2 = FramePacket(1, 43, time.time(), b'', np.asarray([1, 2, 3]).tobytes())

        writer.write(e1)
        writer.write(e2)

        a1 = scanner.next()
        a2 = scanner.next()

        self.assertEqual(e1, a1)
        self.assertEqual(FramePacket(e2.stream_id, e2.frame_id, e2.timestamp, None, e2.data), a2)

    def test_scan_empty_returns_None(self):
        pipe = FakeSocketPipe()

        scanner = SocketFrameScanner(pipe)

        self.assertIsNone(scanner.next())
