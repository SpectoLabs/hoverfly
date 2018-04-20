package models

import (
	"net/url"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/util"
)

type RequestFieldMatchers struct {
	Matcher       string
	Value         interface{}
	ExactMatch    *string
	XmlMatch      *string
	XpathMatch    *string
	JsonMatch     *string
	JsonPathMatch *string
	RegexMatch    *string
	GlobMatch     *string
}

func NewRequestFieldMatchersFromView(matchers *v2.RequestFieldMatchersView) *RequestFieldMatchers {
	if matchers == nil {
		return nil
	}

	var matcher string
	var value interface{}

	if matchers.ExactMatch != nil {
		matcher = "exact"
		value = *matchers.ExactMatch
	}

	if matchers.XmlMatch != nil {
		matcher = "xml"
		value = *matchers.XmlMatch
	}

	if matchers.XpathMatch != nil {
		matcher = "xpath"
		value = *matchers.XpathMatch
	}

	if matchers.JsonMatch != nil {
		matcher = "json"
		value = *matchers.JsonMatch
	}

	if matchers.JsonPathMatch != nil {
		matcher = "jsonpath"
		value = *matchers.JsonPathMatch
	}

	if matchers.RegexMatch != nil {
		matcher = "regex"
		value = *matchers.RegexMatch
	}

	if matchers.GlobMatch != nil {
		matcher = "glob"
		value = *matchers.GlobMatch
	}

	return &RequestFieldMatchers{
		Matcher:       matcher,
		Value:         value,
		ExactMatch:    matchers.ExactMatch,
		XmlMatch:      matchers.XmlMatch,
		XpathMatch:    matchers.XpathMatch,
		JsonMatch:     matchers.JsonMatch,
		JsonPathMatch: matchers.JsonPathMatch,
		RegexMatch:    matchers.RegexMatch,
		GlobMatch:     matchers.GlobMatch,
	}
}

func (this RequestFieldMatchers) BuildView() *v2.RequestFieldMatchersView {
	return &v2.RequestFieldMatchersView{
		ExactMatch:    this.ExactMatch,
		XmlMatch:      this.XmlMatch,
		XpathMatch:    this.XpathMatch,
		JsonMatch:     this.JsonMatch,
		JsonPathMatch: this.JsonPathMatch,
		RegexMatch:    this.RegexMatch,
		GlobMatch:     this.GlobMatch,
	}
}

type RequestMatcherResponsePair struct {
	RequestMatcher RequestMatcher
	Response       ResponseDetails
}

func NewRequestMatcherResponsePairFromView(view *v2.RequestMatcherResponsePairViewV4) *RequestMatcherResponsePair {
	if view.RequestMatcher.Query != nil && view.RequestMatcher.Query.ExactMatch != nil {
		sortedQuery := util.SortQueryString(*view.RequestMatcher.Query.ExactMatch)
		view.RequestMatcher.Query.ExactMatch = &sortedQuery
	}

	var headersWithMatchers map[string]*RequestFieldMatchers
	for key, view := range view.RequestMatcher.HeadersWithMatchers {
		if headersWithMatchers == nil {
			headersWithMatchers = map[string]*RequestFieldMatchers{}
		}
		headersWithMatchers[key] = NewRequestFieldMatchersFromView(view)
	}

	var queriesWithMatchers map[string]*RequestFieldMatchers
	for key, view := range view.RequestMatcher.QueriesWithMatchers {
		if queriesWithMatchers == nil {
			queriesWithMatchers = map[string]*RequestFieldMatchers{}
		}
		queriesWithMatchers[key] = NewRequestFieldMatchersFromView(view)
	}

	return &RequestMatcherResponsePair{
		RequestMatcher: RequestMatcher{
			Path:                NewRequestFieldMatchersFromView(view.RequestMatcher.Path),
			Method:              NewRequestFieldMatchersFromView(view.RequestMatcher.Method),
			Destination:         NewRequestFieldMatchersFromView(view.RequestMatcher.Destination),
			Scheme:              NewRequestFieldMatchersFromView(view.RequestMatcher.Scheme),
			Query:               NewRequestFieldMatchersFromView(view.RequestMatcher.Query),
			Body:                NewRequestFieldMatchersFromView(view.RequestMatcher.Body),
			Headers:             view.RequestMatcher.Headers,
			HeadersWithMatchers: headersWithMatchers,
			QueriesWithMatchers: queriesWithMatchers,
			RequiresState:       view.RequestMatcher.RequiresState,
		},
		Response: NewResponseDetailsFromResponse(view.Response),
	}
}

func (this *RequestMatcherResponsePair) BuildView() v2.RequestMatcherResponsePairViewV4 {

	var path, method, destination, scheme, query, body *v2.RequestFieldMatchersView

	if this.RequestMatcher.Path != nil {
		path = this.RequestMatcher.Path.BuildView()
	}

	if this.RequestMatcher.Method != nil {
		method = this.RequestMatcher.Method.BuildView()
	}

	if this.RequestMatcher.Destination != nil {
		destination = this.RequestMatcher.Destination.BuildView()
	}

	if this.RequestMatcher.Scheme != nil {
		scheme = this.RequestMatcher.Scheme.BuildView()
	}

	if this.RequestMatcher.Query != nil {
		query = this.RequestMatcher.Query.BuildView()
	}

	if this.RequestMatcher.Body != nil {
		body = this.RequestMatcher.Body.BuildView()
	}

	headersWithMatchers := map[string]*v2.RequestFieldMatchersView{}
	for key, matcher := range this.RequestMatcher.HeadersWithMatchers {
		headersWithMatchers[key] = matcher.BuildView()
	}

	queriesWithMatchers := map[string]*v2.RequestFieldMatchersView{}
	for key, matcher := range this.RequestMatcher.QueriesWithMatchers {
		queriesWithMatchers[key] = matcher.BuildView()
	}

	return v2.RequestMatcherResponsePairViewV4{
		RequestMatcher: v2.RequestMatcherViewV4{
			Path:                path,
			Method:              method,
			Destination:         destination,
			Scheme:              scheme,
			Query:               query,
			Body:                body,
			Headers:             this.RequestMatcher.Headers,
			HeadersWithMatchers: headersWithMatchers,
			QueriesWithMatchers: queriesWithMatchers,
			RequiresState:       this.RequestMatcher.RequiresState,
		},
		Response: this.Response.ConvertToResponseDetailsViewV4(),
	}
}

type RequestMatcher struct {
	Path                *RequestFieldMatchers
	Method              *RequestFieldMatchers
	Destination         *RequestFieldMatchers
	Scheme              *RequestFieldMatchers
	Query               *RequestFieldMatchers
	Body                *RequestFieldMatchers
	Headers             map[string][]string
	HeadersWithMatchers map[string]*RequestFieldMatchers
	QueriesWithMatchers map[string]*RequestFieldMatchers
	RequiresState       map[string]string
}

func (this RequestMatcher) IncludesHeaderMatching() bool {
	return (this.Headers != nil && len(this.Headers) > 0) || (this.HeadersWithMatchers != nil && len(this.HeadersWithMatchers) > 0)
}

func (this RequestMatcher) IncludesStateMatching() bool {
	return this.RequiresState != nil && len(this.RequiresState) > 0
}

func (this RequestMatcher) ToEagerlyCachable() *RequestDetails {
	if this.Body == nil || this.Body.ExactMatch == nil ||
		this.Destination == nil || this.Destination.ExactMatch == nil ||
		this.Method == nil || this.Method.ExactMatch == nil ||
		this.Path == nil || this.Path.ExactMatch == nil ||
		this.Query == nil || this.Query.ExactMatch == nil ||
		this.Scheme == nil || this.Scheme.ExactMatch == nil {
		return nil
	}

	if this.Headers != nil && len(this.Headers) > 0 {
		return nil
	}

	if this.RequiresState != nil && len(this.RequiresState) > 0 {
		return nil
	}

	query, _ := url.ParseQuery(*this.Query.ExactMatch)

	return &RequestDetails{
		Body:        *this.Body.ExactMatch,
		Destination: *this.Destination.ExactMatch,
		Headers:     this.Headers,
		Method:      *this.Method.ExactMatch,
		Path:        *this.Path.ExactMatch,
		Query:       query,
		Scheme:      *this.Scheme.ExactMatch,
	}
}

type MatchError struct {
	ClosestMiss *ClosestMiss
	error       string
}

func NewMatchErrorWithClosestMiss(closestMiss *ClosestMiss, error string, isCachable bool) *MatchError {
	return &MatchError{
		ClosestMiss: closestMiss,
		error:       error,
	}
}

func NewMatchError(error string, matchedOnAllButHeadersAtLeastOnce bool) *MatchError {
	return &MatchError{
		error: error,
	}
}

func (this MatchError) Error() string {
	return this.error
}
