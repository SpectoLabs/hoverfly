package models

import (
	"github.com/aymerick/raymond"
)

type CachedResponse struct {
	Request                  RequestDetails
	MatchingPair             *RequestMatcherResponsePair
	ClosestMiss              *ClosestMiss
	ResponseStateTemplates   map[string]*raymond.Template
	ResponseTemplate         *raymond.Template
	ResponseHeadersTemplates map[string][]*raymond.Template
}
