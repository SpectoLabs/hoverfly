package v2

import (
	"github.com/SpectoLabs/hoverfly/core/interfaces"
)

type SimulationViewV5 struct {
	DataViewV5 `json:"data"`
	MetaView   `json:"meta"`
}

type DataViewV5 struct {
	RequestResponsePairs []RequestMatcherResponsePairViewV5 `json:"pairs"`
	GlobalActions        GlobalActionsView                  `json:"globalActions"`
}

type RequestMatcherResponsePairViewV5 struct {
	RequestMatcher RequestMatcherViewV5  `json:"request"`
	Response       ResponseDetailsViewV5 `json:"response"`
}

// RequestDetailsView is used when marshalling and unmarshalling RequestDetails
type RequestMatcherViewV5 struct {
	Path            []MatcherViewV5            `json:"path,omitempty"`
	Method          []MatcherViewV5            `json:"method,omitempty"`
	Destination     []MatcherViewV5            `json:"destination,omitempty"`
	Scheme          []MatcherViewV5            `json:"scheme,omitempty"`
	Body            []MatcherViewV5            `json:"body,omitempty"`
	Headers         map[string][]MatcherViewV5 `json:"headers,omitempty"`
	Query           *QueryMatcherViewV5        `json:"query,omitempty"`
	RequiresState   map[string]string          `json:"requiresState,omitempty"`
	DeprecatedQuery []MatcherViewV5            `json:"deprecatedQuery,omitempty"`
}

type QueryMatcherViewV5 map[string][]MatcherViewV5

type MatcherViewV5 struct {
	Matcher string                 `json:"matcher"`
	Value   interface{}            `json:"value"`
	Config  map[string]interface{} `json:"config,omitempty"`
}

func NewMatcherView(matcher string, value interface{}) MatcherViewV5 {
	return MatcherViewV5{
		Matcher: matcher,
		Value:   value,
	}
}

//Gets Response - required for interfaces.RequestResponsePairView
func (this RequestMatcherResponsePairViewV5) GetResponse() interfaces.Response { return this.Response }

type ResponseDetailsViewV5 struct {
	Status           int                    `json:"status"`
	Body             string                 `json:"body"`
	BodyFile         string              	`json:"bodyFile,omitempty"`
	EncodedBody      bool                   `json:"encodedBody"`
	Headers          map[string][]string    `json:"headers,omitempty"`
	Templated        bool                   `json:"templated"`
	TransitionsState map[string]string      `json:"transitionsState,omitempty"`
	RemovesState     []string               `json:"removesState,omitempty"`
	FixedDelay       int                    `json:"fixedDelay,omitempty"`
	LogNormalDelay   *LogNormalDelayOptions `json:"logNormalDelay,omitempty"`
}

//Gets Status - required for interfaces.Response
func (this ResponseDetailsViewV5) GetStatus() int { return this.Status }

// Gets Body - required for interfaces.Response
func (this ResponseDetailsViewV5) GetBody() string { return this.Body }

// Gets BodyFile - required for interfaces.Response
func (this ResponseDetailsViewV5) GetBodyFile() string { return this.BodyFile }

// Gets EncodedBody - required for interfaces.Response
func (this ResponseDetailsViewV5) GetEncodedBody() bool { return this.EncodedBody }

func (this ResponseDetailsViewV5) GetTemplated() bool { return this.Templated }

func (this ResponseDetailsViewV5) GetTransitionsState() map[string]string {
	return this.TransitionsState
}

func (this ResponseDetailsViewV5) GetRemovesState() []string { return this.RemovesState }

// Gets Headers - required for interfaces.Response
func (this ResponseDetailsViewV5) GetHeaders() map[string][]string { return this.Headers }

// Gets FixedDelay - required for interfaces.Response
func (this ResponseDetailsViewV5) GetFixedDelay() int { return this.FixedDelay }

// Gets LogNormalDelay - required for interfaces.Response
// The trick here to return nil with the right type to compare later.
func (this ResponseDetailsViewV5) GetLogNormalDelay() interfaces.ResponseDelay {
	if this.LogNormalDelay != nil {
		return this.LogNormalDelay
	}

	return nil
}

type LogNormalDelayOptions struct {
	Min    int `json:"min"`
	Max    int `json:"max"`
	Mean   int `json:"mean"`
	Median int `json:"median"`
}

func (l *LogNormalDelayOptions) GetMin() int {
	return l.Min
}

func (l *LogNormalDelayOptions) GetMax() int {
	return l.Max
}

func (l *LogNormalDelayOptions) GetMean() int {
	return l.Mean
}

func (l *LogNormalDelayOptions) GetMedian() int {
	return l.Median
}
