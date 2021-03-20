import json
import logging
import socket
import struct
import time
from typing import Callable, Union

from cogstream.api import StreamSpec
from cogstream.api.engines import StreamMetadata
from cogstream.api.messages import format_from_attributes
from cogstream.engine.channel import JpegReceiveChannel, FrameReceiveChannel, Frame
from cogstream.engine.engine import Engine, EngineResult
from cogstream.engine.io import SocketFrameScanner, socket_recv
from cogstream.engine.transform import build_transformer
from cogstream.typing import deep_from_dict

logger = logging.getLogger(__name__)

Address = Union[tuple, str]

ConnectionHandler = Callable[[socket.socket, StreamSpec], None]
EngineFactory = Callable[[StreamMetadata], Engine]


def serve(address: Address, connection_handler: ConnectionHandler):
    logger.info('starting server on address %s', address)
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    server_socket.bind(address)
    server_socket.listen(1)

    conn = None
    try:
        while True:
            logger.info('waiting for next connection')
            conn, addr = server_socket.accept()
            logger.info('initiating handshake with %s', addr)

            # TODO: multiple connections
            try:
                stream_spec = _socket_recv_stream_spec(conn)
                logger.info('stream spec received: %s', stream_spec)
                connection_handler(conn, stream_spec)
            except:
                logger.exception("exception while handling connection %s", conn)
            finally:
                logger.info('closing connection %s', addr)
                conn.close()
    except KeyboardInterrupt:
        pass
    finally:
        if conn:
            conn.close()

        server_socket.close()


def _socket_recv_stream_spec(sock) -> StreamSpec:
    header = socket_recv(sock, 4)
    data_len = struct.unpack('<I', header)[0]

    if data_len <= 0:
        raise ValueError('error receiving header')

    data = socket_recv(sock, data_len)
    doc = json.loads(data)
    return deep_from_dict(doc, StreamSpec)


def _start_stream(engine: Engine, channel: FrameReceiveChannel, metadata: StreamMetadata):
    transform = build_transformer(metadata.client_format, metadata.engine_format)

    while True:
        then = time.time()
        try:
            frame = channel.recv()

            if frame is None:
                logger.debug('stopping stream')
                break
            logger.debug('received frame id=%s,ts=%.4f', frame.frame_id, frame.timestamp)

            image = transform(frame.image)
            tf_frame = Frame(image, frame.frame_id, frame.metadata, frame.timestamp)

            # call engine
            processed = engine.process(tf_frame)
            result = EngineResult(tf_frame.frame_id, time.time(), processed)

            # TODO: actually return
            logger.debug('engine result: %s', result)

        except ConnectionResetError:
            logger.debug('stopping stream due to ConnectionResetError')
            break

        logger.debug('receiving packet bytes took %.2fms', ((time.time() - then) * 1000))


def _stream_handler(sock, engine_factory: EngineFactory, spec: StreamSpec):
    # initialize stream metadata and setup engine
    client_format = format_from_attributes(spec.attributes)
    metadata = StreamMetadata(spec, client_format)
    engine = engine_factory(metadata)
    metadata.engine_format = engine.get_descriptor().specification.input_format

    # todo: determine from stream metadata
    channel = JpegReceiveChannel(SocketFrameScanner(sock))
    try:
        engine.setup()
        _start_stream(engine, channel, metadata)
    finally:
        engine.close()


def serve_engine(address: Address, engine_factory: EngineFactory):
    def handler(sock, spec: StreamSpec):
        _stream_handler(sock, engine_factory, spec)

    return serve(address, handler)
