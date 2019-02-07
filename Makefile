hoverfly-test:
	cd core && \
	go test -v $$(go list ./... | grep -v -E 'vendor')

hoverctl-test:
	cd hoverctl && \
	go test -v $$(go list ./... | grep -v -E 'vendor')

hoverfly-build: hoverfly-test
	cd core/cmd/hoverfly && \
	go build -o ../../../target/hoverfly

hoverctl-build: hoverctl-test
	cd hoverctl && \
	go build -ldflags "-X main.hoverctlVersion=$(GIT_TAG_NAME)" -o ../target/hoverctl

CORE_FUNCTIONAL_TESTS = $(shell cd functional-tests/core && go list ./...)

hoverfly-functional-test: hoverfly-build
	cp target/hoverfly functional-tests/core/bin/hoverfly
	cd functional-tests/core && \
	go test -v $(CORE_FUNCTIONAL_TESTS)

hoverctl-functional-test:
	cp target/hoverfly functional-tests/hoverctl/bin/hoverfly
	cd functional-tests/hoverctl && \
	go test -v $(go list ./... | grep -v -E 'vendor')

test: hoverfly-functional-test hoverctl-test hoverctl-functional-test

build:
	cd core/cmd/hoverfly && \
	go build -o ../../../target/hoverfly

	cd hoverctl && \
	go build -ldflags "-X main.hoverctlVersion=$(GIT_TAG_NAME)" -o ../target/hoverctl

build-ui:
	wget https://github.com/SpectoLabs/hoverfly-ui/releases/download/$(GIT_TAG_NAME)/$(GIT_TAG_NAME).zip
	unzip $(GIT_TAG_NAME).zip -d hoverfly-ui	
	cd core && \
	statik -src=../hoverfly-ui
	rm -rf $(GIT_TAG_NAME).zip
	rm -rf hoverfly-ui

benchmark:
	cd core && \
	go test -bench=BenchmarkProcessRequest -run=XXX -cpuprofile profile_cpu.out -memprofile profile_mem.out --benchtime=20s

fmt:
	go fmt $$(go list ./... | grep -v -E 'vendor')

update-dependencies:
	godep save -t ./...

update-version:
	awk \
		-v line=$$(awk '/hoverfly.version/{print NR; exit}' core/hoverfly.go) \
		-v version=${VERSION} \
		'{ if (NR == line) print "	hoverfly.version = \"${VERSION}\""; else print $0}' core/hoverfly.go > core/hoverfly2.go
	rm -rf core/hoverfly.go
	mv core/hoverfly2.go core/hoverfly.go
	git add core/hoverfly.go
	awk \
		-v line=$$(awk '/version/{print NR; exit}' docs/conf.py) \
		-v version=${VERSION} \
		'{ if (NR == line) print "version = \x27${VERSION}\x27"; else print $0}' docs/conf.py > docs/conf2.py
	rm -rf docs/conf.py
	mv docs/conf2.py docs/conf.py
	git add docs/conf.py
	target/hoverctl > docs/pages/reference/hoverctl/hoverctl.output
	git add docs/pages/reference/hoverctl/hoverctl.output
	git commit -m "Updated hoverfly version to ${VERSION}"
