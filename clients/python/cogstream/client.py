import json
import logging
import socket
import time

import cv2
import numpy as np
from websocket import create_connection

import cogstream.protocol as protocol

logger = logging.getLogger(__name__)


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

        jpg: np.ndarray = cv2.imencode('.jpg', frame)[1]

        client.request(jpg)

        delay = ia - (time.time() - start)
        if delay >= 0:
            time.sleep(delay)

        key = cv2.waitKey(1)
        if key == ord('q'):
            break

        logger.info('fps: %.2f' % (1 / (time.time() - start)))


class MediatorClient:
    def __init__(self, host, port):
        self._ws = create_connection(f"ws://{host}:{port}")

    def request_operation(self, op_spec):
        self._ws.send(json.dumps({
            "type": 2,
            "content": op_spec
        }))
        return json.loads(self._ws.recv())

    def establish_format(self, stream_format):
        self._ws.send(json.dumps({
            "type": 4,
            "content": stream_format
        }))
        return json.loads(self._ws.recv())

    def close(self):
        self._ws.close()


class EngineClient:
    def __init__(self, address) -> None:
        super().__init__()
        self.address = address
        self.sock = None
        self.handshake = False

    def open(self, stream_spec):
        if self.sock is not None:
            raise ValueError('already connected')

        address = self.address
        logger.info('connecting to server at %s', address)
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.connect(address)

        protocol.send_packet(sock, json.dumps(stream_spec).encode('UTF-8'))

        self.sock = sock
        self.handshake = True

    def close(self):
        if self.sock is None:
            return

        logger.info('closing socket')
        self.sock.close()

    def request(self, frame):
        payload = protocol.serialize_frame(frame)
        protocol.send_packet(self.sock, payload)
        return 'ok'
