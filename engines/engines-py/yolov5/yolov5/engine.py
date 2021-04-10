import time
from typing import Any

import torch
from cogstream.api import EngineDescriptor, Specification
from cogstream.api.format import Format, ColorMode
from cogstream.engine import Engine, Frame


class Yolov5(Engine):

    def __init__(self, device=None) -> None:
        super().__init__()

        self.device = device
        self.model = None
        self.names = []
        self.colors = []

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

    def process(self, frame: Frame) -> Any:
        self.inference(frame.image)
        pass

    def setup(self):
        if self.device is None:
            self.device = torch.device("cuda") if torch.cuda.is_available() else torch.device("cpu")

        self.model = torch.hub.load('ultralytics/yolov5', 'yolov5s', pretrained=True).autoshape()
        self.model.to(self.device)
        self.names = self.model.module.names if hasattr(self.model, 'module') else self.model.names

    @torch.no_grad()
    def inference(self, frame_rgb):
        results = self.model(frame_rgb, 320 + 32 * 4)  # includes NMS
        return results

    def visualize(self, detections: 'models.common.Detections', frame, conf_th=None):
        import cv2

        tl = 2
        tf = 1
        white = (255, 255, 255)
        black = (0, 0, 0)

        points = detections.xyxy[0].cpu().numpy()

        for point in points:
            xyxy, conf, cls = point[:4], point[4], int(point[5])

            if conf_th is not None and conf < conf_th:
                continue

            c1, c2 = (int(xyxy[0]), int(xyxy[1])), (int(xyxy[2]), int(xyxy[3]))

            print(f'{time.time():0.4f},{detections.names[cls]},{conf:.4f},{c1[0]},{c1[1]},{c2[0]},{c2[1]}')

            # draw rectangle
            cv2.rectangle(frame, c1, c2, self.colors[cls], thickness=tl, lineType=cv2.LINE_AA)

            # draw label
            label = f'{detections.names[cls]} ({conf * 100:.0f}%)'
            t_size = cv2.getTextSize(label, 0, fontScale=tl / 3, thickness=tf)[0]
            c3 = c1[0] + t_size[0], c1[1] - t_size[1] - 3
            cv2.rectangle(frame, c1, c3, black, - 1, cv2.LINE_AA)  # label background
            cv2.putText(frame, label, (c1[0], c1[1] - 2), 0, tl / 3, white, thickness=tf, lineType=cv2.LINE_AA)

        return frame
