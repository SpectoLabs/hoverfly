FROM golang:1.11.4 AS build-env
WORKDIR /usr/local/go/src/github.com/SpectoLabs/hoverfly
COPY . /usr/local/go/src/github.com/SpectoLabs/hoverfly    
RUN cd core/cmd/hoverfly && CGO_ENABLED=0 GOOS=linux go install -ldflags "-s -w"

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build-env /usr/local/go/bin/hoverfly /bin/hoverfly
ENTRYPOINT ["/bin/hoverfly", "-listen-on-host=0.0.0.0"]
CMD [""]

EXPOSE 8500 8888
