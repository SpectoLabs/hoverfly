package matching

import (
	"github.com/SpectoLabs/hoverfly/core/models"
)

func StrongestMatchStrategy(req models.RequestDetails, webserver bool, simulation *models.Simulation, state map[string]string) *MatchingResult {
	var requestMatch *models.RequestMatcherResponsePair
	var err *models.MatchError
	var cachable bool

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

		fieldMatch = HeaderMatching(requestMatcher, req.Headers)
		if !fieldMatch.Matched {
			matched = false
			matchedOnAllButState = false
			missedFields = append(missedFields, "headers")
		}
		matchScore += fieldMatch.MatchScore

		fieldMatch = QueryMatching(requestMatcher, req.Query)
		if !fieldMatch.Matched {
			matched = false
			matchedOnAllButState = false
			missedFields = append(missedFields, "queries")
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

	return &MatchingResult{
		Pair:     requestMatch,
		Error:    err,
		Cachable: cachable,
	}
}
