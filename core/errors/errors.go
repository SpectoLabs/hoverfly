package errors

import "github.com/SpectoLabs/hoverfly/core/models"

type HoverflyError struct {
	Message    string
	StatusCode int
}

func (err HoverflyError) Error() string {
	return err.Message
}

func NoCacheSetError() *HoverflyError {
	return &HoverflyError{
		Message:    "No cache set",
		StatusCode: 412,
	}
}

func RecordedRequestNotInCacheError() *HoverflyError {
	return &HoverflyError{
		Message:    "Could not find recorded request in cache",
		StatusCode: 412,
	}
}

func DecodePayloadError() *HoverflyError {
	return &HoverflyError{
		Message:    "Failed to decode payload from cache",
		StatusCode: 500,
	}
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

func MiddlewareNotSetError() *HoverflyError {
	return &HoverflyError{
		Message: "Cannot execute middleware as middleware has not been correctly set",
	}
}
