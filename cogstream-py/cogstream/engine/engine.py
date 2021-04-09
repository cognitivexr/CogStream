import abc
from typing import Any, NamedTuple

import numpy as np

from cogstream.api.engines import EngineDescriptor


class Frame(NamedTuple):
    image: np.ndarray
    frame_id: int = None
    metadata: bytes = None
    timestamp: float = None


class EngineResult(NamedTuple):
    frame_id: int
    timestamp: float
    result: Any


class EngineResultWriter(abc.ABC):
    @abc.abstractmethod
    def write(self, result: EngineResult): ...


class Engine(abc.ABC):

    def get_descriptor(self) -> EngineDescriptor:
        raise NotImplementedError

    def process(self, frame: 'Frame', results: EngineResultWriter):
        raise NotImplementedError

    def setup(self):
        pass

    def close(self):
        pass
