import logging
import time

import cv2
import pkg_resources
from cogstream.api import EngineDescriptor, Specification
from cogstream.api.format import AnyFormat
from cogstream.engine import Engine, EngineResultWriter, EngineResult

logger = logging.getLogger(__name__)


def create_engine(metadata) -> Engine:
    return FacesEngine()


def load_classifier():
    path = pkg_resources.resource_filename('cv2', 'data/haarcascade_frontalface_default.xml')
    return cv2.CascadeClassifier(path)


class FacesEngine(Engine):
    classifier: cv2.CascadeClassifier

    def get_descriptor(self) -> EngineDescriptor:
        return EngineDescriptor('facescv', Specification('analyze', AnyFormat))

    def setup(self):
        self.classifier = load_classifier()

    def process(self, frame: 'Frame', results: EngineResultWriter):
        gray = cv2.cvtColor(frame.image, cv2.COLOR_BGR2GRAY)

        faces = self.classifier.detectMultiScale(
            gray,
            scaleFactor=1.2,
            minNeighbors=5,
            minSize=(30, 30),
            flags=cv2.CASCADE_SCALE_IMAGE
        )
        # faces is either an empty tuple, or a numpy array of bounding box tuples: (x, y, w, h)

        payload = faces.tolist() if len(faces) else []

        result = EngineResult(frame.frame_id, time.time(), payload)
        results.write(result)
