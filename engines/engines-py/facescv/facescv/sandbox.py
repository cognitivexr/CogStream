import cv2
from cogstream.engine import EngineResultWriter, EngineResult, Frame

from facescv.engine import FacesEngine


def main():
    engine = FacesEngine()
    engine.setup()

    class LocalResultWriter(EngineResultWriter):

        last_result: EngineResult

        def write(self, result: EngineResult):
            self.last_result = result

    cap = cv2.VideoCapture(0)
    results = LocalResultWriter()

    try:
        while True:
            more, frame = cap.read()
            if not more:
                break

            engine.process(Frame(frame), results)

            faces = results.last_result

            if faces is not None and faces.result is not None:
                i = 0
                for (x, y, w, h) in faces.result:
                    i += 1
                    cv2.rectangle(frame, (x, y), (x + w, y + h), (0, 255, 0), 2)
                    face = frame[y:y + h, x:x + w]
                    cv2.imshow(f'face-{i}', face)

            cv2.imshow('faces', frame)

            if cv2.waitKey(1) & 0xFF == ord('q'):
                return
    finally:
        engine.close()
        cap.release()
        cv2.destroyAllWindows()


if __name__ == '__main__':
    main()
