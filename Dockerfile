FROM golang:1.5

MAINTAINER karolis.rusenas@gmail.com

ADD . /go/src/github.com/rusenask/genproxy

# provide redis connection details
# ENV RedisAddress=redis
# ENV RedisPassword=redis_pass

ENV GO15VENDOREXPERIMENT 1

RUN go install github.com/rusenask/genproxy

ENTRYPOINT /go/bin/genproxy

EXPOSE 8500 8888

