FROM golang:1.11 as builder

LABEL MAINTAINER="telia-oss"

ADD . /go/src/github.com/telia-oss/appsync-resource

WORKDIR /go/src/github.com/telia-oss/appsync-resource

ENV TARGET linux
ENV ARCH amd64

RUN make build

FROM alpine:3.8 AS resource
COPY --from=builder /go/src/github.com/telia-oss/appsync-resource/check /opt/resource/check
COPY --from=builder /go/src/github.com/telia-oss/appsync-resource/in /opt/resource/in
COPY --from=builder /go/src/github.com/telia-oss/appsync-resource/out /opt/resource/out
RUN chmod +x /opt/resource/out

FROM resource
