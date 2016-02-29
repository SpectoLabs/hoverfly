deps:
	go get github.com/Sirupsen/logrus
	go get github.com/elazarl/goproxy
	go get github.com/meatballhat/negroni-logrus
	go get github.com/codegangsta/negroni
	go get github.com/go-zoo/bone
	go get github.com/boltdb/bolt
	go get github.com/rakyll/statik
	go get github.com/rcrowley/go-metrics
	go get github.com/gorilla/websocket
	go get github.com/dgrijalva/jwt-go
	go get github.com/rusenask/goproxy

test: deps
	go test

build: deps
	cd cmd/hoverfly/ && go build
