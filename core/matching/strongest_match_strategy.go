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
		var score int
		matched := true
		matchedOnAllButHeaders := true
		matchedOnAllButState := true

		requestMatcher := matchingPair.RequestMatcher

		fieldMatch := FieldMatcher(requestMatcher.Body, req.Body)
		if !fieldMatch.Matched {
			matchedOnAllButHeaders = false
			matchedOnAllButState = false
			matched = false
			missedFields = append(missedFields, "body")
		}
		score += fieldMatch.Score

		if !webserver {
			match := FieldMatcher(requestMatcher.Destination, req.Destination)
			if !match.Matched {
				matchedOnAllButHeaders = false
				matchedOnAllButState = false
				matched = false
				missedFields = append(missedFields, "destination")
			}
			score += match.Score
		}

		fieldMatch = FieldMatcher(requestMatcher.Path, req.Path)
		if !fieldMatch.Matched {
			matchedOnAllButHeaders = false
			matchedOnAllButState = false
			matched = false
			missedFields = append(missedFields, "path")
		}
		score += fieldMatch.Score

		fieldMatch = FieldMatcher(requestMatcher.Query, req.QueryString())
		if !fieldMatch.Matched {
			matchedOnAllButHeaders = false
			matchedOnAllButState = false
			matched = false
			missedFields = append(missedFields, "query")
		}
		score += fieldMatch.Score

		fieldMatch = FieldMatcher(requestMatcher.Method, req.Method)
		if !fieldMatch.Matched {
			matchedOnAllButHeaders = false
			matchedOnAllButState = false
			matched = false
			missedFields = append(missedFields, "method")
		}
		score += fieldMatch.Score

		fieldMatch = HeaderMatching(requestMatcher, req.Headers)
		if !fieldMatch.Matched {
			matched = false
			matchedOnAllButState = false
			missedFields = append(missedFields, "headers")
		}
		score += fieldMatch.Score

		fieldMatch = QueryMatching(requestMatcher, req.Query)
		if !fieldMatch.Matched {
			matched = false
			matchedOnAllButState = false
			missedFields = append(missedFields, "queries")
		}
		score += fieldMatch.Score

		fieldMatch = StateMatcher(state, requestMatcher.RequiresState)
		if !fieldMatch.Matched {
			matched = false
			matchedOnAllButHeaders = false
			missedFields = append(missedFields, "state")
		}
		score += fieldMatch.Score

		// This only counts if there was actually a matcher for headers
		if matchedOnAllButHeaders && requestMatcher.Headers != nil && len(requestMatcher.Headers) > 0 {
			matchedOnAllButHeadersAtLeastOnce = true
		}

		// This only counts of there was actually a matcher for state
		if matchedOnAllButState && requestMatcher.RequiresState != nil && len(requestMatcher.RequiresState) > 0 {
			matchedOnAllButStateAtLeastOnce = true
		}

		if matched == true && score >= strongestMatchScore {
			requestMatch = &models.RequestMatcherResponsePair{
				RequestMatcher: requestMatcher,
				Response:       matchingPair.Response,
			}
			strongestMatchScore = score
			closestMiss = nil
		} else if matched == false && requestMatch == nil && score >= closestMissScore {
			closestMissScore = score
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
