package cors

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Configs struct {
	Enabled 			bool
	AllowOrigin 		string
	AllowMethods 		string
	AllowHeaders 		string
	PreflightMaxAge 	int64
	AllowCredentials 	bool
	ExposeHeaders 		string
}

func DefaultCORSConfigs() *Configs {
	return &Configs {
		Enabled: true,
		AllowOrigin: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS",
		AllowHeaders: "Content-Type,Origin,Accept,Authorization,Content-Length,X-Requested-With",
		PreflightMaxAge: 1800,
		AllowCredentials: true,
		ExposeHeaders: "",
	}
}

func (c *Configs) InterceptPreflightRequest(r *http.Request) *http.Response {
	if r.Method != http.MethodOptions || r.Header.Get("Origin") == "" || r.Header.Get("Access-Control-Request-Methods") == "" {
		return nil
	}
	resp := &http.Response{}
	resp.Request = r
	resp.Header = make(http.Header)
	resp.Header.Add("Access-Control-Allow-Origin", c.AllowOrigin)
	resp.Header.Add("Access-Control-Allow-Methods", c.AllowMethods)
	resp.Header.Add("Access-Control-Max-Age", strconv.FormatInt(c.PreflightMaxAge, 10))
	if r.Header.Get("Access-Control-Request-Headers") == "" {
		resp.Header.Add("Access-Control-Allow-Headers", c.AllowHeaders)
	} else {
		resp.Header.Add("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
	}
	if c.AllowCredentials {
		resp.Header.Add("Access-Control-Allow-Credentials", "true")
	}
	resp.StatusCode = http.StatusOK
	buf := bytes.NewBufferString("")
	resp.ContentLength = 0
	resp.Body = ioutil.NopCloser(buf)
	return resp
}

func (c *Configs) AddCORSHeaders(r *http.Request, resp *http.Response) {
	if r.Header.Get("Origin") == "" {
		return
	}

	if resp.Header == nil {
		resp.Header = make(http.Header)
	}
	resp.Header.Add("Access-Control-Allow-Origin", c.AllowOrigin)
	resp.Header.Add("Access-Control-Expose-Headers", c.ExposeHeaders)
	if c.AllowCredentials {
		resp.Header.Add("Access-Control-Allow-Credentials", "true")
	}
}
