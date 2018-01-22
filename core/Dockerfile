FROM debian:jessie

MAINTAINER karolis.rusenas@opencredo.com

# you already built the binaries on linux, right?
COPY target/hoverfly /go/bin/hoverfly

ENTRYPOINT ["/go/bin/hoverfly", "-listen-on-host=0.0.0.0"]
CMD [""]

EXPOSE 8500 8888
