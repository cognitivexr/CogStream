import argparse

import cv2

from cogstream.client import MediatorClient, EngineClient, stream_camera


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

    engine_spec = mediator.request_operation({
        "code": "record",
        "attributes": {
            "foo": [
                "bar",
                "baz"
            ]
        }
    })
    print(engine_spec)

    # todo: handle engine specs and use them for establishing format

    stream_spec = mediator.establish_format({
        "attributes": {
            "la": [
                "le",
                "lu"
            ]
        }
    })
    print(stream_spec)
    mediator.close()

    # todo: use agreement for connection
    address = stream_spec["content"]["engineAddress"].split(":")
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
