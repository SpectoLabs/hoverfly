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
	Arguments map[string]string `json:"arguments,omitempty"`
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

type CacheView struct {
	RequestResponsePairs []RequestResponsePairViewV1 `json:"cache"`
}

type LogsView struct {
	Logs []map[string]interface{} `json:"logs"`
}
