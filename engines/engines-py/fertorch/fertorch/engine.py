import json
import logging
import os
import tempfile
import time
import urllib.request as request

import cv2
import pkg_resources
from cogstream.api import EngineDescriptor, Specification
from cogstream.api.format import Format, ColorMode, Orientation
from cogstream.engine import Engine, EngineResultWriter, Frame, EngineResult
from FERModel import classes, FERModel, get_default_device
from torchvision import transforms
import torch


logger = logging.getLogger(__name__)

model = FERModel(1, 7)

softmax = torch.nn.Softmax(dim=1)

def create_engine(metadata=None) -> Engine:
    engine = EmotionEngine()
    engine.setup()
    return engine

def img2tensor(x):
    transform = transforms.Compose(
            [transforms.ToTensor(),
             transforms.Normalize((0.5), (0.5))])
    return transform(x)

def predict(x):
    out = model(x[None])
    scaled = softmax(out).cpu().detach().numpy()
    return [{"label": classes[i], "probability":p} for (i,p) in enumerate(scaled[0])]


class EmotionEngine(Engine):

    def __init__(self):
        self.face_detector = None

    def setup(self):

        self.face_detector = cv2.CascadeClassifier('/home/nvidia/CognitiveXR/FER-Pytorch/haarcascade_frontalface_default.xml')


    def get_descriptor(self) -> EngineDescriptor:
        return EngineDescriptor('fertorch', Specification('analyze', Format(
            color_mode=ColorMode.Gray,
            orientation=Orientation.TopLeft
        )))

    def process(self, frame: Frame, results: EngineResultWriter):
        gray = frame.image

        #cv2.imshow('faces', gray)
        #key = cv2.waitKey(1)
        #if key == ord('q'):
        #    raise KeyboardInterrupt

        faces = self.face_detector.detectMultiScale(
            gray,
            scaleFactor=1.2,
            minNeighbors=3,
            minSize=(30, 30),
            flags=cv2.CASCADE_SCALE_IMAGE
        )

        if len(faces) == 0:
            results.write(EngineResult(frame.frame_id, time.time(), []))
            return

        result_payload = []

        i = 0
        for (x, y, w, h) in faces:
            i += 1

            face = gray[y:y + h, x:x + w]
            face_input = cv2.resize(face, (48, 48))
            emotions = predict(img2tensor(face_input))
            labels = [{'probability': e['probability'], 'label': e['label']} for e in emotions]
            result_payload.append({'face': (x, y, w, h), 'emotions': labels})

        results.write(EngineResult(frame.frame_id, time.time(), result_payload))
