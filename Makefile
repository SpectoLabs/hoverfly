hoverctl-dependencies:
	cd hoverctl && \
	glide install

hoverctl-functional-test-dependencies:
	cd functional-tests/hoverctl && \
	glide install

hoverctl-test: hoverctl-dependencies
	cd hoverctl && \
	go test -v $(go list ./... | grep -v -E 'vendor')

hoverctl-functional-test: hoverctl-functional-test-dependencies
	cd functional-tests/hoverctl && \
	go test -v $(go list ./... | grep -v -E 'vendor')

build:
	cd hoverctl && \
	go build -o target/hoverctl



