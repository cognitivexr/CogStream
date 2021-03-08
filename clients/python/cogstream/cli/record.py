import argparse
import logging
import time

import cv2
import numpy as np

from cogstream.client import Client

logger = logging.getLogger(__name__)


def stream_camera(cap, client, show=True):
    goal_fps = 25

    # target frame inter-arrival time
    ia = 1 / goal_fps

    while True:
        start = time.time()

        check, frame = cap.read()
        if not check:
            logger.info('no more frames to read')
            break

        if show:
            cv2.imshow("capture", frame)

        jpg: np.ndarray = cv2.imencode('.jpg', frame)[1]

        client.request(jpg)

        delay = ia - (time.time() - start)
        if delay >= 0:
            time.sleep(delay)

        key = cv2.waitKey(1)
        if key == ord('q'):
            break

        logger.info('fps: %.2f' % (1 / (time.time() - start)))


def main():
    # FIXME: this currently connects directly to the engine. should do the handshake first

    parser = argparse.ArgumentParser(description='CogStream Client')
    parser.add_argument('--host', type=str, help='the address to connect to', default='127.0.0.1')
    parser.add_argument('--port', type=int, help='the port to connect to (default 53210)', default=53210)

    parser.add_argument('--height', type=int, help='camera capture height', default=640)
    parser.add_argument('--width', type=int, help='camera capture width', default=360)

    logging.basicConfig(level=logging.INFO)

    args = parser.parse_args()

    address = (args.host, args.port)

    client = Client(address)

    # FIXME perform handshake with record operation
    width = args.width
    height = args.height

    agreement = {"sessionId": "12345", "config": {"width": width, "height": height, "colorMode": 1}}
    client.open(agreement)

    if not client.handshake:
        print('handshake failed')
        client.close()
        return

    cap = cv2.VideoCapture(0)

    try:
        cap.set(cv2.CAP_PROP_FRAME_WIDTH, width)
        cap.set(cv2.CAP_PROP_FRAME_HEIGHT, height)

        stream_camera(cap, client)
    except KeyboardInterrupt:
        pass
    finally:
        cap.release()
        client.close()


if __name__ == '__main__':
    main()
