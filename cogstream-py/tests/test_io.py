import time
from unittest import TestCase

from cogstream.engine.io import SocketResultWriter, SocketResultScanner, ResultPacket


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
