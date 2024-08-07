package models

import (
	"encoding/json"
	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	"github.com/SpectoLabs/hoverfly/core/util"
)

type RequestFieldMatchers struct {
	Matcher string
	Value   interface{}
	Config  map[string]interface{}
	DoMatch *RequestFieldMatchers
}

func NewRequestFieldMatchersFromView(matchers []v2.MatcherViewV5) []RequestFieldMatchers {
	if matchers == nil {
		return nil
	}
	convertedMatchers := []RequestFieldMatchers{}
	for _, matcher := range matchers {
		doMatch := getDoMatchRequestFromMatcherView(matcher.DoMatch)
		value := getValueFromMatcherView(&matcher)
		convertedMatchers = append(convertedMatchers, RequestFieldMatchers{
			Matcher: matcher.Matcher,
			Value:   value,
			Config:  matcher.Config,
			DoMatch: doMatch,
		})
	}
	return convertedMatchers
}

func getValueFromMatcherView(matcher *v2.MatcherViewV5) interface{} {

	if matcher.Matcher == "form" {
		formFieldsMap, ok := matcher.Value.(map[string]interface{})
		if !ok {
			//return default value incase of any issue
			return matcher.Value
		}
		returnValue := make(map[string][]RequestFieldMatchers)
		for formField, formMatchers := range formFieldsMap {
			marshalledFormMatcherValue, _ := json.Marshal(formMatchers)
			var matchers []RequestFieldMatchers
			err := json.Unmarshal(marshalledFormMatcherValue, &matchers)
			if err != nil {
				//return default value incase of any issue
				return matcher.Value
			}
			returnValue[formField] = matchers
		}
		return returnValue
	} else {
		return matcher.Value
	}
}

func getDoMatchRequestFromMatcherView(matcher *v2.MatcherViewV5) *RequestFieldMatchers {

	if matcher == nil {
		return nil
	}
	matcherValue := *matcher
	return &RequestFieldMatchers{
		Matcher: matcherValue.Matcher,
		Value:   matcherValue.Value,
		Config:  matcherValue.Config,
		DoMatch: getDoMatchRequestFromMatcherView(matcherValue.DoMatch),
	}

}

func NewRequestFieldMatchersFromMapView(mapMatchers map[string][]v2.MatcherViewV5) map[string][]RequestFieldMatchers {
	var matchers map[string][]RequestFieldMatchers
	for key, view := range mapMatchers {
		if matchers == nil {
			matchers = map[string][]RequestFieldMatchers{}
		}
		matchers[key] = NewRequestFieldMatchersFromView(view)
	}
	return matchers
}

func NewQueryRequestFieldMatchersFromMapView(mapMatchers *v2.QueryMatcherViewV5) *QueryRequestFieldMatchers {
	var matchers *QueryRequestFieldMatchers
	if mapMatchers != nil {
		matchers = &QueryRequestFieldMatchers{}
		for key, view := range *mapMatchers {
			if matchers == nil {
				matchers = &QueryRequestFieldMatchers{}
			}
			matchers.Add(key, NewRequestFieldMatchersFromView(view))
		}
	}

	return matchers
}

func (this RequestFieldMatchers) BuildView() v2.MatcherViewV5 {
	doMatch := getViewFromRequestFieldMatcher(this.DoMatch)
	value := getValueFromRequestFieldMatcher(&this)
	return v2.MatcherViewV5{
		Matcher: this.Matcher,
		Value:   value,
		Config:  this.Config,
		DoMatch: doMatch,
	}
}

func getValueFromRequestFieldMatcher(matcher *RequestFieldMatchers) interface{} {

	if matcher.Matcher == "form" {
		formFieldMatchers := matcher.Value.(map[string][]RequestFieldMatchers)
		returnValue := make(map[string][]v2.MatcherViewV5)
		for formField, matchers := range formFieldMatchers {
			var matchersView []v2.MatcherViewV5
			for _, matcher := range matchers {
				matchersView = append(matchersView, matcher.BuildView())
			}
			returnValue[formField] = matchersView
		}
		return returnValue
	} else {
		return matcher.Value
	}
}

func getViewFromRequestFieldMatcher(matcher *RequestFieldMatchers) *v2.MatcherViewV5 {

	if matcher == nil {
		return nil
	}
	matcherValue := *matcher
	return &v2.MatcherViewV5{
		Matcher: matcherValue.Matcher,
		Value:   matcherValue.Value,
		Config:  matcherValue.Config,
		DoMatch: getViewFromRequestFieldMatcher(matcherValue.DoMatch),
	}
}

type RequestMatcherResponsePair struct {
	Labels         []string
	RequestMatcher RequestMatcher
	Response       ResponseDetails
}

func NewRequestMatcherResponsePairFromView(view *v2.RequestMatcherResponsePairViewV5) *RequestMatcherResponsePair {

	return &RequestMatcherResponsePair{
		Labels: view.Labels,
		RequestMatcher: RequestMatcher{
			Path:          NewRequestFieldMatchersFromView(view.RequestMatcher.Path),
			Method:        NewRequestFieldMatchersFromView(view.RequestMatcher.Method),
			Destination:   NewRequestFieldMatchersFromView(view.RequestMatcher.Destination),
			Scheme:        NewRequestFieldMatchersFromView(view.RequestMatcher.Scheme),
			Body:          NewRequestFieldMatchersFromView(view.RequestMatcher.Body),
			Headers:       NewRequestFieldMatchersFromMapView(view.RequestMatcher.Headers),
			Query:         NewQueryRequestFieldMatchersFromMapView(view.RequestMatcher.Query),
			RequiresState: view.RequestMatcher.RequiresState,
		},
		Response: NewResponseDetailsFromResponse(view.Response),
	}
}

func (this *RequestMatcherResponsePair) BuildView() v2.RequestMatcherResponsePairViewV5 {

	var path, method, destination, scheme, body []v2.MatcherViewV5

	if this.RequestMatcher.Path != nil && len(this.RequestMatcher.Path) != 0 {
		views := []v2.MatcherViewV5{}
		for _, matcher := range this.RequestMatcher.Path {
			views = append(views, matcher.BuildView())
		}
		path = views
	}

	if this.RequestMatcher.Method != nil && len(this.RequestMatcher.Method) != 0 {
		views := []v2.MatcherViewV5{}
		for _, matcher := range this.RequestMatcher.Method {
			views = append(views, matcher.BuildView())
		}
		method = views
	}

	if this.RequestMatcher.Destination != nil && len(this.RequestMatcher.Destination) != 0 {
		views := []v2.MatcherViewV5{}
		for _, matcher := range this.RequestMatcher.Destination {
			views = append(views, matcher.BuildView())
		}
		destination = views
	}

	if this.RequestMatcher.Scheme != nil && len(this.RequestMatcher.Scheme) != 0 {
		views := []v2.MatcherViewV5{}
		for _, matcher := range this.RequestMatcher.Scheme {
			views = append(views, matcher.BuildView())
		}
		scheme = views
	}

	if this.RequestMatcher.Body != nil && len(this.RequestMatcher.Body) != 0 {
		views := []v2.MatcherViewV5{}
		for _, matcher := range this.RequestMatcher.Body {
			views = append(views, matcher.BuildView())
		}
		body = views
	}

	headersWithMatchers := map[string][]v2.MatcherViewV5{}
	for key, matchers := range this.RequestMatcher.Headers {
		views := []v2.MatcherViewV5{}
		for _, matcher := range matchers {
			views = append(views, matcher.BuildView())
		}
		headersWithMatchers[key] = views
	}

	var queriesWithMatchers *v2.QueryMatcherViewV5
	if this.RequestMatcher.Query != nil {
		queriesWithMatchers = &v2.QueryMatcherViewV5{}
		for key, matchers := range *this.RequestMatcher.Query {
			views := []v2.MatcherViewV5{}
			for _, matcher := range matchers {
				views = append(views, matcher.BuildView())
			}
			(*queriesWithMatchers)[key] = views
		}
	}

	return v2.RequestMatcherResponsePairViewV5{
		Labels: this.Labels,
		RequestMatcher: v2.RequestMatcherViewV5{
			Path:          path,
			Method:        method,
			Destination:   destination,
			Scheme:        scheme,
			Body:          body,
			Headers:       headersWithMatchers,
			Query:         queriesWithMatchers,
			RequiresState: this.RequestMatcher.RequiresState,
		},
		Response: this.Response.ConvertToResponseDetailsViewV5(),
	}
}

type RequestMatcher struct {
	Path          []RequestFieldMatchers
	Method        []RequestFieldMatchers
	Destination   []RequestFieldMatchers
	Scheme        []RequestFieldMatchers
	Body          []RequestFieldMatchers
	Headers       map[string][]RequestFieldMatchers
	Query         *QueryRequestFieldMatchers
	RequiresState map[string]string
}

type QueryRequestFieldMatchers map[string][]RequestFieldMatchers

func (q *QueryRequestFieldMatchers) Add(k string, v []RequestFieldMatchers) {
	(*q)[k] = v
}

func (q *QueryRequestFieldMatchers) Get(k string) []RequestFieldMatchers {
	return (*q)[k]
}

func (this RequestMatcher) IncludesHeaderMatching() bool {
	return this.Headers != nil && len(this.Headers) > 0
}

func (this RequestMatcher) IncludesStateMatching() bool {
	return this.RequiresState != nil && len(this.RequiresState) > 0
}

func (this RequestMatcher) ToEagerlyCacheable() *RequestDetails {
	if this.Body == nil || len(this.Body) != 1 || this.Body[0].Matcher != matchers.Exact ||
		this.Destination == nil || len(this.Destination) != 1 || this.Destination[0].Matcher != matchers.Exact ||
		this.Method == nil || len(this.Method) != 1 || this.Method[0].Matcher != matchers.Exact ||
		this.Path == nil || len(this.Path) != 1 || this.Path[0].Matcher != matchers.Exact ||
		this.Scheme == nil || len(this.Scheme) != 1 || this.Scheme[0].Matcher != matchers.Exact {
		return nil
	}

	if this.IncludesHeaderMatching() {
		return nil
	}

	if this.IncludesStateMatching() {
		return nil
	}

	query := make(map[string][]string)
	if this.Query != nil && len(*this.Query) > 0 {
		for key, valueMatchers := range *this.Query {
			for _, valueMatcher := range valueMatchers {
				if valueMatcher.Matcher != matchers.Exact && !(valueMatcher.Matcher == matchers.Array && (valueMatcher.Config == nil || len(valueMatcher.Config) == 0)) {
					return nil
				}
				if valueMatcher.Matcher == matchers.Array {
					if value, ok := util.GetStringArray(valueMatcher.Value); ok {
						query[key] = value
					} else {
						//ll hardly the case
						query[key] = []string{}
					}

				} else {
					query[key] = []string{valueMatcher.Value.(string)}
				}

			}
		}
	}

	return &RequestDetails{
		Body:        this.Body[0].Value.(string),
		Destination: this.Destination[0].Value.(string),
		Method:      this.Method[0].Value.(string),
		Path:        this.Path[0].Value.(string),
		Query:       query,
		Scheme:      this.Scheme[0].Value.(string),
	}
}

type MatchError struct {
	ClosestMiss *ClosestMiss
	error       string
}

func NewMatchErrorWithClosestMiss(closestMiss *ClosestMiss, error string) *MatchError {
	return &MatchError{
		ClosestMiss: closestMiss,
		error:       error,
	}
}

func NewMatchError(error string) *MatchError {
	return &MatchError{
		error: error,
	}
}

func (this MatchError) Error() string {
	return this.error
}
