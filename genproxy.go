package main

import (
	"github.com/elazarl/goproxy"
	log "github.com/Sirupsen/logrus"

	"net/http"
	"os"
	"flag"
	"regexp"
	"fmt"
)


const DefaultPort = ":8500"


func main() {
	// getting proxy configuration
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	record := flag.Bool("record", false, "should proxy record")
	externalSystem := flag.String("target", "mirage.readthedocs.org", "destination host to catch")
	flag.Parse()

	// getting settings
	initSettings()


	// getting default database
	port := os.Getenv("ProxyPort")
	if (port == "") {
		port = DefaultPort
	} else {
		port = fmt.Sprintf(":%s", port)
	}


	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.TextFormatter{})


	// app starting
	log.WithFields(log.Fields{
		"RedisAddress": AppConfig.redisAddress,
		"ProxyPort": port,
	}).Info("app is starting")

	// getting redis connection
	d := DBClient{pool: getRedisPool()}

	// creating proxy
	proxy := goproxy.NewProxyHttpServer()

	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*$"))).
	HandleConnect(goproxy.AlwaysMitm)

	// hijacking plain connections
	proxy.OnRequest(goproxy.DstHostIs(*externalSystem)).DoFunc(
		func(r *http.Request,ctx *goproxy.ProxyCtx)(*http.Request,*http.Response) {

			log.Info("connection found......")
			log.Info(fmt.Sprintf("Url path:  %s", r.URL.Path))
			if *record {
				log.Info("*** RECORD ***")
				d.recordRequest(r)
			} else {
				log.Info("*** PLAYBACK ***")
				_ = d.getResponse(r)
			}
			return r,nil
		})


	proxy.Verbose = *verbose
	log.Fatal(http.ListenAndServe(port, proxy))
}