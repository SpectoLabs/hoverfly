FROM golang:1.9.1

MAINTAINER benji.hooper@specto.io

ADD . /go/src/github.com/SpectoLabs/hoverfly

ENV GO15VENDOREXPERIMENT 1

RUN go install github.com/SpectoLabs/hoverfly/cmd/hoverfly/

ENTRYPOINT ["/go/bin/hoverfly"]
CMD [""]

EXPOSE 8500 8888
