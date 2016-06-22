dependencies: hoverctl-dependencies hoverctl-functional-test-dependencies

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

hoverctl-build: hoverctl-test
	cd hoverctl && \
	go build -o ../target/hoverctl



