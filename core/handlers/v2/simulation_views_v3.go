package v2

import (
	"github.com/SpectoLabs/hoverfly/core/interfaces"
)

type SimulationViewV3 struct {
	DataViewV3 `json:"data"`
	MetaView   `json:"meta"`
}

type DataViewV3 struct {
	RequestResponsePairs []RequestMatcherResponsePairViewV3 `json:"pairs"`
	GlobalActions        GlobalActionsView                  `json:"globalActions"`
}

type RequestMatcherResponsePairViewV3 struct {
	Response       ResponseDetailsViewV3 `json:"response"`
	RequestMatcher RequestMatcherViewV3  `json:"request"`
}

// RequestDetailsView is used when marshalling and unmarshalling RequestDetails
type RequestMatcherViewV3 struct {
	Path        *RequestFieldMatchersView `json:"path,omitempty"`
	Method      *RequestFieldMatchersView `json:"method,omitempty"`
	Destination *RequestFieldMatchersView `json:"destination,omitempty"`
	Scheme      *RequestFieldMatchersView `json:"scheme,omitempty"`
	Query       *RequestFieldMatchersView `json:"query,omitempty"`
	Body        *RequestFieldMatchersView `json:"body,omitempty"`
	Headers     map[string][]string       `json:"headers,omitempty"`
}

//Gets Response - required for interfaces.RequestResponsePairView
func (this RequestMatcherResponsePairViewV3) GetResponse() interfaces.Response { return this.Response }

type ResponseDetailsViewV3 struct {
	Status      int                 `json:"status"`
	Body        string              `json:"body"`
	EncodedBody bool                `json:"encodedBody"`
	Headers     map[string][]string `json:"headers,omitempty"`
	Templated   bool                `json:"templated"`
}

//Gets Status - required for interfaces.Response
func (this ResponseDetailsViewV3) GetStatus() int { return this.Status }

// Gets Body - required for interfaces.Response
func (this ResponseDetailsViewV3) GetBody() string { return this.Body }

// Gets EncodedBody - required for interfaces.Response
func (this ResponseDetailsViewV3) GetEncodedBody() bool { return this.EncodedBody }

func (this ResponseDetailsViewV3) GetTemplated() bool { return this.Templated }

// Gets Headers - required for interfaces.Response
func (this ResponseDetailsViewV3) GetHeaders() map[string][]string { return this.Headers }

func (this ResponseDetailsViewV3) GetTransitionsState() map[string]string { return nil }

func (this ResponseDetailsViewV3) GetRemovesState() []string { return nil }
