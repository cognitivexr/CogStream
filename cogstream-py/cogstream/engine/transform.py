from typing import Callable, Tuple

import cv2.cv2 as cv2
import numpy as np

from cogstream.api.format import ColorMode, Orientation, Format

Function = Callable[[np.ndarray], np.ndarray]


def _no_transform(_: np.ndarray) -> np.ndarray:
    return None


def _cv2_cvt_col(code) -> Function:
    def func(src):
        return cv2.cvtColor(src, code=code)

    return func


def _cv2_flip(code) -> Function:
    def func(src):
        return cv2.flip(src, flipCode=code)

    return func


def _cv2_rotate(code) -> Function:
    def func(src):
        return cv2.rotate(src, rotateCode=code)

    return func


def _cv2_resize_with_scale(dim: Tuple[int, int]) -> Function:
    def func(src):
        return cv2.resize(src, dsize=dim)

    return func


def _get_color_transform(from_color: ColorMode, to_color: ColorMode) -> Function:
    if from_color == to_color:
        return NoTransform
    if from_color == ColorMode.UNKNOWN or to_color == ColorMode.UNKNOWN:
        return NoTransform

    code = get_color_conversion_code(from_color, to_color)
    return _cv2_cvt_col(code)


def _pipeline(*fns) -> Function:
    pipeline = [fn for fn in fns if fn is not NoTransform]

    if len(pipeline) == 0:
        return NoTransform
    if len(pipeline) == 1:
        return pipeline[0]

    def fn(src: np.ndarray):
        dst = src

        for f in pipeline:
            dst = f(dst)

        return dst

    return fn


def _get_orientation_transformation(source: Orientation, target: Orientation) -> Function:
    if source == target:
        return NoTransform
    if source == Orientation.UNKNOWN or target == Orientation.UNKNOWN:
        return NoTransform

    rotate, flip = NoTransform, NoTransform

    rotate_flag = get_rotate_flag(source.angle(), target.angle())
    if rotate_flag >= 0:
        rotate = _cv2_rotate(rotate_flag)

    if source.mirrored() != target.mirrored():
        a = int(((target.angle() % 360) / 90) % 2)
        if a == 0:
            flip = _cv2_flip(1)
        else:
            flip = _cv2_flip(0)

    return Pipeline(rotate, flip)


def get_color_conversion_code(from_color: ColorMode, to_color: ColorMode):
    # a bit hacky, ... but works
    var = f'COLOR_{from_color.name.upper()}2{to_color.name.upper()}'

    try:
        return getattr(cv2, var)
    except AttributeError:
        raise ValueError(f'no direct color conversion from {from_color.name} to {to_color.name}')


def get_rotate_flag(from_angle: int, to_angle: int) -> int:
    angle = to_angle - from_angle

    if angle < 0:
        angle += 360

    rotations = (angle % 360) / 90

    if rotations == 0:
        return -1
    elif rotations == 1:
        return cv2.ROTATE_90_CLOCKWISE
    elif rotations == 2:
        return cv2.ROTATE_180
    elif rotations == 3:
        return cv2.ROTATE_90_COUNTERCLOCKWISE

    raise ValueError('cannot rotate %d -> %d' % (from_angle, to_angle))


def build_transformer(source: Format, target: Format) -> Function:
    fns = list()

    if ColorMode.UNKNOWN not in [source.color_mode, target.color_mode]:
        fns.append(_get_color_transform(source.color_mode, target.color_mode))

    fns.append(_get_orientation_transformation(source.orientation, target.orientation))

    if source.width != target.width or source.height != target.height:
        fns.append(_cv2_resize_with_scale((target.width, target.height)))

    return Pipeline(*fns)


NoTransform = _no_transform
ColorTransform = _get_color_transform
OrientationTransform = _get_orientation_transformation
Pipeline = _pipeline
ResizeWithScale = _cv2_resize_with_scale
Transformer = build_transformer
