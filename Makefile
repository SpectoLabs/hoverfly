deps:
	go get

test: deps
	go test

build-jenkins: deps
	cd cmd/hoverfly/ && go build

build: deps
	cd cmd/hoverfly/ && go build -o ${GOPATH}/bin/hoverflyb

build_ci: deps
	go get -u bitbucket.org/tebeka/go2xunit
	go get -u github.com/mitchellh/gox
