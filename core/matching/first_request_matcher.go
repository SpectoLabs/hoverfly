package matching

import (
	"errors"

	"github.com/SpectoLabs/hoverfly/core/models"
)

func FirstMatchRequestMatcher(req models.RequestDetails, webserver bool, simulation *models.Simulation) (*models.RequestMatcherResponsePair, error) {

	for _, matchingPair := range simulation.MatchingPairs {
		// TODO: not matching by default on URL and body - need to enable this
		// TODO: enable matching on scheme

		requestMatcher := matchingPair.RequestMatcher

		if !UnscoredFieldMatcher(requestMatcher.Body, req.Body).Matched {
			continue
		}

		if !webserver {
			if !UnscoredFieldMatcher(requestMatcher.Destination, req.Destination).Matched {
				continue
			}
		}

		if !UnscoredFieldMatcher(requestMatcher.Path, req.Path).Matched {
			continue
		}

		if !UnscoredFieldMatcher(requestMatcher.Query, req.Query).Matched {
			continue
		}

		if !UnscoredFieldMatcher(requestMatcher.Method, req.Method).Matched {
			continue
		}

		if !CountlessHeaderMatcher(requestMatcher.Headers, req.Headers).Matched {
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