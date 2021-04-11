import logging

from cogstream.engine import serve_engine

from yolov5.engine import Yolov5


def main():
    logging.basicConfig(level=logging.INFO)
    serve_engine(('0.0.0.0', 45671), lambda metadata: Yolov5())


if __name__ == '__main__':
    main()
