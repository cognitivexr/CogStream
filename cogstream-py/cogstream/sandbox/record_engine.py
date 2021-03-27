import logging

import cv2.cv2 as cv2

from cogstream.api import EngineDescriptor, StreamMetadata, Specification
from cogstream.api.format import Format, ColorMode, Orientation
from cogstream.engine import serve_engine, Engine, Frame


class InputFormat(object):
    pass


class Recorder(Engine):

    def __init__(self) -> None:
        super().__init__()
        self.format = Format(800, 600, ColorMode(1), Orientation(1))
        self.out = None

    def get_descriptor(self) -> EngineDescriptor:
        return EngineDescriptor('record-py', Specification('record', self.format))

    def setup(self):
        self.out = cv2.VideoWriter('/tmp/after.mp4', cv2.VideoWriter_fourcc('m', 'p', '4', 'v'), 25, (800, 600))

    def process(self, frame: Frame):
        self.out.write(frame.image)

    def close(self):
        self.out.release()


def main():
    def create_recorder(_: StreamMetadata):
        return Recorder()

    logging.basicConfig(level=logging.DEBUG)

    serve_engine(("0.0.0.0", 54322), create_recorder)


if __name__ == '__main__':
    main()
