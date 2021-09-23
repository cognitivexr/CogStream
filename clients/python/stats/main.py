import argparse
import logging

import cv2

from cogstream.api import OperationSpec, ClientFormatSpec, to_attributes, StreamSpec
from cogstream.mediator.client import MediatorClient
from cogstream.engine.client import EngineClient, stream_camera
from cogstream.typing import deep_from_dict


def main():
    stream_spec = StreamSpec('127.0.0.1:53210', to_attributes({
        "format.width": 1920,
        "format.height": 1080,
        "format.orientation": 1,
        "format.colorMode": 1
    }))
    engine_client = EngineClient(stream_spec)
    engine_client.open()

    if not engine_client.acknowledged:
        print('engine stream could not be initiated')
        engine_client.close()
        return

    cap = cv2.VideoCapture("/home/silvio/videos/sample/1920x1080.mp4")

    if not cap.isOpened():
        raise Exception("Error opening video stream or file")

    try:
        stream_camera(cap, engine_client)
    except KeyboardInterrupt:
        pass
    finally:
        cap.release()
        engine_client.close()


if __name__ == '__main__':
    main()
