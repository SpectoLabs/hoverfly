package errors

import "github.com/SpectoLabs/hoverfly/core/models"

type HoverflyError struct {
	Message    string
	StatusCode int
}

func (err HoverflyError) Error() string {
	return err.Message
}

func MatchingFailedError(closestMiss *models.ClosestMiss) *HoverflyError {
	message := "Could not find a match for request, create or record a valid matcher first!"

	if closestMiss != nil {
		message = message + closestMiss.GetMessage()
	}
	return &HoverflyError{
		Message:    message,
		StatusCode: 412,
	}
}
