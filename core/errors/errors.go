package errors

import "github.com/SpectoLabs/hoverfly/core/models"

type HoverflyError struct {
	Message string
}

func (err HoverflyError) Error() string {
	return err.Message
}

func NoCacheSetError() *HoverflyError {
	return &HoverflyError{
		Message: "No cache set",
	}
}

func RecordedRequestNotInCacheError() *HoverflyError {
	return &HoverflyError{
		Message: "Could not find recorded request in cache",
	}
}

func MatchingFailedError(closestMiss *models.ClosestMiss) *HoverflyError {
	message := "Could not find a match for request, create or record a valid matcher first!"

	if closestMiss != nil {
		message = message + closestMiss.GetMessage()
	}
	return &HoverflyError{
		Message: message,
	}
}

func MiddlewareNotSetError() *HoverflyError {
	return &HoverflyError{
		Message: "Cannot execute middleware as middleware has not been correctly set",
	}
}
