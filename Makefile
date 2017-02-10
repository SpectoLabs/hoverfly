hoverfly-test:
	cd core && \
	go test -v $$(go list ./... | grep -v -E 'vendor')

hoverctl-test:
	cd hoverctl && \
	go test -v $$(go list ./... | grep -v -E 'vendor')

hoverfly-build: hoverfly-test
	cd core/cmd/hoverfly && \
	go build -ldflags "-X main.hoverflyVersion=$(GIT_TAG_NAME)" -o ../../../target/hoverfly

hoverctl-build: hoverctl-test
	cd hoverctl && \
	go build -ldflags "-X main.hoverctlVersion=$(GIT_TAG_NAME)" -o ../target/hoverctl

hoverfly-functional-test: hoverfly-build
	cp target/hoverfly functional-tests/core/bin/hoverfly
	cd functional-tests/core && \
	go test -v $(go list ./... | grep -v -E 'vendor')

hoverctl-functional-test: hoverctl-build
	cp target/hoverctl functional-tests/hoverctl/bin/hoverctl
	cp target/hoverfly functional-tests/hoverctl/bin/hoverfly
	cd functional-tests/hoverctl && \
	go test -v $(go list ./... | grep -v -E 'vendor')

test: hoverfly-functional-test hoverctl-functional-test

build:
	cd core/cmd/hoverfly && \
	go build -ldflags "-X main.hoverflyVersion=$(GIT_TAG_NAME)" -o ../../../target/hoverfly

	cd hoverctl && \
	go build -ldflags "-X main.hoverctlVersion=$(GIT_TAG_NAME)" -o ../target/hoverctl

fmt:
	go fmt $$(go list ./... | grep -v -E 'vendor')

update-version:
	awk \
		-v line=$$(awk '/h.version/{print NR; exit}' core/hoverfly.go) \
		-v version=${VERSION} \
		'{ if (NR == line) print "	h.version = \"${VERSION}\""; else print $0}' core/hoverfly.go > core/hoverfly2.go
	rm -rf core/hoverfly.go
	mv core/hoverfly2.go core/hoverfly.go
	git add core/hoverfly.go
	git commit -m "Updated hoverfly version to ${VERSION}"
