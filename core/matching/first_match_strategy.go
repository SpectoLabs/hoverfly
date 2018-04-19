package matching

import (
	"github.com/SpectoLabs/hoverfly/core/models"
)

func FirstMatchStrategy(req models.RequestDetails, webserver bool, simulation *models.Simulation, currentState map[string]string) *MatchingResult {

	matchedOnAllButHeadersAtLeastOnce := false
	matchedOnAllButStateAtLeastOnce := false

	for _, matchingPair := range simulation.GetMatchingPairs() {
		// TODO: not matching by default on URL and body - need to enable this
		// TODO: enable matching on scheme

		requestMatcher := matchingPair.RequestMatcher
		matchedOnAllButHeaders := true
		matchedOnAllButState := true
		isAMatch := true

		if !UnscoredFieldMatcher(requestMatcher.Body, req.Body).Matched {
			matchedOnAllButHeaders = false
			matchedOnAllButState = false
			isAMatch = false
			continue
		}

		if !webserver {
			if !UnscoredFieldMatcher(requestMatcher.Destination, req.Destination).Matched {
				matchedOnAllButHeaders = false
				matchedOnAllButState = false
				isAMatch = false
				continue
			}
		}

		if !UnscoredFieldMatcher(requestMatcher.Path, req.Path).Matched {
			matchedOnAllButHeaders = false
			matchedOnAllButState = false
			isAMatch = false
			continue
		}

		if !UnscoredFieldMatcher(requestMatcher.Query, req.QueryString()).Matched {
			matchedOnAllButHeaders = false
			matchedOnAllButState = false
			isAMatch = false
			continue
		}

		if !UnscoredFieldMatcher(requestMatcher.Method, req.Method).Matched {
			matchedOnAllButHeaders = false
			matchedOnAllButState = false
			isAMatch = false
			continue
		}

		if !HeaderMatching(requestMatcher, req.Headers).Matched {
			matchedOnAllButState = false
			isAMatch = false
		}

		if !QueryMatching(requestMatcher, req.Query).Matched {
			matchedOnAllButState = false
			isAMatch = false
		}

		if !UnscoredStateMatcher(currentState, requestMatcher.RequiresState).Matched {
			matchedOnAllButHeaders = false
			isAMatch = false
		}

		if matchedOnAllButHeaders {
			matchedOnAllButHeadersAtLeastOnce = true
		}

		if matchedOnAllButState {
			matchedOnAllButStateAtLeastOnce = true
		}

		if !isAMatch {
			continue
		}

		// return the first requestMatcher to match
		match := &models.RequestMatcherResponsePair{
			RequestMatcher: requestMatcher,
			Response:       matchingPair.Response,
		}
		return &MatchingResult{
			Pair:     match,
			Error:    nil,
			Cachable: isCachable(match, matchedOnAllButHeadersAtLeastOnce, matchedOnAllButStateAtLeastOnce),
		}
	}

	return &MatchingResult{
		Pair:     nil,
		Error:    models.NewMatchError("No match found", matchedOnAllButHeadersAtLeastOnce),
		Cachable: isCachable(nil, matchedOnAllButHeadersAtLeastOnce, matchedOnAllButStateAtLeastOnce),
	}
}
