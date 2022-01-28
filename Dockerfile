# syntax=docker/dockerfile:1

##
## Build go
##
#FROM golang:1.16-buster AS build
#
#WORKDIR /app
#
# ....
#RUN make engines-go


##
## Build python
##
FROM python:3.8-slim-buster as build-py

RUN apt-get update && apt-get install make

WORKDIR /engines

COPY . .

# RUN make engines-py

WORKDIR /engines/engines-py/debug

RUN make install

CMD ["make", "start" ]

##
## Deploy
##
#FROM python:3.8-slim-buster
#
#WORKDIR /
#
#COPY --from=build-py ...
#
#
#ENTRYPOINT [""] ... mediator

