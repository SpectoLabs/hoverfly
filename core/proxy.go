package hoverfly

import (
	"bufio"
	"encoding/base64"
	"net"
	"net/http"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/goproxy"
	"github.com/SpectoLabs/goproxy/ext/auth"
	"github.com/SpectoLabs/hoverfly/core/authentication"
	"github.com/SpectoLabs/hoverfly/core/authentication/backends"
	"github.com/SpectoLabs/hoverfly/core/util"
)

var disableProxyAuthoriation = false

var XHoverflyAuthorizationHeader = "X-HOVERFLY-AUTHORIZATION"
var ProxyAuthorizationHeader = "Proxy-Authorization"

// Creates goproxy.ProxyHttpServer and configures it to be used as a proxy for Hoverfly
// goproxy is given handlers that use the Hoverfly request processing
func NewProxy(hoverfly *Hoverfly) *goproxy.ProxyHttpServer {
	disableProxyAuthoriation = hoverfly.Cfg.DisableProxyAuthoriation

	// creating proxy
	proxy := goproxy.NewProxyHttpServer()

	proxy.OnRequest(goproxy.UrlMatches(regexp.MustCompile(hoverfly.Cfg.Destination))).
		HandleConnect(goproxy.AlwaysMitm)

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

	// enable curl -p for all hosts on port 80
	proxy.OnRequest(goproxy.UrlMatches(regexp.MustCompile(hoverfly.Cfg.Destination))).
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
	proxy.OnRequest(goproxy.UrlMatches(regexp.MustCompile(hoverfly.Cfg.Destination))).DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			resp := hoverfly.processRequest(r)
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
	proxy.OnResponse(goproxy.UrlMatches(regexp.MustCompile(hoverfly.Cfg.Destination))).DoFunc(
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
		log.Warn("NonproxyHandler")
		resp := hoverfly.processRequest(r)
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

func proxyBasicAndBearer(proxy *goproxy.ProxyHttpServer, realm string, basicFunc func(user, passwd string) bool, bearerFunc func(token string) bool) {

	proxy.OnRequest().Do(goproxy.FuncReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		if !authFromHeader(req, basicFunc, bearerFunc) {
			return nil, auth.BasicUnauthorized(req, realm)
		}
		return req, nil
	}))

	proxy.OnRequest().HandleConnect(goproxy.FuncHttpsHandler(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		if !authFromHeader(ctx.Req, basicFunc, bearerFunc) {
			ctx.Resp = auth.BasicUnauthorized(ctx.Req, realm)
			return goproxy.RejectConnect, host
		}
		return goproxy.MitmConnect, host
	}))
}

func authFromHeader(req *http.Request, basicFunc func(user, passwd string) bool, bearerFunc func(token string) bool) bool {
	headerValue := req.Header.Get(XHoverflyAuthorizationHeader)
	if headerValue == "" && !disableProxyAuthoriation {
		headerValue = req.Header.Get(ProxyAuthorizationHeader)
	}

	authheader := strings.SplitN(headerValue, " ", 2)
	req.Header.Del(ProxyAuthorizationHeader)
	if len(authheader) != 2 {
		return false
	}
	if authheader[0] == "Basic" {
		userpassraw, err := base64.StdEncoding.DecodeString(authheader[1])
		if err != nil {
			return false
		}
		userpass := strings.SplitN(string(userpassraw), ":", 2)
		if len(userpass) != 2 {
			return false
		}
		return basicFunc(userpass[0], userpass[1])
	} else if authheader[0] == "Bearer" {
		return bearerFunc(authheader[1])
	}

	return false
}
