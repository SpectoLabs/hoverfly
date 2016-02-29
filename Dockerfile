FROM golang:1.5

MAINTAINER karolis.rusenas@opencredo.com

ADD . /go/src/github.com/SpectoLabs/hoverfly

ENV GO15VENDOREXPERIMENT 1

RUN go install github.com/SpectoLabs/hoverfly/cmd/hoverfly/

ENTRYPOINT /go/bin/hoverfly

EXPOSE 8500 8888
