package matching

import (
	"strings"

	"github.com/SpectoLabs/hoverfly/core/models"
)

func Match(strongestMatch string, req models.RequestDetails, webserver bool, simulation *models.Simulation, state map[string]string) *MatchingResult {
	if strings.ToLower(strongestMatch) == "strongest" {
		return StrongestMatchRequestMatcher(req, webserver, simulation, state)
	} else {
		return FirstMatchRequestMatcher(req, webserver, simulation, state)
	}
}

type MatchingResult struct {
	Pair     *models.RequestMatcherResponsePair
	Error    *models.MatchError
	Cachable bool
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

type MatchingError struct {
	StatusCode  int
	Description string
}

func (this MatchingError) Error() string {
	return this.Description
}

func MissedError(miss *models.ClosestMiss) *MatchingError {
	description := "Could not find a match for request, create or record a valid matcher first!"

	if miss != nil {
		description = description + miss.GetMessage()
	}
	return &MatchingError{
		StatusCode:  412,
		Description: description,
	}
}
