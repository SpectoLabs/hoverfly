package hoverfly

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/modes"
	"github.com/SpectoLabs/hoverfly/core/templating"
	"github.com/SpectoLabs/hoverfly/core/util"
)

// DoRequest - performs request and returns response that should be returned to client and error
func (hf *Hoverfly) DoRequest(request *http.Request) (*http.Response, error) {

	// We can't have this set. And it only contains "/pkg/net/http/" anyway
	request.RequestURI = ""

	requestBody, _ := ioutil.ReadAll(request.Body)

	request.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	resp, err := hf.HTTP.Do(request)

	request.Body = ioutil.NopCloser(bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}

	resp.Header.Set("hoverfly", "Was-Here")

	return resp, nil

}

// GetResponse returns stored response from cache
func (hf *Hoverfly) GetResponse(requestDetails models.RequestDetails) (*models.ResponseDetails, *matching.MatchingError) {

	cachedResponse, cacheErr := hf.CacheMatcher.GetCachedResponse(&requestDetails)
	if cacheErr == nil && cachedResponse.MatchingPair == nil {
		return nil, matching.MissedError(cachedResponse.ClosestMiss)
	} else if cacheErr == nil {
		if cachedResponse.MatchingPair.Response.TransitionsState != nil {
			hf.TransitionState(cachedResponse.MatchingPair.Response.TransitionsState)
		}
		if cachedResponse.MatchingPair.Response.RemovesState != nil {
			hf.RemoveState(cachedResponse.MatchingPair.Response.RemovesState)
		}
		return &cachedResponse.MatchingPair.Response, nil
	}

	var pair *models.RequestMatcherResponsePair
	var err *models.MatchError
	var cachable bool

	mode := (hf.modeMap[modes.Simulate]).(*modes.SimulateMode)

	strongestMatch := strings.ToLower(mode.MatchingStrategy) == "strongest"

	// Matching
	if strongestMatch {
		pair, err, cachable = matching.StrongestMatchRequestMatcher(requestDetails, hf.Cfg.Webserver, hf.Simulation, hf.state)
	} else {
		pair, err, cachable = matching.FirstMatchRequestMatcher(requestDetails, hf.Cfg.Webserver, hf.Simulation, hf.state)
	}

	if err == nil {
		// Templating
		if pair.Response.Templated == true {
			responseBody, err := templating.ApplyTemplate(&requestDetails, hf.state, pair.Response.Body)
			if err == nil {
				pair.Response.Body = responseBody
			}
		}
		// State transitions
		if pair.Response.TransitionsState != nil {
			hf.TransitionState(pair.Response.TransitionsState)
		}
		if pair.Response.RemovesState != nil {
			hf.RemoveState(pair.Response.RemovesState)
		}
	}

	// Caching
	if cachable {
		hf.CacheMatcher.SaveRequestMatcherResponsePair(requestDetails, pair, err)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":       err.Error(),
			"query":       requestDetails.Query,
			"path":        requestDetails.Path,
			"destination": requestDetails.Destination,
			"method":      requestDetails.Method,
		}).Warn("Failed to find matching request from simulation")

		return nil, matching.MissedError(err.ClosestMiss)
	}

	return &pair.Response, nil
}

func (hf *Hoverfly) TransitionState(transition map[string]string) {
	for k, v := range transition {
		hf.state[k] = v
	}
}

func (hf *Hoverfly) RemoveState(toRemove []string) {
	for _, key := range toRemove {
		delete(hf.state, key)
	}
}

// save gets request fingerprint, extracts request body, status code and headers, then saves it to cache
func (hf *Hoverfly) Save(request *models.RequestDetails, response *models.ResponseDetails, headersWhitelist []string) error {
	body := &models.RequestFieldMatchers{
		ExactMatch: util.StringToPointer(request.Body),
	}
	contentType := util.GetContentTypeFromHeaders(request.Headers)
	if contentType == "json" {
		body = &models.RequestFieldMatchers{
			JsonMatch: util.StringToPointer(request.Body),
		}
	} else if contentType == "xml" {
		body = &models.RequestFieldMatchers{
			XmlMatch: util.StringToPointer(request.Body),
		}
	}

	var headers map[string][]string
	if headersWhitelist == nil {
		headersWhitelist = []string{}
	}

	if len(headersWhitelist) >= 1 && headersWhitelist[0] == "*" {
		headers = request.Headers
	} else {
		headers = map[string][]string{}
		for _, header := range headersWhitelist {
			headerValues := request.Headers[header]
			if len(headerValues) > 0 {
				headers[header] = headerValues
			}
		}
	}

	pair := models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Path: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer(request.Path),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer(request.Method),
			},
			Destination: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer(request.Destination),
			},
			Scheme: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer(request.Scheme),
			},
			Query: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer(request.QueryString()),
			},
			Body:    body,
			Headers: headers,
		},
		Response: *response,
	}

	hf.Simulation.AddRequestMatcherResponsePair(&pair)

	return nil
}

func (this Hoverfly) ApplyMiddleware(pair models.RequestResponsePair) (models.RequestResponsePair, error) {
	if this.Cfg.Middleware.IsSet() {
		return this.Cfg.Middleware.Execute(pair)
	}

	return pair, nil
}
