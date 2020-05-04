package v2

import (
	"github.com/SpectoLabs/hoverfly/core/delay"
	"github.com/SpectoLabs/hoverfly/core/interfaces"
)

type SimulationViewV6 struct {
	DataViewV6 `json:"data"`
	MetaView   `json:"meta"`
}

type DataViewV6 struct {
	RequestResponsePairs []RequestMatcherResponsePairViewV6 `json:"pairs"`
	GlobalActions        GlobalActionsView                  `json:"globalActions"`
}

type RequestMatcherResponsePairViewV6 struct {
	RequestMatcher RequestMatcherViewV6  `json:"request"`
	Response       ResponseDetailsViewV6 `json:"response"`
}

// RequestMatcherViewV6 is used when marshalling and unmarshalling RequestDetails
type RequestMatcherViewV6 struct {
	Path            []MatcherViewV6            `json:"path,omitempty"`
	Method          []MatcherViewV6            `json:"method,omitempty"`
	Destination     []MatcherViewV6            `json:"destination,omitempty"`
	Scheme          []MatcherViewV6            `json:"scheme,omitempty"`
	Body            []MatcherViewV6            `json:"body,omitempty"`
	Headers         map[string][]MatcherViewV6 `json:"headers,omitempty"`
	Query           *QueryMatcherViewV6        `json:"query,omitempty"`
	RequiresState   map[string]string          `json:"requiresState,omitempty"`
	DeprecatedQuery []MatcherViewV6            `json:"deprecatedQuery,omitempty"`
}

type QueryMatcherViewV6 map[string][]MatcherViewV6

type MatcherViewV6 struct {
	Matcher string                 `json:"matcher"`
	Value   interface{}            `json:"value"`
	Config  map[string]interface{} `json:"config,omitempty"`
}

func NewMatcherViewV6(matcher string, value interface{}) MatcherViewV6 {
	return MatcherViewV6{
		Matcher: matcher,
		Value:   value,
	}
}

// Gets Response - required for interfaces.RequestResponsePairView
func (this RequestMatcherResponsePairViewV6) GetResponse() interfaces.Response { return this.Response }

type ResponseDetailsViewV6 struct {
	Status           int                          `json:"status"`
	Body             string                       `json:"body"`
	BodyFile         string                       `json:"bodyFile"`
	EncodedBody      bool                         `json:"encodedBody"`
	Headers          map[string][]string          `json:"headers,omitempty"`
	Templated        bool                         `json:"templated"`
	TransitionsState map[string]string            `json:"transitionsState,omitempty"`
	RemovesState     []string                     `json:"removesState,omitempty"`
	FixedDelay       int                          `json:"fixedDelay"`
	LogNormalDelay   *delay.LogNormalDelayOptions `json:"logNormalDelay,omitempty"`
}

//Gets Status - required for interfaces.Response
func (this ResponseDetailsViewV6) GetStatus() int { return this.Status }

// Gets Body - required for interfaces.Response
func (this ResponseDetailsViewV6) GetBody() string { return this.Body }

// Gets BodyFile - required for interfaces.Response
func (this ResponseDetailsViewV6) GetBodyFile() string { return this.BodyFile }

// Gets EncodedBody - required for interfaces.Response
func (this ResponseDetailsViewV6) GetEncodedBody() bool { return this.EncodedBody }

func (this ResponseDetailsViewV6) GetTemplated() bool { return this.Templated }

func (this ResponseDetailsViewV6) GetTransitionsState() map[string]string {
	return this.TransitionsState
}

func (this ResponseDetailsViewV6) GetRemovesState() []string { return this.RemovesState }

// Gets Headers - required for interfaces.Response
func (this ResponseDetailsViewV6) GetHeaders() map[string][]string { return this.Headers }

// Gets FixedDelay - required for interfaces.Response
func (this ResponseDetailsViewV6) GetFixedDelay() int { return this.FixedDelay }

// Gets LogNormalDelay - required for interfaces.Response
func (this ResponseDetailsViewV6) GetLogNormalDelay() *delay.LogNormalDelayOptions { return this.LogNormalDelay }