from dataclasses import dataclass
from typing import List, Dict, Tuple

from cogstream.api.format import Format, ColorMode, Orientation

Attributes = Dict[str, List[str]]


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


def format_from_attributes(attrs: Attributes) -> Format:
    width = int(attrs.get('format.width', [0])[0])
    height = int(attrs.get('format.height', [0])[0])
    # TODO: consider string values of color mode
    color_mode = ColorMode(int(attrs.get('format.colorMode', [0])[0]))
    orientation = Orientation(int(attrs.get('format.orientation', [0])[0]))

    return Format(width, height, color_mode, orientation)


def format_to_attributes(f: Format, attrs: Attributes):
    attrs['format.width'] = [str(f.width)]
    attrs['format.height'] = [str(f.height)]
    attrs['format.colorMode'] = [str(f.color_mode.value)]
    attrs['format.orientation'] = [str(f.orientation.value)]
