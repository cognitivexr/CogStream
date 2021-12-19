import argparse
import logging
import time

import cv2
from cogstream.api import StreamSpec, to_attributes
from cogstream.engine import EngineResult
from cogstream.engine.client import EngineClient, stream_camera


def show_and_draw(frame, result: EngineResult, conf_th=0.3):
    tl = 2
    tf = 1
    white = (255, 255, 255)
    black = (0, 0, 0)
    green = (0, 255, 0)

    for point in result.result:
        xyxy, conf, label = point['xyxy'], point['conf'], point['label']

        if conf_th is not None and conf < conf_th:
            continue

        c1, c2 = (int(xyxy[0]), int(xyxy[1])), (int(xyxy[2]), int(xyxy[3]))

        print(f'{time.time():0.4f},{label},{conf:.4f},{c1[0]},{c1[1]},{c2[0]},{c2[1]}')

        # draw rectangle
        cv2.rectangle(frame, c1, c2, green, thickness=tl, lineType=cv2.LINE_AA)

        # draw label
        label = f'{label} ({conf * 100:.0f}%)'
        t_size = cv2.getTextSize(label, 0, fontScale=tl / 3, thickness=tf)[0]
        c3 = c1[0] + t_size[0], c1[1] - t_size[1] - 3
        cv2.rectangle(frame, c1, c3, black, - 1, cv2.LINE_AA)  # label background
        cv2.putText(frame, label, (c1[0], c1[1] - 2), 0, tl / 3, white, thickness=tf, lineType=cv2.LINE_AA)

    cv2.imshow("show", frame)

    key = cv2.waitKey(1)
    if key == ord('q'):
        exit(0)


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

        stream_camera(cap, engine_client, show=False, on_result=show_and_draw)
    except KeyboardInterrupt:
        pass
    finally:
        cap.release()
        engine_client.close()


if __name__ == '__main__':
    main()
