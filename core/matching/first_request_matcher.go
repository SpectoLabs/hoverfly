package matching

import (
	"github.com/SpectoLabs/hoverfly/core/models"
)

func FirstMatchRequestMatcher(req models.RequestDetails, webserver bool, simulation *models.Simulation, currentState map[string]string) (*models.RequestMatcherResponsePair, *models.MatchError,  bool) {

	matchedOnAllButHeadersAtLeastOnce := false
	matchedOnAllButStateAtLeastOnce := false

	for _, matchingPair := range simulation.MatchingPairs {
		// TODO: not matching by default on URL and body - need to enable this
		// TODO: enable matching on scheme

		requestMatcher := matchingPair.RequestMatcher
		matchedOnAllButHeaders := true
		matchedOnAllButState := true

		if !UnscoredFieldMatcher(requestMatcher.Body, req.Body).Matched {
			matchedOnAllButHeaders = false
			matchedOnAllButState = false
			continue
		}

		if !webserver {
			if !UnscoredFieldMatcher(requestMatcher.Destination, req.Destination).Matched {
				matchedOnAllButHeaders = false
				matchedOnAllButState = false
				continue
			}
		}

		if !UnscoredFieldMatcher(requestMatcher.Path, req.Path).Matched {
			matchedOnAllButHeaders = false
			matchedOnAllButState = false
			continue
		}

		if !UnscoredFieldMatcher(requestMatcher.Query, req.QueryString()).Matched {
			matchedOnAllButHeaders = false
			matchedOnAllButState = false
			continue
		}

		if !UnscoredFieldMatcher(requestMatcher.Method, req.Method).Matched {
			matchedOnAllButHeaders = false
			matchedOnAllButState = false
			continue
		}

		if !CountlessHeaderMatcher(requestMatcher.Headers, req.Headers).Matched {
			if matchedOnAllButHeaders {
				matchedOnAllButHeadersAtLeastOnce = true
			}
			continue
		}

		if !UnscoredStateMatcher(currentState, requestMatcher.RequiresState).Matched {
			if matchedOnAllButState {
				matchedOnAllButStateAtLeastOnce = true
			}
			continue
		}

		// return the first requestMatcher to match
		match := &models.RequestMatcherResponsePair{
			RequestMatcher: requestMatcher,
			Response:       matchingPair.Response,
		}

		return match, nil, isCachable(match, matchedOnAllButHeadersAtLeastOnce, matchedOnAllButStateAtLeastOnce)
	}

	return nil, models.NewMatchError("No match found", matchedOnAllButHeadersAtLeastOnce), isCachable(nil, matchedOnAllButHeadersAtLeastOnce, matchedOnAllButStateAtLeastOnce)
}
