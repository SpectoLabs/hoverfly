package hoverfly

import (
	"crypto/tls"
	"errors"
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

func GetHttpClient(hf *Hoverfly, host string) (*http.Client, error) {
	if hf.Cfg.PACFile != nil {
		parser := new(gopac.Parser)
		if err := parser.ParseBytes(hf.Cfg.PACFile); err != nil {
			return nil, errors.New("Unable to parse PAC file\n\n" + err.Error())
		}

		result, err := parser.FindProxy("", host)
		if err != nil {
			return nil, errors.New("Unable to parse PAC file\n\n" + err.Error())
		}

		for _, s := range strings.Split(result, ";") {
			if s == "DIRECT" {
				log.Println("DIRECT")
				return GetDefaultHoverflyHTTPClient(hf.Cfg.TLSVerification, ""), nil
			}
			if s[0:6] == "PROXY " {
				return GetDefaultHoverflyHTTPClient(hf.Cfg.TLSVerification, s[6:]), nil
			}
		}
	}
	return hf.HTTP, nil
}
