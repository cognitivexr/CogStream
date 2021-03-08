import json
import logging
import socket

import cogstream.protocol as protocol

logger = logging.getLogger(__name__)


class Client:
    def __init__(self, address) -> None:
        super().__init__()
        self.address = address
        self.sock = None
        self.handshake = False

    def open(self, agreement):
        if self.sock is not None:
            raise ValueError('already connected')

        address = self.address
        logger.info('connecting to server at %s', address)
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.connect(address)

        protocol.send_packet(sock, json.dumps(agreement).encode('UTF-8'))

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
