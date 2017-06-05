package matching

import (
	"errors"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"fmt"
	"encoding/json"
	"strings"
)


func StrongestMatchRequestMatcher(req models.RequestDetails, webserver bool, simulation *models.Simulation) (requestMatch *models.RequestMatcherResponsePair, closestMiss *ClosestMiss, err error) {

	var closestMissScore int
	var strongestMatchScore int

	for _, matchingPair := range simulation.MatchingPairs {
		// TODO: not matching by default on URL and body - need to enable this
		// TODO: enable matching on scheme

		missedFields := make([]string, 0)
		var matchScore int
		matched := true

		requestMatcher := matchingPair.RequestMatcher

		fieldMatch := ScoredFieldMatcher(requestMatcher.Body, req.Body)
		if !fieldMatch.Matched {
			matched = false
			missedFields = append(missedFields, "body")
		}
		matchScore += fieldMatch.MatchScore

		if !webserver {
			match := ScoredFieldMatcher(requestMatcher.Destination, req.Destination)
			if !match.Matched {
				matched = false
				missedFields = append(missedFields, "destination")
			}
			matchScore += match.MatchScore
		}

		fieldMatch = ScoredFieldMatcher(requestMatcher.Path, req.Path)
		if !fieldMatch.Matched {
			matched = false
			missedFields = append(missedFields, "path")
		}
		matchScore += fieldMatch.MatchScore

		fieldMatch = ScoredFieldMatcher(requestMatcher.Query, req.Query)
		if !fieldMatch.Matched {
			matched = false
			missedFields = append(missedFields, "query")
		}
		matchScore += fieldMatch.MatchScore

		fieldMatch = ScoredFieldMatcher(requestMatcher.Method, req.Method)
		if !fieldMatch.Matched {
			matched = false
			missedFields = append(missedFields, "method")
		}
		matchScore += fieldMatch.MatchScore

		fieldMatch = CountingHeaderMatcher(requestMatcher.Headers, req.Headers)
		if !fieldMatch.Matched {
			matched = false
			missedFields = append(missedFields, "headers")
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
			closestMiss = &ClosestMiss{
				RequestDetails: &req,
				RequestMatcher: &view.RequestMatcher,
				Response: &view.Response,
				MissedFields: missedFields,
			}
		}
	}

	if requestMatch == nil {
		err = errors.New("No match found")
	}

	return
}

type ClosestMiss struct {
	RequestDetails * models.RequestDetails
	RequestMatcher * v2.RequestMatcherViewV2
	Response * v2.ResponseDetailsView
	MissedFields []string
}

func (this *ClosestMiss) GetMessage() string {

	requestBytes, _ := json.MarshalIndent(this.RequestDetails, "", "    ")
	matcherBytes, _ := json.MarshalIndent(this.RequestMatcher, "", "    ")
	responseBytes, _ := json.MarshalIndent(this.Response, "", "    ")

	return "\n\nThe following request was made, but was not matched by Hoverfly:\n\n" +
		string(requestBytes) +
		"\n\nThe closest miss was the following matcher:\n\n" +
		string(matcherBytes) +
		"\n\nBut it did not match on the following fields:\n\n" +
		fmt.Sprint("[" + strings.Join(this.MissedFields, ", ") + "]")  +
		"\n\nWhich if hit would have given the following response:\n\n" +
		string(responseBytes)
}