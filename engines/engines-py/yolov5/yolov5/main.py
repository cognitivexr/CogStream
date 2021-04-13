import argparse
import logging

from cogstream.engine import serve_engine

from yolov5.engine import Yolov5


def main():
    logging.basicConfig(level=logging.INFO)

    parser = argparse.ArgumentParser(description='CogStream engine for object detection with PyTorch and YOLOv5')
    parser.add_argument('--host', type=str, help='the host to bind to', default='0.0.0.0')
    parser.add_argument('--port', type=int, help='the port to bind to (default 54321)', default=54321)

    args = parser.parse_args()

    engine = Yolov5()
    engine.setup()

    serve_engine((args.host, args.port), lambda metadata: engine)


if __name__ == '__main__':
    main()
