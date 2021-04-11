import logging
import time
from typing import Any

import torch
from cogstream.api import EngineDescriptor, Specification
from cogstream.api.format import Format, ColorMode
from cogstream.engine import Engine, Frame, EngineResultWriter, EngineResult

logger = logging.getLogger(__name__)


class Yolov5(Engine):

    def __init__(self, device=None, preload=True) -> None:
        super().__init__()

        self.device = device
        self.model = None
        self.names = []
        self.colors = []

        if preload:
            self.setup()

    def get_descriptor(self) -> EngineDescriptor:
        return EngineDescriptor(
            name='yolov5-py',
            specification=Specification(
                operation='analyze',
                input_format=Format(640, 480, ColorMode.RGB),
                attributes={
                    'input_scale': ['int'],
                    'cuda': ['bool']
                }
            )
        )

    def process(self, frame: Frame, writer: EngineResultWriter) -> Any:
        detections = self.inference(frame.image)
        # TODO: to engine res

        objects = list()
        points = detections.xyxy[0].cpu().numpy()

        for point in points:
            xyxy, conf, cls = point[:4], point[4], int(point[5])

            objects.append({'xyxy': list(xyxy), 'conf': conf, 'label': self.names[cls]})

        logger.debug('writing inference result %s', objects)
        writer.write(EngineResult(frame.frame_id, time.time(), objects))

    def setup(self):
        if self.model is not None:
            return

        if self.device is None:
            self.device = torch.device("cuda") if torch.cuda.is_available() else torch.device("cpu")

        self.model = torch.hub.load('ultralytics/yolov5', 'yolov5s', pretrained=True).autoshape()
        self.model.to(self.device)
        self.names = self.model.module.names if hasattr(self.model, 'module') else self.model.names

    @torch.no_grad()
    def inference(self, frame_rgb) -> 'models.common.Detections':
        results = self.model(frame_rgb, 320 + 32 * 4)  # includes NMS
        return results

    def inference_visualize(self, frame_rgb) -> 'np.ndarray':
        import cv2
        frame = frame_rgb

        detections = self.inference(frame)

        tl = 2
        tf = 1
        white = (255, 255, 255)
        black = (0, 0, 0)

        points = detections.xyxy[0].cpu().numpy()

        for point in points:
            xyxy, conf, cls = point[:4], point[4], int(point[5])

            c1, c2 = (int(xyxy[0]), int(xyxy[1])), (int(xyxy[2]), int(xyxy[3]))

            print(f'{time.time():0.4f},{detections.names[cls]},{conf:.4f},{c1[0]},{c1[1]},{c2[0]},{c2[1]}')

            # draw rectangle
            cv2.rectangle(frame, c1, c2, 255, thickness=tl, lineType=cv2.LINE_AA)

            # draw label
            label = f'{detections.names[cls]} ({conf * 100:.0f}%)'
            t_size = cv2.getTextSize(label, 0, fontScale=tl / 3, thickness=tf)[0]
            c3 = c1[0] + t_size[0], c1[1] - t_size[1] - 3
            cv2.rectangle(frame, c1, c3, black, - 1, cv2.LINE_AA)  # label background
            cv2.putText(frame, label, (c1[0], c1[1] - 2), 0, tl / 3, white, thickness=tf, lineType=cv2.LINE_AA)

        return frame
