package matching

import (
	"errors"
	"strings"

	"github.com/SpectoLabs/hoverfly/core/models"
	glob "github.com/ryanuber/go-glob"
)

type TemplateMatcher struct{}

func (t TemplateMatcher) Match(req models.RequestDetails, webserver bool, simulation *models.Simulation) (*models.RequestMatcherResponsePair, error) {

	for _, matchingPair := range simulation.MatchingPairs {
		// TODO: not matching by default on URL and body - need to enable this
		// TODO: need to enable regex matches
		// TODO: enable matching on scheme

		template := matchingPair.RequestMatcher

		if !FieldMatcher(template.Body, req.Body) {
			continue
		}

		if !webserver {
			if !FieldMatcher(template.Destination, req.Destination) {
				continue
			}
		}

		if !FieldMatcher(template.Path, req.Path) {
			continue
		}

		if !FieldMatcher(template.Query, req.Query) {
			continue
		}

		if !FieldMatcher(template.Method, req.Method) {
			continue
		}

		if !HeaderMatch(template.Headers, req.Headers) {
			continue
		}

		// return the first template to match
		return &models.RequestMatcherResponsePair{
			RequestMatcher: template,
			Response:       matchingPair.Response,
		}, nil
	}
	return nil, errors.New("No match found")
}

/**
Check keys and corresponding values in template headers are also present in request headers
*/
func HeaderMatch(templateHeaders, requestHeaders map[string][]string) bool {

	for templateHeaderKey, templateHeaderValues := range templateHeaders {
		for requestHeaderKey, requestHeaderValues := range requestHeaders {
			delete(requestHeaders, requestHeaderKey)
			requestHeaders[strings.ToLower(requestHeaderKey)] = requestHeaderValues

		}

		requestTemplateValues, templateHeaderMatched := requestHeaders[strings.ToLower(templateHeaderKey)]
		if !templateHeaderMatched {
			return false
		}

		for _, templateHeaderValue := range templateHeaderValues {
			templateValueMatched := false
			for _, requestHeaderValue := range requestTemplateValues {
				if glob.Glob(strings.ToLower(templateHeaderValue), strings.ToLower(requestHeaderValue)) {
					templateValueMatched = true
				}
			}

			if !templateValueMatched {
				return false
			}
		}
	}
	return true
}
