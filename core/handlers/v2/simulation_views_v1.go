package v2

import (
	"github.com/SpectoLabs/hoverfly/core/interfaces"
	"github.com/SpectoLabs/hoverfly/core/util"
)

type SimulationViewV1 struct {
	DataViewV1 `json:"data"`
	MetaView   `json:"meta"`
}

func (this SimulationViewV1) Upgrade() SimulationViewV3 {
	var pairs []RequestMatcherResponsePairViewV3
	for _, pairV1 := range this.RequestResponsePairViewV1 {

		var schemeMatchers, methodMatchers, destinationMatchers, pathMatchers, queryMatchers, bodyMatchers *RequestFieldMatchersView
		var headers map[string][]string

		isNotRecording := pairV1.Request.RequestType != nil && *pairV1.Request.RequestType != "recording"

		if isNotRecording {
			headers = pairV1.Request.Headers
		}
		if pairV1.Request.Scheme != nil {

			if isNotRecording {
				schemeMatchers = &RequestFieldMatchersView{
					GlobMatch: pairV1.Request.Scheme,
				}
			} else {
				schemeMatchers = &RequestFieldMatchersView{
					ExactMatch: pairV1.Request.Scheme,
				}
			}
		}

		if pairV1.Request.Method != nil {

			if isNotRecording {
				methodMatchers = &RequestFieldMatchersView{
					GlobMatch: pairV1.Request.Method,
				}
			} else {
				methodMatchers = &RequestFieldMatchersView{
					ExactMatch: pairV1.Request.Method,
				}
			}
		}

		if pairV1.Request.Destination != nil {
			if isNotRecording {
				destinationMatchers = &RequestFieldMatchersView{
					GlobMatch: pairV1.Request.Destination,
				}
			} else {
				destinationMatchers = &RequestFieldMatchersView{
					ExactMatch: pairV1.Request.Destination,
				}
			}
		}

		if pairV1.Request.Path != nil {
			if isNotRecording {
				pathMatchers = &RequestFieldMatchersView{
					GlobMatch: pairV1.Request.Path,
				}
			} else {
				pathMatchers = &RequestFieldMatchersView{
					ExactMatch: pairV1.Request.Path,
				}
			}
		}

		if pairV1.Request.Query != nil {
			if isNotRecording {
				queryMatchers = &RequestFieldMatchersView{
					GlobMatch: pairV1.Request.Query,
				}
			} else {
				queryMatchers = &RequestFieldMatchersView{
					ExactMatch: pairV1.Request.Query,
				}
			}
		}

		if pairV1.Request.Body != nil {
			if isNotRecording {
				bodyMatchers = &RequestFieldMatchersView{
					GlobMatch: pairV1.Request.Body,
				}
			} else {
				bodyMatchers = &RequestFieldMatchersView{
					ExactMatch: pairV1.Request.Body,
				}
			}
		}

		pair := RequestMatcherResponsePairViewV3{
			RequestMatcher: RequestMatcherViewV3{
				Scheme:      schemeMatchers,
				Method:      methodMatchers,
				Destination: destinationMatchers,
				Path:        pathMatchers,
				Query:       queryMatchers,
				Body:        bodyMatchers,
				Headers:     headers,
			},
			Response: ResponseDetailsViewV3{
				Body:        pairV1.Response.Body,
				EncodedBody: pairV1.Response.EncodedBody,
				Headers:     pairV1.Response.Headers,
				Status:      pairV1.Response.Status,
				Templated:   false,
			},
		}
		pairs = append(pairs, pair)
	}

	return SimulationViewV3{
		DataViewV3{
			RequestResponsePairs: pairs,
		},
		MetaView{
			SchemaVersion:   "v3",
			HoverflyVersion: this.HoverflyVersion,
			TimeExported:    this.TimeExported,
		},
	}
}

type DataViewV1 struct {
	RequestResponsePairViewV1 []RequestResponsePairViewV1 `json:"pairs"`
	GlobalActions             GlobalActionsView           `json:"globalActions"`
}

type RequestResponsePairViewV1 struct {
	Response ResponseDetailsView `json:"response"`
	Request  RequestDetailsView  `json:"request"`
}

//Gets Response - required for interfaces.RequestResponsePairView
func (this RequestResponsePairViewV1) GetResponse() interfaces.Response { return this.Response }

//Gets RequestMatcher - required for interfaces.RequestResponsePairView
func (this RequestResponsePairViewV1) GetRequest() interfaces.Request { return this.Request }

// ResponseDetailsView is used when marshalling and
// unmarshalling requests. This struct's Body may be Base64
// encoded based on the EncodedBody field.
type ResponseDetailsView struct {
	Status      int                 `json:"status"`
	Body        string              `json:"body"`
	EncodedBody bool                `json:"encodedBody"`
	Headers     map[string][]string `json:"headers,omitempty"`
}

//Gets Status - required for interfaces.Response
func (this ResponseDetailsView) GetStatus() int { return this.Status }

// Gets Body - required for interfaces.Response
func (this ResponseDetailsView) GetBody() string { return this.Body }

// Gets EncodedBody - required for interfaces.Response
func (this ResponseDetailsView) GetEncodedBody() bool { return this.EncodedBody }

func (this ResponseDetailsView) GetTemplated() bool { return false }

// Gets Headers - required for interfaces.Response
func (this ResponseDetailsView) GetHeaders() map[string][]string { return this.Headers }

// RequestDetailsView is used when marshalling and unmarshalling RequestDetails
type RequestDetailsView struct {
	RequestType *string             `json:"requestType,omitempty"`
	Path        *string             `json:"path"`
	Method      *string             `json:"method"`
	Destination *string             `json:"destination"`
	Scheme      *string             `json:"scheme"`
	Query       *string             `json:"query"`
	Body        *string             `json:"body"`
	Headers     map[string][]string `json:"headers"`
}

//Gets Path - required for interfaces.RequestMatcher
func (this RequestDetailsView) GetPath() *string { return this.Path }

//Gets Method - required for interfaces.RequestMatcher
func (this RequestDetailsView) GetMethod() *string { return this.Method }

//Gets Destination - required for interfaces.RequestMatcher
func (this RequestDetailsView) GetDestination() *string { return this.Destination }

//Gets Scheme - required for interfaces.RequestMatcher
func (this RequestDetailsView) GetScheme() *string { return this.Scheme }

//Gets Query - required for interfaces.RequestMatcher
func (this RequestDetailsView) GetQuery() *string {
	if this.Query == nil {
		return this.Query
	}
	queryString := util.SortQueryString(*this.Query)
	return &queryString
}

//Gets Body - required for interfaces.RequestMatcher
func (this RequestDetailsView) GetBody() *string { return this.Body }

//Gets Headers - required for interfaces.RequestMatcher
func (this RequestDetailsView) GetHeaders() map[string][]string { return this.Headers }
