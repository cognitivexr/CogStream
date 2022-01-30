FROM golang

WORKDIR /cogstream

COPY mediator/ /cogstream/mediator
COPY api/ /cogstream/api

RUN cd mediator; go build ./cmd/mediator
RUN cp ./mediator/mediator /usr/bin
RUN mkdir /cogstream/engines

EXPOSE 8191

CMD mediator -plugins /cogstream/engines
