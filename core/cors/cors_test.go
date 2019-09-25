package cors

import (
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"testing"
)

func Test_InterceptPreflightRequest_ReturnSuccessResponseWithDefaultHeaders(t *testing.T) {
	RegisterTestingT(t)

	unit := DefaultCORSConfigs()

	r, err := http.NewRequest(http.MethodOptions, "http://somehost.com", nil)
	Expect(err).To(BeNil())
	r.Header.Set("Origin", "http://originhost.com")
	r.Header.Set("Access-Control-Request-Method", "PUT")
	resp := unit.InterceptPreflightRequest(r)

	Expect(resp).ToNot(BeNil())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))
	Expect(resp.Header.Get("Access-Control-Allow-Origin")).To(Equal("http://originhost.com"))
	Expect(resp.Header.Get("Access-Control-Allow-Methods")).To(Equal("GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS"))
	Expect(resp.Header.Get("Access-Control-Max-Age")).To(Equal("1800"))
	Expect(resp.Header.Get("Access-Control-Allow-Credentials")).To(Equal("true"))
	Expect(resp.Header.Get("Access-Control-Allow-Headers")).To(Equal("Content-Type,Origin,Accept,Authorization,Content-Length,X-Requested-With"))
	responseBody, err := ioutil.ReadAll(resp.Body)
	Expect(string(responseBody)).To(Equal(""))
}

func Test_InterceptPreflightRequest_ReturnNilIfRequestIsNotOptions(t *testing.T) {
	RegisterTestingT(t)

	unit := DefaultCORSConfigs()

	r, err := http.NewRequest(http.MethodGet, "http://somehost.com", nil)
	Expect(err).To(BeNil())
	r.Header.Set("Origin", "http://originhost.com")
	r.Header.Set("Access-Control-Request-Method", "PUT")
	resp := unit.InterceptPreflightRequest(r)

	Expect(resp).To(BeNil())
}

func Test_InterceptPreflightRequest_ReturnNilIfRequestHasNoOriginHeader(t *testing.T) {
	RegisterTestingT(t)

	unit := DefaultCORSConfigs()

	r, err := http.NewRequest(http.MethodOptions, "http://somehost.com", nil)
	Expect(err).To(BeNil())
	r.Header.Set("Access-Control-Request-Method", "PUT")
	resp := unit.InterceptPreflightRequest(r)

	Expect(resp).To(BeNil())
}

func Test_InterceptPreflightRequest_ReturnNilIfRequestHasNoRequestMethodsHeader(t *testing.T) {
	RegisterTestingT(t)

	unit := DefaultCORSConfigs()

	r, err := http.NewRequest(http.MethodOptions, "http://somehost.com", nil)
	Expect(err).To(BeNil())
	r.Header.Set("Origin", "http://originhost.com")
	resp := unit.InterceptPreflightRequest(r)

	Expect(resp).To(BeNil())
}

func Test_InterceptPreflightRequest_ShouldEchoAllowHeaders(t *testing.T) {
	RegisterTestingT(t)

	unit := DefaultCORSConfigs()

	r, err := http.NewRequest(http.MethodOptions, "http://somehost.com", nil)
	Expect(err).To(BeNil())
	r.Header.Set("Origin", "http://originhost.com")
	r.Header.Set("Access-Control-Request-Method", "PUT")
	r.Header.Set("Access-Control-Request-Headers", "X-PINGOTHER,Content-Type")
	resp := unit.InterceptPreflightRequest(r)

	Expect(resp).ToNot(BeNil())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))
	Expect(resp.Header.Get("Access-Control-Allow-Headers")).To(Equal("X-PINGOTHER,Content-Type"))
}

func Test_InterceptPreflightRequest_ShouldNotSetAllowCredentialsHeaderIfFalse(t *testing.T) {
	RegisterTestingT(t)

	unit := DefaultCORSConfigs()
	unit.AllowCredentials = false

	r, err := http.NewRequest(http.MethodOptions, "http://somehost.com", nil)
	Expect(err).To(BeNil())
	r.Header.Set("Origin", "http://originhost.com")
	r.Header.Set("Access-Control-Request-Method", "PUT")
	resp := unit.InterceptPreflightRequest(r)

	Expect(resp).ToNot(BeNil())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))
	Expect(resp.Header.Get("Access-Control-Allow-Credentials")).To(Equal(""))
}

func Test_InterceptPreflightRequest_UseDefaultAllowOriginValueIfAllowCredentialsHeaderIsFalse(t *testing.T) {
	RegisterTestingT(t)

	unit := DefaultCORSConfigs()
	unit.AllowCredentials = false

	r, err := http.NewRequest(http.MethodOptions, "http://somehost.com", nil)
	Expect(err).To(BeNil())
	r.Header.Set("Origin", "http://originhost.com")
	r.Header.Set("Access-Control-Request-Method", "PUT")
	resp := unit.InterceptPreflightRequest(r)

	Expect(resp).ToNot(BeNil())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))
	Expect(resp.Header.Get("Access-Control-Allow-Origin")).To(Equal("*"))
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
	Expect(resp.Header.Get("Access-Control-Allow-Origin")).To(Equal("http://originhost.com"))
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

func Test_AddCORSHeaders_ShouldSetExposeHeadersIfPresent(t *testing.T) {
	RegisterTestingT(t)

	unit := DefaultCORSConfigs()
	unit.ExposeHeaders = "Content-Type"

	r, err := http.NewRequest(http.MethodGet, "http://somehost.com", nil)
	Expect(err).To(BeNil())
	r.Header.Set("Origin", "http://originhost.com")

	resp := &http.Response{}
	unit.AddCORSHeaders(r, resp)

	Expect(resp).ToNot(BeNil())
	Expect(resp.Header.Get("Access-Control-Allow-Origin")).To(Equal("http://originhost.com"))
	Expect(resp.Header.Get("Access-Control-Allow-Credentials")).To(Equal("true"))
	Expect(resp.Header.Get("Access-Control-Expose-Headers")).To(Equal("Content-Type"))
}

func Test_AddCORSHeaders_ShouldPreserveExistingCORSHeaders(t *testing.T) {
	RegisterTestingT(t)

	unit := DefaultCORSConfigs()

	r, err := http.NewRequest(http.MethodGet, "http://somehost.com", nil)
	Expect(err).To(BeNil())
	r.Header.Set("Origin", "http://originhost.com")

	resp := &http.Response{}
	resp.Header = make(http.Header)
	resp.Header.Set("Access-Control-Allow-Origin", "*")

	unit.AddCORSHeaders(r, resp)

	Expect(resp).ToNot(BeNil())
	Expect(resp.Header.Get("Access-Control-Allow-Origin")).To(Equal("*"))
	Expect(resp.Header.Get("Access-Control-Allow-Credentials")).To(Equal(""))
	Expect(resp.Header.Get("Access-Control-Expose-Headers")).To(Equal(""))
}
