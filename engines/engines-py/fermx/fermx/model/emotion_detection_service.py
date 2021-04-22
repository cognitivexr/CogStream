# Copyright 2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.
# Licensed under the Apache License, Version 2.0 (the "License").
# You may not use this file except in compliance with the License.
# A copy of the License is located at
#     http://www.apache.org/licenses/LICENSE-2.0
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.
import logging
import time

import mxnet as mx
import numpy as np

from .mxnet_model_service import MXNetModelService
from .mxnet_utils import ndarray

logger = logging.getLogger(__name__)


class FERPlusService(MXNetModelService):
    """
    FERPlus implements emotion detection on facial images.
    The underlying model is based on a 2016 paper by Barsoum et al. (https://arxiv.org/abs/1608.01041)
    It was built and trained using Microsoft Cognitive Toolkit (CNTK),
    and the code is available on the FERPlus repository (https://github.com/Microsoft/FERPlus)
    """

    def preprocess(self, data):
        then = time.time()

        input_image = data
        # # opencv
        # input_image = cv2.resize(data, (64, 64))
        # input_image = cv2.cvtColor(data, gray)
        # # pillow
        # input_image = input_image.resize((64, 64))
        # input_image = input_image.convert('L')

        # input_array = np.asarray(input_image.getdata(), dtype=np.float64)
        input_array = np.asarray(input_image, dtype=np.float64)
        input_array = (input_array - 127.5) / 127.5

        if logger.isEnabledFor(logging.DEBUG):
            logging.debug('preprocess took %.2f', ((time.time() - then) * 1000))

        input_tensor = mx.nd.reshape(mx.nd.array(input_array), shape=(1, 1, 64, 64))
        return [input_tensor]

    def postprocess(self, data):
        """
        Post-process the inference result to take through a softmax, sort by
        probability and return classes mapped to probabilities.
        """

        emotion_classes = ['Neutral', 'Happy', 'Surprise', 'Sad', 'Anger', 'Disgust', 'Fear', 'Contempt']
        softmax_output = data[0].softmax()

        return [ndarray.top_probability(softmax_output, emotion_classes)]


_service = FERPlusService()


def handle(data, context):
    if not _service.initialized:
        _service.initialize(context)

    if data is None:
        return None

    return _service.handle(data, context)
