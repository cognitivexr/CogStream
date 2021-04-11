import os.path
from unittest import TestCase

import cv2
from cogstream.engine import EngineResultWriter, EngineResult, Frame

from yolov5.engine import Yolov5


class Results(EngineResultWriter):

    def __init__(self) -> None:
        super().__init__()
        self.items = list()

    def write(self, result: EngineResult):
        self.items.append(result)


class TestYolov5(TestCase):
    def test(self):
        engine = Yolov5()
        engine.setup()

        results = Results()

        image = cv2.imread(os.path.join(os.path.dirname(__file__), 'kitten.jpg'))
        image = cv2.cvtColor(image, cv2.COLOR_BGR2RGB)
        frame = Frame(image=image)

        engine.process(frame, results)

        self.assertEqual(1, len(results.items))
        self.assertEqual('cat', results.items[0].result[0]['label'])

        # frame_inference = engine.inference_visualize(frame.image)
        # frame_inference = cv2.cvtColor(frame_inference, cv2.COLOR_RGB2BGR)
        # cv2.imwrite('kitten_inference.jpg', frame_inference)
