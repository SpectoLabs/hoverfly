package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

// synthesizeResponse calls middleware to populate response data, nothing gets pass proxy
func synthesizeResponse(req *http.Request, middleware string) *http.Response {

	// this is mainly for testing, since when you create
	if req.Body == nil {
		req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("")))
	}
	responseBody, err := ioutil.ReadAll(req.Body)
	req.Body.Close()

	var bodyStr string
	if err != nil {
		log.WithFields(log.Fields{
			"middleware": middleware,
			"error":      err.Error(),
		}).Error("Failed to read request body when synthesizing response")
	} else {
		bodyStr = string(responseBody)
	}

	request := requestDetails{
		Path:        req.URL.Path,
		Method:      req.Method,
		Destination: req.Host,
		Query:       req.URL.RawQuery,
		Body:        bodyStr,
		RemoteAddr:  req.RemoteAddr,
		Headers:     req.Header,
	}
	payload := Payload{Request: request}

	log.WithFields(log.Fields{
		"middleware":  middleware,
		"body":        bodyStr,
		"destination": request.Destination,
	}).Debug("Synthesizing new response")

	c := NewConstructor(req, payload)

	if middleware != "" {
		err := c.ApplyMiddleware(middleware)
		if err != nil {
			var errorPayload Payload
			errorPayload.Response.Status = 503
			errorPayload.Response.Body = fmt.Sprintf("Middleware error: %s", err.Error())
			c.payload = errorPayload
		}
	} else {
		c.payload.Response.Body = "Precondition failed: middleware not provided."
		c.payload.Response.Status = 428
	}

	response := c.reconstructResponse()
	return response

}
