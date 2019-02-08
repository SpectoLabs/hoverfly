package models

import (
	"github.com/aymerick/raymond"
)

type CachedResponse struct {
	Request      RequestDetails
	MatchingPair *RequestMatcherResponsePair
	ClosestMiss  *ClosestMiss
	ResponseTemplate *raymond.Template
}
