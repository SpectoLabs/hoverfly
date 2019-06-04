package cors

import (
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"testing"
)

func Test_InterceptNewPreflightRequest_ReturnSuccessResponseWithDefaultHeaders(t *testing.T) {
	RegisterTestingT(t)

	unit := DefaultCORSConfigs()

	r, err := http.NewRequest(http.MethodOptions, "http://somehost.com", nil)
	Expect(err).To(BeNil())
	r.Header.Set("Origin", "http://originhost.com")
	r.Header.Set("Access-Control-Request-Methods", "PUT,POST")
	resp := unit.InterceptNewPreflightRequest(r)

	Expect(resp).ToNot(BeNil())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))
	Expect(resp.Header.Get("Access-Control-Allow-Origin")).To(Equal("*"))
	Expect(resp.Header.Get("Access-Control-Allow-Methods")).To(Equal("GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS"))
	Expect(resp.Header.Get("Access-Control-Max-Age")).To(Equal("1800"))
	Expect(resp.Header.Get("Access-Control-Allow-Credentials")).To(Equal("true"))
	Expect(resp.Header.Get("Access-Control-Allow-Headers")).To(Equal("Content-Type,Origin,Accept,Authorization,Content-Length,X-Requested-With"))
	responseBody, err := ioutil.ReadAll(resp.Body)
	Expect(string(responseBody)).To(Equal(""))
}

func Test_InterceptNewPreflightRequest_ReturnNilIfRequestIsNotOptions(t *testing.T) {
	RegisterTestingT(t)

	unit := DefaultCORSConfigs()

	r, err := http.NewRequest(http.MethodGet, "http://somehost.com", nil)
	Expect(err).To(BeNil())
	r.Header.Set("Origin", "http://originhost.com")
	r.Header.Set("Access-Control-Request-Methods", "PUT,POST")
	resp := unit.InterceptNewPreflightRequest(r)

	Expect(resp).To(BeNil())
}

func Test_InterceptNewPreflightRequest_ReturnNilIfRequestHasNoOriginHeader(t *testing.T) {
	RegisterTestingT(t)

	unit := DefaultCORSConfigs()

	r, err := http.NewRequest(http.MethodOptions, "http://somehost.com", nil)
	Expect(err).To(BeNil())
	r.Header.Set("Access-Control-Request-Methods", "PUT,POST")
	resp := unit.InterceptNewPreflightRequest(r)

	Expect(resp).To(BeNil())
}

func Test_InterceptNewPreflightRequest_ReturnNilIfRequestHasNoRequestMethodsHeader(t *testing.T) {
	RegisterTestingT(t)

	unit := DefaultCORSConfigs()

	r, err := http.NewRequest(http.MethodOptions, "http://somehost.com", nil)
	Expect(err).To(BeNil())
	r.Header.Set("Origin", "http://originhost.com")
	resp := unit.InterceptNewPreflightRequest(r)

	Expect(resp).To(BeNil())
}

func Test_InterceptNewPreflightRequest_ShouldEchoAllowHeaders(t *testing.T) {
	RegisterTestingT(t)

	unit := DefaultCORSConfigs()

	r, err := http.NewRequest(http.MethodOptions, "http://somehost.com", nil)
	Expect(err).To(BeNil())
	r.Header.Set("Origin", "http://originhost.com")
	r.Header.Set("Access-Control-Request-Methods", "PUT,POST")
	r.Header.Set("Access-Control-Request-Headers", "X-PINGOTHER,Content-Type")
	resp := unit.InterceptNewPreflightRequest(r)

	Expect(resp).ToNot(BeNil())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))
	Expect(resp.Header.Get("Access-Control-Allow-Headers")).To(Equal("X-PINGOTHER,Content-Type"))
}

func Test_InterceptNewPreflightRequest_ShouldNotSetAllowCredentialsHeaderIfFalse(t *testing.T) {
	RegisterTestingT(t)

	unit := DefaultCORSConfigs()
	unit.AllowCredentials = false

	r, err := http.NewRequest(http.MethodOptions, "http://somehost.com", nil)
	Expect(err).To(BeNil())
	r.Header.Set("Origin", "http://originhost.com")
	r.Header.Set("Access-Control-Request-Methods", "PUT,POST")
	resp := unit.InterceptNewPreflightRequest(r)

	Expect(resp).ToNot(BeNil())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))
	Expect(resp.Header.Get("Access-Control-Allow-Credentials")).To(Equal(""))
}


func Test_AddCORSHeaders_ShouldAddDefaultCORSHeadersToResponse(t *testing.T) {
	RegisterTestingT(t)

	unit := DefaultCORSConfigs()

	r, err := http.NewRequest(http.MethodGet, "http://somehost.com", nil)
	Expect(err).To(BeNil())
	r.Header.Set("Origin", "http://originhost.com")

	resp := &http.Response{}
	unit.AddCORSHeaders(r, resp)

	Expect(resp).ToNot(BeNil())
	Expect(resp.Header.Get("Access-Control-Allow-Origin")).To(Equal("*"))
	Expect(resp.Header.Get("Access-Control-Allow-Credentials")).To(Equal("true"))
	Expect(resp.Header.Get("Access-Control-Expose-Headers")).To(Equal(""))
}

func Test_AddCORSHeaders_ShouldNotAddCORSHeadersIfRequestHasNoOriginHeader(t *testing.T) {
	RegisterTestingT(t)

	unit := DefaultCORSConfigs()

	r, err := http.NewRequest(http.MethodGet, "http://somehost.com", nil)
	Expect(err).To(BeNil())

	resp := &http.Response{}
	unit.AddCORSHeaders(r, resp)

	Expect(resp).ToNot(BeNil())
	Expect(resp.Header.Get("Access-Control-Allow-Origin")).To(Equal(""))
	Expect(resp.Header.Get("Access-Control-Allow-Credentials")).To(Equal(""))
	Expect(resp.Header.Get("Access-Control-Expose-Headers")).To(Equal(""))
}

func Test_AddCORSHeaders_ShouldNotSetAllowCredentialsHeaderIfFalse(t *testing.T) {
	RegisterTestingT(t)

	unit := DefaultCORSConfigs()
	unit.AllowCredentials = false

	r, err := http.NewRequest(http.MethodGet, "http://somehost.com", nil)
	Expect(err).To(BeNil())
	r.Header.Set("Origin", "http://originhost.com")

	resp := &http.Response{}
	unit.AddCORSHeaders(r, resp)

	Expect(resp).ToNot(BeNil())
	Expect(resp.Header.Get("Access-Control-Allow-Origin")).To(Equal("*"))
	Expect(resp.Header.Get("Access-Control-Allow-Credentials")).To(Equal(""))
	Expect(resp.Header.Get("Access-Control-Expose-Headers")).To(Equal(""))
}