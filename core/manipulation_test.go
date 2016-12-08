package hoverfly

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
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
	pair := models.RequestResponsePair{Request: request}

	c := NewConstructor(req, pair)
	newRequest, err := c.ReconstructRequest()
	Expect(err).To(BeNil())
	Expect(newRequest.Method).To(Equal("POST"))
	Expect(newRequest.URL.Path).To(Equal("/random-path"))
	Expect(newRequest.Host).To(Equal("changed.destination.com"))
	Expect(newRequest.URL.RawQuery).To(Equal("?foo=bar"))
}

func TestReconstructRequestBodyRequestResponsePair(t *testing.T) {
	RegisterTestingT(t)

	req, _ := http.NewRequest("GET", "http://example.com", nil)

	emptyPair := models.RequestResponsePair{}
	c := NewConstructor(req, emptyPair)
	c.requestResponsePair.Request.Method = "OPTIONS"
	c.requestResponsePair.Request.Destination = "newdestination"
	c.requestResponsePair.Request.Body = "new request body here"

	newRequest, err := c.ReconstructRequest()

	Expect(err).To(BeNil())
	Expect(newRequest.Method).To(Equal("OPTIONS"))
	Expect(newRequest.Host).To(Equal("newdestination"))

	body, err := ioutil.ReadAll(newRequest.Body)

	Expect(err).To(BeNil())
	Expect(string(body)).To(Equal("new request body here"))
}

func TestReconstructRequestHeadersInPair(t *testing.T) {
	RegisterTestingT(t)

	req, _ := http.NewRequest("GET", "http://example.com", nil)

	req.Header.Set("Header", "ValueX")

	emptyPair := models.RequestResponsePair{}
	c := NewConstructor(req, emptyPair)
	c.requestResponsePair.Request.Headers = req.Header
	c.requestResponsePair.Request.Destination = "destination.com"

	newRequest, err := c.ReconstructRequest()
	Expect(err).To(BeNil())
	Expect(newRequest.Header.Get("Header")).To(Equal("ValueX"))
}

func TestReconstructResponseHeadersInPair(t *testing.T) {
	RegisterTestingT(t)

	req, _ := http.NewRequest("GET", "http://example.com", nil)

	pair := models.RequestResponsePair{}

	pair.Response.Status = 201
	pair.Response.Body = "body here"

	headers := make(map[string][]string)
	headers["Header"] = []string{"one"}

	pair.Response.Headers = headers

	c := NewConstructor(req, pair)

	response := c.ReconstructResponse()

	Expect(response.Header.Get("Header")).To(Equal(headers["Header"][0]))

}

func TestReconstructionFailure(t *testing.T) {
	RegisterTestingT(t)

	req, _ := http.NewRequest("GET", "http://example.com", nil)

	emptyPair := models.RequestResponsePair{}
	c := NewConstructor(req, emptyPair)
	c.requestResponsePair.Request.Method = "GET"
	c.requestResponsePair.Request.Body = "new request body here"

	_, err := c.ReconstructRequest()
	Expect(err).ToNot(BeNil())
}
