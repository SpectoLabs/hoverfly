package hoverfly

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/SpectoLabs/goproxy"
	"github.com/SpectoLabs/goproxy/ext/auth"
	"github.com/SpectoLabs/hoverfly/core/authentication"
	"github.com/SpectoLabs/hoverfly/core/authentication/backends"
	"github.com/SpectoLabs/hoverfly/core/util"
)

var ProxyAuthorizationHeader string

// Creates goproxy.ProxyHttpServer and configures it to be used as a proxy for Hoverfly
// goproxy is given handlers that use the Hoverfly request processing
func NewProxy(hoverfly *Hoverfly) *goproxy.ProxyHttpServer {
	ProxyAuthorizationHeader = hoverfly.Cfg.ProxyAuthorizationHeader

	// creating proxy
	proxy := goproxy.NewProxyHttpServer()

	proxy.OnRequest(matchesFilter(hoverfly.Cfg.Destination)).
		HandleConnect(goproxy.FuncHttpsHandler(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
			if hoverfly.Cfg.PlainHttpTunneling && !strings.HasSuffix(host, ":443") {
				return goproxy.HTTPMitmConnect, host
			}
			return goproxy.MitmConnect, host
		}))

	if hoverfly.Cfg.HttpsOnly {
		log.Info("Disabling HTTP")
		proxy.DisableNonTls(true)
	}

	if hoverfly.Cfg.AuthEnabled {
		log.Info("Enabling proxy authentication")
		proxyBasicAndBearer(proxy, "hoverfly", func(user, password string) bool {

			proxyUser := &backends.User{
				Username: user,
				Password: password,
			}

			responseStatus, _ := authentication.Login(proxyUser, hoverfly.Authentication, nil, 0)

			return responseStatus == http.StatusOK
		}, func(headerToken string) bool {
			return authentication.IsJwtTokenValid(headerToken, hoverfly.Authentication, hoverfly.Cfg.SecretKey, hoverfly.Cfg.JWTExpirationDelta)
		})
	}

	// processing connections
	proxy.OnRequest(matchesFilter(hoverfly.Cfg.Destination)).DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			startTime := time.Now()
			resp := hoverfly.processRequest(r)
			hoverfly.Journal.NewEntry(r, resp, hoverfly.Cfg.Mode, startTime)
			return r, resp
		})

	if hoverfly.Cfg.Verbose {
		proxy.OnRequest().DoFunc(
			func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
				log.WithFields(log.Fields{
					"destination": r.Host,
					"path":        r.URL.Path,
					"query":       r.URL.RawQuery,
					"method":      r.Method,
					"mode":        hoverfly.Cfg.GetMode(),
				}).Debug("got request..")
				return r, nil
			})
	}

	// intercepts response
	proxy.OnResponse(matchesFilter(hoverfly.Cfg.Destination)).DoFunc(
		func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			hoverfly.Counter.Count(hoverfly.Cfg.GetMode())
			return resp
		})

	proxy.Verbose = hoverfly.Cfg.Verbose
	// proxy starting message
	log.WithFields(log.Fields{
		"Destination": hoverfly.Cfg.Destination,
		"ProxyPort":   hoverfly.Cfg.ProxyPort,
		"Mode":        hoverfly.Cfg.GetMode(),
	}).Info("Proxy prepared...")

	return proxy
}

// Creates goproxy.ProxyHttpServer and configures it to be used as a webserver for Hoverfly
// goproxy is given a non proxy handler that uses the Hoverfly request processing
func NewWebserverProxy(hoverfly *Hoverfly) *goproxy.ProxyHttpServer {
	// creating proxy
	proxy := goproxy.NewProxyHttpServer()
	proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		r.URL.Scheme = "http"
		resp := hoverfly.processRequest(r)
		hoverfly.Journal.NewEntry(r, resp, hoverfly.Cfg.Mode, startTime)
		body, err := util.GetResponseBody(resp)

		if err != nil {
			log.Error("Error reading response body")
			w.WriteHeader(500)
			return
		}

		for name, values := range resp.Header {
			name = strings.ToLower(name)

			for _, value := range values {
				w.Header().Add(name, value)
			}
		}

		w.Header().Set("Req", r.RequestURI)
		w.Header().Set("Resp", resp.Header.Get("Content-Length"))
		w.WriteHeader(resp.StatusCode)
		w.Write([]byte(body))

		hoverfly.Counter.Count(hoverfly.Cfg.GetMode())
	})

	if hoverfly.Cfg.Verbose {
		proxy.OnRequest().DoFunc(
			func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
				log.WithFields(log.Fields{
					"destination": r.Host,
					"path":        r.URL.Path,
					"query":       r.URL.RawQuery,
					"method":      r.Method,
					"mode":        hoverfly.Cfg.GetMode(),
				}).Debug("got request..")
				return r, nil
			})
	}

	log.WithFields(log.Fields{
		"Destination":   hoverfly.Cfg.Destination,
		"WebserverPort": hoverfly.Cfg.ProxyPort,
		"Mode":          hoverfly.Cfg.GetMode(),
	}).Info("Webserver prepared...")

	return proxy
}

func unauthorizedError(request *http.Request, realm, message string) *http.Response {
	response := auth.BasicUnauthorized(request, realm)
	response.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(message)))
	response.ContentLength = int64(len(message))

	return response
}

func proxyBasicAndBearer(proxy *goproxy.ProxyHttpServer, realm string, basicFunc func(user, passwd string) bool, bearerFunc func(token string) bool) {

	proxy.OnRequest().Do(goproxy.FuncReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		err := authFromHeader(req, basicFunc, bearerFunc)
		if err != nil {
			return nil, unauthorizedError(req, realm, err.Error())
		}
		return req, nil
	}))

	proxy.OnRequest().HandleConnect(goproxy.FuncHttpsHandler(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		err := authFromHeader(ctx.Req, basicFunc, bearerFunc)
		if err != nil {
			ctx.Resp = unauthorizedError(ctx.Req, realm, err.Error())
			return goproxy.RejectConnect, host
		}
		return goproxy.MitmConnect, host
	}))
}

func authFromHeader(req *http.Request, basicFunc func(user, passwd string) bool, bearerFunc func(token string) bool) error {
	headerValue := req.Header.Get(ProxyAuthorizationHeader)

	if ProxyAuthorizationHeader != "Proxy-Authorization" && req.Header.Get("Proxy-Authorization") != "" {
		return fmt.Errorf("407 `Proxy-Authorization` header is disabled, use `X-HOVERFLY-AUTHORIZATION` instead")
	}

	authheader := strings.SplitN(headerValue, " ", 2)
	req.Header.Del(ProxyAuthorizationHeader)
	if len(authheader) != 2 {
		return fmt.Errorf("407 Proxy authentication required")
	}
	if authheader[0] == "Basic" {
		userpassraw, err := base64.StdEncoding.DecodeString(authheader[1])
		if err != nil {
			return fmt.Errorf("407 Proxy authentication required")
		}
		userpass := strings.SplitN(string(userpassraw), ":", 2)
		if len(userpass) != 2 {
			return fmt.Errorf("407 Proxy authentication required")
		}
		result := basicFunc(userpass[0], userpass[1])
		if result == false {
			return fmt.Errorf("407 Proxy authentication required")
		}
	} else if authheader[0] == "Bearer" {
		result := bearerFunc(authheader[1])
		if result == false {
			return fmt.Errorf("407 Proxy authentication required")
		}
	} else {
		return fmt.Errorf("407 Unknown authentication type `%v`, only `Basic` or `Bearer` are supported", authheader[0])
	}

	return nil
}

func matchesFilter(filter string) goproxy.ReqConditionFunc {
	re := regexp.MustCompile(filter)
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {

		scheme := req.URL.Scheme
		host := req.Host		// relative URL does not have host value, and this is a safer way to get the hostname from request struct
		path := req.URL.Path

		if scheme == "https" || strings.Contains(host, ":443") {
			scheme = "https"
			host = strings.Replace(host, ":443", "", 1)
		}

		if req.Method == http.MethodConnect {
			parsed, _ := url.Parse(filter)
			re = regexp.MustCompile(parsed.Scheme + "://" + parsed.Host)
		}

		return re.MatchString(path) ||
			re.MatchString(host+path) ||
			re.MatchString(scheme+"://"+host+path)

	}
}
