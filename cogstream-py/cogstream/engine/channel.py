import abc
import time
from typing import NamedTuple, Optional

import cv2
import numpy as np

from cogstream.engine.io import FrameWriter, FramePacket, FrameScanner

maxint = 2 ** 32 - 1


class Frame(NamedTuple):
    image: np.ndarray
    frame_id: int = None
    metadata: bytes = None
    timestamp: float = None


class FrameReceiveChannel(abc.ABC):
    @abc.abstractmethod
    def recv(self) -> Optional[Frame]: ...


class FrameSendChannel(abc.ABC):
    @abc.abstractmethod
    def send(self, frame: Frame): ...


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
