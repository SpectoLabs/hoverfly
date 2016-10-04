package v2

import (
	"github.com/SpectoLabs/hoverfly/core/metrics"
	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
)

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

type HoverflyView struct {
	DestinationView
	MiddlewareView
	ModeView
	UsageView
}

type SimulationView struct {
	DataView `json:"data"`
	MetaView `json:"meta"`
}

type DataView struct {
	RequestResponsePairs []RequestResponsePairView `json:"pairs"`
	GlobalActions GlobalActionsView `json:"globalActions"`
}

type RequestResponsePairView struct {
	Response ResponseDetailsView `json:"response"`
	Request  RequestDetailsView  `json:"request"`
}

// PayloadView is used when marshalling and unmarshalling payloads.
type RequestDetailsView struct {
	RequestType *string             `json:"requestType"`
	Path        *string             `json:"path"`
	Method      *string             `json:"method"`
	Destination *string             `json:"destination"`
	Scheme      *string             `json:"scheme"`
	Query       *string             `json:"query"`
	Body        *string             `json:"body"`
	Headers     map[string][]string `json:"headers"`
}

// RequestDetailsView is used when marshalling and unmarshalling RequestDetails

// ResponseDetailsView is used when marshalling and
// unmarshalling requests. This struct's Body may be Base64
// encoded based on the EncodedBody field.
type ResponseDetailsView struct {
	Status      int                 `json:"status"`
	Body        string              `json:"body"`
	EncodedBody bool                `json:"encodedBody"`
	Headers     map[string][]string `json:"headers"`
}

type GlobalActionsView struct {
	Delays []v1.ResponseDelayView `json:"delays"`
}

type MetaView struct {
	SchemaVersion string `json:"schemaVersion"`
	HoverflyVersion string `json:"hoverflyVersion"`
	TimeExported string `json:"timeExported"`
}