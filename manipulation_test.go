package hoverfly

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestReconstructRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	// changing payload so we don't have to call middleware
	request := RequestDetails{
		Path:        "/random-path",
		Method:      "POST",
		Query:       "?foo=bar",
		Destination: "changed.destination.com",
	}
	payload := Payload{Request: request}

	c := NewConstructor(req, payload)
	newRequest, err := c.reconstructRequest()
	expect(t, err, nil)
	expect(t, newRequest.Method, "POST")
	expect(t, newRequest.URL.Path, "/random-path")
	expect(t, newRequest.Host, "changed.destination.com")
	expect(t, newRequest.URL.RawQuery, "?foo=bar")
}

func TestReconstructRequestBodyPayload(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	payload := Payload{}
	c := NewConstructor(req, payload)
	c.payload.Request.Method = "OPTIONS"
	c.payload.Request.Destination = "newdestination"
	c.payload.Request.Body = "new request body here"

	newRequest, err := c.reconstructRequest()

	expect(t, err, nil)
	expect(t, newRequest.Method, "OPTIONS")
	expect(t, newRequest.Host, "newdestination")

	body, err := ioutil.ReadAll(newRequest.Body)

	expect(t, err, nil)
	expect(t, string(body), "new request body here")
}

func TestReconstructRequestHeadersPayload(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	req.Header.Set("Header", "ValueX")

	payload := Payload{}
	c := NewConstructor(req, payload)
	c.payload.Request.Headers = req.Header
	c.payload.Request.Destination = "destination.com"

	newRequest, err := c.reconstructRequest()
	expect(t, err, nil)
	expect(t, newRequest.Header.Get("Header"), "ValueX")
}

func TestReconstructResponseHeadersPayload(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	payload := Payload{}

	payload.Response.Status = 201
	payload.Response.Body = "body here"

	headers := make(map[string][]string)
	headers["Header"] = []string{"one"}

	payload.Response.Headers = headers

	c := NewConstructor(req, payload)

	response := c.reconstructResponse()

	expect(t, response.Header.Get("Header"), headers["Header"][0])

}

func TestReconstructionFailure(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	payload := Payload{}
	c := NewConstructor(req, payload)
	c.payload.Request.Method = "GET"
	c.payload.Request.Body = "new request body here"

	_, err := c.reconstructRequest()
	refute(t, err, nil)
}
