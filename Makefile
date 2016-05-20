dependencies:
	glide up

test:
	go test -v $(go list ./... | grep -v -E 'vendor')

build:
	go build -o hoverfly-cli
