package matching

import (
	"github.com/SpectoLabs/hoverfly/core/models"
)

func StrongestMatchRequestMatcher(req models.RequestDetails, webserver bool, simulation *models.Simulation, state map[string]string) (requestMatch *models.RequestMatcherResponsePair, err *models.MatchError, cachable bool) {

	var closestMissScore int
	var strongestMatchScore int
	var closestMiss *models.ClosestMiss
	matchedOnAllButHeadersAtLeastOnce := false
	matchedOnAllButStateAtLeastOnce := false
	cachable = true

	for _, matchingPair := range simulation.GetMatchingPairs() {
		// TODO: not matching by default on URL and body - need to enable this
		// TODO: enable matching on scheme

		missedFields := make([]string, 0)
		var matchScore int
		matched := true
		matchedOnAllButHeaders := true
		matchedOnAllButState := true

		requestMatcher := matchingPair.RequestMatcher

		fieldMatch := ScoredFieldMatcher(requestMatcher.Body, req.Body)
		if !fieldMatch.Matched {
			matchedOnAllButHeaders = false
			matchedOnAllButState = false
			matched = false
			missedFields = append(missedFields, "body")
		}
		matchScore += fieldMatch.MatchScore

		if !webserver {
			match := ScoredFieldMatcher(requestMatcher.Destination, req.Destination)
			if !match.Matched {
				matchedOnAllButHeaders = false
				matchedOnAllButState = false
				matched = false
				missedFields = append(missedFields, "destination")
			}
			matchScore += match.MatchScore
		}

		fieldMatch = ScoredFieldMatcher(requestMatcher.Path, req.Path)
		if !fieldMatch.Matched {
			matchedOnAllButHeaders = false
			matchedOnAllButState = false
			matched = false
			missedFields = append(missedFields, "path")
		}
		matchScore += fieldMatch.MatchScore

		fieldMatch = ScoredFieldMatcher(requestMatcher.Query, req.QueryString())
		if !fieldMatch.Matched {
			matchedOnAllButHeaders = false
			matchedOnAllButState = false
			matched = false
			missedFields = append(missedFields, "query")
		}
		matchScore += fieldMatch.MatchScore

		fieldMatch = ScoredFieldMatcher(requestMatcher.Method, req.Method)
		if !fieldMatch.Matched {
			matchedOnAllButHeaders = false
			matchedOnAllButState = false
			matched = false
			missedFields = append(missedFields, "method")
		}
		matchScore += fieldMatch.MatchScore

		fieldMatch = CountingHeaderMatcher(requestMatcher.Headers, req.Headers)
		if !fieldMatch.Matched {
			matched = false
			matchedOnAllButState = false
			missedFields = append(missedFields, "headers")
		}
		matchScore += fieldMatch.MatchScore

		fieldMatch = ScoredStateMatcher(state, requestMatcher.RequiresState)
		if !fieldMatch.Matched {
			matched = false
			matchedOnAllButHeaders = false
			missedFields = append(missedFields, "state")
		}
		matchScore += fieldMatch.MatchScore

		// This only counts if there was actually a matcher for headers
		if matchedOnAllButHeaders && requestMatcher.Headers != nil && len(requestMatcher.Headers) > 0 {
			matchedOnAllButHeadersAtLeastOnce = true
		}

		// This only counts of there was actually a matcher for state
		if matchedOnAllButState && requestMatcher.RequiresState != nil && len(requestMatcher.RequiresState) > 0 {
			matchedOnAllButStateAtLeastOnce = true
		}

		if matched == true && matchScore >= strongestMatchScore {
			requestMatch = &models.RequestMatcherResponsePair{
				RequestMatcher: requestMatcher,
				Response:       matchingPair.Response,
			}
			strongestMatchScore = matchScore
			closestMiss = nil
		} else if matched == false && requestMatch == nil && matchScore >= closestMissScore {
			closestMissScore = matchScore
			view := matchingPair.BuildView()
			closestMiss = &models.ClosestMiss{
				RequestDetails: req,
				RequestMatcher: view.RequestMatcher,
				Response:       view.Response,
				MissedFields:   missedFields,
				State:          state,
			}
		}
	}

	cachable = isCachable(requestMatch, matchedOnAllButHeadersAtLeastOnce, matchedOnAllButStateAtLeastOnce)

	if requestMatch == nil {
		err = models.NewMatchErrorWithClosestMiss(closestMiss, "No match found", matchedOnAllButHeadersAtLeastOnce)
	}

	return
}

func isCachable(requestMatch *models.RequestMatcherResponsePair, matchedOnAllButHeadersAtLeastOnce bool, matchedOnAllButStateAtLeastOnce bool) bool {
	// Do not cache misses if the only thing they missed on was headers/state because a subsequent request which is the same
	// but with different headers/state will need to go through matching
	if requestMatch == nil && (matchedOnAllButHeadersAtLeastOnce || matchedOnAllButStateAtLeastOnce) {
		return false
	} else if requestMatch != nil {

		// And do not cache hits if they matched on headers because a subsequent request which is the same
		// but with different headers wouldn't match
		if requestMatch.RequestMatcher.IncludesHeaderMatching() {
			return false
		}

		// And do not cache hits if another request matched on all but headers, as it could be stronger match
		if matchedOnAllButHeadersAtLeastOnce {
			return false
		}

		// And do not cache hits if they matched on state because a subsequent request which is the same
		// but with different state wouldn't match
		if requestMatch.RequestMatcher.IncludesStateMatching() {
			return false
		}

		// And don't cache hits if another matcher matched on everything apart from state, as we would potentially hit
		// that matcher in the future if state changed
		if matchedOnAllButStateAtLeastOnce {
			return false
		}
	}

	return true
}
