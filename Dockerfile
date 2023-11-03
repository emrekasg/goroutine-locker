FROM golang:alpine

VOLUME /cpuaffinity

RUN apk add --no-cache bash

RUN apk add --no-cache wget git sudo gcc build-base

ENV GOPATH /go

CMD ["/bin/bash", "-c", "while true; do sleep 100000; done"]