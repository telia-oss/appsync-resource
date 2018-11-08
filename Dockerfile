FROM golang:1.11 as builder

LABEL MAINTAINER="telia-oss"

ADD . /go/src/github.com/telia-oss/appsync-resource

WORKDIR /go/src/github.com/telia-oss/appsync-resource

ENV TARGET linux
ENV ARCH amd64

RUN make build

FROM alpine/git:latest as resource
COPY --from=builder /go/src/github.com/telia-oss/appsync-resource/out /opt/resource/out
RUN chmod +x /opt/resource/*

FROM resource
