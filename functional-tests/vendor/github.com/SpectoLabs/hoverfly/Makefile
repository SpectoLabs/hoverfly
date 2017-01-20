dependencies: hoverfly-dependencies hoverfly-functional-test-dependencies hoverctl-dependencies hoverctl-functional-test-dependencies

hoverfly-dependencies:
	cd core && \
	glide install --strip-vcs

hoverctl-dependencies:
	cd hoverctl && \
	glide install --strip-vcs

hoverfly-functional-test-dependencies:
	cd functional-tests/core && \
	glide install --strip-vcs

hoverctl-functional-test-dependencies:
	cd functional-tests/hoverctl && \
	glide install --strip-vcs

hoverfly-test: hoverfly-dependencies
	cd core && \
	go test -v $$(go list ./... | grep -v -E 'vendor')

hoverctl-test: hoverctl-dependencies
	cd hoverctl && \
	go test -v $$(go list ./... | grep -v -E 'vendor')

hoverfly-build: hoverfly-test
	cd core/cmd/hoverfly && \
	go build -ldflags "-X main.hoverflyVersion=$(GIT_TAG_NAME)" -o ../../../target/hoverfly

hoverctl-build: hoverctl-test
	cd hoverctl && \
	go build -ldflags "-X main.hoverctlVersion=$(GIT_TAG_NAME)" -o ../target/hoverctl

hoverfly-functional-test: hoverfly-functional-test-dependencies hoverfly-build
	cp target/hoverfly functional-tests/core/bin/hoverfly
	cd functional-tests/core && \
	go test -v $(go list ./... | grep -v -E 'vendor')

hoverctl-functional-test: hoverctl-functional-test-dependencies hoverctl-build
	cp target/hoverctl functional-tests/hoverctl/bin/hoverctl
	cp target/hoverfly functional-tests/hoverctl/bin/hoverfly
	cd functional-tests/hoverctl && \
	go test -v $(go list ./... | grep -v -E 'vendor')

test: hoverfly-functional-test hoverctl-functional-test

build: test

gox-build: 
	rm -rf target/*
	cd core/cmd/hoverfly && \
	$(GOPATH)/bin/gox -ldflags "-X main.hoverflyVersion=$(GIT_TAG_NAME)"
	mv core/cmd/hoverfly/hoverfly_* target/
	cd hoverctl && \
	$(GOPATH)/bin/gox -ldflags "-X main.hoverctlVersion=$(GIT_TAG_NAME)"
	mv hoverctl/hoverctl_* target/

version-binaries:
	for f in target/*; \
	do \
 		mv $$f $${f:0:15}_$(GIT_TAG_NAME)$${f:15}; \
	done;

rename-darwin-binaries:
	mv target/hoverfly_darwin_386 target/hoverfly_OSX_386
	mv target/hoverfly_darwin_amd64 target/hoverfly_OSX_amd64
	mv target/hoverctl_darwin_386 target/hoverctl_OSX_386
	mv target/hoverctl_darwin_amd64 target/hoverctl_OSX_amd64

build-release: gox-build rename-darwin-binaries version-binaries

fmt:
	go fmt $$(go list ./... | grep -v -E 'vendor')
