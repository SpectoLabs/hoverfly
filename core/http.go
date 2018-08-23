package hoverfly

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/benjih/gopac"
)

func GetDefaultHoverflyHTTPClient(tlsVerification bool, upstreamProxy string) *http.Client {

	var proxyURL func(*http.Request) (*url.URL, error)
	if upstreamProxy == "" {
		proxyURL = http.ProxyURL(nil)
	} else {
		if upstreamProxy[0:4] != "http" {
			upstreamProxy = "http://" + upstreamProxy
		}
		u, err := url.Parse(upstreamProxy)
		if err != nil {
			log.Fatalf("Could not parse upstream proxy: ", err.Error())
		}
		proxyURL = http.ProxyURL(u)
	}

	return &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}, Transport: &http.Transport{
		Proxy: proxyURL,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !tlsVerification,
			Renegotiation:      tls.RenegotiateFreelyAsClient,
		},
	}}
}

func GetHttpClient(hf *Hoverfly, host string) *http.Client {
	if hf.Cfg.PACFile != nil {
		parser := new(gopac.Parser)
		if err := parser.ParseBytes(hf.Cfg.PACFile); err != nil {
			log.Fatalf("Failed to parse PAC (%s)", err)
		}

		result, err := parser.FindProxy("", host)

		if err != nil {
			log.Fatalf("Failed to find proxy entry (%s)", err)
		}

		for _, s := range strings.Split(result, ";") {
			if s == "DIRECT" {
				log.Println("DIRECT")
				return GetDefaultHoverflyHTTPClient(hf.Cfg.TLSVerification, "")
			}
			if s[0:6] == "PROXY " {
				return GetDefaultHoverflyHTTPClient(hf.Cfg.TLSVerification, s[6:])
			}
		}
	}
	return hf.HTTP
}
