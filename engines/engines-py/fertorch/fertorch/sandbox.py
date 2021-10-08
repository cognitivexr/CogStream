import logging

import cv2
from cogstream.engine import EngineResultWriter, EngineResult, Frame

from fermx.engine import create_engine


class LocalResultWriter(EngineResultWriter):
    last_result: EngineResult = None

    def write(self, result: EngineResult):
        self.last_result = result


def main():
    logging.basicConfig(level=logging.INFO)
    engine = create_engine()
    engine.setup()

    cap = cv2.VideoCapture(0)
    cap.set(cv2.CAP_PROP_FRAME_WIDTH, 1280)
    cap.set(cv2.CAP_PROP_FRAME_HEIGHT, 720)
    results = LocalResultWriter()

    try:
        while True:
            more, frame = cap.read()
            if not more:
                break

            engine.process(Frame(frame), results)

            result = results.last_result

            if result is not None:
                i = 0

                for item in result.result:
                    i += 1
                    x, y, h, w = item['face']
                    emotion = item['emotions'][0]
                    cv2.rectangle(frame, (x, y), (x + w, y + h), (0, 255, 0), 2)
                    face = frame[y:y + h, x:x + w]
                    cv2.imshow(f'face-{i}', face)

                    label = emotion['label']
                    cv2.putText(frame, label, (x, y - 5), 0, 1, (20, 255, 20), thickness=2, lineType=cv2.LINE_AA)

                cv2.imshow('result', frame)

            if cv2.waitKey(1) & 0xFF == ord('q'):
                return
    finally:
        engine.close()
        cap.release()
        cv2.destroyAllWindows()


if __name__ == '__main__':
    main()
