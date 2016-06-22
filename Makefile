dependencies:
	cd hoverctl && \
	glide install

	cd hoverctl-ft && \
	glide install
test:
	cd hoverctl && \
	go test -v $(go list ./... | grep -v -E 'vendor')

build:
	cd hoverctl && \
	go build -o target/hoverctl



