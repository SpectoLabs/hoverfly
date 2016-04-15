deps:
	go get

test: deps
	go test

build-drone: deps
	go get -u all
	cd cmd/hoverfly/ && go build

build: deps
	cd cmd/hoverfly/ && go build -o ${GOPATH}/bin/hoverflyb

build-ami:
	packer build -var 'aws_access_key=${AWS_ACCESS_KEY}' -var 'aws_secret_key=${AWS_SECRET_KEY}' packer.json
