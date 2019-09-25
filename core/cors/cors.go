package cors

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Configs struct {
	Enabled          bool
	AllowOrigin      string
	AllowMethods     string
	AllowHeaders     string
	PreflightMaxAge  int64
	AllowCredentials bool
	ExposeHeaders    string
}

func DefaultCORSConfigs() *Configs {
	return &Configs{
		Enabled:          true,
		AllowOrigin:      "*",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS",
		AllowHeaders:     "Content-Type,Origin,Accept,Authorization,Content-Length,X-Requested-With",
		PreflightMaxAge:  1800,
		AllowCredentials: true,
		ExposeHeaders:    "",
	}
}

// TODO provide config to pass through OPTIONS call
// Intercept pre-flight request and return 200 response with CORS headers
func (c *Configs) InterceptPreflightRequest(r *http.Request) *http.Response {
	if r.Method != http.MethodOptions || r.Header.Get("Origin") == "" || r.Header.Get("Access-Control-Request-Method") == "" {
		return nil
	}
	resp := &http.Response{}
	resp.Request = r
	resp.Header = make(http.Header)
	resp.Header.Set("Access-Control-Allow-Origin", c.getAllowOrigin(r))
	resp.Header.Set("Access-Control-Allow-Methods", c.AllowMethods)
	resp.Header.Set("Access-Control-Max-Age", strconv.FormatInt(c.PreflightMaxAge, 10))
	allowHeaders := r.Header.Get("Access-Control-Request-Headers")

	if allowHeaders == "" {
		allowHeaders = c.AllowHeaders
	}
	resp.Header.Set("Access-Control-Allow-Headers", allowHeaders)

	if c.AllowCredentials {
		resp.Header.Set("Access-Control-Allow-Credentials", "true")
	}

	resp.StatusCode = http.StatusOK
	buf := bytes.NewBufferString("")
	resp.ContentLength = 0
	resp.Body = ioutil.NopCloser(buf)
	return resp
}

// Add CORS headers to the response if the request contains Origin header
func (c *Configs) AddCORSHeaders(r *http.Request, resp *http.Response) {
	if r.Header.Get("Origin") == "" || (resp.Header != nil && resp.Header.Get("Access-Control-Allow-Origin") != "") {
		return
	}

	if resp.Header == nil {
		resp.Header = make(http.Header)
	}

	resp.Header.Set("Access-Control-Allow-Origin", c.getAllowOrigin(r))

	if c.AllowCredentials {
		resp.Header.Set("Access-Control-Allow-Credentials", "true")
	}

	if c.ExposeHeaders != "" {
		resp.Header.Set("Access-Control-Expose-Headers", c.ExposeHeaders)
	}
}

// Safely get Allow-Origin header value. The value cannot be a wildcard while Allow-Credentials is enabled
func (c *Configs) getAllowOrigin(r *http.Request) string {
	allowOrigin := c.AllowOrigin
	if c.AllowCredentials && allowOrigin == "*" {
		allowOrigin = r.Header.Get("Origin")
	}
	return allowOrigin
}
