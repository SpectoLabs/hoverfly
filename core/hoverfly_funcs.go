package hoverfly

import (
	"github.com/aymerick/raymond"
	"net/http"
	"strings"

	"github.com/SpectoLabs/hoverfly/core/errors"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/modes"
	"github.com/SpectoLabs/hoverfly/core/util"
	log "github.com/sirupsen/logrus"
)

// DoRequest - performs request and returns response that should be returned to client and error
func (hf *Hoverfly) DoRequest(request *http.Request) (*http.Response, error) {

	// We can't have this set. And it only contains "/pkg/net/http/" anyway
	request.RequestURI = ""

	client, err := GetHttpClient(hf, request.Host)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	resp.Header.Set("hoverfly", "Was-Here")

	return resp, nil

}

// GetResponse returns stored response from cache
func (hf *Hoverfly) GetResponse(requestDetails models.RequestDetails) (*models.ResponseDetails, *errors.HoverflyError) {

	var response models.ResponseDetails
	var cachedResponse *models.CachedResponse

	cachedResponse, cacheErr := hf.CacheMatcher.GetCachedResponse(&requestDetails)

	// Get the cached response and return if there is a miss
	if cacheErr == nil && cachedResponse.MatchingPair == nil {
		return nil, errors.MatchingFailedError(cachedResponse.ClosestMiss)
		// If it's cached, use that response
	} else if cacheErr == nil {
		response = cachedResponse.MatchingPair.Response
		//If it's not cached, perform matching to find a hit
	} else {
		mode := (hf.modeMap[modes.Simulate]).(*modes.SimulateMode)

		// Matching
		result := matching.Match(mode.MatchingStrategy, requestDetails, hf.Cfg.Webserver, hf.Simulation, hf.state)

		// Cache result
		if result.Cachable {
			cachedResponse, _ = hf.CacheMatcher.SaveRequestMatcherResponsePair(requestDetails, result.Pair, result.Error)
		}

		// If we miss, just return
		if result.Error != nil {
			log.WithFields(log.Fields{
				"error":       result.Error.Error(),
				"query":       requestDetails.Query,
				"path":        requestDetails.Path,
				"destination": requestDetails.Destination,
				"method":      requestDetails.Method,
			}).Warn("Failed to find matching request from simulation")

			return nil, errors.MatchingFailedError(result.Error.ClosestMiss)
		} else {
			response = result.Pair.Response
		}
	}

	// Templating applies at the end, once we have loaded a response. Comes BEFORE state transitions,
	// as we use the current state in templates
	if response.Templated == true {

		var template *raymond.Template
		if cachedResponse != nil && cachedResponse.ResponseTemplate != nil {
			template = cachedResponse.ResponseTemplate
		} else {
			// Parse and cache the template
			template, _ = hf.templator.ParseTemplate(response.Body)
			if cachedResponse != nil {
				cachedResponse.ResponseTemplate = template
			}
		}

		responseBody, err :=  hf.templator.RenderTemplate(template, &requestDetails, hf.state.State)

		if err == nil {
			response.Body = responseBody
		} else {
			log.Warnf("Failed to render response template: %s", err.Error())
		}
	}

	// State transitions after we have the response
	if response.TransitionsState != nil {
		hf.state.PatchState(response.TransitionsState)
	}
	if response.RemovesState != nil {
		hf.state.RemoveState(response.RemovesState)
	}

	return &response, nil
}

// save gets request fingerprint, extracts request body, status code and headers, then saves it to cache
func (hf *Hoverfly) Save(request *models.RequestDetails, response *models.ResponseDetails, headersWhitelist []string, recordSequence bool) error {
	body := []models.RequestFieldMatchers{
		{
			Matcher: matchers.Exact,
			Value:   request.Body,
		},
	}
	contentType := util.GetContentTypeFromHeaders(request.Headers)
	if contentType == "json" {
		body = []models.RequestFieldMatchers{
			{
				Matcher: matchers.Json,
				Value:   request.Body,
			},
		}
	} else if contentType == "xml" {
		body = []models.RequestFieldMatchers{
			{
				Matcher: matchers.Xml,
				Value:   request.Body,
			},
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

	requestHeaders := map[string][]models.RequestFieldMatchers{}
	for headerKey, headerValues := range headers {
		requestHeaders[headerKey] = []models.RequestFieldMatchers{
			{
				Matcher: matchers.Exact,
				Value:   strings.Join(headerValues, ";"),
			},
		}
	}

	queries := &models.QueryRequestFieldMatchers{}
	for key, values := range request.Query {
		queries.Add(key, []models.RequestFieldMatchers{
			{
				Matcher: matchers.Exact,
				Value:   strings.Join(values, ";"),
			},
		})
	}

	pair := models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   request.Path,
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   request.Method,
				},
			},
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   request.Destination,
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   request.Scheme,
				},
			},
			Query:   queries,
			Body:    body,
			Headers: requestHeaders,
		},
		Response: *response,
	}
	if recordSequence {
		hf.Simulation.AddPairInSequence(&pair, hf.state)
	} else {
		hf.Simulation.AddPair(&pair)
	}

	return nil
}

func (this Hoverfly) ApplyMiddleware(pair models.RequestResponsePair) (models.RequestResponsePair, error) {
	if this.Cfg.Middleware.IsSet() {
		return this.Cfg.Middleware.Execute(pair)
	}

	return pair, nil
}
