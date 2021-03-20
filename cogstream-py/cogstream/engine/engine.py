import abc
from typing import Any, NamedTuple

from cogstream.api.engines import EngineDescriptor
from cogstream.engine.channel import Frame


class EngineResult(NamedTuple):
    frame_id: int
    timestamp: float
    result: Any


class Engine(abc.ABC):

    def get_descriptor(self) -> EngineDescriptor:
        raise NotImplementedError

    def process(self, frame: Frame) -> Any:
        raise NotImplementedError

    def setup(self):
        pass

    def close(self):
        pass
