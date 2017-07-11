package matching

import (
	"github.com/SpectoLabs/hoverfly/core/models"
)

func StrongestMatchRequestMatcher(req models.RequestDetails, webserver bool, simulation *models.Simulation) (requestMatch *models.RequestMatcherResponsePair, err *models.MatchError) {

	var closestMissScore int
	var strongestMatchScore int
	var closestMiss *models.ClosestMiss
	matchedOnAllButHeadersAtLeastOnce := false

	for _, matchingPair := range simulation.MatchingPairs {
		// TODO: not matching by default on URL and body - need to enable this
		// TODO: enable matching on scheme

		missedFields := make([]string, 0)
		var matchScore int
		matched := true
		matchedOnAllButHeaders := true

		requestMatcher := matchingPair.RequestMatcher

		fieldMatch := ScoredFieldMatcher(requestMatcher.Body, req.Body)
		if !fieldMatch.Matched {
			matchedOnAllButHeaders = false
			matched = false
			missedFields = append(missedFields, "body")
		}
		matchScore += fieldMatch.MatchScore

		if !webserver {
			match := ScoredFieldMatcher(requestMatcher.Destination, req.Destination)
			if !match.Matched {
				matchedOnAllButHeaders = false
				matched = false
				missedFields = append(missedFields, "destination")
			}
			matchScore += match.MatchScore
		}

		fieldMatch = ScoredFieldMatcher(requestMatcher.Path, req.Path)
		if !fieldMatch.Matched {
			matchedOnAllButHeaders = false
			matched = false
			missedFields = append(missedFields, "path")
		}
		matchScore += fieldMatch.MatchScore

		fieldMatch = ScoredFieldMatcher(requestMatcher.Query, req.QueryString())
		if !fieldMatch.Matched {
			matchedOnAllButHeaders = false
			matched = false
			missedFields = append(missedFields, "query")
		}
		matchScore += fieldMatch.MatchScore

		fieldMatch = ScoredFieldMatcher(requestMatcher.Method, req.Method)
		if !fieldMatch.Matched {
			matchedOnAllButHeaders = false
			matched = false
			missedFields = append(missedFields, "method")
		}
		matchScore += fieldMatch.MatchScore

		fieldMatch = CountingHeaderMatcher(requestMatcher.Headers, req.Headers)
		if !fieldMatch.Matched {
			matched = false
			missedFields = append(missedFields, "headers")
			if matchedOnAllButHeaders {
				matchedOnAllButHeadersAtLeastOnce = true
			}
		}
		matchScore += fieldMatch.MatchScore

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
			}
		}
	}

	if requestMatch == nil {
		err = models.NewMatchErrorWithClosestMiss(closestMiss, "No match found", matchedOnAllButHeadersAtLeastOnce)
	}

	return
}
