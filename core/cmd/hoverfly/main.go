// Copyright 2015 SpectoLabs. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// hoverfly is an HTTP/s proxy configurable via flags/environment variables/admin HTTP API
//
// this proxy can be dynamically configured through HTTP calls when it's running, to change modes,
// export and import requests.

package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/SpectoLabs/goproxy"
	hv "github.com/SpectoLabs/hoverfly/core"
	"github.com/SpectoLabs/hoverfly/core/authentication/backends"
	"github.com/SpectoLabs/hoverfly/core/cache"
	hvc "github.com/SpectoLabs/hoverfly/core/certs"
	cs "github.com/SpectoLabs/hoverfly/core/cors"
	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/SpectoLabs/hoverfly/core/matching"
	mw "github.com/SpectoLabs/hoverfly/core/middleware"
	"github.com/SpectoLabs/hoverfly/core/modes"
	"github.com/SpectoLabs/hoverfly/core/util"
	log "github.com/sirupsen/logrus"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var importFlags arrayFlags
var destinationFlags arrayFlags
var logOutputFlags arrayFlags
var responseBodyFilesPath string
var responseBodyFilesAllowedOriginFlags arrayFlags

const boltBackend = "boltdb"
const inmemoryBackend = "memory"

var (
	version       = flag.Bool("version", false, "Get the version of hoverfly")
	verbose       = flag.Bool("v", false, "Should every proxy request be logged to stdout")
	logLevelFlag  = flag.String("log-level", "info", "Set log level (panic, fatal, error, warn, info or debug)")
	capture       = flag.Bool("capture", false, "Start Hoverfly in capture mode - transparently intercepts and saves requests/response")
	synthesize    = flag.Bool("synthesize", false, "Start Hoverfly in synthesize mode (middleware is required)")
	modify        = flag.Bool("modify", false, "Start Hoverfly in modify mode - applies middleware (required) to both outgoing and incoming HTTP traffic")
	spy           = flag.Bool("spy", false, "Start Hoverfly in spy mode, similar to simulate but calls real server when cache miss")
	diff          = flag.Bool("diff", false, "Start Hoverfly in diff mode - calls real server and compares the actual response with the expected simulation config if present")
	middleware    = flag.String("middleware", "", "Set middleware by passing the name of the binary and the path of the middleware script separated by space. (i.e. '-middleware \"python script.py\"')")
	proxyPort     = flag.String("pp", "", "Proxy port - run proxy on another port (i.e. '-pp 9999' to run proxy on port 9999)")
	adminPort     = flag.String("ap", "", "Admin port - run admin interface on another port (i.e. '-ap 1234' to run admin UI on port 1234)")
	listenOnHost  = flag.String("listen-on-host", "", "Specify which network interface to bind to, eg. 0.0.0.0 will bind to all interfaces. By default hoverfly will only bind ports to loopback interface")
	metrics       = flag.Bool("metrics", false, "Enable metrics logging to stdout")
	dev           = flag.Bool("dev", false, "Enable CORS headers to allow Hoverfly Admin UI development")
	devCorsOrigin = flag.String("dev-cors-origin", "http://localhost:4200", "Custom CORS origin for dev mode")
	destination   = flag.String("destination", ".", "Control which URLs Hoverfly should intercept and process, it can be string or regex")
	webserver     = flag.Bool("webserver", false, "Start Hoverfly in webserver mode (simulate mode)")

	addNew          = flag.Bool("add", false, "Add new user '-add -username hfadmin -password hfpass'")
	addUser         = flag.String("username", "", "Username for new user")
	addPassword     = flag.String("password", "", "Password for new user")
	addPasswordHash = flag.String("password-hash", "", "Password hash for new user instead of password")
	isAdmin         = flag.Bool("admin", true, "Supply '-admin=false' to make this non admin user")
	authEnabled     = flag.Bool("auth", false, "Enable authentication")

	generateCA = flag.Bool("generate-ca-cert", false, "Generate CA certificate and private key for MITM")
	certName   = flag.String("cert-name", "hoverfly.proxy", "Cert name")
	certOrg    = flag.String("cert-org", "Hoverfly Authority", "Organisation name for new cert")
	cert       = flag.String("cert", "", "CA certificate used to sign MITM certificates")
	key        = flag.String("key", "", "Private key of the CA used to sign MITM certificates")

	tlsVerification    = flag.Bool("tls-verification", true, "Turn on/off tls verification for outgoing requests (will not try to verify certificates)")
	plainHttpTunneling = flag.Bool("plain-http-tunneling", false, "Use plain http tunneling to host with non-443 port")

	upstreamProxy = flag.String("upstream-proxy", "", "Specify an upstream proxy for hoverfly to route traffic through")

	databasePath = flag.String("db-path", "", "A path to a BoltDB file with persisted user and token data for authentication (DEPRECATED)")
	database     = flag.String("db", inmemoryBackend, "Storage to use - 'boltdb' or 'memory' which will not write anything to disk (DEPRECATED)")
	disableCache = flag.Bool("disable-cache", false, "Disable the request/response cache (the cache that sits in front of matching)")

	logsFormat = flag.String("logs", "plaintext", "Specify format for logs, options are \"plaintext\" and \"json\"")
	logsSize   = flag.Int("logs-size", 1000, "Set the amount of logs to be stored in memory")
	logsFile   = flag.String("logs-file", "hoverfly.log", "Specify log file name for output logs")
	logNoColor = flag.Bool("log-no-color", false, "Disable colors for logging")

	journalSize   = flag.Int("journal-size", 1000, "Set the size of request/response journal")
	cacheSize     = flag.Int("cache-size", 1000, "Set the size of request/response cache")
	cors          = flag.Bool("cors", false, "Enable CORS support")
	noImportCheck = flag.Bool("no-import-check", false, "Skip duplicate request check when importing simulations")

	pacFile = flag.String("pac-file", "", "Path to the pac file to be imported on startup")

	clientAuthenticationDestination = flag.String("client-authentication-destination", "", "Regular expression of destination with client authentication")
	clientAuthenticationClientCert  = flag.String("client-authentication-client-cert", "", "Path to the client certification file used for authentication")
	clientAuthenticationClientKey   = flag.String("client-authentication-client-key", "", "Path to the client key file used for authentication")
	clientAuthenticationCACert      = flag.String("client-authentication-ca-cert", "", "Path to the ca cert file used for authentication")
)

var CA_CERT = []byte(`-----BEGIN CERTIFICATE-----
MIIDbTCCAlWgAwIBAgIVAPFUKC/hDKXSN4nF4Gh/fG7Oby4KMA0GCSqGSIb3DQEB
CwUAMDYxGzAZBgNVBAoTEkhvdmVyZmx5IEF1dGhvcml0eTEXMBUGA1UEAxMOaG92
ZXJmbHkucHJveHkwHhcNMjIwMzI3MjE0OTA4WhcNMzIwMzI0MjE0OTA4WjA2MRsw
GQYDVQQKExJIb3ZlcmZseSBBdXRob3JpdHkxFzAVBgNVBAMTDmhvdmVyZmx5LnBy
b3h5MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA6mqE6O8H14tsul0B
UGhLuxFYdSFsjtWtcR4v5PJDK118pdoYC3hgvejcQHdjzuYfVJybo2UHxhEyomhu
r3KrpcjC0VnGfeibNXY01JDWMVxC2QgutGZb92/wChMBfOKYq5z4MhK+5gdiBkz2
C8/1Q724sw14iIcQB+POY6lVBj3YI5Ja+hjSm6SWVvMVRk1uMgx2CcW6zbgErkNg
xvDnDPHlRl5aIIHNyDMlSczVtq0SlBrTtExmjSg2Edzo1v25DG1LBzV58zYE5/cr
Yh+Dm1XKB88sSBb8bUoAjZCWsl3Dkn8eR8wOdabZZxU/STP/g9yxTM1fcnc4v4e/
QF0TnQIDAQABo3IwcDAOBgNVHQ8BAf8EBAMCAqQwEwYDVR0lBAwwCgYIKwYBBQUH
AwEwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUgrvEDpBhKH6SOUAC2fs8oTri
5XMwGQYDVR0RBBIwEIIOaG92ZXJmbHkucHJveHkwDQYJKoZIhvcNAQELBQADggEB
AOXmvQtdsH4kBGAnMI87SlFbAskrbeY/Kqr1PQyDTt2MVj/SjpsVxNEoIsS5ghcI
EyvhD/3t2q15D3XNc+wixSu8jCTe9N1CGXdiolfZ09SqiBtItvOh9R7pdkCcquh+
69JJayOMInoSnmaf+ic+gbzLiEgfW+Dv/OR2Bmuelrs1zOnHdXhY45bN6PRQFrWt
+Wkr7OqTfoCAz6NGgSWcKrXymTtErX7ZJGYwSc2+nHQznl7RBdyL2BfQAVWaWmhI
s+IfxcKlYBr/nKWOkhD81VrNXFEj6R5kEYOdXYe9ovRmQhKWSz4cpCcMqtx1ye7K
1iBU2wVABfQZhp/3eiIyF9w=
-----END CERTIFICATE-----`)

var CA_KEY = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEA6mqE6O8H14tsul0BUGhLuxFYdSFsjtWtcR4v5PJDK118pdoY
C3hgvejcQHdjzuYfVJybo2UHxhEyomhur3KrpcjC0VnGfeibNXY01JDWMVxC2Qgu
tGZb92/wChMBfOKYq5z4MhK+5gdiBkz2C8/1Q724sw14iIcQB+POY6lVBj3YI5Ja
+hjSm6SWVvMVRk1uMgx2CcW6zbgErkNgxvDnDPHlRl5aIIHNyDMlSczVtq0SlBrT
tExmjSg2Edzo1v25DG1LBzV58zYE5/crYh+Dm1XKB88sSBb8bUoAjZCWsl3Dkn8e
R8wOdabZZxU/STP/g9yxTM1fcnc4v4e/QF0TnQIDAQABAoIBAEC3SZw5KXQTVOAa
fxtgv8+UWVR09tB0I18AU36kd3DIbXooPM0l3adwWyYdD9v14h5s4fb5FG1VICKA
LFaZlNO/GjHL1CW8iuT2jl1E4y1baEUcojBBthAYwi810gpVUIrIWikQzc0ZqrFM
m/zk27Ro803TYTxn9UAIX1laTVPcRbjc6OwLiDoQC/Ftrec7SRU6lGqV60BI0sz0
2j8aLwqOHIIdt9gt0QFqmHau+4PXV++OrQapu2YzL5yTT+jePF9d1TFXOTEe2kjH
OY5qmSYC9+K4XgHKi8qDZTSyMpkzyGse3gFjstgR5/eUZKafWrW7nm1NwuxJJdUn
87qBK40CgYEA7mzeNN5fSXfr4kWTGAkQaaETvMKqsGYcsrVe+ZZzVtJfXnmVjRms
V5KT/mwoZvSURDp+rStQQuXVahpOHwCn0tDow/XJZUuQObemjTdH4g2RvScKGtHb
pB2RsgF7VKTfT6thGdLsejF3M289uIJq124YskzTAs9eaJr745dUPGMCgYEA+7H/
TEFFJCYb3D6hEgooSPETwOp7ErF2zNUHcpSBJ73Cxi9NLjuLShp2BaPEH7WlziAD
r35noRYlKyMqmxYEx298B7zY1mE2ctUsWVKkP2SInH73rvez2+Iapa0Z84w44wul
b9Sf56LgJNYelIA8+gAHJUKeHtyDVxCABFGCb/8CgYA99Y69QHiUuBRVpez21wwr
1w8xA4ml87NLgbSfuchZbKwZ+hCyLVTLIS1Sdbr+HlsVa/oVeGcQK3gNba6Vge8a
6u1CV3Ix37QoO6CNnCsTBKG1/Ro0JAsnGAQPtTDeq0XZB1lhg52ul4I5nJP2ifXH
7DWAyFQhq9AF8Ri6aU4brwKBgBKSm+ggmN2GAmBKLtCJ91cKkw6VPueuOLn8rkQC
OVWZZxoAu41Bz5F0Smk4IGzGlqmTKzJz/WmhnLSGL8qp4UhmLZzUjpujKMVofZFJ
y9zxqjMCG3zJwnfjQ1weXd/e5QO8BEUwR2xsVGXjdvY2UEmSXvSc6dYVJ4vxJ8Ep
0po5AoGAFbb0DOncdG+UbCvdrPXag8R9OB96+wxfjpeFfIHz8KdQC3cWYN5S9uLx
C5Ix+d+Y+XH1nQjP8XFZ8B2d9jKZ74uHyKeaQJwVZ3sFuxy1beZxS9b29u2QiplA
OeW30zIcJM7P+3uZEHo6GUMvh+WJdayxZmPUAEKdHmHudw/JCoE=
-----END RSA PRIVATE KEY-----`)

func init() {
	// overriding default goproxy certificate
	tlsc, err := tls.X509KeyPair(CA_CERT, CA_KEY)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Failed to load certificate and key pair")
	}
	goproxy.GoproxyCa = tlsc
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func main() {
	hoverfly := hv.NewHoverfly()

	flag.Var(&importFlags, "import", "Import from file or from URL (i.e. '-import my_service.json' or '-import http://mypage.com/service_x.json'")
	flag.Var(&destinationFlags, "dest", "Specify which hosts to process (i.e. '-dest fooservice.org -dest barservice.org -dest catservice.org') - other hosts will be ignored will passthrough'")
	flag.Var(&logOutputFlags, "logs-output", "Specify locations for output logs, options are \"console\" and \"file\" (default \"console\")")
	flag.StringVar(&responseBodyFilesPath, "response-body-files-path", "", "When a response contains a relative bodyFile, it will be resolved against this path (default is CWD)")
	flag.Var(&responseBodyFilesAllowedOriginFlags, "response-body-files-allow-origin", "When a response contains a url in bodyFile, it will be loaded only if the origin is allowed")

	flag.Parse()

	if *logsFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			ForceColors:      true,
			DisableTimestamp: false,
			FullTimestamp:    true,
			DisableColors:    *logNoColor,
		})
	}

	if *version {
		fmt.Println(hv.NewHoverfly().GetVersion())
		os.Exit(0)
	}

	if *journalSize < 0 {
		*journalSize = 0
	}

	if *logsSize < 0 {
		*logsSize = 0
	}

	if *cacheSize <= 0 {
		log.WithFields(log.Fields{
			"cache-size": *logsSize,
		}).Fatal("Cache size must be a positive number, alternatively use the disable-cache flag")
	}

	hoverfly.StoreLogsHook.LogsLimit = *logsSize
	hoverfly.Journal.EntryLimit = *journalSize

	// getting settings
	cfg := hv.InitSettings()

	logLevel, err := log.ParseLevel(*logLevelFlag)
	if err != nil {
		log.WithFields(log.Fields{
			"log-level": *logLevelFlag,
		}).Fatal("Unknown log-level value")
	}
	log.SetLevel(logLevel)

	if len(logOutputFlags) == 0 {
		// default logging on console when no flag given
		log.SetOutput(os.Stdout)
	} else {

		// remove duplicates
		logOutputMap := map[string]string{}
		for _, val := range logOutputFlags {
			logOutputMap[val] = val
		}
		_, isLogFile := logOutputMap["file"]
		if !isLogFile && isFlagPassed("logs-file") {
			log.WithFields(log.Fields{
				"logs-file": *logsFile,
			}).Fatal("-logs-file is not allowed unless -logs-output is set to 'file'.")
		}

		writers := make([]io.Writer, 0)
		for _, logsOutput := range logOutputFlags {
			if logsOutput == "file" {
				var formatter log.Formatter
				if *logsFormat == "json" {
					formatter = &log.JSONFormatter{}
				} else {
					formatter = &log.TextFormatter{
						ForceColors:      true,
						DisableTimestamp: false,
						FullTimestamp:    true,
						DisableColors:    true,
					}
				}
				logFileHook, err := util.NewLogFileHook(util.LogFileConfig{
					Filename:  *logsFile,
					Level:     logLevel,
					Formatter: formatter,
				})
				if err == nil {
					// add hook to write logs into file
					log.AddHook(logFileHook)
				} else {
					log.Fatal("Failed to write log file:" + *logsFile)
				}
			} else if logsOutput == "console" {
				writers = append(writers, os.Stdout)
			} else {
				log.WithFields(log.Fields{
					"logs-output": logsOutput,
				}).Fatal("Unknown logs output type")
			}
		}
		log.SetOutput(io.MultiWriter(writers...))
	}

	if *verbose {
		// Only log the warning severity or above.
		log.SetLevel(log.DebugLevel)
		log.Info("Log level set to verbose")
	}
	cfg.Verbose = *verbose

	if *dev {
		handlers.EnableCors = true
		handlers.CorsOrigin = *devCorsOrigin

		log.WithField("allowOrigin", *devCorsOrigin).Warn("Dev mode is enabled")
	}

	if *generateCA {
		tlsc, err := hvc.GenerateAndSave(*certName, *certOrg, 10*365*24*time.Hour)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed to generate certificate")
		}
		goproxy.GoproxyCa = *tlsc

	} else if *cert != "" && *key != "" {
		tlsc, err := tls.LoadX509KeyPair(*cert, *key)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed to load certificate and key pair")
		}

		goproxy.GoproxyCa = tlsc

		log.WithFields(log.Fields{
			"certificate": *cert,
			"key":         *key,
		}).Info("Default keys have been overwritten")
	}

	// overriding environment variables (proxy and admin ports)
	if *proxyPort != "" {
		cfg.ProxyPort = *proxyPort

		log.WithFields(log.Fields{
			"port": *proxyPort,
		}).Info("Default proxy port has been overwritten")
	}
	if *adminPort != "" {
		cfg.AdminPort = *adminPort

		log.WithFields(log.Fields{
			"port": *adminPort,
		}).Info("Default admin port has been overwritten")
	}

	if *listenOnHost != "" {
		cfg.ListenOnHost = *listenOnHost

		log.WithFields(log.Fields{
			"host": *listenOnHost,
		}).Info("Listen on specific interface")
	}

	// overriding environment variable (external proxy)
	if *upstreamProxy != "" {
		cfg.SetUpstreamProxy(*upstreamProxy)

		log.WithFields(log.Fields{
			"url": *upstreamProxy,
		}).Info("Upstream proxy has been set")
	}

	cfg.PlainHttpTunneling = *plainHttpTunneling

	if *cors {
		cfg.CORS = *cs.DefaultCORSConfigs()
		log.Info("CORS has been enabled")
	}

	if *noImportCheck {
		cfg.NoImportCheck = *noImportCheck
		log.Info("Import check has been disabled")
	}

	cfg.ClientAuthenticationDestination = *clientAuthenticationDestination
	cfg.ClientAuthenticationClientCert = *clientAuthenticationClientCert
	cfg.ClientAuthenticationClientKey = *clientAuthenticationClientKey
	cfg.ClientAuthenticationCACert = *clientAuthenticationCACert

	// overriding default middleware setting
	newMiddleware, err := mw.ConvertToNewMiddleware(*middleware)
	if err != nil {
		log.Error(err.Error())
	}
	cfg.Middleware = *newMiddleware

	mode := getInitialMode(cfg)

	// setting mode
	cfg.SetMode(mode)

	// disabling authentication if no-auth for auth disabled env variable
	if *authEnabled {
		cfg.AuthEnabled = true
	}

	// disabling tls verification if flag or env variable is set to 'false' (defaults to true)
	if !cfg.TLSVerification || !*tlsVerification {
		cfg.TLSVerification = false

		log.Info("TLS certificate verification has been disabled")
	}

	if len(destinationFlags) > 0 {
		cfg.Destination = strings.Join(destinationFlags[:], "|")

	} else {
		//  setting destination regexp
		cfg.Destination = *destination
	}

	cfg.ResponsesBodyFilesPath = responseBodyFilesPath

	for _, allowedOrigin := range responseBodyFilesAllowedOriginFlags {
		if !util.IsURL(allowedOrigin) {
			log.WithFields(log.Fields{"origin": allowedOrigin}).Fatal("Origin is not a valid url")
		}
	}

	cfg.ResponsesBodyFilesAllowedOrigins = responseBodyFilesAllowedOriginFlags

	var requestCache cache.FastCache
	var tokenCache cache.Cache
	var userCache cache.Cache

	if *databasePath != "" {
		cfg.DatabasePath = *databasePath
	}

	if *database == boltBackend {
		db := cache.GetDB(cfg.DatabasePath)
		defer db.Close()
		tokenCache = cache.NewBoltDBCache(db, []byte(backends.TokenBucketName))
		userCache = cache.NewBoltDBCache(db, []byte(backends.UserBucketName))

		log.Info("Using boltdb backend")
	} else if *database == inmemoryBackend {
		tokenCache = cache.NewInMemoryCache()
		userCache = cache.NewInMemoryCache()

		log.Info("Using memory backend")
	} else {
		log.WithFields(log.Fields{
			"database": *database,
		}).Fatal("Unknown database type")
	}
	cfg.DisableCache = *disableCache
	cfg.CacheSize = *cacheSize
	if cfg.DisableCache {
		log.Info("Request cache has been disabled")
	} else {
		// Request cache is always in-memory
		requestCache, err = cache.NewLRUCache(cfg.CacheSize)
		if err != nil {
			log.WithFields(log.Fields{
				"error":      err.Error(),
				"cache-size": cfg.CacheSize,
			}).Fatal("Failed to create cache")
		}
	}

	authBackend := backends.NewCacheBasedAuthBackend(tokenCache, userCache)

	hoverfly.Cfg = cfg
	hoverfly.CacheMatcher = matching.CacheMatcher{
		RequestCache: requestCache,
		Webserver:    cfg.Webserver,
	}
	hoverfly.Authentication = authBackend
	hoverfly.HTTP = hv.GetDefaultHoverflyHTTPClient(hoverfly.Cfg.TLSVerification, hoverfly.Cfg.UpstreamProxy)

	// if add new user supplied - adding it to database
	if *addNew || *authEnabled {
		var err error
		if *addPasswordHash != "" {
			err = hoverfly.Authentication.AddUserHashedPassword(*addUser, *addPasswordHash, *isAdmin)
		} else {
			err = hoverfly.Authentication.AddUser(*addUser, *addPassword, *isAdmin)
		}
		if err != nil {
			log.WithFields(log.Fields{
				"error":    err.Error(),
				"username": *addUser,
			}).Fatal("Failed to add new user")
		} else {
			log.WithFields(log.Fields{
				"username": *addUser,
			}).Info("User added successfully")
		}
		cfg.AuthEnabled = true
	}
	if cfg.AuthEnabled {
		if os.Getenv(hv.HoverflyAdminUsernameEV) != "" && os.Getenv(hv.HoverflyAdminPasswordEV) != "" {
			hoverfly.Authentication.AddUser(
				os.Getenv(hv.HoverflyAdminUsernameEV),
				os.Getenv(hv.HoverflyAdminPasswordEV),
				true)
		}

		// checking if there are any users
		users, err := hoverfly.Authentication.GetAllUsers()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed when retrieving users")
		}
		if len(users) < 1 {
			createSuperUser(hoverfly)
		}
	}

	// importing records if environment variable is set
	ev := os.Getenv(hv.HoverflyImportRecordsEV)
	if ev != "" {
		err := hoverfly.Import(ev)
		if err != nil {
			log.WithFields(log.Fields{
				"error":  err.Error(),
				"import": ev,
			}).Fatal("Environment variable for importing was set but failed to import this resource")
		}
	}

	// importing stuff
	if len(importFlags) > 0 {
		for _, v := range importFlags {
			if v != "" {
				log.WithFields(log.Fields{
					"import": v,
				}).Debug("Importing given resource")
				err := hoverfly.Import(v)
				if err != nil {
					log.WithFields(log.Fields{
						"error":  err.Error(),
						"import": v,
					}).Fatal("Failed to import given resource")
				}
			}
		}
		hoverfly.CacheMatcher.PreloadCache(hoverfly.Simulation)
	}

	// start metrics registry flush
	if *metrics {
		hoverfly.Counter.Init()
	}

	cfg.Webserver = *webserver

	if *pacFile != "" {
		pacFileContent, err := ioutil.ReadFile(*pacFile)
		if err != nil {
			log.WithFields(log.Fields{"error": err.Error(), "pacFile": *pacFile}).
				Fatal("Failed to import pac file")
		}

		log.WithField("pacFile", *pacFile).Infoln("Using provided pac file")
		hoverfly.SetPACFile(pacFileContent)
	}

	err = hoverfly.StartProxy()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Failed to start proxy")
	}

	// starting admin interface, this is blocking
	adminApi := hv.AdminApi{}
	adminApi.StartAdminInterface(hoverfly)
}

func createSuperUser(h *hv.Hoverfly) {
	reader := bufio.NewReader(os.Stdin)
	// Prompt and read
	fmt.Println("No users found in the database, please create initial user.")
	fmt.Print("Enter username (default hf): ")
	username, err := reader.ReadString('\n')
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Failed retrieving username input")
	}
	fmt.Print("Enter password (default hf): ")
	password, err := reader.ReadString('\n')
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Failed retrieving password input")
	}
	// Trim whitespace and use defaults if nothing entered
	username = strings.TrimSpace(username)
	if username == "" {
		username = "hf"
	}
	password = strings.TrimSpace(password)
	if password == "" {
		password = "hf"
	}
	err = h.Authentication.AddUser(username, password, true)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Failed to create user")
	} else {
		log.WithFields(log.Fields{
			"username": username,
		}).Info("User created")
	}
}

func getInitialMode(cfg *hv.Configuration) string {
	if *webserver {
		return modes.Simulate
	}

	if *capture {
		// checking whether user supplied other modes
		if *synthesize == true || *modify == true || *spy == true || *diff == true {
			log.Fatal("Two or more modes supplied, check your flags")
		}

		return modes.Capture

	} else if *synthesize {

		if !cfg.Middleware.IsSet() {
			log.Fatal("Synthesize mode chosen although middleware not supplied")
		}

		if *capture == true || *modify == true || *spy == true || *diff == true {
			log.Fatal("Two or more modes supplied, check your flags")
		}

		return modes.Synthesize

	} else if *modify {
		if !cfg.Middleware.IsSet() {
			log.Fatal("Modify mode chosen although middleware not supplied")
		}

		if *capture == true || *synthesize == true || *spy == true || *diff == true {
			log.Fatal("Two or more modes supplied, check your flags")
		}

		return modes.Modify
	} else if *spy {
		if *capture == true || *synthesize == true || *modify == true || *diff == true {
			log.Fatal("Two or more modes supplied, check your flags")
		}

		return modes.Spy
	} else if *diff {
		if *capture == true || *synthesize == true || *modify == true || *spy == true {
			log.Fatal("Two or more modes supplied, check your flags")
		}

		return modes.Diff
	}

	return modes.Simulate
}
