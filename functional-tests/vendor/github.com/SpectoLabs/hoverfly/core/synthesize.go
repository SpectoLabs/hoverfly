package hoverfly

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/models"
)

// SynthesizeResponse calls middleware to populate response data, nothing gets pass proxy
func SynthesizeResponse(req *http.Request, middleware string) (*http.Response, error) {

	// this is mainly for testing, since when you create a request during tests
	// its body will be nil, that results in bad things during read
	if req.Body == nil {
		req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("")))
	}
	defer req.Body.Close()
	requestBody, err := ioutil.ReadAll(req.Body)

	var bodyStr string
	if err != nil {
		log.WithFields(log.Fields{
			"middleware": middleware,
			"error":      err.Error(),
		}).Error("Failed to read request body when synthesizing response")

		// creating new error with more info
		return nil, fmt.Errorf("Synthesize failed, could not read request body - %s", err.Error())
	}

	bodyStr = string(requestBody)

	request := models.RequestDetails{
		Path:        req.URL.Path,
		Method:      req.Method,
		Destination: req.Host,
		Query:       req.URL.RawQuery,
		Body:        bodyStr,
		Headers:     req.Header,
	}
	pair := models.RequestResponsePair{Request: request}

	log.WithFields(log.Fields{
		"middleware":  middleware,
		"body":        bodyStr,
		"destination": request.Destination,
	}).Debug("Synthesizing new response")

	c := NewConstructor(req, pair)

	if middleware != "" {
		err := c.ApplyMiddleware(middleware)
		if err != nil {
			return nil, fmt.Errorf("Synthesize failed, middleware error - %s", err.Error())
		}
	} else {
		return nil, fmt.Errorf("Synthesize failed, middleware not provided")

	}

	response := c.ReconstructResponse()
	return response, nil

}
