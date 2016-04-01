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
	"crypto/tls"
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	hv "github.com/SpectoLabs/hoverfly"
	"github.com/SpectoLabs/hoverfly/authentication/backends"
	"github.com/SpectoLabs/hoverfly/cache"
	hvc "github.com/SpectoLabs/hoverfly/certs"
	"github.com/rusenask/goproxy"
	"strings"
	"time"
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

var (
	verbose         = flag.Bool("v", false, "should every proxy request be logged to stdout")
	capture         = flag.Bool("capture", false, "start Hoverfly in capture mode - transparently intercepts and saves requests/response")
	synthesize      = flag.Bool("synthesize", false, "start Hoverfly in synthesize mode (middleware is required)")
	modify          = flag.Bool("modify", false, "start Hoverfly in modify mode - applies middleware (required) to both outgoing and incomming HTTP traffic")
	middleware      = flag.String("middleware", "", "should proxy use middleware")
	proxyPort       = flag.String("pp", "", "proxy port - run proxy on another port (i.e. '-pp 9999' to run proxy on port 9999)")
	adminPort       = flag.String("ap", "", "admin port - run admin interface on another port (i.e. '-ap 1234' to run admin UI on port 1234)")
	metrics         = flag.Bool("metrics", false, "supply -metrics flag to enable metrics logging to stdout")
	dev             = flag.Bool("dev", false, "supply -dev flag to serve directly from ./static/dist instead from statik binary")
	destination     = flag.String("destination", ".", "destination URI to catch")
	addNew          = flag.Bool("add", false, "add new user '-add -username hfadmin -password hfpass'")
	addUser         = flag.String("username", "", "username for new user")
	addPassword     = flag.String("password", "", "password for new user")
	isAdmin         = flag.Bool("admin", true, "supply '-admin false' to make this non admin user (defaults to 'true') ")
	authEnabled     = flag.Bool("auth", false, "enable authentication, currently it is disabled by default")
	generateCA      = flag.Bool("generate-ca-cert", false, "generate CA certificate and private key for MITM")
	certName        = flag.String("cert-name", "hoverfly.proxy", "cert name")
	certOrg         = flag.String("cert-org", "Hoverfly Authority", "organisation name for new cert")
	cert            = flag.String("cert", "", "CA certificate used to sign MITM certificates")
	key             = flag.String("key", "", "private key of the CA used to sign MITM certificates")
	tlsVerification = flag.Bool("tls-verification", true, "turn on/off tls verification for outgoing requests (will not try to verify certificates) - defaults to true")
	databasePath    = flag.String("db-dir", "", "database location - supply it if you want to provide specific to database (will be created there if it doesn't exist)")
	database        = flag.String("db", "boltdb", "Persistance storage to use - 'boltdb' or 'memory' which will not write anything to disk")
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	flag.Var(&importFlags, "import", "import from file or from URL (i.e. '-import my_service.json' or '-import http://mypage.com/service_x.json'")
	flag.Var(&destinationFlags, "dest", "specify which hosts to process (i.e. '-dest fooservice.org -dest barservice.org -dest catservice.org') - other hosts will be ignored will passthrough'")
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

		goproxy.GoproxyCa = tlsc

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

	// setting mode
	cfg.SetMode(mode)

	// enabling authentication if flag or env variable is set to 'true'
	if cfg.AuthEnabled || *authEnabled {
		cfg.AuthEnabled = true
	}

	// disabling tls verification if flag or env variable is set to 'false' (defaults to true)
	if !cfg.TlsVerification || !*tlsVerification {
		cfg.TlsVerification = false
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

	if *database == "boltdb" {
		log.Info("Creating bolt db backend...")
		db := cache.GetDB(*databasePath)
		defer db.Close()
		requestCache = cache.NewBoltDBCache(db, []byte("requestsBucket"))
		metadataCache = cache.NewBoltDBCache(db, []byte("metadataBucket"))
		tokenCache = cache.NewBoltDBCache(db, []byte(backends.TokenBucketName))
		userCache = cache.NewBoltDBCache(db, []byte(backends.UserBucketName))
	} else {
		log.Info("Creating in memory map backend...")
		requestCache = cache.NewInMemoryCache()
		metadataCache = cache.NewInMemoryCache()
		tokenCache = cache.NewInMemoryCache()
		userCache = cache.NewInMemoryCache()
	}

	hoverfly := hv.NewHoverfly(cfg, requestCache, metadataCache, backends.NewAuthBackend(tokenCache, userCache))

	// if add new user supplied - adding it to database
	if *addNew {
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
		return
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

	err := hoverfly.StartProxy()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("failed to start proxy...")
	}

	// starting admin interface, this is blocking
	hoverfly.StartAdminInterface()
}
