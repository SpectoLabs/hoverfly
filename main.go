package main

import (
	"github.com/elazarl/goproxy"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"os"
	"flag"
	"regexp"
)

func orPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// Initial structure of configuration that is expected from conf.json file
type Configuration struct {
	MirageEndpoint string
	ExternalSystem string
}

// AppConfig stores application configuration
var AppConfig Configuration

func main() {
	// getting proxy configuration
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr := flag.String("addr", ":8080", "proxy listen address")
	externalSystem := flag.String("target", "/", "URL path to catch")
	flag.Parse()


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

	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*$"))).
	HandleConnect(goproxy.AlwaysMitm)

	// hijacking plain connections
	proxy.OnRequest(goproxy.DstHostIs(*externalSystem)).DoFunc(
		func(r *http.Request,ctx *goproxy.ProxyCtx)(*http.Request,*http.Response) {



			log.Info("connection found......")
			return r,nil
		})


	proxy.Verbose = *verbose
	log.Fatal(http.ListenAndServe(*addr, proxy))
}