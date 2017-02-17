package matching

import (
	"errors"
	"strings"

	"github.com/SpectoLabs/hoverfly/core/models"
	glob "github.com/ryanuber/go-glob"
)

type TemplateMatcher struct{}

func (t TemplateMatcher) Match(req models.RequestDetails, webserver bool, simulation *models.Simulation) (*models.ResponseDetails, error) {
	// iterate through the request templates, looking for template to match request
	for _, entry := range simulation.Templates {
		// TODO: not matching by default on URL and body - need to enable this
		// TODO: need to enable regex matches
		// TODO: enable matching on scheme

		template := entry.RequestTemplate

		if template.Body != nil && !glob.Glob(*template.Body, req.Body) {
			continue
		}

		if !webserver {
			if template.Destination != nil && !glob.Glob(*template.Destination, req.Destination) {
				continue
			}
		}
		if template.Path != nil && !glob.Glob(*template.Path, req.Path) {
			continue
		}
		if template.Query != nil && !glob.Glob(*template.Query, req.Query) {
			continue
		}
		if !headerMatch(template.Headers, req.Headers) {
			continue
		}
		if template.Method != nil && !glob.Glob(*template.Method, req.Method) {
			continue
		}

		// return the first template to match
		return &entry.Response, nil
	}
	return nil, errors.New("No match found")
}

/**
Check keys and corresponding values in template headers are also present in request headers
*/
func headerMatch(templateHeaders, requestHeaders map[string][]string) bool {

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
