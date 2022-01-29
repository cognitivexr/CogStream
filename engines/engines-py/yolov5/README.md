# Object Detection with PyTorch and YOLOv5

This CogStream engine performs object detection using pre-trained YOLOv5 models on PyTorch.

## Build Docker images

### CPU

    docker build -f Dockerfile -t cognitivexr/engine-yolov5:cpu .

### GPU acceleration (CUDA 11.0)

    docker build -f Dockerfile.cuda110 -t cognitivexr/engine-yolov5:cuda-110 .
