import time
from unittest import TestCase

from cogstream.engine import EngineResult
from cogstream.engine.channel import JsonResultSendChannel, JsonResultReceiveChannel
from cogstream.engine.io import ResultWriter, ResultScanner, ResultPacket


class ResultWriterScannerPipe(ResultWriter, ResultScanner):

    def __init__(self):
        self.q = list()

    def write(self, packet: ResultPacket):
        self.q.append(packet)

    def next(self) -> ResultPacket:
        return self.q.pop()


class TestJsonResultChannel(TestCase):
    def test_send_recv_result(self):
        pipe = ResultWriterScannerPipe()

        sender = JsonResultSendChannel(0, pipe)
        receiver = JsonResultReceiveChannel(0, pipe)

        e = EngineResult(1, time.time(), 'foobar')

        sender.send_result(e)

        a = receiver.recv_result()

        self.assertEqual(e, a)
