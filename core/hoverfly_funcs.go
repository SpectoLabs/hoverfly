package hoverfly

import (
	"fmt"
	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/aymerick/raymond"
	"io/ioutil"
	"net/http"
	"path/filepath"
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

	resp.Header.Set("Hoverfly", "Was-Here")

	if hf.Cfg.Mode == "spy" {
		resp.Header.Add("Hoverfly", "Forwarded")
	}

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

		responseBody, err := hf.applyBodyTemplating(&requestDetails, &response, cachedResponse)

		if err == nil {
			response.Body = responseBody
		} else {
			log.Warnf("Failed to applying body templating: %s", err.Error())
		}

		responseHeaders, err := hf.applyHeadersTemplating(&requestDetails, &response, cachedResponse)

		if err == nil {
			response.Headers = responseHeaders
		} else {
			log.Warnf("Failed to applying headers templating: %s", err.Error())
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

func (hf *Hoverfly) readResponseBodyFiles(pairs []v2.RequestMatcherResponsePairViewV5) v2.SimulationImportResult {
	result := v2.SimulationImportResult{}

	for i, pair := range pairs {
		if len(pair.Response.GetBody()) > 0 && len(pair.Response.GetBodyFile()) > 0 {
			result.AddBodyAndBodyFileWarning(i)
			continue
		}

		if len(pair.Response.GetBody()) == 0 && len(pair.Response.GetBodyFile()) > 0 {
			var content string
			var err error

			bodyFile := pair.Response.GetBodyFile()

			if util.IsURL(bodyFile) {
				content, err = hf.readResponseBodyURL(bodyFile)
			} else {
				content, err = hf.readResponseBodyFile(bodyFile)
			}

			if err != nil {
				result.SetError(fmt.Errorf("data.pairs[%d].response %s", i, err.Error()))
				return result
			}

			pairs[i].Response.Body = content
		}
	}

	return result
}

func (hf *Hoverfly) readResponseBodyURL(fileURL string) (string, error) {
	isAllowed := false
	for _, allowedOrigin := range hf.Cfg.ResponsesBodyFilesAllowedOrigins {
		if strings.HasPrefix(fileURL, allowedOrigin) {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return "", fmt.Errorf("bodyFile %s is not allowed. To allow this origin run hoverfly with -response-body-files-allow-origin", fileURL)
	}

	resp, err := http.DefaultClient.Get(fileURL)
	if err != nil {
		err := fmt.Errorf("bodyFile %s cannot be downloaded: %s", fileURL, err.Error())
		return "", err
	}

	content, err := util.GetResponseBody(resp)
	if err != nil {
		err := fmt.Errorf("response from bodyFile %s cannot be read: %s", fileURL, err.Error())
		return "", err
	}

	return content, nil
}

func (hf *Hoverfly) readResponseBodyFile(filePath string) (string, error) {
	if filepath.IsAbs(filePath) {
		return "", fmt.Errorf("bodyFile contains absolute path (%s). only relative is supported", filePath)
	}

	fileContents, err := ioutil.ReadFile(filepath.Join(hf.Cfg.ResponsesBodyFilesPath, filePath))
	if err != nil {
		return "", err
	}

	return string(fileContents[:]), nil
}

func (hf *Hoverfly) applyBodyTemplating(requestDetails *models.RequestDetails, response *models.ResponseDetails, cachedResponse *models.CachedResponse) (string, error) {
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

	return hf.templator.RenderTemplate(template, requestDetails, hf.state.State)
}

func (hf *Hoverfly) applyHeadersTemplating(requestDetails *models.RequestDetails, response *models.ResponseDetails, cachedResponse *models.CachedResponse) (map[string][]string, error) {
	var headersTemplates map[string][]*raymond.Template
	if cachedResponse != nil && cachedResponse.ResponseHeadersTemplates != nil {
		headersTemplates = cachedResponse.ResponseHeadersTemplates
	} else {
		var header []*raymond.Template
		headersTemplates = map[string][]*raymond.Template{}
		// Parse and cache headers templates
		for k, v := range response.Headers {
			header = make([]*raymond.Template, len(v))
			for i, h := range v {
				header[i], _ = hf.templator.ParseTemplate(h)
			}

			headersTemplates[k] = header
		}

		if cachedResponse != nil {
			cachedResponse.ResponseHeadersTemplates = headersTemplates
		}
	}

	var (
		header []string
		err    error
	)
	headers := map[string][]string{}

	// Render headers templates
	for k, v := range headersTemplates {
		header = make([]string, len(v))
		for i, h := range v {
			header[i], err = hf.templator.RenderTemplate(h, requestDetails, hf.state.State)

			if err != nil {
				return nil, err
			}
		}
		headers[k] = header
	}

	return headers, nil
}

// save gets request fingerprint, extracts request body, status code and headers, then saves it to cache
func (hf *Hoverfly) Save(request *models.RequestDetails, response *models.ResponseDetails, modeArgs *modes.ModeArguments) error {
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

	if len(modeArgs.Headers) >= 1 {
		if modeArgs.Headers[0] == "*" {
			headers = request.Headers
		} else {
			headers = map[string][]string{}
			for _, header := range modeArgs.Headers {
				headerValues := request.Headers[header]
				if len(headerValues) > 0 {
					headers[header] = headerValues
				}
			}
		}
	}

	var requestHeaders map[string][]models.RequestFieldMatchers
	if len(headers) > 0 {
		requestHeaders = map[string][]models.RequestFieldMatchers{}
		for headerKey, headerValues := range headers {
			requestHeaders[headerKey] = []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   strings.Join(headerValues, ";"),
				},
			}
		}
	}

	var queries *models.QueryRequestFieldMatchers
	if len(request.Query) > 0 {
		queries = &models.QueryRequestFieldMatchers{}
		for key, values := range request.Query {
			queries.Add(key, []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   strings.Join(values, ";"),
				},
			})
		}
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
	if modeArgs.Stateful {
		hf.Simulation.AddPairInSequence(&pair, hf.state)
	} else if modeArgs.OverwriteDuplicate {
		hf.Simulation.AddPairWithOverwritingDuplicate(&pair)
	} else {
		hf.Simulation.AddPair(&pair)
	}

	return nil
}

func (hf Hoverfly) ApplyMiddleware(pair models.RequestResponsePair) (models.RequestResponsePair, error) {
	if hf.Cfg.Middleware.IsSet() {
		return hf.Cfg.Middleware.Execute(pair)
	}

	return pair, nil
}
