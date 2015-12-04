package main

import (
	"bytes"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

type Constructor struct {
	request *http.Request
	payload Payload
}

func NewConstructor(req *http.Request, payload Payload) Constructor {
	c := Constructor{request: req, payload: payload}
	return c
}

func (c *Constructor) ApplyMiddleware(middleware string) error {

	newPayload, err := ExecuteMiddleware(middleware, c.payload)

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err.Error(),
			"middleware": AppConfig.middleware,
		}).Error("Error during middleware transformation, not modifying payload!")

		return err
	} else {

		log.WithFields(log.Fields{
			"middleware": AppConfig.middleware,
			"newPayload": newPayload,
		}).Info("Middleware transformation complete!")
		// override payload with transformed new payload
		c.payload = newPayload

		return nil
	}
}

func (c *Constructor) reconstructResponse() *http.Response {
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
