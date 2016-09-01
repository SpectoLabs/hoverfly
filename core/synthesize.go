package hoverfly

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/models"
)

// SynthesizeResponse calls middleware to populate response data, nothing gets pass proxy
func SynthesizeResponse(req *http.Request, requestDetails models.RequestDetails, middleware string) (*http.Response, error) {
	pair := models.RequestResponsePair{Request: requestDetails}

	log.WithFields(log.Fields{
		"middleware":  middleware,
		"body":        requestDetails.Body,
		"destination": requestDetails.Destination,
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