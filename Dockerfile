FROM golang:1.11 as builder
MAINTAINER telia-oss
ADD . /go/src/github.com/telia-oss/appsync-resource
WORKDIR /go/src/github.com/telia-oss/appsync-resource
ENV TARGET linux
ENV ARCH amd64
RUN make build

FROM alpine:3.6 as resource
COPY --from=builder /go/src/github.com/telia-oss/appsync-resource/check /opt/resource/check
COPY --from=builder /go/src/github.com/telia-oss/appsync-resource/in /opt/resource/in
COPY --from=builder /go/src/github.com/telia-oss/appsync-resource/out /opt/resource/out
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN chmod +x /opt/resource/*

FROM resource