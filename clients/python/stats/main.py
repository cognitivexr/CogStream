import argparse
import logging

import cv2

from cogstream.api import OperationSpec, ClientFormatSpec, to_attributes, StreamSpec
from cogstream.mediator.client import MediatorClient
from cogstream.engine.client import EngineClient, stream_camera
from cogstream.typing import deep_from_dict


def main():
    stream_spec = deep_from_dict({"engineAddress": "127.0.0.1:53210", "attributes": {}}, StreamSpec)
    engine_client = EngineClient(stream_spec)
    engine_client.open()

    if not engine_client.acknowledged:
        print('engine stream could not be initiated')
        engine_client.close()
        return

    cap = cv2.VideoCapture("/home/vader/Videos/sample_1920x1080.mp4")

    try:
        cap.set(cv2.CAP_PROP_FRAME_WIDTH, 1920)
        cap.set(cv2.CAP_PROP_FRAME_HEIGHT, 1080)

        stream_camera(cap, engine_client)
    except KeyboardInterrupt:
        pass
    finally:
        cap.release()
        engine_client.close()


if __name__ == '__main__':
    main()
