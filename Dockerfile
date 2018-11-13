FROM golang:1.11 as builder

LABEL MAINTAINER="telia-oss"

ADD . /go/src/github.com/telia-oss/appsync-resource

WORKDIR /go/src/github.com/telia-oss/appsync-resource

ENV TARGET linux
ENV ARCH amd64

RUN make build

FROM alpine:3.8 as resource
RUN apk update && apk add ca-certificates
COPY --from=builder /go/src/github.com/telia-oss/appsync-resource/bin/check /opt/resource/check
COPY --from=builder /go/src/github.com/telia-oss/appsync-resource/bin/in /opt/resource/in
COPY --from=builder /go/src/github.com/telia-oss/appsync-resource/bin/out /opt/resource/out
RUN chmod +x /opt/resource/*

FROM resource