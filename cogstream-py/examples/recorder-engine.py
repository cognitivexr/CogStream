import logging

import cv2.cv2 as cv2

from cogstream.api import EngineDescriptor, StreamMetadata, Specification
from cogstream.api.format import Format, Orientation, ColorMode
from cogstream.engine import serve_engine, Engine, Frame


class Recorder(Engine):

    def __init__(self, target_file, recording_format: Format, recording_fps: int, codec='mp4v') -> None:
        super().__init__()
        self.target_file = target_file
        self.format = recording_format
        self.fps = recording_fps
        self.codec = codec
        self.out = None

    def get_descriptor(self) -> EngineDescriptor:
        return EngineDescriptor('record-py', Specification('record', self.format))

    def setup(self):
        dim = (self.format.width, self.format.height)
        self.out = cv2.VideoWriter(self.target_file, cv2.VideoWriter_fourcc(*self.codec), self.fps, dim)

    def process(self, frame: Frame):
        self.out.write(frame.image)

    def close(self):
        self.out.release()


def main():
    address = ("0.0.0.0", 54322)
    recording_format = Format(800, 600, ColorMode.RGB, Orientation.TopLeft)
    recording_fps = 22
    file = '/tmp/recorder-engine.mp4'

    def engine_factory(_: StreamMetadata) -> Engine:
        return Recorder(target_file=file, recording_format=recording_format, recording_fps=recording_fps)

    logging.basicConfig(level=logging.DEBUG)
    serve_engine(address, engine_factory)


if __name__ == '__main__':
    main()
