package matching

import (
	"errors"
	"github.com/SpectoLabs/hoverfly/core/models"
	"fmt"
)


func StrongestMatchRequestMatcher(req models.RequestDetails, webserver bool, simulation *models.Simulation) (requestMatch, closestMatch *models.RequestMatcherResponsePair, err error) {

	var closestMatchTotalMatches int
	var strongestMatchTotalMatches int

	for _, matchingPair := range simulation.MatchingPairs {
		// TODO: not matching by default on URL and body - need to enable this
		// TODO: enable matching on scheme

		var totalMatches int
		matched := true

		requestMatcher := matchingPair.RequestMatcher

		fieldMatch := CountingFieldMatcher(requestMatcher.Body, req.Body)
		if !fieldMatch.Matched {
			fmt.Println("Did not match on body")
			matched = false
		}
		totalMatches += fieldMatch.TotalMatches

		if !webserver {
			match := CountingFieldMatcher(requestMatcher.Destination, req.Destination)
			if !match.Matched {
				fmt.Println("Did not match on destination")
				matched = false
			}
			totalMatches += match.TotalMatches
		}

		fieldMatch = CountingFieldMatcher(requestMatcher.Path, req.Path)
		if !fieldMatch.Matched {
			fmt.Println("Did not match on path")
			matched = false
		}
		totalMatches += fieldMatch.TotalMatches

		fieldMatch = CountingFieldMatcher(requestMatcher.Query, req.Query)
		if !fieldMatch.Matched {
			fmt.Println("Did not match on query")
			matched = false
		}
		totalMatches += fieldMatch.TotalMatches

		fieldMatch = CountingFieldMatcher(requestMatcher.Method, req.Method)
		if !fieldMatch.Matched {
			fmt.Println("Did not match on method")
			matched = false
		}
		totalMatches += fieldMatch.TotalMatches

		if !HeaderMatcher(requestMatcher.Headers, req.Headers) {
			fmt.Println("Did not match on headers")
			matched = false
		}

		if matched == true && totalMatches >= strongestMatchTotalMatches {
			requestMatch = &models.RequestMatcherResponsePair{
				RequestMatcher: requestMatcher,
				Response:       matchingPair.Response,
			}
			strongestMatchTotalMatches = totalMatches
		} else if matched == false && requestMatch == nil && totalMatches >= closestMatchTotalMatches {
			closestMatchTotalMatches = totalMatches
			closestMatch = &models.RequestMatcherResponsePair{}
			*closestMatch = matchingPair
			fmt.Println("Not Matched")
		}
	}

	if requestMatch == nil {
		err = errors.New("No match found")
	}

	return
}
