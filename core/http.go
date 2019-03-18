package hoverfly

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/jackwakefield/gopac"
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
			log.Fatalf("Could not parse upstream proxy: %s", err.Error())
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
		if client := parsePACFileResult(result, hf.Cfg.TLSVerification); client != nil {
			return client, nil
		}

	}

	if hf.Cfg.ClientAuthenticationDestination != "" {

		re := regexp.MustCompile(hf.Cfg.ClientAuthenticationDestination)

		if re.MatchString(host) {

			// Load client cert
			cert, err := tls.LoadX509KeyPair(
				hf.Cfg.ClientAuthenticationClientCert,
				hf.Cfg.ClientAuthenticationClientKey,
			)

			if err != nil {
				return nil, errors.New("Unable to load client certs file\n\n" + err.Error())
			}

			caCertPool := x509.NewCertPool()

			var tlsConfig *tls.Config

			if hf.Cfg.ClientAuthenticationCACert != "" {
				// Load CA cert
				caCert, err := ioutil.ReadFile(hf.Cfg.ClientAuthenticationCACert)

				if err != nil {
					return nil, errors.New("Unable to load ca certs file\n\n" + err.Error())
				}

				caCertPool.AppendCertsFromPEM(caCert)

				tlsConfig = &tls.Config{
					Certificates: []tls.Certificate{cert},
					RootCAs:      caCertPool,
				}
			} else {
				tlsConfig = &tls.Config{
					Certificates:       []tls.Certificate{cert},
					RootCAs:            caCertPool,
					InsecureSkipVerify: true,
				}
			}

			tlsConfig.BuildNameToCertificate()

			transport := &http.Transport{TLSClientConfig: tlsConfig}
			client := &http.Client{Transport: transport}

			return client, nil
		}
	}

	return hf.HTTP, nil
}

func parsePACFileResult(result string, tlsVerification bool) *http.Client {
	for _, s := range strings.Split(result, ";") {
		if s == "DIRECT" {
			return GetDefaultHoverflyHTTPClient(tlsVerification, "")
		}
		if s[0:6] == "PROXY " {
			return GetDefaultHoverflyHTTPClient(tlsVerification, s[6:])
		}
	}
	return nil
}
