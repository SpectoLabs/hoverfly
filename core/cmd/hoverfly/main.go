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
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/goproxy"
	hv "github.com/SpectoLabs/hoverfly/core"
	"github.com/SpectoLabs/hoverfly/core/authentication/backends"
	"github.com/SpectoLabs/hoverfly/core/cache"
	hvc "github.com/SpectoLabs/hoverfly/core/certs"
	"github.com/SpectoLabs/hoverfly/core/modes"
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

const boltBackend = "boltdb"
const inmemoryBackend = "memory"

var (
	hoverflyVersion string

	version     = flag.Bool("version", false, "get the version of hoverfly")
	verbose     = flag.Bool("v", false, "should every proxy request be logged to stdout")
	capture     = flag.Bool("capture", false, "start Hoverfly in capture mode - transparently intercepts and saves requests/response")
	synthesize  = flag.Bool("synthesize", false, "start Hoverfly in synthesize mode (middleware is required)")
	modify      = flag.Bool("modify", false, "start Hoverfly in modify mode - applies middleware (required) to both outgoing and incomming HTTP traffic")
	middleware  = flag.String("middleware", "", "should proxy use middleware")
	proxyPort   = flag.String("pp", "", "proxy port - run proxy on another port (i.e. '-pp 9999' to run proxy on port 9999)")
	adminPort   = flag.String("ap", "", "admin port - run admin interface on another port (i.e. '-ap 1234' to run admin UI on port 1234)")
	metrics     = flag.Bool("metrics", false, "supply -metrics flag to enable metrics logging to stdout")
	dev         = flag.Bool("dev", false, "supply -dev flag to serve directly from ./static/dist instead from statik binary")
	destination = flag.String("destination", ".", "destination URI to catch")
	webserver   = flag.Bool("webserver", false, "start Hoverfly in webserver mode (simulate mode)")

	addNew      = flag.Bool("add", false, "add new user '-add -username hfadmin -password hfpass'")
	addUser     = flag.String("username", "", "username for new user")
	addPassword = flag.String("password", "", "password for new user")
	isAdmin     = flag.Bool("admin", true, "supply '-admin false' to make this non admin user (defaults to 'true') ")
	authEnabled = flag.Bool("auth", false, "enable authentication, currently it is disabled by default")

	generateCA = flag.Bool("generate-ca-cert", false, "generate CA certificate and private key for MITM")
	certName   = flag.String("cert-name", "hoverfly.proxy", "cert name")
	certOrg    = flag.String("cert-org", "Hoverfly Authority", "organisation name for new cert")
	cert       = flag.String("cert", "", "CA certificate used to sign MITM certificates")
	key        = flag.String("key", "", "private key of the CA used to sign MITM certificates")

	tlsVerification = flag.Bool("tls-verification", true, "turn on/off tls verification for outgoing requests (will not try to verify certificates) - defaults to true")

	upstreamProxy = flag.String("upstream-proxy", "", "specify an upstream proxy for hoverfly to route traffic through")

	databasePath = flag.String("db-path", "", "database location - supply it to provide specific database location (will be created there if it doesn't exist)")
	database     = flag.String("db", inmemoryBackend, "Persistance storage to use - 'boltdb' or 'memory' which will not write anything to disk")
)

var CA_CERT = []byte(`-----BEGIN CERTIFICATE-----
MIIDkDCCAnigAwIBAgIVAPVhYkM0BEM/yrYlqXluHt7cc6l1MA0GCSqGSIb3DQEB
CwUAMDYxGzAZBgNVBAoTEkhvdmVyZmx5IEF1dGhvcml0eTEXMBUGA1UEAxMOaG92
ZXJmbHkucHJveHkwHhcNMTUwMzI1MTM1NjA4WhcNMTcwMzI0MTM1NjA4WjA2MRsw
GQYDVQQKExJIb3ZlcmZseSBBdXRob3JpdHkxFzAVBgNVBAMTDmhvdmVyZmx5LnBy
b3h5MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyVHPS3AoW7GSExp4
F4b6rofOpCFCk9oyALOqifcLgMqfa+xjzHa9HH7yraT5EPKieTBm+XrJUnWUih+g
klKKvQYUWSx+W/+5LvFI/ZeOzBnBx9ZRlZNGSu613G430GZp3ydbY18wyhDlH3Xc
EmhDEHxBX+OmSj1cLMPFqYhbsA5I79evpSafHQ6vIUcy8tZqIj7vGgpssULLq3K9
Fnbexf8AFkaaRwx/iz3XBXfubrAzjYhr+B57/davpJGu3qkiRBhgWkMO0OHJmFIt
iIuE8Mg6yEyYd1gJdS0zQFa7FRBOAtJbiEnZfR/MFS3DdptIgyOH9f3iFn/Ad9Lv
JzWg3wIDAQABo4GUMIGRMA4GA1UdDwEB/wQEAwICpDATBgNVHSUEDDAKBggrBgEF
BQcDATAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBRMLsxq3r2ownouli5b4BAD
vjFC9DAfBgNVHSMEGDAWgBRMLsxq3r2ownouli5b4BADvjFC9DAZBgNVHREEEjAQ
gg5ob3ZlcmZseS5wcm94eTANBgkqhkiG9w0BAQsFAAOCAQEAYPu97vNZekWN80yA
zxTrakcp0ymcPraZT9mv2+tpicZ9rEa5QVIA4npACFaYmynhO/lyYgHjmBOpy+tX
KhhO7R6tYJaodVY55/B6/yj/bnAUa67kdMjtcb/lZCZX+cSBaUyuQ3xMHq1JF4sk
q09xF61TphpHdjTApQsjocJpNCxw2Ou3ctCUnSsuq6oC597CsnzKZTFJsqs4LO03
lZ7F4bw8LyZl52zzwVwbKnNszAq6ClUyfjVonO85DpBQhq+gFWlSVz4CCC62mig9
nnOsZC8mHNQrE0gHgTmKQlNwUQE7c+cUpfa9sjdlCG6eDkF4OaLPrG975WFuVXTM
V6/wDA==
-----END CERTIFICATE-----`)

var CA_KEY = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAyVHPS3AoW7GSExp4F4b6rofOpCFCk9oyALOqifcLgMqfa+xj
zHa9HH7yraT5EPKieTBm+XrJUnWUih+gklKKvQYUWSx+W/+5LvFI/ZeOzBnBx9ZR
lZNGSu613G430GZp3ydbY18wyhDlH3XcEmhDEHxBX+OmSj1cLMPFqYhbsA5I79ev
pSafHQ6vIUcy8tZqIj7vGgpssULLq3K9Fnbexf8AFkaaRwx/iz3XBXfubrAzjYhr
+B57/davpJGu3qkiRBhgWkMO0OHJmFItiIuE8Mg6yEyYd1gJdS0zQFa7FRBOAtJb
iEnZfR/MFS3DdptIgyOH9f3iFn/Ad9LvJzWg3wIDAQABAoIBAASQlFCzlFav6g4A
1aRC7UAz2B2km2va0LNvX3iNX3dmIMNDsueZ8aPJxRrm2LbnqYNx84PIovP5soqH
OQ7YTEkI8EEtXxga7koAMpV9cEF0fA5Z77OiiT99tiXvYdiZ2eCzdcEFEYgjZe6W
r4zDTHH9P0Y7VTPtvD9PmRXE/784KWLfREf6pwuICtBf2cLXPrq7mIZaV5Tf2od1
k6PHEDSaVZ5rLWgUGzmPyFiyC7phcFpck8JEGM4mp5YEdVwL3j/F8n8pcgDYQd1R
jf2CKtzeGNeQ3mTca8cAlqHmmbB6fUXz6SGK91x1O60tXWMJHv2A6n6wK6k/oOFn
q54ZzOECgYEA8R+JfS1/0lAlxJR8nl/kE0YNmuTpc5ZmHI0jm1Tfz3WipWSWCktI
BzgpygcHdjk9YzGO7d5lO5XqBccP9ntI9+6/ULwpTgusYR+O+OEHuDo50DZACgUP
wqOnsfx37cXd9z6Dp6+xR9c82pqcxEIjBRgZ7gdRRPRndR1lmXTn7JkCgYEA1b2W
WpCkjmfFBgrIJSpf+t6MavTksDILTQpOmRTZj4svhlJ0/EzSBwZqNDCtCYj5D8pq
I67wMeRmVa3zh2AQLdc1hqP3yVWHohWh50rN0FMpKztFwizjNrlXDcHDTwS+RMYr
t2K+UEhpS1XdY76L5fAMs9d5j3OH5/HvDbmjrDcCgYEAiR+cOtnjNSFrOQ4QiKiT
tfpCxnGj6Z4AWABT3YQ4+2w0oMZBJX2GasSfz0qMDcmjhYOres7c1zP8MGjyRQP7
jTPzDODUxJOS5nDiB9tBXp2OP0B6zrfuLIyRU4D2WvwJrQ+aI4Sg1vAqpU8EFABg
lgcMx/bVWtd69nlPTCPVuRECgYADFdN/xyq464KKjclJ0AzGoEPCn3pVmMNU/1sX
Fpf1XHr5I2OQ6ML3Wv5ZdoJo6tM9iRxzG2lYLwXTIsmrIJXbM4oQQXmoLFXi3xER
N6E06p5jg12EagV1msNI7Y0WLOlaMMocwY4htonejoS9ldiLHyXvyqJ0kaRaksFy
n0VfjQKBgQDe8rz2fM4F4ZeaCi75LCik8XCpA//8DquYrtwz+ojM9fK7Z28N8Vir
G/COvEx6J0CycRYFzUUxNWOpFIONCgLQkNEaGppBPkZ/aLqZzeOsakv9dGxVsm2W
gYLP4o2Hv9odPkRDyOasEJ5wIjnk8Aj2fmD34TAVKiSFkwguW9QbHg==
-----END RSA PRIVATE KEY-----
`)

func init() {
	// overriding default goproxy certificate
	tlsc, err := tls.X509KeyPair(CA_CERT, CA_KEY)
	if err != nil {
		log.Fatalf("Failed to load certifiate and key pair, got error: %s", err.Error())
	}
	goproxy.GoproxyCa = tlsc
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	flag.Var(&importFlags, "import", "import from file or from URL (i.e. '-import my_service.json' or '-import http://mypage.com/service_x.json'")
	flag.Var(&destinationFlags, "dest", "specify which hosts to process (i.e. '-dest fooservice.org -dest barservice.org -dest catservice.org') - other hosts will be ignored will passthrough'")
	flag.Parse()

	if *version {
		fmt.Println(hoverflyVersion)
		os.Exit(0)
	}

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

	if *generateCA {
		tlsc, err := hvc.GenerateAndSave(*certName, *certOrg, 365*24*time.Hour)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("failed to generate certificate.")
		}
		goproxy.GoproxyCa = *tlsc

	} else if *cert != "" && *key != "" {
		tlsc, err := tls.LoadX509KeyPair(*cert, *key)
		if err != nil {
			log.Fatalf("Failed to load certifiate and key pair, got error: %s", err.Error())
		}

		log.WithFields(log.Fields{
			"certificate": *cert,
			"key":         *key,
		}).Info("Default keys have been overwritten")

		goproxy.GoproxyCa = tlsc

	}

	// overriding environment variables (proxy and admin ports)
	if *proxyPort != "" {
		cfg.ProxyPort = *proxyPort
	}
	if *adminPort != "" {
		cfg.AdminPort = *adminPort
	}

	// overriding environment variable (external proxy)
	if *upstreamProxy != "" {
		cfg.SetUpstreamProxy(*upstreamProxy)
	}

	// development settings
	cfg.Development = *dev

	// overriding default middleware setting
	newMiddleware, err := hv.ConvertToNewMiddleware(*middleware)
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
		log.Info("tls certificate verification is now turned off!")
	}

	if len(destinationFlags) > 0 {
		cfg.Destination = strings.Join(destinationFlags[:], "|")

	} else {
		//  setting destination regexp
		cfg.Destination = *destination
	}

	var requestCache cache.Cache
	var metadataCache cache.Cache
	var tokenCache cache.Cache
	var userCache cache.Cache

	if *databasePath != "" {
		cfg.DatabasePath = *databasePath
	}

	if *database == boltBackend {
		log.Info("Creating bolt db backend...")
		db := cache.GetDB(cfg.DatabasePath)
		defer db.Close()
		requestCache = cache.NewBoltDBCache(db, []byte("requestsBucket"))
		metadataCache = cache.NewBoltDBCache(db, []byte("metadataBucket"))
		tokenCache = cache.NewBoltDBCache(db, []byte(backends.TokenBucketName))
		userCache = cache.NewBoltDBCache(db, []byte(backends.UserBucketName))
	} else if *database == inmemoryBackend {
		log.Info("Creating in memory map backend...")
		log.Warn("Turning off authentication...")

		requestCache = cache.NewInMemoryCache()
		metadataCache = cache.NewInMemoryCache()
		tokenCache = cache.NewInMemoryCache()
		userCache = cache.NewInMemoryCache()
	} else {
		log.Fatalf("unknown database type chosen: %s", *database)
	}

	authBackend := backends.NewCacheBasedAuthBackend(tokenCache, userCache)

	hoverfly := hv.GetNewHoverfly(cfg, requestCache, metadataCache, authBackend)

	// if add new user supplied - adding it to database
	if *addNew || *authEnabled {
		err := hoverfly.Authentication.AddUser(*addUser, *addPassword, *isAdmin)
		if err != nil {
			log.WithFields(log.Fields{
				"error":    err.Error(),
				"username": *addUser,
			}).Fatal("failed to add new user")
		} else {
			log.WithFields(log.Fields{
				"username": *addUser,
			}).Info("user added successfuly")
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
				"error": err,
			}).Fatal("got error while trying to get all users")
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
		} else {
			err = hoverfly.MetadataCache.Set([]byte("import_from_env_variable"), []byte(ev))
		}
	}

	// importing stuff
	if len(importFlags) > 0 {
		for i, v := range importFlags {
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
				} else {
					err = hoverfly.MetadataCache.Set([]byte(fmt.Sprintf("import_%d", i+1)), []byte(v))
				}
			}
		}
	}

	// start metrics registry flush
	if *metrics {
		hoverfly.Counter.Init()
	}

	cfg.Webserver = *webserver

	err = hoverfly.StartProxy()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("failed to start proxy...")
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
			"error": err,
		}).Fatal("error while getting username input")
	}
	fmt.Print("Enter password (default hf): ")
	password, err := reader.ReadString('\n')
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("error while getting password input")
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
			"error": err,
		}).Fatal("failed to create user.")
	} else {
		log.Infof("User: '%s' created.\n", username)
	}
}

func getInitialMode(cfg *hv.Configuration) string {
	if *webserver {
		return modes.Simulate
	}

	if *capture {
		// checking whether user supplied other modes
		if *synthesize == true || *modify == true {
			log.Fatal("Two or more modes supplied, check your flags")
		}

		return modes.Capture

	} else if *synthesize {

		if !cfg.Middleware.IsSet() {
			log.Fatal("Synthesize mode chosen although middleware not supplied")
		}

		if *capture == true || *modify == true {
			log.Fatal("Two or more modes supplied, check your flags")
		}

		return modes.Synthesize

	} else if *modify {
		if !cfg.Middleware.IsSet() {
			log.Fatal("Modify mode chosen although middleware not supplied")
		}

		if *capture == true || *synthesize == true {
			log.Fatal("Two or more modes supplied, check your flags")
		}

		return modes.Modify
	}

	return modes.Simulate
}
