import abc
import socket
import struct
from typing import NamedTuple, Tuple, Optional


class FramePacket(NamedTuple):
    # header
    stream_id: int
    frame_id: int
    timestamp: float
    # payload
    metadata: bytes
    data: bytes

    def split_time_fields(self) -> Tuple[int, int]:
        s = int(self.timestamp)
        ns = int((self.timestamp - s) * 1e7)
        return s, ns

    def __str__(self):
        return f'FramePacket(' \
               f'stream_id={self.stream_id}, ' \
               f'frame_id={self.frame_id}, ' \
               f'timestamp={self.timestamp}, ' \
               f'metadata=len({len(self.metadata) if self.metadata else 0}), ' \
               f'data=len({len(self.data)}))'


class FrameWriter(abc.ABC):
    @abc.abstractmethod
    def write(self, frame: FramePacket): ...


class FrameScanner(abc.ABC):
    @abc.abstractmethod
    def next(self) -> FramePacket: ...


class SocketFrameWriter(FrameWriter):

    def __init__(self, sock: socket.socket) -> None:
        super().__init__()
        self.sock = sock

    def write(self, frame: FramePacket):
        s, ns = frame.split_time_fields()
        metadata_len = len(frame.metadata) if frame.metadata is not None else 0
        data_len = len(frame.data)

        arr = struct.pack('<IIIIII', frame.stream_id, frame.frame_id, s, ns, metadata_len, data_len)

        if metadata_len == 0:
            self.sock.sendall(arr + frame.data)
        else:
            self.sock.sendall(arr + frame.metadata + frame.data)


def socket_recv(sock, n):
    """
    Helper function that reads from the socket until n bytes have been received and returns the data, or return None if
    EOF is hit.
    
    :param sock: the socket to receive data from
    :param n: the number of bytes to read
    :return: the data received
    """
    data = b''
    while len(data) < n:
        packet = sock.recv(n - len(data))
        if not packet:
            return None
        data += packet
    return data


class SocketFrameScanner(FrameScanner):

    def __init__(self, sock: socket.socket) -> None:
        super().__init__()
        self.sock = sock

    def next(self) -> FramePacket:
        return SocketFrameScanner.recv_packet(self.sock)

    @staticmethod
    def recv_packet(sock) -> Optional[FramePacket]:
        header = socket_recv(sock, 24)

        if header is None:
            return None

        stream_id, frame_id, t_s, t_ns, metadata_len, data_len = struct.unpack('<IIIIII', header)
        timestamp = t_s + (float(t_ns) / 1e7)

        if metadata_len > 0:
            metadata = socket_recv(sock, metadata_len)
            if metadata is None:
                return None
        else:
            metadata = None

        data = socket_recv(sock, data_len)

        if data is None:
            return None

        return FramePacket(stream_id, frame_id, timestamp, metadata, data)
