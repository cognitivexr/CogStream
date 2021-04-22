import argparse
import logging

from cogstream.engine import serve_engine

from facescv.engine import create_engine


def main():
    logging.basicConfig(level=logging.INFO)

    parser = argparse.ArgumentParser(description='CogStream engine for face detection using OpenCV')
    parser.add_argument('--host', type=str, help='the host to bind to', default='0.0.0.0')
    parser.add_argument('--port', type=int, help='the port to bind to (default random port)', default=54321)

    args = parser.parse_args()

    serve_engine((args.host, args.port), create_engine)


if __name__ == '__main__':
    main()
