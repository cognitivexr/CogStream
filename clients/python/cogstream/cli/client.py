import argparse

import cv2

from cogstream.client import MediatorClient, EngineClient, stream_camera, OperationSpec, ClientFormatSpec


def main():
    parser = argparse.ArgumentParser(description='CogStream Client')
    parser.add_argument('--host', type=str, help='the address to connect to', default='127.0.0.1')
    parser.add_argument('--port', type=int, help='the port to connect to (default 8191)', default=8191)
    parser.add_argument('--operation', type=str, help='the operation type', required=True)

    args = parser.parse_args()

    mediator = MediatorClient(args.host, args.port)

    # todo: remove this eventually
    if args.operation != 'record':
        return -1

    op_spec = OperationSpec(args.operation, {"format.width": ["640"], "format.height": ["360"], "codec": ["xvid"]})
    available_engines = mediator.request_operation(op_spec)

    if not available_engines.engines:
        return -1

    # todo: select for real

    selection = available_engines.engines[0]

    client_format = ClientFormatSpec(selection.name,
                                     {"format.width": ["640"], "format.height": ["360"], "codec": ["xvid"]})
    stream_spec = mediator.establish_format(client_format)
    print(stream_spec)
    mediator.close()

    # todo: use agreement for connection
    address = stream_spec.engineAddress.split(":")
    tup_addr = (address[0], int(address[1]))

    engine_client = EngineClient(tup_addr)
    agreement = {"sessionId": "12345", "config": {"width": 640, "height": 360, "colorMode": 1}}
    engine_client.open(agreement)

    if not engine_client.handshake:
        print('handshake failed')
        engine_client.close()
        return

    cap = cv2.VideoCapture(0)

    try:
        cap.set(cv2.CAP_PROP_FRAME_WIDTH, 640)
        cap.set(cv2.CAP_PROP_FRAME_HEIGHT, 360)

        stream_camera(cap, engine_client)
    except KeyboardInterrupt:
        pass
    finally:
        cap.release()
        engine_client.close()


if __name__ == '__main__':
    main()
