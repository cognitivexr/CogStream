from typing import NamedTuple

from cogstream.api.messages import StreamSpec, Attributes
from cogstream.api.format import Format


class Specification(NamedTuple):
    operation: str
    input_format: Format
    attributes: Attributes = dict()


class EngineDescriptor(NamedTuple):
    name: str
    specification: Specification


class StreamMetadata:
    spec: StreamSpec
    client_format: Format
    engine_format: Format

    def __init__(self, spec, client_format, engine_format=None) -> None:
        super().__init__()
        self.spec = spec
        self.client_format = client_format
        self.engine_format = engine_format
