# syntax=docker/dockerfile:1

FROM debian:buster

RUN apt-get update

# set noninteractive installation
ENV DEBIAN_FRONTEND "noninteractive"
# install tzdata package
RUN apt-get install -y tzdata
# set your timezone
RUN ln -fs /usr/share/zoneinfo/America/New_York /etc/localtime
RUN dpkg-reconfigure --frontend noninteractive tzdata

RUN apt-get install -y \
    bzip2 \
    g++ \
    git \
    libgl1-mesa-glx \
    libhdf5-dev \
    openmpi-bin \
    wget \
    python3 \
    python3-dev \
    python3-pip \
    python3-tk \
    python3-opencv

RUN wget https://go.dev/dl/go1.17.6.linux-amd64.tar.gz

RUN tar -xvf go1.17.6.linux-amd64.tar.gz 

RUN mv go /usr/local

ENV GOROOT "/usr/local/go"

ENV GOPATH "$HOME/go"

ENV PATH "$GOPATH/bin:$GOROOT/bin:$PATH"

RUN go version

RUN go env