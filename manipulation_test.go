package hoverfly

import (
	"github.com/SpectoLabs/hoverfly/core/testutil"
	"io/ioutil"
	"net/http"
	"testing"
	"github.com/SpectoLabs/hoverfly/core/models"
)

func TestReconstructRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	// changing payload so we don't have to call middleware
	request := models.RequestDetails{
		Path:        "/random-path",
		Method:      "POST",
		Query:       "?foo=bar",
		Destination: "changed.destination.com",
	}
	payload := models.Payload{Request: request}

	c := NewConstructor(req, payload)
	newRequest, err := c.ReconstructRequest()
	testutil.Expect(t, err, nil)
	testutil.Expect(t, newRequest.Method, "POST")
	testutil.Expect(t, newRequest.URL.Path, "/random-path")
	testutil.Expect(t, newRequest.Host, "changed.destination.com")
	testutil.Expect(t, newRequest.URL.RawQuery, "?foo=bar")
}

func TestReconstructRequestBodyPayload(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	payload := models.Payload{}
	c := NewConstructor(req, payload)
	c.payload.Request.Method = "OPTIONS"
	c.payload.Request.Destination = "newdestination"
	c.payload.Request.Body = "new request body here"

	newRequest, err := c.ReconstructRequest()

	testutil.Expect(t, err, nil)
	testutil.Expect(t, newRequest.Method, "OPTIONS")
	testutil.Expect(t, newRequest.Host, "newdestination")

	body, err := ioutil.ReadAll(newRequest.Body)

	testutil.Expect(t, err, nil)
	testutil.Expect(t, string(body), "new request body here")
}

func TestReconstructRequestHeadersPayload(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	req.Header.Set("Header", "ValueX")

	payload := models.Payload{}
	c := NewConstructor(req, payload)
	c.payload.Request.Headers = req.Header
	c.payload.Request.Destination = "destination.com"

	newRequest, err := c.ReconstructRequest()
	testutil.Expect(t, err, nil)
	testutil.Expect(t, newRequest.Header.Get("Header"), "ValueX")
}

func TestReconstructResponseHeadersPayload(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	payload := models.Payload{}

	payload.Response.Status = 201
	payload.Response.Body = "body here"

	headers := make(map[string][]string)
	headers["Header"] = []string{"one"}

	payload.Response.Headers = headers

	c := NewConstructor(req, payload)

	response := c.ReconstructResponse()

	testutil.Expect(t, response.Header.Get("Header"), headers["Header"][0])

}

func TestReconstructionFailure(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	payload := models.Payload{}
	c := NewConstructor(req, payload)
	c.payload.Request.Method = "GET"
	c.payload.Request.Body = "new request body here"

	_, err := c.ReconstructRequest()
	testutil.Refute(t, err, nil)
}
