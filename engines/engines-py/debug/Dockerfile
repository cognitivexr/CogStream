FROM ubuntu:20.04

RUN apt-get update

# configure tzdata (which otherwise involves interaction via stdin)
RUN export DEBIAN_FRONTEND=noninteractive
RUN apt-get install -y tzdata
RUN ln -fs /usr/share/zoneinfo/America/New_York /etc/localtime
RUN dpkg-reconfigure --frontend noninteractive tzdata

# install base engine dependencies
RUN apt-get install -y \
    python3 \
    python3-pip \
    python3-opencv

COPY ./requirements.txt ./requirements.txt
RUN sed -i "s/..\/..\/..\///g" requirements.txt
RUN pip3 install -r requirements.txt

# default engine port
EXPOSE 54321

# engine code
COPY debug ./debug
CMD python3 -m debug.main
