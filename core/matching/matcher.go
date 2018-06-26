package matching

import (
	"strings"

	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/state"
)

func Match(strongestMatch string, req models.RequestDetails, webserver bool, simulation *models.Simulation, state *state.State) *MatchingResult {
	if strings.ToLower(strongestMatch) == "strongest" {
		return MatchingStrategyRunner(req, webserver, simulation, state, &StrongestMatchStrategy{})
	} else {
		return MatchingStrategyRunner(req, webserver, simulation, state, &FirstMatchStrategy{})
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
