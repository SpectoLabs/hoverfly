FROM golang:1.25.4 AS build-env
WORKDIR /usr/local/go/src/github.com/SpectoLabs/hoverfly
COPY . /usr/local/go/src/github.com/SpectoLabs/hoverfly

# Support multi-arch builds with buildx by honoring TARGETOS/TARGETARCH
ARG TARGETOS=linux
ARG TARGETARCH
RUN cd core/cmd/hoverfly \
    && CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
       go install -ldflags "-s -w"

# Final minimal image based on Alpine, without running apk to avoid QEMU trigger issues
FROM alpine:3.20
# Provide CA certificates by copying from the builder (PEM bundle is arch-independent)
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# Copy the statically linked hoverfly binary
COPY --from=build-env /usr/local/go/bin/hoverfly /bin/hoverfly

ENTRYPOINT ["/bin/hoverfly", "-listen-on-host=0.0.0.0"]
EXPOSE 8500 8888
