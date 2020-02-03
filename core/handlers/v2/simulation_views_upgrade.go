package v2

import (
	"net/url"
	"strings"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
)

func upgradeV1(originalSimulation SimulationViewV1) SimulationViewV6 {
	var pairs []RequestMatcherResponsePairViewV6
	for _, pairV1 := range originalSimulation.RequestResponsePairViewV1 {

		schemeMatchers := []MatcherViewV6{}
		methodMatchers := []MatcherViewV6{}
		destinationMatchers := []MatcherViewV6{}
		pathMatchers := []MatcherViewV6{}
		queryMatchers := []MatcherViewV6{}
		bodyMatchers := []MatcherViewV6{}

		isNotRecording := pairV1.Request.RequestType != nil && *pairV1.Request.RequestType != "recording"

		if pairV1.Request.Scheme != nil {

			if isNotRecording {
				schemeMatchers = append(schemeMatchers, MatcherViewV6{
					Matcher: matchers.Glob,
					Value:   *pairV1.Request.Scheme,
				})
			} else {
				schemeMatchers = append(schemeMatchers, MatcherViewV6{
					Matcher: matchers.Exact,
					Value:   *pairV1.Request.Scheme,
				})
			}
		}

		if pairV1.Request.Method != nil {
			if isNotRecording {
				methodMatchers = append(methodMatchers, MatcherViewV6{
					Matcher: matchers.Glob,
					Value:   *pairV1.Request.Method,
				})
			} else {
				methodMatchers = append(methodMatchers, MatcherViewV6{
					Matcher: matchers.Exact,
					Value:   *pairV1.Request.Method,
				})
			}
		}

		if pairV1.Request.Destination != nil {
			if isNotRecording {
				destinationMatchers = append(destinationMatchers, MatcherViewV6{
					Matcher: matchers.Glob,
					Value:   *pairV1.Request.Destination,
				})
			} else {
				destinationMatchers = append(destinationMatchers, MatcherViewV6{
					Matcher: matchers.Exact,
					Value:   *pairV1.Request.Destination,
				})
			}
		}

		if pairV1.Request.Path != nil {
			if isNotRecording {
				pathMatchers = append(pathMatchers, MatcherViewV6{
					Matcher: matchers.Glob,
					Value:   *pairV1.Request.Path,
				})
			} else {
				pathMatchers = append(pathMatchers, MatcherViewV6{
					Matcher: matchers.Exact,
					Value:   *pairV1.Request.Path,
				})
			}
		}

		if pairV1.Request.Query != nil {
			query, _ := url.QueryUnescape(*pairV1.Request.Query)
			if isNotRecording {
				queryMatchers = append(queryMatchers, MatcherViewV6{
					Matcher: matchers.Glob,
					Value:   query,
				})
			} else {
				queryMatchers = append(queryMatchers, MatcherViewV6{
					Matcher: matchers.Exact,
					Value:   query,
				})
			}
		}

		if pairV1.Request.Body != nil {
			if isNotRecording {
				bodyMatchers = append(bodyMatchers, MatcherViewV6{
					Matcher: matchers.Glob,
					Value:   *pairV1.Request.Body,
				})
			} else {
				bodyMatchers = append(bodyMatchers, MatcherViewV6{
					Matcher: matchers.Exact,
					Value:   *pairV1.Request.Body,
				})
			}
		}

		headersWithMatchers := getMatchersFromRequestHeaders(pairV1.Request.Headers)

		pair := RequestMatcherResponsePairViewV6{
			RequestMatcher: RequestMatcherViewV6{
				Scheme:          schemeMatchers,
				Method:          methodMatchers,
				Destination:     destinationMatchers,
				Path:            pathMatchers,
				Body:            bodyMatchers,
				Headers:         headersWithMatchers,
				RequiresState:   nil,
				DeprecatedQuery: queryMatchers,
			},
			Response: ResponseDetailsViewV6{
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

	return SimulationViewV6{
		DataViewV6{
			RequestResponsePairs: pairs,
		},
		newMetaView(originalSimulation.MetaView),
	}
}

func upgradeV2(originalSimulation SimulationViewV2) SimulationViewV6 {
	requestResponsePairs := []RequestMatcherResponsePairViewV6{}

	for _, requestResponsePairV2 := range originalSimulation.DataViewV2.RequestResponsePairs {
		schemeMatchers := []MatcherViewV6{}
		methodMatchers := []MatcherViewV6{}
		destinationMatchers := []MatcherViewV6{}
		pathMatchers := []MatcherViewV6{}
		queryMatchers := []MatcherViewV6{}
		bodyMatchers := []MatcherViewV6{}

		schemeMatchers = v2GetMatchersFromRequestFieldMatchersView(requestResponsePairV2.RequestMatcher.Scheme)
		methodMatchers = v2GetMatchersFromRequestFieldMatchersView(requestResponsePairV2.RequestMatcher.Method)
		destinationMatchers = v2GetMatchersFromRequestFieldMatchersView(requestResponsePairV2.RequestMatcher.Destination)
		pathMatchers = v2GetMatchersFromRequestFieldMatchersView(requestResponsePairV2.RequestMatcher.Path)
		bodyMatchers = v2GetMatchersFromRequestFieldMatchersView(requestResponsePairV2.RequestMatcher.Body)

		if requestResponsePairV2.RequestMatcher.Query != nil {
			if requestResponsePairV2.RequestMatcher.Query.ExactMatch != nil {
				unescapedQuery, _ := url.QueryUnescape(*requestResponsePairV2.RequestMatcher.Query.ExactMatch)
				queryMatchers = append(queryMatchers, MatcherViewV6{
					Matcher: matchers.Exact,
					Value:   unescapedQuery,
				})
			}
			if requestResponsePairV2.RequestMatcher.Query.GlobMatch != nil {
				unescapedQuery, _ := url.QueryUnescape(*requestResponsePairV2.RequestMatcher.Query.GlobMatch)
				queryMatchers = append(queryMatchers, MatcherViewV6{
					Matcher: matchers.Glob,
					Value:   unescapedQuery,
				})
			}
		}

		headersWithMatchers := getMatchersFromRequestHeaders(requestResponsePairV2.RequestMatcher.Headers)

		requestResponsePair := RequestMatcherResponsePairViewV6{
			RequestMatcher: RequestMatcherViewV6{
				Destination:     destinationMatchers,
				Headers:         headersWithMatchers,
				Method:          methodMatchers,
				Path:            pathMatchers,
				Scheme:          schemeMatchers,
				Body:            bodyMatchers,
				RequiresState:   nil,
				DeprecatedQuery: queryMatchers,
			},
			Response: ResponseDetailsViewV6{
				Body:             requestResponsePairV2.Response.Body,
				EncodedBody:      requestResponsePairV2.Response.EncodedBody,
				Headers:          requestResponsePairV2.Response.Headers,
				Status:           requestResponsePairV2.Response.Status,
				Templated:        false,
				TransitionsState: nil,
				RemovesState:     nil,
			},
		}

		requestResponsePairs = append(requestResponsePairs, requestResponsePair)
	}

	return SimulationViewV6{
		DataViewV6{
			RequestResponsePairs: requestResponsePairs,
			GlobalActions:        originalSimulation.GlobalActions,
		},
		newMetaView(originalSimulation.MetaView),
	}
}

func upgradeV4(originalSimulation SimulationViewV4) SimulationViewV6 {
	requestResponsePairs := []RequestMatcherResponsePairViewV6{}

	for _, requestResponsePairV2 := range originalSimulation.DataViewV4.RequestResponsePairs {
		schemeMatchers := []MatcherViewV6{}
		methodMatchers := []MatcherViewV6{}
		destinationMatchers := []MatcherViewV6{}
		pathMatchers := []MatcherViewV6{}
		queryMatchers := []MatcherViewV6{}
		bodyMatchers := []MatcherViewV6{}

		schemeMatchers = v2GetMatchersFromRequestFieldMatchersView(requestResponsePairV2.RequestMatcher.Scheme)
		methodMatchers = v2GetMatchersFromRequestFieldMatchersView(requestResponsePairV2.RequestMatcher.Method)
		destinationMatchers = v2GetMatchersFromRequestFieldMatchersView(requestResponsePairV2.RequestMatcher.Destination)
		pathMatchers = v2GetMatchersFromRequestFieldMatchersView(requestResponsePairV2.RequestMatcher.Path)
		bodyMatchers = v2GetMatchersFromRequestFieldMatchersView(requestResponsePairV2.RequestMatcher.Body)

		if requestResponsePairV2.RequestMatcher.Query != nil {
			if requestResponsePairV2.RequestMatcher.Query.ExactMatch != nil {
				unescapedQuery, _ := url.QueryUnescape(*requestResponsePairV2.RequestMatcher.Query.ExactMatch)
				queryMatchers = append(queryMatchers, MatcherViewV6{
					Matcher: matchers.Exact,
					Value:   unescapedQuery,
				})
			}
			if requestResponsePairV2.RequestMatcher.Query.GlobMatch != nil {
				unescapedQuery, _ := url.QueryUnescape(*requestResponsePairV2.RequestMatcher.Query.GlobMatch)
				queryMatchers = append(queryMatchers, MatcherViewV6{
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

		var queriesWithMatchers *QueryMatcherViewV6
		if requestResponsePairV2.RequestMatcher.QueriesWithMatchers != nil {
			queriesWithMatchers = &QueryMatcherViewV6{}
			for key, value := range *requestResponsePairV2.RequestMatcher.QueriesWithMatchers {
				(*queriesWithMatchers)[key] = v2GetMatchersFromRequestFieldMatchersView(value)
			}
		}

		requestResponsePair := RequestMatcherResponsePairViewV6{
			RequestMatcher: RequestMatcherViewV6{
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
			Response: ResponseDetailsViewV6{
				Body:             requestResponsePairV2.Response.Body,
				EncodedBody:      requestResponsePairV2.Response.EncodedBody,
				Headers:          requestResponsePairV2.Response.Headers,
				Status:           requestResponsePairV2.Response.Status,
				Templated:        requestResponsePairV2.Response.Templated,
				TransitionsState: requestResponsePairV2.Response.TransitionsState,
				RemovesState:     requestResponsePairV2.Response.RemovesState,
			},
		}

		requestResponsePairs = append(requestResponsePairs, requestResponsePair)
	}

	return SimulationViewV6{
		DataViewV6{
			RequestResponsePairs: requestResponsePairs,
			GlobalActions:        originalSimulation.GlobalActions,
		},
		newMetaView(originalSimulation.MetaView),
	}
}

func upgradeV5(originalSimulation SimulationViewV5) SimulationViewV6 {
	requestResponsePairs := []RequestMatcherResponsePairViewV6{}
	makeV6Matcher := func(v5Matchers []MatcherViewV5) []MatcherViewV6 {
		matcherViews := []MatcherViewV6{}

		for _, v5Matcher := range v5Matchers {
			matcherViews = append(matcherViews, MatcherViewV6{
				Matcher: v5Matcher.Matcher,
				Value:   v5Matcher.Value,
				Config:  v5Matcher.Config,
			})
		}

		return matcherViews
	}

	for _, requestResponsePairV5 := range originalSimulation.DataViewV5.RequestResponsePairs {
		schemeMatchers := []MatcherViewV6{}
		methodMatchers := []MatcherViewV6{}
		destinationMatchers := []MatcherViewV6{}
		pathMatchers := []MatcherViewV6{}
		bodyMatchers := []MatcherViewV6{}
		deprecatedQueryMatchers := []MatcherViewV6{}
		headersMatchers := map[string][]MatcherViewV6{}

		for v5HeaderKey, v5HeaderMatchers := range requestResponsePairV5.RequestMatcher.Headers {
			headersMatchers[v5HeaderKey] = makeV6Matcher(v5HeaderMatchers)
		}

		var queriesMatchers *QueryMatcherViewV6
		if requestResponsePairV5.RequestMatcher.Query != nil {
			queriesMatchers = &QueryMatcherViewV6{}
			for key, value := range *requestResponsePairV5.RequestMatcher.Query {
				(*queriesMatchers)[key] = makeV6Matcher(value)
			}
		}

		schemeMatchers = makeV6Matcher(requestResponsePairV5.RequestMatcher.Scheme)
		methodMatchers = makeV6Matcher(requestResponsePairV5.RequestMatcher.Method)
		destinationMatchers = makeV6Matcher(requestResponsePairV5.RequestMatcher.Destination)
		pathMatchers = makeV6Matcher(requestResponsePairV5.RequestMatcher.Path)
		bodyMatchers = makeV6Matcher(requestResponsePairV5.RequestMatcher.Body)
		deprecatedQueryMatchers = makeV6Matcher(requestResponsePairV5.RequestMatcher.DeprecatedQuery)

		requestResponsePair := RequestMatcherResponsePairViewV6{
			RequestMatcher: RequestMatcherViewV6{
				Destination:     destinationMatchers,
				Method:          methodMatchers,
				Path:            pathMatchers,
				Scheme:          schemeMatchers,
				Body:            bodyMatchers,
				Headers:         headersMatchers,
				Query:           queriesMatchers,
				RequiresState:   requestResponsePairV5.RequestMatcher.RequiresState,
				DeprecatedQuery: deprecatedQueryMatchers,
			},
			Response: ResponseDetailsViewV6{
				Body:             requestResponsePairV5.Response.Body,
				EncodedBody:      requestResponsePairV5.Response.EncodedBody,
				Headers:          requestResponsePairV5.Response.Headers,
				Status:           requestResponsePairV5.Response.Status,
				Templated:        requestResponsePairV5.Response.Templated,
				TransitionsState: requestResponsePairV5.Response.TransitionsState,
				RemovesState:     requestResponsePairV5.Response.RemovesState,
			},
		}

		requestResponsePairs = append(requestResponsePairs, requestResponsePair)
	}

	return SimulationViewV6{
		DataViewV6{
			RequestResponsePairs: requestResponsePairs,
			GlobalActions:        originalSimulation.GlobalActions,
		},
		newMetaView(originalSimulation.MetaView),
	}
}

func getMatchersFromRequestHeaders(headers map[string][]string) map[string][]MatcherViewV6 {
	requestHeaders := map[string][]MatcherViewV6{}
	for headerKey, headerValues := range headers {
		requestHeaders[headerKey] = []MatcherViewV6{
			{
				Matcher: matchers.Glob,
				Value:   strings.Join(headerValues, ";"),
			},
		}
	}
	return requestHeaders
}

func v2GetMatchersFromRequestFieldMatchersView(requestFieldMatchers *RequestFieldMatchersView) []MatcherViewV6 {
	matcherViews := []MatcherViewV6{}
	if requestFieldMatchers != nil {
		if requestFieldMatchers.ExactMatch != nil {
			matcherViews = append(matcherViews, MatcherViewV6{
				Matcher: matchers.Exact,
				Value:   *requestFieldMatchers.ExactMatch,
			})
		}
		if requestFieldMatchers.GlobMatch != nil {
			matcherViews = append(matcherViews, MatcherViewV6{
				Matcher: matchers.Glob,
				Value:   *requestFieldMatchers.GlobMatch,
			})
		}
		if requestFieldMatchers.JsonMatch != nil {
			matcherViews = append(matcherViews, MatcherViewV6{
				Matcher: matchers.Json,
				Value:   *requestFieldMatchers.JsonMatch,
			})
		}
		if requestFieldMatchers.JsonPathMatch != nil {
			matcherViews = append(matcherViews, MatcherViewV6{
				Matcher: matchers.JsonPath,
				Value:   *requestFieldMatchers.JsonPathMatch,
			})
		}
		if requestFieldMatchers.RegexMatch != nil {
			matcherViews = append(matcherViews, MatcherViewV6{
				Matcher: matchers.Regex,
				Value:   *requestFieldMatchers.RegexMatch,
			})
		}
		if requestFieldMatchers.XmlMatch != nil {
			matcherViews = append(matcherViews, MatcherViewV6{
				Matcher: matchers.Xml,
				Value:   *requestFieldMatchers.XmlMatch,
			})
		}
		if requestFieldMatchers.XpathMatch != nil {
			matcherViews = append(matcherViews, MatcherViewV6{
				Matcher: matchers.Xpath,
				Value:   *requestFieldMatchers.XpathMatch,
			})
		}
	}
	return matcherViews
}

func newMetaView(originalMeta MetaView) MetaView {
	return MetaView{
		SchemaVersion:   "v6",
		HoverflyVersion: originalMeta.HoverflyVersion,
		TimeExported:    originalMeta.TimeExported,
	}
}
