package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/elazarl/goproxy"

	"bufio"
	"flag"
	"fmt"
	"net"
	"net/http"
	"regexp"
)

// VirtualizeMode - default mode when Hoverfly looks for captured requests to respond
const VirtualizeMode = "virtualize"

// SynthesizeMode - all requests are sent to middleware to create response
const SynthesizeMode = "synthesize"

// ModifyMode - middleware is applied to outgoing and incoming traffic
const ModifyMode = "modify"

// CaptureMode - requests are captured and stored in cache
const CaptureMode = "capture"

// orPanic - wrapper for logging errors
func orPanic(err error) {
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Panic("Got error.")
	}
}

func main() {
	// Output to stderr instead of stdout, could also be a file.
	//	log.SetOutput(os.Stderr)
	//	log.SetFormatter(&log.TextFormatter{})
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

	flag.Parse()

	// getting settings
	cfg := InitSettings()

	if *verbose {
		// Only log the warning severity or above.
		log.SetLevel(log.DebugLevel)
	}
	cfg.verbose = *verbose

	// overriding environment variables (proxy and admin ports)
	if *proxyPort != "" {
		cfg.proxyPort = *proxyPort
	}
	if *adminPort != "" {
		cfg.adminPort = *adminPort
	}

	// overriding default middleware setting
	cfg.middleware = *middleware

	// setting default mode
	mode := VirtualizeMode

	if *capture {
		mode = CaptureMode
		// checking whether user supplied other modes
		if *synthesize == true || *modify == true {
			log.Fatal("Two or more modes supplied, check your flags")
		}
	} else if *synthesize {
		mode = SynthesizeMode

		if cfg.middleware == "" {
			log.Fatal("Synthesize mode chosen although middleware not supplied")
		}

		if *capture == true || *modify == true {
			log.Fatal("Two or more modes supplied, check your flags")
		}
	} else if *modify {
		mode = ModifyMode

		if cfg.middleware == "" {
			log.Fatal("Modify mode chosen although middleware not supplied")
		}

		if *capture == true || *synthesize == true {
			log.Fatal("Two or more modes supplied, check your flags")
		}
	}

	// overriding default settings
	cfg.mode = mode

	// overriding destination
	cfg.destination = *destination

	proxy, dbClient := getNewHoverfly(cfg)
	defer dbClient.cache.db.Close()

	log.Warn(http.ListenAndServe(fmt.Sprintf(":%s", cfg.proxyPort), proxy))
}

// getNewHoverfly returns a configured ProxyHttpServer and DBClient, also starts admin interface on configured port
func getNewHoverfly(cfg *Configuration) (*goproxy.ProxyHttpServer, DBClient) {

	// getting boltDB
	db := getDB(cfg.databaseName)

	cache := Cache{
		db:             db,
		requestsBucket: []byte(requestsBucketName),
	}

	// getting connections
	d := DBClient{
		cache: cache,
		http:  &http.Client{},
		cfg:   cfg,
	}

	// creating proxy
	proxy := goproxy.NewProxyHttpServer()

	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile(d.cfg.destination))).
		HandleConnect(goproxy.AlwaysMitm)

	// enable curl -p for all hosts on port 80
	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile(d.cfg.destination))).
		HijackConnect(func(req *http.Request, client net.Conn, ctx *goproxy.ProxyCtx) {
		defer func() {
			if e := recover(); e != nil {
				ctx.Logf("error connecting to remote: %v", e)
				client.Write([]byte("HTTP/1.1 500 Cannot reach destination\r\n\r\n"))
			}
			client.Close()
		}()
		clientBuf := bufio.NewReadWriter(bufio.NewReader(client), bufio.NewWriter(client))
		remote, err := net.Dial("tcp", req.URL.Host)
		orPanic(err)
		remoteBuf := bufio.NewReadWriter(bufio.NewReader(remote), bufio.NewWriter(remote))
		for {
			req, err := http.ReadRequest(clientBuf.Reader)
			orPanic(err)
			orPanic(req.Write(remoteBuf))
			orPanic(remoteBuf.Flush())
			resp, err := http.ReadResponse(remoteBuf.Reader, req)

			orPanic(err)
			orPanic(resp.Write(clientBuf.Writer))
			orPanic(clientBuf.Flush())
		}
	})

	// processing connections
	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile(cfg.destination))).DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			return d.processRequest(r)
		})

	go d.startAdminInterface()

	proxy.Verbose = d.cfg.verbose
	// proxy starting message
	log.WithFields(log.Fields{
		"Destination": d.cfg.destination,
		"ProxyPort":   d.cfg.proxyPort,
		"Mode":        d.cfg.GetMode(),
	}).Info("Proxy prepared...")

	return proxy, d
}

func hoverflyError(req *http.Request, err error, msg string, statusCode int) *http.Response {
	return goproxy.NewResponse(req,
		goproxy.ContentTypeText, statusCode,
		fmt.Sprintf("Hoverfly Error! %s. Got error: %s \n", msg, err.Error()))
}

// processRequest - processes incoming requests and based on proxy state (record/playback)
// returns HTTP response.
func (d *DBClient) processRequest(req *http.Request) (*http.Request, *http.Response) {

	mode := d.cfg.GetMode()

	if mode == CaptureMode {
		newResponse, err := d.captureRequest(req)

		if err != nil {
			return req, hoverflyError(req, err, "Could not capture request", http.StatusServiceUnavailable)
		}
		log.WithFields(log.Fields{
			"mode":        mode,
			"middleware":  d.cfg.middleware,
			"path":        req.URL.Path,
			"rawQuery":    req.URL.RawQuery,
			"method":      req.Method,
			"destination": req.Host,
		}).Info("request and response captured")

		return req, newResponse

	} else if mode == SynthesizeMode {
		response, err := synthesizeResponse(req, d.cfg.middleware)

		if err != nil {
			return req, hoverflyError(req, err, "Could not create synthetic response!", http.StatusServiceUnavailable)
		}

		log.WithFields(log.Fields{
			"mode":        mode,
			"middleware":  d.cfg.middleware,
			"path":        req.URL.Path,
			"rawQuery":    req.URL.RawQuery,
			"method":      req.Method,
			"destination": req.Host,
		}).Info("synthetic response created successfuly")

		return req, response

	} else if mode == ModifyMode {
		response, err := d.modifyRequestResponse(req, d.cfg.middleware)

		if err != nil {
			log.WithFields(log.Fields{
				"error":      err.Error(),
				"middleware": d.cfg.middleware,
			}).Error("Got error when performing request modification")
			return req, hoverflyError(
				req,
				err,
				fmt.Sprintf("Middleware (%s) failed or something else happened!", d.cfg.middleware),
				http.StatusServiceUnavailable)
		}
		// returning modified response
		return req, response
	}

	newResponse := d.getResponse(req)
	return req, newResponse

}
