package v2

import (
	"net/url"
	"strings"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
)

func upgradeV1(originalSimulation SimulationViewV1) SimulationViewV5 {
	var pairs []RequestMatcherResponsePairViewV5
	for _, pairV1 := range originalSimulation.RequestResponsePairViewV1 {

		schemeMatchers := []MatcherViewV5{}
		methodMatchers := []MatcherViewV5{}
		destinationMatchers := []MatcherViewV5{}
		pathMatchers := []MatcherViewV5{}
		queryMatchers := []MatcherViewV5{}
		bodyMatchers := []MatcherViewV5{}

		isNotRecording := pairV1.Request.RequestType != nil && *pairV1.Request.RequestType != "recording"

		if pairV1.Request.Scheme != nil {

			if isNotRecording {
				schemeMatchers = append(schemeMatchers, MatcherViewV5{
					Matcher: matchers.Glob,
					Value:   *pairV1.Request.Scheme,
				})
			} else {
				schemeMatchers = append(schemeMatchers, MatcherViewV5{
					Matcher: matchers.Exact,
					Value:   *pairV1.Request.Scheme,
				})
			}
		}

		if pairV1.Request.Method != nil {
			if isNotRecording {
				methodMatchers = append(methodMatchers, MatcherViewV5{
					Matcher: matchers.Glob,
					Value:   *pairV1.Request.Method,
				})
			} else {
				methodMatchers = append(methodMatchers, MatcherViewV5{
					Matcher: matchers.Exact,
					Value:   *pairV1.Request.Method,
				})
			}
		}

		if pairV1.Request.Destination != nil {
			if isNotRecording {
				destinationMatchers = append(destinationMatchers, MatcherViewV5{
					Matcher: matchers.Glob,
					Value:   *pairV1.Request.Destination,
				})
			} else {
				destinationMatchers = append(destinationMatchers, MatcherViewV5{
					Matcher: matchers.Exact,
					Value:   *pairV1.Request.Destination,
				})
			}
		}

		if pairV1.Request.Path != nil {
			if isNotRecording {
				pathMatchers = append(pathMatchers, MatcherViewV5{
					Matcher: matchers.Glob,
					Value:   *pairV1.Request.Path,
				})
			} else {
				pathMatchers = append(pathMatchers, MatcherViewV5{
					Matcher: matchers.Exact,
					Value:   *pairV1.Request.Path,
				})
			}
		}

		if pairV1.Request.Query != nil {
			query, _ := url.QueryUnescape(*pairV1.Request.Query)
			if isNotRecording {
				queryMatchers = append(queryMatchers, MatcherViewV5{
					Matcher: matchers.Glob,
					Value:   query,
				})
			} else {
				queryMatchers = append(queryMatchers, MatcherViewV5{
					Matcher: matchers.Exact,
					Value:   query,
				})
			}
		}

		if pairV1.Request.Body != nil {
			if isNotRecording {
				bodyMatchers = append(bodyMatchers, MatcherViewV5{
					Matcher: matchers.Glob,
					Value:   *pairV1.Request.Body,
				})
			} else {
				bodyMatchers = append(bodyMatchers, MatcherViewV5{
					Matcher: matchers.Exact,
					Value:   *pairV1.Request.Body,
				})
			}
		}

		headersWithMatchers := getMatchersFromRequestHeaders(pairV1.Request.Headers)

		pair := RequestMatcherResponsePairViewV5{
			RequestMatcher: RequestMatcherViewV5{
				Scheme:          schemeMatchers,
				Method:          methodMatchers,
				Destination:     destinationMatchers,
				Path:            pathMatchers,
				Body:            bodyMatchers,
				Headers:         headersWithMatchers,
				RequiresState:   nil,
				DeprecatedQuery: queryMatchers,
			},
			Response: ResponseDetailsViewV5{
				Body:             pairV1.Response.Body,
				EncodedBody:      pairV1.Response.EncodedBody,
				Headers:          pairV1.Response.Headers,
				Status:           pairV1.Response.Status,
				Templated:        false,
				TransitionsState: nil,
				RemovesState:     nil,
			},
		}
		pairs = append(pairs, pair)
	}

	return SimulationViewV5{
		DataViewV5{
			RequestResponsePairs: pairs,
		},
		newMetaView(originalSimulation.MetaView),
	}
}

func upgradeV2(originalSimulation SimulationViewV2) SimulationViewV5 {
	requestReponsePairs := []RequestMatcherResponsePairViewV5{}

	for _, requestResponsePairV2 := range originalSimulation.DataViewV2.RequestResponsePairs {
		schemeMatchers := []MatcherViewV5{}
		methodMatchers := []MatcherViewV5{}
		destinationMatchers := []MatcherViewV5{}
		pathMatchers := []MatcherViewV5{}
		queryMatchers := []MatcherViewV5{}
		bodyMatchers := []MatcherViewV5{}

		schemeMatchers = v2GetMatchersFromRequestFieldMatchersView(requestResponsePairV2.RequestMatcher.Scheme)
		methodMatchers = v2GetMatchersFromRequestFieldMatchersView(requestResponsePairV2.RequestMatcher.Method)
		destinationMatchers = v2GetMatchersFromRequestFieldMatchersView(requestResponsePairV2.RequestMatcher.Destination)
		pathMatchers = v2GetMatchersFromRequestFieldMatchersView(requestResponsePairV2.RequestMatcher.Path)
		bodyMatchers = v2GetMatchersFromRequestFieldMatchersView(requestResponsePairV2.RequestMatcher.Body)

		if requestResponsePairV2.RequestMatcher.Query != nil {
			if requestResponsePairV2.RequestMatcher.Query.ExactMatch != nil {
				unescapedQuery, _ := url.QueryUnescape(*requestResponsePairV2.RequestMatcher.Query.ExactMatch)
				queryMatchers = append(queryMatchers, MatcherViewV5{
					Matcher: matchers.Exact,
					Value:   unescapedQuery,
				})
			}
			if requestResponsePairV2.RequestMatcher.Query.GlobMatch != nil {
				unescapedQuery, _ := url.QueryUnescape(*requestResponsePairV2.RequestMatcher.Query.GlobMatch)
				queryMatchers = append(queryMatchers, MatcherViewV5{
					Matcher: matchers.Glob,
					Value:   unescapedQuery,
				})
			}
		}

		headersWithMatchers := getMatchersFromRequestHeaders(requestResponsePairV2.RequestMatcher.Headers)

		requestResponsePair := RequestMatcherResponsePairViewV5{
			RequestMatcher: RequestMatcherViewV5{
				Destination:     destinationMatchers,
				Headers:         headersWithMatchers,
				Method:          methodMatchers,
				Path:            pathMatchers,
				Scheme:          schemeMatchers,
				Body:            bodyMatchers,
				RequiresState:   nil,
				DeprecatedQuery: queryMatchers,
			},
			Response: ResponseDetailsViewV5{
				Body:             requestResponsePairV2.Response.Body,
				EncodedBody:      requestResponsePairV2.Response.EncodedBody,
				Headers:          requestResponsePairV2.Response.Headers,
				Status:           requestResponsePairV2.Response.Status,
				Templated:        false,
				TransitionsState: nil,
				RemovesState:     nil,
			},
		}

		requestReponsePairs = append(requestReponsePairs, requestResponsePair)
	}

	return SimulationViewV5{
		DataViewV5{
			RequestResponsePairs: requestReponsePairs,
			GlobalActions:        originalSimulation.GlobalActions,
		},
		newMetaView(originalSimulation.MetaView),
	}
}

func upgradeV4(originalSimulation SimulationViewV4) SimulationViewV5 {
	requestReponsePairs := []RequestMatcherResponsePairViewV5{}

	for _, requestResponsePairV2 := range originalSimulation.DataViewV4.RequestResponsePairs {
		schemeMatchers := []MatcherViewV5{}
		methodMatchers := []MatcherViewV5{}
		destinationMatchers := []MatcherViewV5{}
		pathMatchers := []MatcherViewV5{}
		queryMatchers := []MatcherViewV5{}
		bodyMatchers := []MatcherViewV5{}

		schemeMatchers = v2GetMatchersFromRequestFieldMatchersView(requestResponsePairV2.RequestMatcher.Scheme)
		methodMatchers = v2GetMatchersFromRequestFieldMatchersView(requestResponsePairV2.RequestMatcher.Method)
		destinationMatchers = v2GetMatchersFromRequestFieldMatchersView(requestResponsePairV2.RequestMatcher.Destination)
		pathMatchers = v2GetMatchersFromRequestFieldMatchersView(requestResponsePairV2.RequestMatcher.Path)
		bodyMatchers = v2GetMatchersFromRequestFieldMatchersView(requestResponsePairV2.RequestMatcher.Body)

		if requestResponsePairV2.RequestMatcher.Query != nil {
			if requestResponsePairV2.RequestMatcher.Query.ExactMatch != nil {
				unescapedQuery, _ := url.QueryUnescape(*requestResponsePairV2.RequestMatcher.Query.ExactMatch)
				queryMatchers = append(queryMatchers, MatcherViewV5{
					Matcher: matchers.Exact,
					Value:   unescapedQuery,
				})
			}
			if requestResponsePairV2.RequestMatcher.Query.GlobMatch != nil {
				unescapedQuery, _ := url.QueryUnescape(*requestResponsePairV2.RequestMatcher.Query.GlobMatch)
				queryMatchers = append(queryMatchers, MatcherViewV5{
					Matcher: matchers.Glob,
					Value:   unescapedQuery,
				})
			}
		}

		headersWithMatchers := getMatchersFromRequestHeaders(requestResponsePairV2.RequestMatcher.Headers)

		for key, value := range requestResponsePairV2.RequestMatcher.HeadersWithMatchers {
			values := v2GetMatchersFromRequestFieldMatchersView(value)
			if headersWithMatchers[key] == nil {
				headersWithMatchers[key] = values
			} else {
				headersWithMatchers[key] = append(headersWithMatchers[key], values...)
			}
		}

		var queriesWithMatchers *QueryMatcherViewV5
		if requestResponsePairV2.RequestMatcher.QueriesWithMatchers != nil {
			queriesWithMatchers = &QueryMatcherViewV5{}
			for key, value := range *requestResponsePairV2.RequestMatcher.QueriesWithMatchers {
				(*queriesWithMatchers)[key] = v2GetMatchersFromRequestFieldMatchersView(value)
			}
		}

		requestResponsePair := RequestMatcherResponsePairViewV5{
			RequestMatcher: RequestMatcherViewV5{
				Destination:     destinationMatchers,
				Method:          methodMatchers,
				Path:            pathMatchers,
				Scheme:          schemeMatchers,
				Body:            bodyMatchers,
				Headers:         headersWithMatchers,
				Query:           queriesWithMatchers,
				RequiresState:   requestResponsePairV2.RequestMatcher.RequiresState,
				DeprecatedQuery: queryMatchers,
			},
			Response: ResponseDetailsViewV5{
				Body:             requestResponsePairV2.Response.Body,
				EncodedBody:      requestResponsePairV2.Response.EncodedBody,
				Headers:          requestResponsePairV2.Response.Headers,
				Status:           requestResponsePairV2.Response.Status,
				Templated:        requestResponsePairV2.Response.Templated,
				TransitionsState: requestResponsePairV2.Response.TransitionsState,
				RemovesState:     requestResponsePairV2.Response.RemovesState,
			},
		}

		requestReponsePairs = append(requestReponsePairs, requestResponsePair)
	}

	return SimulationViewV5{
		DataViewV5{
			RequestResponsePairs: requestReponsePairs,
			GlobalActions:        originalSimulation.GlobalActions,
		},
		newMetaView(originalSimulation.MetaView),
	}
}

func getMatchersFromRequestHeaders(headers map[string][]string) map[string][]MatcherViewV5 {
	requestHeaders := map[string][]MatcherViewV5{}
	for headerKey, headerValues := range headers {
		requestHeaders[headerKey] = []MatcherViewV5{
			{
				Matcher: matchers.Glob,
				Value:   strings.Join(headerValues, ";"),
			},
		}
	}
	return requestHeaders
}

func v2GetMatchersFromRequestFieldMatchersView(requestFieldMatchers *RequestFieldMatchersView) []MatcherViewV5 {
	matcherViews := []MatcherViewV5{}
	if requestFieldMatchers != nil {
		if requestFieldMatchers.ExactMatch != nil {
			matcherViews = append(matcherViews, MatcherViewV5{
				Matcher: matchers.Exact,
				Value:   *requestFieldMatchers.ExactMatch,
			})
		}
		if requestFieldMatchers.GlobMatch != nil {
			matcherViews = append(matcherViews, MatcherViewV5{
				Matcher: matchers.Glob,
				Value:   *requestFieldMatchers.GlobMatch,
			})
		}
		if requestFieldMatchers.JsonMatch != nil {
			matcherViews = append(matcherViews, MatcherViewV5{
				Matcher: matchers.Json,
				Value:   *requestFieldMatchers.JsonMatch,
			})
		}
		if requestFieldMatchers.JsonPathMatch != nil {
			matcherViews = append(matcherViews, MatcherViewV5{
				Matcher: matchers.JsonPath,
				Value:   *requestFieldMatchers.JsonPathMatch,
			})
		}
		if requestFieldMatchers.RegexMatch != nil {
			matcherViews = append(matcherViews, MatcherViewV5{
				Matcher: matchers.Regex,
				Value:   *requestFieldMatchers.RegexMatch,
			})
		}
		if requestFieldMatchers.XmlMatch != nil {
			matcherViews = append(matcherViews, MatcherViewV5{
				Matcher: matchers.Xml,
				Value:   *requestFieldMatchers.XmlMatch,
			})
		}
		if requestFieldMatchers.XpathMatch != nil {
			matcherViews = append(matcherViews, MatcherViewV5{
				Matcher: matchers.Xpath,
				Value:   *requestFieldMatchers.XpathMatch,
			})
		}
	}
	return matcherViews
}

func newMetaView(originalMeta MetaView) MetaView {
	return MetaView{
		SchemaVersion:   "v5",
		HoverflyVersion: originalMeta.HoverflyVersion,
		TimeExported:    originalMeta.TimeExported,
	}
}
