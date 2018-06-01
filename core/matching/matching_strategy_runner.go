package matching

import (
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/state"
)

type MatchingStrategy interface {
	PreMatching()
	Matching(*FieldMatch, string)
	PostMatching(models.RequestDetails, models.RequestMatcher, models.RequestMatcherResponsePair, *state.State) *MatchingResult
	Result() *MatchingResult
}

func MatchingStrategyRunner(req models.RequestDetails, webserver bool, simulation *models.Simulation, state *state.State, strategy MatchingStrategy) *MatchingResult {
	for _, matchingPair := range simulation.GetMatchingPairs() {
		requestMatcher := matchingPair.RequestMatcher
		strategy.PreMatching()

		strategy.Matching(FieldMatcher(requestMatcher.Body, req.Body), "body")

		if !webserver {
			strategy.Matching(FieldMatcher(requestMatcher.Destination, req.Destination), "destination")
		}

		strategy.Matching(FieldMatcher(requestMatcher.Path, req.Path), "path")

		strategy.Matching(FieldMatcher(requestMatcher.DeprecatedQuery, req.QueryString()), "query")

		strategy.Matching(FieldMatcher(requestMatcher.Method, req.Method), "method")

		strategy.Matching(HeaderMatching(requestMatcher, req.Headers), "headers")

		strategy.Matching(QueryMatching(requestMatcher, req.Query), "queries")

		strategy.Matching(StateMatcher(state, requestMatcher.RequiresState), "state")

		if result := strategy.PostMatching(req, requestMatcher, matchingPair, state); result != nil {
			return result
		}
	}

	return strategy.Result()
}
