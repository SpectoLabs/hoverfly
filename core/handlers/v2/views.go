package v2

import "github.com/SpectoLabs/hoverfly/core/metrics"

type DestinationView struct {
	Destination string `json:"destination"`
}

type UsageView struct {
	Usage metrics.Stats `json:"usage"`
}

type MiddlewareView struct {
	Middleware string `json:"middleware"`
}

type ModeView struct {
	Mode string `json:"mode"`
}
