import json
import logging
import os
import tempfile
import time
import urllib.request as request
import zipfile

import cv2
import pkg_resources
from cogstream.api import EngineDescriptor, Specification
from cogstream.api.format import AnyFormat
from cogstream.engine import Engine, EngineResultWriter, Frame, EngineResult

mar_url = 'https://s3.amazonaws.com/model-server/model_archive_1.0/FERPlus.mar'
model_path = os.path.join(os.path.expanduser('~'), '.cache', 'cogstream', 'models', 'fermx')

logger = logging.getLogger(__name__)


def create_engine(metadata=None) -> Engine:
    engine = EmotionEngine()
    engine.setup()
    return engine


def loade_face_detector():
    path = pkg_resources.resource_filename('cv2', 'data/haarcascade_frontalface_default.xml')
    return cv2.CascadeClassifier(path)


def download_ferplus_model(target_dir):
    if not os.path.isdir(target_dir):
        logger.info('creating directory to store model %s', target_dir)
        os.makedirs(target_dir, exist_ok=True)

    with tempfile.TemporaryDirectory(suffix='_cogstream') as tmpdir:
        mar = os.path.join(tmpdir, 'FERPlus.mar')

        logger.info('downloading ferplus model %s -> %s', mar_url, mar)
        request.urlretrieve(mar_url, filename=mar)

        logger.info('extracting mar file %s -> %s', mar, target_dir)
        with zipfile.ZipFile(mar, 'r') as zip_ref:
            zip_ref.extract('FERPlus-0000.params', target_dir)
            zip_ref.extract('FERPlus-symbol.json', target_dir)
            zip_ref.extract('signature.json', target_dir)

        with open(os.path.join(target_dir, 'downloaded'), mode='w') as fd:
            fd.write(str(time.time()))
            fd.write('\n')


def _load_ferplus_service(model_dir) -> 'FERPlusService':
    from .model import FERPlusService, Context

    with open(os.path.join(os.path.dirname(__file__), 'model/MAR-INF/MANIFEST.json'), 'r') as fd:
        manifest = json.load(fd)

    ctx = Context('FERPlus', model_dir, manifest, batch_size=1, gpu=0, mms_version='0.0')
    service = FERPlusService()
    service.initialize(ctx)
    return service


def require_ferplus_service():
    if not os.path.exists(os.path.join(model_path, 'downloaded')):
        download_ferplus_model(target_dir=model_path)

    return _load_ferplus_service(model_path)


class EmotionEngine(Engine):

    def __init__(self):
        self._do_inference = None
        self.face_detector = None

    def setup(self):
        if self._do_inference is not None:
            return

        service = require_ferplus_service()

        def inference(data):
            return service.handle(data, service._context)

        self._do_inference = inference

        self.face_detector = loade_face_detector()

    def _do_inference(self, img):
        ...

    def get_descriptor(self) -> EngineDescriptor:
        return EngineDescriptor('fermx', Specification('analyze', AnyFormat))

    def process(self, frame: Frame, results: EngineResultWriter):
        gray = cv2.cvtColor(frame.image, cv2.COLOR_BGR2GRAY)

        faces = self.face_detector.detectMultiScale(
            gray,
            scaleFactor=1.2,
            minNeighbors=7,
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
            face_input = cv2.resize(face, (64, 64))
            emotions = self._do_inference(face_input)
            # cv2.imshow(f'emotion-{i}', face_input)
            result_payload.append({'face': (x, y, w, h), 'emotions': emotions[0]})

        results.write(EngineResult(frame.frame_id, time.time(), result_payload))
