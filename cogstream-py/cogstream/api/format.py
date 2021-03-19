from enum import Enum
from typing import NamedTuple


class ColorMode(Enum):
    UNKNOWN = 0
    RGB = 1
    RGBA = 2
    Gray = 3
    BGR = 4
    BGRA = 5
    HLS = 6
    Lab = 7
    Luv = 8
    Bayer = 9


class Orientation(Enum):
    UNKNOWN = 0
    TopLeft = 1
    TopRight = 2
    BottomRight = 3
    BottomLeft = 4
    LeftTop = 5
    RightTop = 6
    RightBottom = 7
    LeftBottom = 8

    def angle(self) -> int:
        return _orientation_angles[self]

    def mirrored(self) -> int:
        return _orientation_mirrored[self]


_orientation_angles = {
    Orientation.TopLeft: 0,
    Orientation.TopRight: 0,
    Orientation.BottomRight: 180,
    Orientation.BottomLeft: 180,
    Orientation.LeftTop: 90,
    Orientation.RightTop: 90,
    Orientation.RightBottom: 270,
    Orientation.LeftBottom: 270,
}

_orientation_mirrored = {
    Orientation.TopLeft: False,
    Orientation.TopRight: True,
    Orientation.BottomRight: False,
    Orientation.BottomLeft: True,
    Orientation.LeftTop: True,
    Orientation.RightTop: False,
    Orientation.RightBottom: True,
    Orientation.LeftBottom: False,
}


class Format(NamedTuple):
    width: int = 0
    height: int = 0
    color_mode: ColorMode = ColorMode.UNKNOWN
    orientation: Orientation = Orientation.UNKNOWN


AnyFormat = Format(0, 0, ColorMode.UNKNOWN, Orientation.UNKNOWN)
