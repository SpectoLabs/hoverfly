deps:
	go get

test: deps
	go test

build-drone: deps
	go get -u all
	cd cmd/hoverfly/ && go build

#build: deps
#	cd cmd/hoverfly/ && go build -o ${GOPATH}/bin/hoverflyb

build-ami:
	packer build -var 'aws_access_key=${AWS_ACCESS_KEY}' -var 'aws_secret_key=${AWS_SECRET_KEY}' packer.json

dependencies: hoverctl-dependencies hoverctl-functional-test-dependencies

hoverfly-dependencies:
	cd core && \
	glide install

hoverctl-dependencies:
	cd hoverctl && \
	glide install

hoverfly-functional-test-dependencies:
	cd functional-tests/core && \
	glide install

hoverctl-functional-test-dependencies:
	cd functional-tests/hoverctl && \
	glide install

hoverctl-test: hoverctl-dependencies
	cd hoverctl && \
	go test -v $(go list ./... | grep -v -E 'vendor')

hoverctl-build: hoverctl-test
	cd hoverctl && \
	go build -o ../target/hoverctl

hoverctl-functional-test: hoverctl-functional-test-dependencies hoverctl-build
	cp target/hoverctl functional-tests/hoverctl/bin/hoverctl
	cd functional-tests/hoverctl && \
	go test -v $(go list ./... | grep -v -E 'vendor')

build: hoverctl-functional-test
