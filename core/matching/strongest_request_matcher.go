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

		var totalMatches int
		matched := true

		requestMatcher := matchingPair.RequestMatcher

		fieldMatch := CountingFieldMatcher(requestMatcher.Body, req.Body)
		if !fieldMatch.Matched {
			matched = false
		}
		totalMatches += fieldMatch.TotalMatches

		if !webserver {
			match := CountingFieldMatcher(requestMatcher.Destination, req.Destination)
			if !match.Matched {
				matched = false
			}
			totalMatches += match.TotalMatches
		}

		fieldMatch = CountingFieldMatcher(requestMatcher.Path, req.Path)
		if !fieldMatch.Matched {
			matched = false
		}
		totalMatches += fieldMatch.TotalMatches

		fieldMatch = CountingFieldMatcher(requestMatcher.Query, req.Query)
		if !fieldMatch.Matched {
			matched = false
		}
		totalMatches += fieldMatch.TotalMatches

		fieldMatch = CountingFieldMatcher(requestMatcher.Method, req.Method)
		if !fieldMatch.Matched {
			matched = false
		}
		totalMatches += fieldMatch.TotalMatches

		fieldMatch = CountingHeaderMatcher(requestMatcher.Headers, req.Headers)
		if !fieldMatch.Matched {
			matched = false
		}
		totalMatches += fieldMatch.TotalMatches

		if matched == true && totalMatches >= strongestMatchTotalMatches {
			requestMatch = &models.RequestMatcherResponsePair{
				RequestMatcher: requestMatcher,
				Response:       matchingPair.Response,
			}
			strongestMatchTotalMatches = totalMatches
			closestMiss = nil
		} else if matched == false && requestMatch == nil && totalMatches >= closestMissTotalMatches {
			closestMissTotalMatches = totalMatches
			closestMiss = &models.RequestMatcherResponsePair{}
			*closestMiss = matchingPair
		}
	}

	if requestMatch == nil {
		err = errors.New("No match found")
	}

	return
}
