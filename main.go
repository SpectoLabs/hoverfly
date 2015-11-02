package main

import (
	"github.com/elazarl/goproxy"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"os"
)

// Initial structure of configuration that is expected from conf.json file
type Configuration struct {
	MirageEndpoint string
	ExternalSystem string
}

// AppConfig stores application configuration
var AppConfig Configuration

func main() {
	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.TextFormatter{})

	// getting app config
	mirageEndpoint := os.Getenv("MirageEndpoint")

	AppConfig.MirageEndpoint = mirageEndpoint

	// app starting
	log.WithFields(log.Fields{
		"MirageEndpoint": AppConfig.MirageEndpoint,
	}).Info("app is starting")

	proxy := goproxy.NewProxyHttpServer()

	proxy.Verbose = true
	log.Fatal(http.ListenAndServe(":8080", proxy))
}