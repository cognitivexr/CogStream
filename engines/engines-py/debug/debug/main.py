import argparse
import logging
import socket
import struct
import time
from typing import Optional

import numpy as np
from cogstream.api import StreamSpec, StreamMetadata, EngineDescriptor, Specification
from cogstream.api.format import AnyFormat
from cogstream.engine import srv, Engine, EngineResultWriter, EngineResult, Frame, FrameReceiveChannel
from cogstream.engine.channel import EngineChannel, JsonResultSendChannel
from cogstream.engine.io import SocketFrameScanner, SocketResultWriter, FrameScanner, FramePacket, socket_recv

logger = logging.getLogger(__name__)


class DebugReceiveChannel(FrameReceiveChannel):
    scanner: FrameScanner

    def __init__(self, scanner: FrameScanner) -> None:
        super().__init__()
        self.scanner = scanner

    def recv(self) -> Optional[Frame]:
        logger.debug('reading next frame packet ...')
        packet = self.scanner.next()

        if packet is None:
            logger.debug('frame packet is None, EOF?')
            return None

        logger.debug('frame packet received: %s', packet)

        return Frame(np.array([]), packet.frame_id, packet.metadata, packet.timestamp)


class DebugSocketFrameScanner(FrameScanner):

    def __init__(self, sock: socket.socket) -> None:
        super().__init__()
        self.sock = sock

    def next(self) -> FramePacket:
        return self.recv_packet(self.sock)

    @staticmethod
    def recv_packet(sock) -> Optional[FramePacket]:
        logger.debug('receiving 24 byte FramePacket header...')
        header = socket_recv(sock, 24)

        if header is None:
            return None

        stream_id, frame_id, t_s, t_ns, metadata_len, data_len = struct.unpack('<IIIIII', header)
        logger.debug('deserialized header %d,%d,%d,%d,%d,%d', stream_id, frame_id, t_s, t_ns, metadata_len, data_len)

        timestamp = t_s + (float(t_ns) / 1e7)

        if metadata_len > 0:
            metadata = socket_recv(sock, metadata_len)
            if metadata is None:
                return None
            logger.debug('received %d FramePacket metadata bytes', len(metadata))
        else:
            metadata = None

        data = socket_recv(sock, data_len)

        logger.debug('received %d FramePacket data bytes', len(data))

        if data is None:
            return None

        return FramePacket(stream_id, frame_id, timestamp, metadata, data)


class DebugEngine(Engine):

    def get_descriptor(self) -> EngineDescriptor:
        return EngineDescriptor('debug', Specification('analyze', AnyFormat, dict()))

    def process(self, frame: Frame, results: EngineResultWriter):
        logger.info('DebugEngine received frame %s', frame)
        result = EngineResult(frame.frame_id, time.time(), {'frame_id': frame.frame_id, 'answer': 42})
        logger.info('DebugEngine sending result %s', result)
        results.write(result)


def _create_channel(sock, _spec: StreamSpec, _metadata: StreamMetadata):
    in_channel = DebugReceiveChannel(DebugSocketFrameScanner(sock))
    out_channel = JsonResultSendChannel(0, SocketResultWriter(sock))

    return EngineChannel(in_channel, out_channel)


def main():
    logging.basicConfig(level=logging.DEBUG)

    parser = argparse.ArgumentParser(description='CogStream engine for debugging network clients')
    parser.add_argument('--host', type=str, help='the host to bind to', default='0.0.0.0')
    parser.add_argument('--port', type=int, help='the port to bind to (default 54321)', default=54321)

    args = parser.parse_args()

    # overwrite channel to inject debug channels that don't decode image data
    srv._create_channel = _create_channel

    try:
        srv.serve_engine((args.host, args.port), lambda metadata: DebugEngine())
    except KeyboardInterrupt:
        print('ok, bye')


if __name__ == '__main__':
    main()
