package v2

import (
	"net/url"
)

type SimulationViewV2 struct {
	DataViewV2 `json:"data"`
	MetaView   `json:"meta"`
}

func (this SimulationViewV2) Upgrade() SimulationViewV4 {
	requestReponsePairsV3 := []RequestMatcherResponsePairViewV4{}

	for _, requestResponsePairV2 := range this.DataViewV2.RequestResponsePairs {
		if requestResponsePairV2.RequestMatcher.Query != nil {
			if requestResponsePairV2.RequestMatcher.Query.ExactMatch != nil {
				unescapedQuery, _ := url.QueryUnescape(*requestResponsePairV2.RequestMatcher.Query.ExactMatch)
				requestResponsePairV2.RequestMatcher.Query.ExactMatch = &unescapedQuery
			}
			if requestResponsePairV2.RequestMatcher.Query.GlobMatch != nil {
				unescapedQuery, _ := url.QueryUnescape(*requestResponsePairV2.RequestMatcher.Query.GlobMatch)
				requestResponsePairV2.RequestMatcher.Query.GlobMatch = &unescapedQuery
			}
		}
		requestResponsePairV3 := RequestMatcherResponsePairViewV4{
			RequestMatcher: RequestMatcherViewV4{
				Body:        requestResponsePairV2.RequestMatcher.Body,
				Destination: requestResponsePairV2.RequestMatcher.Destination,
				Headers:     requestResponsePairV2.RequestMatcher.Headers,
				Method:      requestResponsePairV2.RequestMatcher.Method,
				Path:        requestResponsePairV2.RequestMatcher.Path,
				Query:       requestResponsePairV2.RequestMatcher.Query,
				Scheme:      requestResponsePairV2.RequestMatcher.Scheme,
			},
			Response: ResponseDetailsViewV4{
				Body:        requestResponsePairV2.Response.Body,
				EncodedBody: requestResponsePairV2.Response.EncodedBody,
				Headers:     requestResponsePairV2.Response.Headers,
				Status:      requestResponsePairV2.Response.Status,
				Templated:   false,
			},
		}

		requestReponsePairsV3 = append(requestReponsePairsV3, requestResponsePairV3)
	}

	return SimulationViewV4{
		DataViewV4: DataViewV4{
			RequestResponsePairs: requestReponsePairsV3,
			GlobalActions:        this.GlobalActions,
		},
		MetaView: this.MetaView,
	}
}

type DataViewV2 struct {
	RequestResponsePairs []RequestMatcherResponsePairViewV2 `json:"pairs"`
	GlobalActions        GlobalActionsView                  `json:"globalActions"`
}

type RequestMatcherResponsePairViewV2 struct {
	Response       ResponseDetailsView  `json:"response"`
	RequestMatcher RequestMatcherViewV2 `json:"request"`
}

type RequestFieldMatchersView struct {
	ExactMatch    *string `json:"exactMatch,omitempty"`
	XmlMatch      *string `json:"xmlMatch,omitempty"`
	XpathMatch    *string `json:"xpathMatch,omitempty"`
	JsonMatch     *string `json:"jsonMatch,omitempty"`
	JsonPathMatch *string `json:"jsonPathMatch,omitempty"`
	RegexMatch    *string `json:"regexMatch,omitempty"`
	GlobMatch     *string `json:"globMatch,omitempty"`
}

// RequestDetailsView is used when marshalling and unmarshalling RequestDetails
type RequestMatcherViewV2 struct {
	Path        *RequestFieldMatchersView `json:"path,omitempty"`
	Method      *RequestFieldMatchersView `json:"method,omitempty"`
	Destination *RequestFieldMatchersView `json:"destination,omitempty"`
	Scheme      *RequestFieldMatchersView `json:"scheme,omitempty"`
	Query       *RequestFieldMatchersView `json:"query,omitempty"`
	Body        *RequestFieldMatchersView `json:"body,omitempty"`
	Headers     map[string][]string       `json:"headers,omitempty"`
}
