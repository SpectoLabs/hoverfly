package v2

import (
	"github.com/SpectoLabs/hoverfly/core/metrics"
)

type DestinationView struct {
	Destination string `json:"destination"`
}

type UsageView struct {
	Usage metrics.Stats `json:"usage"`
}

type MiddlewareView struct {
	Binary string `json:"binary"`
	Script string `json:"script"`
	Remote string `json:"remote"`
}

type ModeView struct {
	Mode      string            `json:"mode"`
	Arguments ModeArgumentsView `json:"arguments,omitempty"`
}

type ModeArgumentsView struct {
	Headers          []string `json:"headersWhitelist,omitempty"`
	MatchingStrategy *string  `json:"matchingStrategy,omitempty"`
}

type IsWebServerView struct {
	IsWebServer bool `json:"isWebServer"`
}

type VersionView struct {
	Version string `json:"version"`
}

type UpstreamProxyView struct {
	UpstreamProxy string `json:"upstream-proxy"`
}

type HoverflyView struct {
	DestinationView
	MiddlewareView `json:"middleware"`
	ModeView
	IsWebServerView
	UsageView
	VersionView
	UpstreamProxyView
}

type LogsView struct {
	Logs []map[string]interface{} `json:"logs"`
}

type CacheView struct {
	Cache []CachedResponseView `json:"cache"`
}

type CachedResponseView struct {
	Key          string                            `json:"key"`
	MatchingPair *RequestMatcherResponsePairViewV4 `json:"matchingPair,omitempty"`
	HeaderMatch  bool                              `json:"headerMatch"`
	ClosestMiss  *ClosestMissView                  `json:"closestMiss"`
}

type ClosestMissView struct {
	Response       ResponseDetailsViewV4 `json:"response"`
	RequestMatcher RequestMatcherViewV4  `json:"requestMatcher"`
	MissedFields   []string              `json:"missedFields"`
}

type JournalView struct {
	Journal []JournalEntryView `json:"journal"`
}

type JournalEntryView struct {
	Request     RequestDetailsView  `json:"request"`
	Response    ResponseDetailsView `json:"response"`
	Mode        string              `json:"mode"`
	TimeStarted string              `json:"timeStarted"`
	Latency     float64             `json:"latency"`
}

type JournalEntryFilterView struct {
	Request *RequestMatcherViewV2 `json:"request"`
}

type StateView struct {
	State map[string]string `json:"state"`
}
