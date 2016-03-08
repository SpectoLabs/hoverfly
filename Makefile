deps:
	go get -u github.com/Sirupsen/logrus
	go get -u github.com/elazarl/goproxy
	go get -u github.com/meatballhat/negroni-logrus
	go get -u github.com/codegangsta/negroni
	go get -u github.com/go-zoo/bone
	go get -u github.com/boltdb/bolt
	go get -u github.com/rakyll/statik
	go get -u github.com/rcrowley/go-metrics
	go get -u github.com/gorilla/websocket
	go get -u github.com/dgrijalva/jwt-go
	go get -u github.com/rusenask/goproxy
	go get -u github.com/SpectoLabs/hoverfly

test: deps
	go test

build: deps
	cd cmd/hoverfly/ && go build

build_ci: deps
	go get -u bitbucket.org/tebeka/go2xunit
	go get -u github.com/mitchellh/gox
