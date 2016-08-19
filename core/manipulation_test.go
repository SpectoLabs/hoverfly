package hoverfly

import (
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestReconstructRequest(t *testing.T) {
	RegisterTestingT(t)

	req, _ := http.NewRequest("GET", "http://example.com", nil)

	// changing payload so we don't have to call middleware
	request := models.RequestDetails{
		Path:        "/random-path",
		Method:      "POST",
		Query:       "?foo=bar",
		Destination: "changed.destination.com",
	}
	payload := models.RequestResponsePair{Request: request}

	c := NewConstructor(req, payload)
	newRequest, err := c.ReconstructRequest()
	Expect(err).To(BeNil())
	Expect(newRequest.Method).To(Equal("POST"))
	Expect(newRequest.URL.Path).To(Equal("/random-path"))
	Expect(newRequest.Host).To(Equal("changed.destination.com"))
	Expect(newRequest.URL.RawQuery).To(Equal("?foo=bar"))
}

func TestReconstructRequestBodyPayload(t *testing.T) {
	RegisterTestingT(t)

	req, _ := http.NewRequest("GET", "http://example.com", nil)

	payload := models.RequestResponsePair{}
	c := NewConstructor(req, payload)
	c.payload.Request.Method = "OPTIONS"
	c.payload.Request.Destination = "newdestination"
	c.payload.Request.Body = "new request body here"

	newRequest, err := c.ReconstructRequest()

	Expect(err).To(BeNil())
	Expect(newRequest.Method).To(Equal("OPTIONS"))
	Expect(newRequest.Host).To(Equal("newdestination"))

	body, err := ioutil.ReadAll(newRequest.Body)

	Expect(err).To(BeNil())
	Expect(string(body)).To(Equal("new request body here"))
}

func TestReconstructRequestHeadersPayload(t *testing.T) {
	RegisterTestingT(t)

	req, _ := http.NewRequest("GET", "http://example.com", nil)

	req.Header.Set("Header", "ValueX")

	payload := models.RequestResponsePair{}
	c := NewConstructor(req, payload)
	c.payload.Request.Headers = req.Header
	c.payload.Request.Destination = "destination.com"

	newRequest, err := c.ReconstructRequest()
	Expect(err).To(BeNil())
	Expect(newRequest.Header.Get("Header")).To(Equal("ValueX"))
}

func TestReconstructResponseHeadersPayload(t *testing.T) {
	RegisterTestingT(t)

	req, _ := http.NewRequest("GET", "http://example.com", nil)

	payload := models.RequestResponsePair{}

	payload.Response.Status = 201
	payload.Response.Body = "body here"

	headers := make(map[string][]string)
	headers["Header"] = []string{"one"}

	payload.Response.Headers = headers

	c := NewConstructor(req, payload)

	response := c.ReconstructResponse()

	Expect(response.Header.Get("Header")).To(Equal(headers["Header"][0]))

}

func TestReconstructionFailure(t *testing.T) {
	RegisterTestingT(t)

	req, _ := http.NewRequest("GET", "http://example.com", nil)

	payload := models.RequestResponsePair{}
	c := NewConstructor(req, payload)
	c.payload.Request.Method = "GET"
	c.payload.Request.Body = "new request body here"

	_, err := c.ReconstructRequest()
	Expect(err).ToNot(BeNil())
}

func TestIsMiddlewareLocal_WithNonHttpString(t *testing.T) {
	RegisterTestingT(t)

	Expect(isMiddlewareLocal("python middleware.py")).To(BeTrue())
}

func TestIsMiddlewareLocal_WithHttpString(t *testing.T) {
	RegisterTestingT(t)

	Expect(isMiddlewareLocal("http://remotemiddleware.com/process")).To(BeFalse())
}

func TestIsMiddlewareLocal_WithHttpsString(t *testing.T) {
	RegisterTestingT(t)

	Expect(isMiddlewareLocal("http://remotemiddleware.com/process")).To(BeFalse())
}
