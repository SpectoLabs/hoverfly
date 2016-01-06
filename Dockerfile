FROM golang:1.5

MAINTAINER karolis.rusenas@opencredo.com

ADD . /go/src/github.com/spectolabs/hoverfly

ENV GO15VENDOREXPERIMENT 1

RUN go install github.com/spectolabs/hoverfly

ENTRYPOINT /go/bin/hoverfly

EXPOSE 8500 8888
