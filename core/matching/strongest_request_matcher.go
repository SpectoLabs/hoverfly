package matching

import (
	"errors"
	"github.com/SpectoLabs/hoverfly/core/models"
)


func StrongestMatchRequestMatcher(req models.RequestDetails, webserver bool, simulation *models.Simulation) (requestMatch, closestMiss *models.RequestMatcherResponsePair, err error) {

	var closestMissTotalMatches int
	var strongestMatchTotalMatches int

	for _, matchingPair := range simulation.MatchingPairs {
		// TODO: not matching by default on URL and body - need to enable this
		// TODO: enable matching on scheme

		var matchScore int
		matched := true

		requestMatcher := matchingPair.RequestMatcher

		fieldMatch := CountingFieldMatcher(requestMatcher.Body, req.Body)
		if !fieldMatch.Matched {
			matched = false
		}
		matchScore += fieldMatch.MatchScore

		if !webserver {
			match := CountingFieldMatcher(requestMatcher.Destination, req.Destination)
			if !match.Matched {
				matched = false
			}
			matchScore += match.MatchScore
		}

		fieldMatch = CountingFieldMatcher(requestMatcher.Path, req.Path)
		if !fieldMatch.Matched {
			matched = false
		}
		matchScore += fieldMatch.MatchScore

		fieldMatch = CountingFieldMatcher(requestMatcher.Query, req.Query)
		if !fieldMatch.Matched {
			matched = false
		}
		matchScore += fieldMatch.MatchScore

		fieldMatch = CountingFieldMatcher(requestMatcher.Method, req.Method)
		if !fieldMatch.Matched {
			matched = false
		}
		matchScore += fieldMatch.MatchScore

		fieldMatch = CountingHeaderMatcher(requestMatcher.Headers, req.Headers)
		if !fieldMatch.Matched {
			matched = false
		}
		matchScore += fieldMatch.MatchScore

		if matched == true && matchScore >= strongestMatchTotalMatches {
			requestMatch = &models.RequestMatcherResponsePair{
				RequestMatcher: requestMatcher,
				Response:       matchingPair.Response,
			}
			strongestMatchTotalMatches = matchScore
			closestMiss = nil
		} else if matched == false && requestMatch == nil && matchScore >= closestMissTotalMatches {
			closestMissTotalMatches = matchScore
			closestMiss = &models.RequestMatcherResponsePair{}
			*closestMiss = matchingPair
		}
	}

	if requestMatch == nil {
		err = errors.New("No match found")
	}

	return
}
