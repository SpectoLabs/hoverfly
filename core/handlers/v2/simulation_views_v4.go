package v2

import (
	"github.com/SpectoLabs/hoverfly/core/interfaces"
)

type SimulationViewV4 struct {
	DataViewV4 `json:"data"`
	MetaView   `json:"meta"`
}

type DataViewV4 struct {
	RequestResponsePairs []RequestMatcherResponsePairViewV4 `json:"pairs"`
	GlobalActions        GlobalActionsView                  `json:"globalActions"`
}

type RequestMatcherResponsePairViewV4 struct {
	RequestMatcher RequestMatcherViewV4  `json:"request"`
	Response       ResponseDetailsViewV4 `json:"response"`
}

// RequestDetailsView is used when marshalling and unmarshalling RequestDetails
type RequestMatcherViewV4 struct {
	Path                *RequestFieldMatchersView            `json:"path,omitempty"`
	Method              *RequestFieldMatchersView            `json:"method,omitempty"`
	Destination         *RequestFieldMatchersView            `json:"destination,omitempty"`
	Scheme              *RequestFieldMatchersView            `json:"scheme,omitempty"`
	Query               *RequestFieldMatchersView            `json:"query,omitempty"`
	Body                *RequestFieldMatchersView            `json:"body,omitempty"`
	Headers             map[string][]string                  `json:"headers,omitempty"`
	HeadersWithMatchers map[string]*RequestFieldMatchersView `json:"headersWithMatchers,omitempty"`
	QueriesWithMatchers *QueryMatcherViewV4                  `json:"queriesWithMatchers,omitempty"`
	RequiresState       map[string]string                    `json:"requiresState,omitempty"`
}

type QueryMatcherViewV4 map[string]*RequestFieldMatchersView

//Gets Response - required for interfaces.RequestResponsePairView
func (this RequestMatcherResponsePairViewV4) GetResponse() interfaces.Response { return this.Response }

type ResponseDetailsViewV4 struct {
	Status           int                 `json:"status"`
	Body             string              `json:"body"`
	EncodedBody      bool                `json:"encodedBody"`
	Headers          map[string][]string `json:"headers,omitempty"`
	Templated        bool                `json:"templated"`
	TransitionsState map[string]string   `json:"transitionsState,omitempty"`
	RemovesState     []string            `json:"removesState,omitempty"`
}

//Gets Status - required for interfaces.Response
func (this ResponseDetailsViewV4) GetStatus() int { return this.Status }

// Gets Body - required for interfaces.Response
func (this ResponseDetailsViewV4) GetBody() string { return this.Body }

// Gets EncodedBody - required for interfaces.Response
func (this ResponseDetailsViewV4) GetEncodedBody() bool { return this.EncodedBody }

func (this ResponseDetailsViewV4) GetTemplated() bool { return this.Templated }

func (this ResponseDetailsViewV4) GetTransitionsState() map[string]string {
	return this.TransitionsState
}

func (this ResponseDetailsViewV4) GetRemovesState() []string { return this.RemovesState }

// Gets Headers - required for interfaces.Response
func (this ResponseDetailsViewV4) GetHeaders() map[string][]string { return this.Headers }
