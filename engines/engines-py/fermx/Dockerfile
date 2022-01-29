FROM ubuntu:20.04

RUN apt-get update

# configure tzdata (which otherwise involves interaction via stdin)
RUN export DEBIAN_FRONTEND=noninteractive
RUN apt-get install -y tzdata
RUN ln -fs /usr/share/zoneinfo/America/New_York /etc/localtime
RUN dpkg-reconfigure --frontend noninteractive tzdata

# install base engine dependencies
RUN apt-get install -y \
    bzip2 \
    g++ \
    git \
    wget \
    python3 \
    python3-dev \
    python3-pip \
    python3-tk \
    python3-opencv

RUN pip install cogstream-py mxnet

RUN apt-get autoremove -y && apt-get autoclean -y

# default engine port
EXPOSE 54321

# engine code
COPY fermx ./fermx
CMD python3 -m fermx.main
