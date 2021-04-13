import argparse
import logging

import cv2
from cogstream.api import StreamSpec, to_attributes
from cogstream.engine.client import EngineClient, stream_camera


def main():
    logging.basicConfig(level=logging.INFO)
    parser = argparse.ArgumentParser(description='CogStream Client for YOLOv5 engine')

    parser.add_argument('--capture-width', type=int, help='camera capture height', default=800)
    parser.add_argument('--capture-height', type=int, help='camera capture width', default=600)

    args = parser.parse_args()

    stream_spec = StreamSpec('127.0.0.1:54321', to_attributes({
        "format.width": args.capture_width,
        "format.height": args.capture_height,
        "format.orientation": 1,
        "format.colorMode": 1
    }))

    engine_client = EngineClient(stream_spec)
    engine_client.open()

    if not engine_client.acknowledged:
        print('engine stream could not be initiated')
        engine_client.close()
        return

    cap = cv2.VideoCapture(0)

    try:
        cap.set(cv2.CAP_PROP_FRAME_WIDTH, args.capture_width)
        cap.set(cv2.CAP_PROP_FRAME_HEIGHT, args.capture_height)

        stream_camera(cap, engine_client, show=False)
    except KeyboardInterrupt:
        pass
    finally:
        cap.release()
        engine_client.close()


if __name__ == '__main__':
    main()
