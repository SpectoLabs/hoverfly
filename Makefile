deps:
	go get

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

hoverfly-test: hoverfly-dependencies
	cd core && \
	go test -v $(go list ./.. | grep -v -E 'vendor')

hoverctl-test: hoverctl-dependencies
	cd hoverctl && \
	go test -v $(go list ./... | grep -v -E 'vendor')

hoverfly-build: hoverfly-test
	cd core/cmd/hoverfly && \
	go build -o ../../../target/hoverfly

hoverctl-build: hoverctl-test
	cd hoverctl && \
	go build -o ../target/hoverctl

hoverfly-functional-test: hoverfly-functional-test-dependencies hoverfly-build
	cp target/hoverfly functional-tests/core/bin/hoverfly
	cd functional-tests/core && \
	go test -v $(go list ./.. | grep -v -E 'vendor')

hoverctl-functional-test: hoverctl-functional-test-dependencies hoverctl-build
	cp target/hoverctl functional-tests/hoverctl/bin/hoverctl
	cd functional-tests/hoverctl && \
	go test -v $(go list ./... | grep -v -E 'vendor')

test: hoverfly-functional-test hoverctl-functional-test
