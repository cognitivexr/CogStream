# Facial expression detection wit OpenCV and MXNet

This CogStream engine performs facial expression detection using OpenCV and an MXNet model trained on the FER+ dataset.
It detects each face in the frame and performs emotion detection on individual faces instead of the entire frame.

## Enable GPU Support

To use the MXNet GPU support you need CUDA, cuDNN and the mxnet cuda pip package.
Follow the [instructions to install cuDNN](https://docs.nvidia.com/deeplearning/cudnn/archives/cudnn_765/cudnn-install/index.html#install-linux).
As of April 2021, the cuDNN version needed was v7.
Make sure the CUDA version matches the CUDA version in the pip package of,
e.g., `mxnet-cu101` for cuda version 10.1.
There seems to be an [issue with mxnet 1.8.0](https://stackoverflow.com/questions/66786887/getting-oserror-libnccl-so-2-while-importing-mxnet),
which is solved by installing 1.7 with `pip install mxnet-cu101==1.7`, 

```bash
cat /usr/local/cuda/version.txt
# OR
nvcc --version

# mxnet
pip install mxnet-cu<version>
```

## Copyright

The copyright of the packaged MXNet code (`fermx/model`) is held by Amazon.com under the Apache License, Version 2.0.
