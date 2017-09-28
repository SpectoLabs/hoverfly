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

	var response models.ResponseDetails

	cachedResponse, cacheErr := hf.CacheMatcher.GetCachedResponse(&requestDetails)

	// Get the cached response and return if there is a miss
	if cacheErr == nil && cachedResponse.MatchingPair == nil {
		return nil, matching.MissedError(cachedResponse.ClosestMiss)
		// If it's cached, use that response
	} else if cacheErr == nil {
		response = cachedResponse.MatchingPair.Response
		//If it's not cached, perform matching to find a hit
	} else {
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

		// Cache result
		if cachable {
			hf.CacheMatcher.SaveRequestMatcherResponsePair(requestDetails, pair, err)
		}

		// If we miss, just return
		if err != nil {
			log.WithFields(log.Fields{
				"error":       err.Error(),
				"query":       requestDetails.Query,
				"path":        requestDetails.Path,
				"destination": requestDetails.Destination,
				"method":      requestDetails.Method,
			}).Warn("Failed to find matching request from simulation")

			return nil, matching.MissedError(err.ClosestMiss)
		} else {
			response = pair.Response
		}
	}

	// Templating applies at the end, once we have loaded a response. Comes BEFORE state transitions,
	// as we use the current state in templates
	if response.Templated == true {
		responseBody, err := hf.templator.ApplyTemplate(&requestDetails, hf.state, response.Body)
		if err == nil {
			response.Body = responseBody
		}
	}

	// State transitions after we have the response
	if response.TransitionsState != nil {
		hf.TransitionState(response.TransitionsState)
	}
	if response.RemovesState != nil {
		hf.RemoveState(response.RemovesState)
	}

	return &response, nil
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
