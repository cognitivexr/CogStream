import argparse
import logging

import cv2

from cogstream.api import OperationSpec, ClientFormatSpec, StreamSpec, to_attributes
from cogstream.mediator.client import MediatorClient
from cogstream.engine.client import EngineClient, stream_camera


def main():
    logging.basicConfig(level=logging.INFO)
    parser = argparse.ArgumentParser(description='CogStream Client for recording videos')
    parser.add_argument('--host', type=str, help='the mediator address', default='127.0.0.1')
    parser.add_argument('--port', type=int, help='the mediator port (default 8191)', default=8191)

    parser.add_argument('--capture-width', type=int, help='camera capture height', default=800)
    parser.add_argument('--capture-height', type=int, help='camera capture width', default=600)
    parser.add_argument('--record-width', type=int, help='video height', default=640)
    parser.add_argument('--record-height', type=int, help='video width', default=480)

    args = parser.parse_args()

    # mediator = MediatorClient(args.host, args.port)

    # op_spec = OperationSpec('record', to_attributes({
    #     "format.width": args.record_width,
    #     "format.height": args.record_height,
    #     "codec": "xvid"
    # }))
    # available_engines = mediator.request_operation(op_spec)

    # if not available_engines.engines:
    #     print("error: no available engine")
    #     return -1

    # # todo: select for real
    # selection = available_engines.engines[0]

    # client_format = ClientFormatSpec(selection.name, to_attributes({
    #     "format.width": args.capture_width,
    #     "format.height": args.capture_height,
    #     "format.colorMode": "RGB",
    # }))

    # stream_spec = mediator.establish_format(client_format)
    # mediator.close()
    
    stream_spec = StreamSpec('127.0.0.1:53210', to_attributes({
        "format.width": args.record_width,
        "format.height": args.record_height,
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

        stream_camera(cap, engine_client)
    except KeyboardInterrupt:
        pass
    finally:
        cap.release()
        engine_client.close()


if __name__ == '__main__':
    main()
