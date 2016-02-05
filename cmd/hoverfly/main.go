package main

import (
	log "github.com/Sirupsen/logrus"
	hv "github.com/SpectoLabs/hoverfly"

	"flag"
	"fmt"
	"net/http"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	// getting proxy configuration
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	// modes
	capture := flag.Bool("capture", false, "should proxy capture requests")
	synthesize := flag.Bool("synthesize", false, "should proxy capture requests")
	modify := flag.Bool("modify", false, "should proxy only modify requests")

	destination := flag.String("destination", ".", "destination URI to catch")
	middleware := flag.String("middleware", "", "should proxy use middleware")

	// proxy port
	proxyPort := flag.String("pp", "", "proxy port - run proxy on another port (i.e. '-pp 9999' to run proxy on port 9999)")
	// admin port
	adminPort := flag.String("ap", "", "admin port - run admin interface on another port (i.e. '-ap 1234' to run admin UI on port 1234)")

	// metrics
	metrics := flag.Bool("metrics", false, "supply -metrics flag to enable metrics logging to stdout")

	// development
	dev := flag.Bool("dev", false, "supply -dev flag to serve directly from ./static/dist instead from statik binary")

	flag.Parse()

	// getting settings
	cfg := hv.InitSettings()

	if *verbose {
		// Only log the warning severity or above.
		log.SetLevel(log.DebugLevel)
	}
	cfg.Verbose = *verbose

	if *dev {
		// making text pretty
		log.SetFormatter(&log.TextFormatter{})
	}

	// overriding environment variables (proxy and admin ports)
	if *proxyPort != "" {
		cfg.ProxyPort = *proxyPort
	}
	if *adminPort != "" {
		cfg.AdminPort = *adminPort
	}

	// development settings
	cfg.Development = *dev

	// overriding default middleware setting
	cfg.Middleware = *middleware

	// setting default mode
	mode := hv.VirtualizeMode

	if *capture {
		mode = hv.CaptureMode
		// checking whether user supplied other modes
		if *synthesize == true || *modify == true {
			log.Fatal("Two or more modes supplied, check your flags")
		}
	} else if *synthesize {
		mode = hv.SynthesizeMode

		if cfg.Middleware == "" {
			log.Fatal("Synthesize mode chosen although middleware not supplied")
		}

		if *capture == true || *modify == true {
			log.Fatal("Two or more modes supplied, check your flags")
		}
	} else if *modify {
		mode = hv.ModifyMode

		if cfg.Middleware == "" {
			log.Fatal("Modify mode chosen although middleware not supplied")
		}

		if *capture == true || *synthesize == true {
			log.Fatal("Two or more modes supplied, check your flags")
		}
	}

	// overriding default settings
	cfg.Mode = mode

	// overriding destination
	cfg.Destination = *destination

	proxy, dbClient := hv.GetNewHoverfly(cfg)
	defer dbClient.Cache.DS.Close()

	// starting admin interface
	go dbClient.StartAdminInterface()

	// start metrics registry flush
	if *metrics {
		go dbClient.Counter.Init()
	}

	log.Warn(http.ListenAndServe(fmt.Sprintf(":%s", cfg.ProxyPort), proxy))
}
