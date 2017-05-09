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

hoverfly-functional-test: hoverfly-build
	cp target/hoverfly functional-tests/core/bin/hoverfly
	cd functional-tests/core && \
	go test -v $(go list ./... | grep -v -E 'vendor')

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
	cd core/static/admin && \
	npm run compile && \
	npm run deploy
	cd core && \
	statik -src=./static/admin/dist

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
