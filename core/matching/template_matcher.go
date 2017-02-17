package matching

import (
	"strings"

	glob "github.com/ryanuber/go-glob"
)

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
