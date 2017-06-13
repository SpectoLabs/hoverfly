package v2

import (
	"time"

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
	MatchingPair *RequestMatcherResponsePairViewV2 `json:"matchingPair,omitempty"`
	HeaderMatch  bool                              `json:"headerMatch"`
	ClosestMiss  *ClosestMissView                  `json:"closestMiss"`
}

type JournalEntryView struct {
	Request     RequestDetailsViewV1 `json:"request"`
	Response    ResponseDetailsView  `json:"response"`
	Mode        string               `json:"mode"`
	TimeStarted time.Time            `json:"timeStarted"`
	Latency     time.Duration        `json:"latency"`
}
