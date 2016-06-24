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

hoverfly-test:
	cd core && \
	go test -v $(go list ./.. | grep -v -E 'vendor')

hoverctl-test:
	cd hoverctl && \
	go test -v $(go list ./... | grep -v -E 'vendor')

hoverfly-build: hoverfly-test
	cd core/cmd/hoverfly && \
	go build -o ../../../target/hoverfly

hoverctl-build: hoverctl-test
	cd hoverctl && \
	go build -o ../target/hoverctl

hoverfly-functional-test: hoverfly-build
	cp target/hoverfly functional-tests/core/bin/hoverfly
	cd functional-tests/core && \
	go test -v $(go list ./.. | grep -v -E 'vendor')

hoverctl-functional-test: hoverctl-build
	cp target/hoverctl functional-tests/hoverctl/bin/hoverctl
	cd functional-tests/hoverctl && \
	go test -v $(go list ./... | grep -v -E 'vendor')

test: hoverfly-functional-test hoverctl-functional-test

build: test

gox-build: test
	rm -rf target/*
	cd core/cmd/hoverfly && \
	$(GOPATH)/bin/gox
	mv core/cmd/hoverfly/hoverfly_* target/
	cd hoverctl && \
	$(GOPATH)/bin/gox
	mv hoverctl/hoverctl_* target/
