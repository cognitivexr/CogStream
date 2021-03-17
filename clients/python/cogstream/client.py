import json
import logging
import socket
import time
from dataclasses import dataclass
from typing import List, Dict, Tuple
from cogstream.typing import deep_to_dict, deep_from_dict

import cv2
import numpy as np
from websocket import create_connection

import cogstream.protocol as protocol

logger = logging.getLogger(__name__)

Attributes = Dict[str, List[str]]


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


@dataclass
class OperationSpec:
    code: str
    attributes: Attributes


@dataclass
class EngineSpec:
    name: str
    attributes: Attributes


@dataclass
class AvailableEngines:
    engines: List[EngineSpec]


@dataclass
class ClientFormatSpec:
    engine: str
    attributes: Attributes


@dataclass
class StreamSpec:
    engineAddress: str
    attributes: Attributes

    def get_socket_address(self) -> Tuple[str, int]:
        address = self.engineAddress.split(":")
        return address[0], int(address[1])


def to_attributes(dictionary: Dict) -> Attributes:
    return AttributeBuilder().update(dictionary).build()


class AttributeBuilder:
    """
    Creates attribute objects for CogStream handshake messages. Can be called in various ways.
    For example:

        b = AttributeBuilder().set('format.height', 360).update({'codecs': ('xvid', 'mpeg')})
        b['format.width'] = '640'
        print(b.build())

    Will output: {'format.height': ['360'], 'codecs': ['xvid', 'mpeg'], 'format.width': ['640']}
    """

    def __init__(self):
        self._attributes = {}

    def __setitem__(self, key, value):
        if isinstance(value, (list, tuple)):
            self._attributes[key] = [str(v) for v in value]
            return
        self._attributes[key] = [str(value)]

    def set(self, key, value):
        self[key] = value
        return self

    def update(self, doc: Dict):
        for k, v in doc.items():
            self[k] = v
        return self

    def build(self) -> Attributes:
        return self._attributes


class MediatorClient:
    def __init__(self, host, port):
        self._ws = create_connection(f"ws://{host}:{port}")

    def request_operation(self, operation_spec: OperationSpec) -> AvailableEngines:
        self._ws.send(json.dumps({
            "type": 2,
            "content": deep_to_dict(operation_spec)
        }))
        # todo handle wrong return message
        raw = json.loads(self._ws.recv())
        return deep_from_dict(raw['content'], AvailableEngines)

    def establish_format(self, client_format_spec: ClientFormatSpec) -> StreamSpec:
        self._ws.send(json.dumps({
            "type": 4,
            "content": deep_to_dict(client_format_spec)
        }))
        # todo handle wrong return message
        raw = json.loads(self._ws.recv())
        return deep_from_dict(raw['content'], StreamSpec)

    def close(self):
        self._ws.close()


class EngineClient:
    def __init__(self, stream_spec: StreamSpec) -> None:
        super().__init__()
        self.stream_spec = stream_spec
        self.address = stream_spec.get_socket_address()
        self.sock = None
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
        protocol.send_packet(sock, doc.encode('UTF-8'))

        self.sock = sock
        self.acknowledged = True

    def request(self, frame):
        payload = protocol.serialize_frame(frame)
        protocol.send_packet(self.sock, payload)
        return 'ok'

    def close(self):
        if self.sock is None:
            return

        logger.info('closing socket')
        self.sock.close()
