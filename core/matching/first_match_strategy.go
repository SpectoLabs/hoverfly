package matching

import (
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/state"
)

type FirstMatchStrategy struct {
	matchedOnAllButHeaders            bool
	matchedOnAllButState              bool
	matched                           bool
	matchedOnAllButHeadersAtLeastOnce bool
	matchedOnAllButStateAtLeastOnce   bool
	matchingPair                      *models.RequestMatcherResponsePair
}

func (s *FirstMatchStrategy) PreMatching() {
	s.matchedOnAllButHeaders = true
	s.matchedOnAllButState = true
	s.matched = true
}

func (s *FirstMatchStrategy) Matching(fieldMatch *FieldMatch, field string) {
	if !fieldMatch.Matched {

		if field != "headers" {
			s.matchedOnAllButState = false

		}
		if field != "state" {
			s.matchedOnAllButHeaders = false

		}
		s.matched = false
	}
}

func (s *FirstMatchStrategy) PostMatching(req models.RequestDetails, requestMatcher models.RequestMatcher, matchingPair models.RequestMatcherResponsePair, state *state.State) *MatchingResult {
	if s.matchedOnAllButHeaders {
		s.matchedOnAllButHeadersAtLeastOnce = true
	}

	if s.matchedOnAllButState {
		s.matchedOnAllButStateAtLeastOnce = true
	}
	if s.matched && s.matchingPair == nil {
		s.matchingPair = &matchingPair
		return s.Result()
	}

	return nil
}

func (s *FirstMatchStrategy) Result() *MatchingResult {
	if s.matchedOnAllButHeaders {
		s.matchedOnAllButHeadersAtLeastOnce = true
	}

	if s.matchedOnAllButState {
		s.matchedOnAllButStateAtLeastOnce = true
	}
	if s.matchingPair != nil {

		return &MatchingResult{
			Pair:     s.matchingPair,
			Error:    nil,
			Cachable: isCachable(s.matchingPair, s.matchedOnAllButHeadersAtLeastOnce, s.matchedOnAllButStateAtLeastOnce),
		}
	}

	return &MatchingResult{
		Pair:     nil,
		Error:    models.NewMatchError("No match found"),
		Cachable: isCachable(nil, s.matchedOnAllButHeadersAtLeastOnce, s.matchedOnAllButStateAtLeastOnce),
	}
}
