import json
import logging
import socket
import struct
import time

import cv2
import numpy as np

from cogstream.mediator.client import StreamSpec
from cogstream.engine.channel import JpegSendChannel, FrameSendChannel, Frame
from cogstream.engine.io import SocketFrameWriter
from cogstream.typing import deep_to_dict

logger = logging.getLogger(__name__)


def _send_stream_spec(sock, data: bytes):
    """
    Sends the given data over the socket and prefixes the packet with an appropriate packet header.

    :param sock: the socket to send the data over
    :param data: the data to send
    :return:
    """
    # Prefix each message with a 4-byte length (little endian)
    arr = struct.pack('<i', len(data)) + data
    sock.sendall(arr)


class EngineClient:
    def __init__(self, stream_spec: StreamSpec) -> None:
        super().__init__()
        self.stream_spec = stream_spec
        self.address = stream_spec.get_socket_address()
        self.sock = None
        self.channel: FrameSendChannel = None
        self.acknowledged = False

    def open(self):
        if self.sock is not None:
            raise ValueError('already connected')

        address = self.address
        logger.info('connecting to engine at %s', address)
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.connect(address)

        doc = json.dumps(deep_to_dict(self.stream_spec))
        logger.info("initializing stream with stream spec: %s", doc)
        _send_stream_spec(sock, doc.encode('UTF-8'))

        self.sock = sock
        self.channel = JpegSendChannel(0, SocketFrameWriter(sock))
        self.acknowledged = True

    def request(self, frame: np.ndarray):
        self.channel.send(Frame(frame))
        return 'ok'

    def close(self):
        if self.sock is None:
            return

        logger.info('closing socket')
        self.sock.close()

    def stream_camera(self, cam, show=True):
        stream_camera(cam, self, show)


def stream_camera(cam, client, show=True):
    goal_fps = 25

    # target frame inter-arrival time
    ia = 1 / goal_fps

    while True:
        start = time.time()

        check, frame = cam.read()
        if not check:
            logger.info('no more frames to read')
            break

        if show:
            cv2.imshow("capture", frame)

        client.request(frame)

        delay = ia - (time.time() - start)
        if delay >= 0:
            time.sleep(delay)

        key = cv2.waitKey(1)
        if key == ord('q'):
            break

        logger.info('fps: %.2f' % (1 / (time.time() - start)))
