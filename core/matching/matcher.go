package matching

import (
	"strings"

	"github.com/SpectoLabs/hoverfly/core/models"
)

func Match(strongestMatch string, req models.RequestDetails, webserver bool, simulation *models.Simulation, state map[string]string) (requestMatch *models.RequestMatcherResponsePair, err *models.MatchError, cachable bool) {
	if strings.ToLower(strongestMatch) == "strongest" {
		return StrongestMatchRequestMatcher(req, webserver, simulation, state)
	} else {
		return FirstMatchRequestMatcher(req, webserver, simulation, state)
	}
}

type MatchingError struct {
	StatusCode  int
	Description string
}

func (this MatchingError) Error() string {
	return this.Description
}

func MissedError(miss *models.ClosestMiss) *MatchingError {
	description := "Could not find a match for request, create or record a valid matcher first!"

	if miss != nil {
		description = description + miss.GetMessage()
	}
	return &MatchingError{
		StatusCode:  412,
		Description: description,
	}
}
