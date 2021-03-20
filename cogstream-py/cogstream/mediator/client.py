import json
import logging

from websocket import create_connection

from cogstream.api import OperationSpec, AvailableEngines, ClientFormatSpec, StreamSpec
from cogstream.typing import deep_to_dict, deep_from_dict

logger = logging.getLogger(__name__)


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
