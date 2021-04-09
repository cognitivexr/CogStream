import abc
import json
import time
from typing import Optional

import cv2
import numpy as np

from cogstream.engine.engine import Frame, EngineResult, EngineResultWriter
from cogstream.engine.io import FrameWriter, FramePacket, FrameScanner, ResultWriter, ResultPacket, ResultScanner
from cogstream.typing import deep_to_dict

maxint = 2 ** 32 - 1


class FrameReceiveChannel(abc.ABC):
    @abc.abstractmethod
    def recv(self) -> Optional[Frame]: ...


class FrameSendChannel(abc.ABC):
    @abc.abstractmethod
    def send(self, frame: Frame): ...


class ResultReceiveChannel(abc.ABC):
    @abc.abstractmethod
    def recv_result(self) -> Optional[EngineResult]: ...


class ResultSendChannel(EngineResultWriter, abc.ABC):
    @abc.abstractmethod
    def send_result(self, result: EngineResult): ...

    def write(self, result: EngineResult):
        self.send_result(result)


class EngineChannel(FrameReceiveChannel, ResultSendChannel):
    in_channel: FrameReceiveChannel
    out_channel: ResultSendChannel

    def __init__(self, in_channel: FrameReceiveChannel, out_channel: ResultSendChannel) -> None:
        super().__init__()
        self.in_channel = in_channel
        self.out_channel = out_channel

    def recv(self) -> Optional[Frame]:
        return self.in_channel.recv()

    def send_result(self, result: EngineResult):
        return self.out_channel.send_result(result)


class ClientChannel(ResultReceiveChannel, FrameSendChannel):
    in_channel: ResultReceiveChannel
    out_channel: FrameSendChannel

    def __init__(self, in_channel: ResultReceiveChannel, out_channel: FrameSendChannel) -> None:
        super().__init__()
        self.in_channel = in_channel
        self.out_channel = out_channel

    def recv_result(self) -> Optional[EngineResult]:
        return self.in_channel.recv_result()

    def send(self, frame: Frame):
        return self.out_channel.send(frame)


class JsonResultSendChannel(ResultSendChannel):
    writer: ResultWriter

    def __init__(self, stream_id, writer: ResultWriter) -> None:
        super().__init__()
        self.stream_id = stream_id
        self.writer = writer

    def send_result(self, result: EngineResult):
        obj = result.result
        doc = deep_to_dict(obj)
        data = json.dumps(doc).encode('UTF-8')
        packet = ResultPacket(self.stream_id, result.frame_id, result.timestamp, data)
        self.writer.write(packet)


class JsonResultReceiveChannel(ResultReceiveChannel):
    scanner: ResultScanner

    def __init__(self, stream_id, scanner: ResultScanner) -> None:
        super().__init__()
        self.stream_id = stream_id
        self.scanner = scanner

    def recv_result(self) -> Optional[EngineResult]:
        packet = self.scanner.next()
        if not packet:
            return None

        doc = json.loads(packet.data.decode('UTF-8'))
        # TODO: deserialize into structured type: where to store the type information?
        # obj = deep_from_dict(doc, ?)
        obj = doc

        return EngineResult(frame_id=packet.frame_id, timestamp=packet.timestamp, result=obj)


class JpegReceiveChannel(FrameReceiveChannel):
    scanner: FrameScanner

    def __init__(self, scanner: FrameScanner) -> None:
        super().__init__()
        self.scanner = scanner

    def recv(self) -> Optional[Frame]:
        packet = self.scanner.next()
        if packet is None:
            return None

        arr = np.frombuffer(packet.data, dtype=np.uint8)
        img = cv2.imdecode(arr, flags=1)

        return Frame(img, packet.frame_id, packet.metadata, packet.timestamp)


class JpegSendChannel(FrameSendChannel):
    """
    A JpegChannel encodes frames as jpegs before sending them over the wire
    """
    stream_id: int
    frame_counter: int
    writer: FrameWriter

    def __init__(self, stream_id: int, writer: FrameWriter) -> None:
        super().__init__()
        self.stream_id = stream_id
        self.writer = writer
        self.frame_counter = 0

    def send(self, frame: Frame):
        _, jpg = cv2.imencode('.jpg', frame.image)
        arr = np.asarray(jpg, dtype=np.uint8)
        data = arr.tobytes()

        timestamp = frame.timestamp
        if timestamp is None:
            timestamp = time.time()
        metadata = frame.metadata

        packet = FramePacket(self.stream_id, self.frame_counter, timestamp, metadata, data)
        self.frame_counter = (self.frame_counter + 1) % maxint

        self.writer.write(packet)
