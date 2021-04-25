import argparse
import logging

import cv2
from cogstream.api import StreamSpec, to_attributes
from cogstream.engine import EngineResult
from cogstream.engine.client import EngineClient, stream_camera


def on_result(frame, result: EngineResult):
    """
    Main client code. Uses the EngineResult of the fermx engine to display a bounding box around the captured faces and
    display their emotion.

    :param frame: the frame capture from the camera
    :param result: the EngineResult
    """

    if result is not None:
        i = 0

        for item in result.result:
            i += 1
            x, y, h, w = item['face']
            emotion = item['emotions'][0]
            cv2.rectangle(frame, (x, y), (x + w, y + h), (0, 255, 0), 2)
            # face = frame[y:y + h, x:x + w]
            # cv2.imshow(f'face-{i}', face)

            label = emotion['class']
            cv2.putText(frame, label, (x, y - 5), 0, 1, (20, 255, 20), thickness=2, lineType=cv2.LINE_AA)

    cv2.imshow('faces', frame)

    key = cv2.waitKey(1)
    if key == ord('q'):
        raise KeyboardInterrupt


def main():
    logging.basicConfig(level=logging.INFO)
    parser = argparse.ArgumentParser(description='CogStream client for fermx')

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
        stream_camera(cap, engine_client, show=False, on_result=on_result)
    except KeyboardInterrupt:
        pass
    finally:
        cap.release()
        engine_client.close()


if __name__ == '__main__':
    main()
