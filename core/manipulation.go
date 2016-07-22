package hoverfly

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/models"
	"strings"
)

// Constructor - holds information about original request (which is needed to create response
// and also holds payload
type Constructor struct {
	request *http.Request
	payload models.Payload
}

// NewConstructor - returns constructor instance
func NewConstructor(req *http.Request, payload models.Payload) *Constructor {
	c := &Constructor{request: req, payload: payload}
	return c
}

// ApplyMiddleware - activates given middleware, middleware should be passed as string to executable, can be
// full path.
func (c *Constructor) ApplyMiddleware(middleware string) error {
	var newPayload models.Payload
	var err error

	if isMiddlewareLocal(middleware) {
		newPayload, err = ExecuteMiddlewareLocally(middleware, c.payload)
	} else {
		newPayload, err = ExecuteMiddlewareRemotely(middleware, c.payload)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err.Error(),
			"middleware": middleware,
		}).Error("Error during middleware transformation, not modifying payload!")

		return err
	}

	log.WithFields(log.Fields{
		"middleware": middleware,
	}).Debug("Middleware transformation complete!")
	// override payload with transformed new payload
	c.payload = newPayload

	return nil

}

func isMiddlewareLocal(middleware string) (bool) {
	return !strings.HasPrefix(middleware, "http")
}

// ReconstructResponse changes original response with details provided in Constructor Payload.Response
func (c *Constructor) ReconstructResponse() *http.Response {
	response := &http.Response{}
	response.Request = c.request

	// adding headers
	response.Header = make(http.Header)

	// applying payload
	if len(c.payload.Response.Headers) > 0 {
		for k, values := range c.payload.Response.Headers {
			// headers is a map, appending each value
			for _, v := range values {
				response.Header.Add(k, v)
			}

		}
	}
	// adding body, length, status code
	buf := bytes.NewBufferString(c.payload.Response.Body)
	response.ContentLength = int64(buf.Len())
	response.Body = ioutil.NopCloser(buf)
	response.StatusCode = c.payload.Response.Status

	return response
}

// ReconstructRequest replaces original request with details provided in Constructor Payload.Request
func (c *Constructor) ReconstructRequest() (*http.Request, error) {
	// let's default to what was given
	if c.payload.Request.Scheme == "" {
		c.payload.Request.Scheme = c.request.URL.Scheme
	}

	if c.payload.Request.Destination == "" {
		return nil, fmt.Errorf("failed to reconstruct request, destination not specified")
	}

	newRequest, err := http.NewRequest(
		c.payload.Request.Method,
		fmt.Sprintf("%s://%s", c.payload.Request.Scheme, c.payload.Request.Destination),
		ioutil.NopCloser(bytes.NewBuffer([]byte(c.payload.Request.Body))))

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Request reconstruction failed...")
		return nil, err
	}

	newRequest.Method = c.payload.Request.Method
	newRequest.URL.Path = c.payload.Request.Path
	newRequest.URL.RawQuery = c.payload.Request.Query
	newRequest.Header = c.payload.Request.Headers

	// overriding original request
	c.request = newRequest

	return newRequest, nil
}
