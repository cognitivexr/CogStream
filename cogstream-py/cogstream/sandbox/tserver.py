import logging
import socket
import time

import cv2

from cogstream.engine.channel import JpegReceiveChannel
from cogstream.engine.io import SocketFrameScanner

logger = logging.getLogger(__name__)


def start_stream(conn):
    channel = JpegReceiveChannel(SocketFrameScanner(conn))

    while True:
        then = time.time()
        try:
            frame = channel.recv()
            logger.debug('received frame %s', frame)

            if frame is None:
                logger.debug('stopping stream')
                break

            cv2.imshow('server', frame.image)

            key = cv2.waitKey(1)
            if key == ord('q'):
                break
        except ConnectionResetError:
            logger.debug('stopping stream due to ConnectionResetError')
            break

        logger.debug('receiving packet bytes took %.2fms', ((time.time() - then) * 1000))


def serve(address):
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

            try:
                logger.info('client-server handshake successful, starting stream')
                start_stream(conn)
            finally:
                logger.info('closing connection %s', addr)
                conn.close()
    except KeyboardInterrupt:
        pass
    finally:
        if conn:
            conn.close()

        server_socket.close()


def main():
    logging.basicConfig(level=logging.DEBUG)
    serve(('0.0.0.0', 54321))


if __name__ == '__main__':
    main()
