package main

import (
	"github.com/Abramovic/logrus_influxdb"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/elazarl/goproxy"
	_ "github.com/influxdb/influxdb/client/v2"
	"github.com/meatballhat/negroni-logrus"

	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"
)

const DefaultPort = ":8500"

func main() {
	// getting proxy configuration
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	record := flag.Bool("record", false, "should proxy record")
	destination := flag.String("destination", ".", "destination URI to catch")
	flag.Parse()

	// getting settings
	initSettings()

	// adding influxdb hook
	err := addInfluxLoggingHook()

	if err != nil {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Error("Failed to add InfluxDB hook")
	}

	// overriding default settings
	AppConfig.recordState = *record

	// overriding destination
	AppConfig.destination = *destination

	// getting default database
	port := os.Getenv("ProxyPort")
	if port == "" {
		port = DefaultPort
	} else {
		port = fmt.Sprintf(":%s", port)
	}

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.TextFormatter{})

	redisPool := getRedisPool()
	defer redisPool.Close()

	cache := Cache{pool: redisPool,
		prefix: AppConfig.cachePrefix}

	// getting connections
	d := DBClient{
		cache: cache,
		http:  &http.Client{},
	}

	// creating proxy
	proxy := goproxy.NewProxyHttpServer()

	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*$"))).
		HandleConnect(goproxy.AlwaysMitm)

	// just helper handler to know where request hits proxy or no
	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			log.WithFields(log.Fields{
				"destination": r.URL.Host,
			}).Info("Got request")

			return r, nil
		})

	// hijacking plain connections
	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile(*destination))).DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {

			log.Info("connection found......")
			log.Info(fmt.Sprintf("Url path:  %s", r.URL.Path))

			return d.processRequest(r)
		})

	go d.startAdminInterface()

	proxy.Verbose = *verbose

	// proxy starting message
	log.WithFields(log.Fields{
		"RedisAddress": AppConfig.redisAddress,
		"Destination":  *destination,
		"ProxyPort":    port,
	}).Info("Proxy is starting...")

	log.Warn(http.ListenAndServe(port, proxy))
}

// processRequest - processes incoming requests and based on proxy state (record/playback)
// returns HTTP response.
func (d *DBClient) processRequest(req *http.Request) (*http.Request, *http.Response) {

	if AppConfig.recordState {
		log.Info("*** RECORD ***")
		newResponse, err := d.recordRequest(req)
		if err != nil {
			// something bad happened, passing through
			return req, nil
		} else {
			// discarding original requests and returns supplied response
			return req, newResponse
		}

	} else {
		log.Info("*** PLAYBACK ***")
		newResponse := d.getResponse(req)
		return req, newResponse
	}
}

func (d *DBClient) startAdminInterface() {
	// starting admin interface
	mux := getBoneRouter(*d)
	n := negroni.Classic()
	n.Use(negronilogrus.NewMiddleware())
	n.UseHandler(mux)

	// admin interface starting message
	log.WithFields(log.Fields{
		"RedisAddress": AppConfig.redisAddress,
		"AdminPort":    AppConfig.adminInterface,
	}).Info("Admin interface is starting...")

	n.Run(AppConfig.adminInterface)
}

func addInfluxLoggingHook() error {
	// checking whether app should send logs to influxdb
	influxdbAddress := os.Getenv("InfluxAddress")

	if influxdbAddress != "" {

		// getting default events database
		influxDatabaseName := os.Getenv("InfluxDBName")
		if influxDatabaseName == "" {
			influxDatabaseName = "events"
		}

		var maxRetries = 10
		var errMaxRetriesReached = errors.New("exceeded retry limit")
		var err error
		//		var cont bool

		attempt := 1
		for {
			hook, err := logrus_influxdb.NewInfluxDBHook("192.168.59.103", "logrus", nil)
			log.Info("Hook created, next step - adding to logrus")
			if err == nil {
				log.AddHook(hook)
				log.Info("Hook to InfluxDB added successfuly")
				break
			}
			attempt++
			log.Warn("Failed to connect to InfluxDB, maybe it is not running yet? Waiting...")
			time.Sleep(5 * time.Second)
			if attempt > maxRetries {
				log.WithFields(log.Fields{
					"Error":              err.Error(),
					"InfluxDB":           influxdbAddress,
					"InfluxDatabaseName": influxDatabaseName,
				}).Error("Unable to add InfluxDB hook")
				return errMaxRetriesReached
			}
		}
		return err

	} else {
		return nil
	}
}
