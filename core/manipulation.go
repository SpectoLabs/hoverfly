package hoverfly

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/models"
)

// Constructor - holds information about original request (which is needed to create response
// and also holds payload
type Constructor struct {
	request             *http.Request
	requestResponsePair models.RequestResponsePair
}

// NewConstructor - returns constructor instance
func NewConstructor(req *http.Request, pair models.RequestResponsePair) *Constructor {
	c := &Constructor{request: req, requestResponsePair: pair}
	return c
}

// ApplyMiddleware - activates given middleware, middleware should be passed as string to executable, can be
// full path.
func (c *Constructor) ApplyMiddleware(middleware string) error {
	var newPair models.RequestResponsePair
	var err error

	if isMiddlewareLocal(middleware) {
		newPair, err = ExecuteMiddlewareLocally(middleware, c.requestResponsePair)
	} else {
		middlewareObject := &Middleware{
			Script: middleware,
		}
		newPair, err = middlewareObject.ExecuteMiddlewareRemotely(c.requestResponsePair)
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
	c.requestResponsePair = newPair

	return nil

}

func isMiddlewareLocal(middleware string) bool {
	return !strings.HasPrefix(middleware, "http")
}

// ReconstructResponse changes original response with details provided in Constructor Payload.Response
func (c *Constructor) ReconstructResponse() *http.Response {
	response := &http.Response{}
	response.Request = c.request

	// adding headers
	response.Header = make(http.Header)

	// applying payload
	if len(c.requestResponsePair.Response.Headers) > 0 {
		for k, values := range c.requestResponsePair.Response.Headers {
			// headers is a map, appending each value
			for _, v := range values {
				response.Header.Add(k, v)
			}

		}
	}
	// adding body, length, status code
	buf := bytes.NewBufferString(c.requestResponsePair.Response.Body)
	response.ContentLength = int64(buf.Len())
	response.Body = ioutil.NopCloser(buf)
	response.StatusCode = c.requestResponsePair.Response.Status

	return response
}

// ReconstructRequest replaces original request with details provided in Constructor Payload.Request
func (c *Constructor) ReconstructRequest() (*http.Request, error) {
	// let's default to what was given
	if c.requestResponsePair.Request.Scheme == "" {
		c.requestResponsePair.Request.Scheme = c.request.URL.Scheme
	}

	if c.requestResponsePair.Request.Destination == "" {
		return nil, fmt.Errorf("failed to reconstruct request, destination not specified")
	}

	newRequest, err := http.NewRequest(
		c.requestResponsePair.Request.Method,
		fmt.Sprintf("%s://%s", c.requestResponsePair.Request.Scheme, c.requestResponsePair.Request.Destination),
		ioutil.NopCloser(bytes.NewBuffer([]byte(c.requestResponsePair.Request.Body))))

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Request reconstruction failed...")
		return nil, err
	}

	newRequest.Method = c.requestResponsePair.Request.Method
	newRequest.URL.Path = c.requestResponsePair.Request.Path
	newRequest.URL.RawQuery = c.requestResponsePair.Request.Query
	newRequest.Header = c.requestResponsePair.Request.Headers

	// overriding original request
	c.request = newRequest

	return newRequest, nil
}
