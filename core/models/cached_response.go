package models

type CachedResponse struct {
	Request      RequestDetails
	MatchingPair *RequestMatcherResponsePair
	ClosestMiss  *ClosestMiss
}
