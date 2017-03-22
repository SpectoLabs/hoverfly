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

		if !HeaderMatcher(template.Headers, req.Headers) {
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
