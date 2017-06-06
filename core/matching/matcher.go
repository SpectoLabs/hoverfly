package matching

import "github.com/SpectoLabs/hoverfly/core/models"

type MatchingError struct {
	StatusCode  int
	Description string
}

func (this MatchingError) Error() string {
	return this.Description
}

func MissedError(miss * models.ClosestMiss) * MatchingError {
	description := "Could not find a match for request, create or record a valid matcher first!"

	if miss != nil {
		description = description + miss.GetMessage()
	}
	return &MatchingError{
		StatusCode:  412,
		Description: description,
	}
}
