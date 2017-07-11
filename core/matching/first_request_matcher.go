package matching

import (
	"github.com/SpectoLabs/hoverfly/core/models"
)

func FirstMatchRequestMatcher(req models.RequestDetails, webserver bool, simulation *models.Simulation) (*models.RequestMatcherResponsePair, *models.MatchError) {

	matchedOnAllButHeadersAtLeastOnce := false

	for _, matchingPair := range simulation.MatchingPairs {
		// TODO: not matching by default on URL and body - need to enable this
		// TODO: enable matching on scheme

		requestMatcher := matchingPair.RequestMatcher
		matchedOnAllButHeaders := true

		if !UnscoredFieldMatcher(requestMatcher.Body, req.Body).Matched {
			matchedOnAllButHeaders = false
			continue
		}

		if !webserver {
			if !UnscoredFieldMatcher(requestMatcher.Destination, req.Destination).Matched {
				matchedOnAllButHeaders = false
				continue
			}
		}

		if !UnscoredFieldMatcher(requestMatcher.Path, req.Path).Matched {
			matchedOnAllButHeaders = false
			continue
		}

		if !UnscoredFieldMatcher(requestMatcher.Query, req.QueryString()).Matched {
			matchedOnAllButHeaders = false
			continue
		}

		if !UnscoredFieldMatcher(requestMatcher.Method, req.Method).Matched {
			matchedOnAllButHeaders = false
			continue
		}

		if !CountlessHeaderMatcher(requestMatcher.Headers, req.Headers).Matched {
			if matchedOnAllButHeaders {
				matchedOnAllButHeadersAtLeastOnce = true
			}
			continue
		}

		// return the first requestMatcher to match
		return &models.RequestMatcherResponsePair{
			RequestMatcher: requestMatcher,
			Response:       matchingPair.Response,
		}, nil
	}
	return nil, models.NewMatchError("No match found", matchedOnAllButHeadersAtLeastOnce)
}
