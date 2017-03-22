package matching

import (
	"errors"

	"github.com/SpectoLabs/hoverfly/core/models"
)

func RequestMatcher(req models.RequestDetails, webserver bool, simulation *models.Simulation) (*models.RequestMatcherResponsePair, error) {

	for _, matchingPair := range simulation.MatchingPairs {
		// TODO: not matching by default on URL and body - need to enable this
		// TODO: need to enable regex matches
		// TODO: enable matching on scheme

		requestMatcher := matchingPair.RequestMatcher

		if !FieldMatcher(requestMatcher.Body, req.Body) {
			continue
		}

		if !webserver {
			if !FieldMatcher(requestMatcher.Destination, req.Destination) {
				continue
			}
		}

		if !FieldMatcher(requestMatcher.Path, req.Path) {
			continue
		}

		if !FieldMatcher(requestMatcher.Query, req.Query) {
			continue
		}

		if !FieldMatcher(requestMatcher.Method, req.Method) {
			continue
		}

		if !HeaderMatcher(requestMatcher.Headers, req.Headers) {
			continue
		}

		// return the first requestMatcher to match
		return &models.RequestMatcherResponsePair{
			RequestMatcher: requestMatcher,
			Response:       matchingPair.Response,
		}, nil
	}
	return nil, errors.New("No match found")
}
