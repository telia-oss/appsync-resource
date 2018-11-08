FROM golang:1.11 as builder

LABEL MAINTAINER="telia-oss"

ADD . /go/src/github.com/telia-oss/appsync-resource

WORKDIR /go/src/github.com/telia-oss/appsync-resource

ENV TARGET linux
ENV ARCH amd64

RUN make build

COPY out /opt/resource/out

RUN chmod +x /opt/resource/out
