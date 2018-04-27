package v2

type SimulationViewV2 struct {
	DataViewV2 `json:"data"`
	MetaView   `json:"meta"`
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
