package matching

import (
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/state"
	"github.com/SpectoLabs/hoverfly/core/util"
)

type StrongestMatchStrategy struct {
	matchedOnAllButHeaders            bool
	matchedOnAllButState              bool
	matchedOnAllButHeadersAtLeastOnce bool
	matchedOnAllButStateAtLeastOnce   bool
	matched                           bool
	score                             int
	strongestMatchScore               int
	closestMissScore                  int
	closestMiss                       *models.ClosestMiss
	missedFields                      []string
	requestMatch                      *models.RequestMatcherResponsePair
}

func (s *StrongestMatchStrategy) PreMatching() {
	s.matched = true
	s.missedFields = make([]string, 0)
	s.matchedOnAllButHeaders = true
	s.matchedOnAllButState = true
	s.score = 0
}

func (s *StrongestMatchStrategy) Matching(fieldMatch *FieldMatch, field string) {
	if !fieldMatch.Matched {
		if field != "headers" {
			s.matchedOnAllButHeaders = false
		}
		if field != "state" {
			s.matchedOnAllButState = false
		}
		s.matched = false
		s.missedFields = append(s.missedFields, field)
	}
	s.score += fieldMatch.Score
}

func (s *StrongestMatchStrategy) PostMatching(req models.RequestDetails, requestMatcher models.RequestMatcher, matchingPair models.RequestMatcherResponsePair, state *state.State) *MatchingResult {
	// This only counts if there was actually a matcher for headers
	if s.matchedOnAllButHeaders && requestMatcher.Headers != nil && len(requestMatcher.Headers) > 0 {
		s.matchedOnAllButHeadersAtLeastOnce = true
	}

	// This only counts of there was actually a matcher for state
	if s.matchedOnAllButState && requestMatcher.RequiresState != nil && len(requestMatcher.RequiresState) > 0 {
		s.matchedOnAllButStateAtLeastOnce = true
	}

	if s.matched == true && s.score >= s.strongestMatchScore {
		s.requestMatch = &models.RequestMatcherResponsePair{
			RequestMatcher: requestMatcher,
			Response:       matchingPair.Response,
		}
		s.strongestMatchScore = s.score
		s.closestMiss = nil
	} else if s.matched == false && s.requestMatch == nil && s.score >= s.closestMissScore {
		s.closestMissScore = s.score
		view := matchingPair.BuildView()
		state.RWMutex.RLock()
		copyState := util.CopyMap(state.State)
		state.RWMutex.RUnlock()
		s.closestMiss = &models.ClosestMiss{
			RequestDetails: req,
			RequestMatcher: view.RequestMatcher,
			Response:       view.Response,
			MissedFields:   s.missedFields,
			State:          copyState,
		}
	}

	return nil
}

func (s *StrongestMatchStrategy) Result() *MatchingResult {
	cachable := isCachable(s.requestMatch, s.matchedOnAllButHeadersAtLeastOnce, s.matchedOnAllButStateAtLeastOnce)
	var err *models.MatchError
	if s.requestMatch == nil {
		err = models.NewMatchErrorWithClosestMiss(s.closestMiss, "No match found")
	}

	return &MatchingResult{
		Pair:     s.requestMatch,
		Error:    err,
		Cachable: cachable,
	}
}
