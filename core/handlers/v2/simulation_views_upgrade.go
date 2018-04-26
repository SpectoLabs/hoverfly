package v2

import (
	"fmt"
	"net/url"
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
		var headers map[string][]string

		isNotRecording := pairV1.Request.RequestType != nil && *pairV1.Request.RequestType != "recording"

		if isNotRecording {
			headers = pairV1.Request.Headers
		}
		if pairV1.Request.Scheme != nil {

			if isNotRecording {
				schemeMatchers = append(schemeMatchers, MatcherViewV5{
					Matcher: "glob",
					Value:   *pairV1.Request.Scheme,
				})
			} else {
				schemeMatchers = append(schemeMatchers, MatcherViewV5{
					Matcher: "exact",
					Value:   *pairV1.Request.Scheme,
				})
			}
		}

		if pairV1.Request.Method != nil {
			if isNotRecording {
				methodMatchers = append(methodMatchers, MatcherViewV5{
					Matcher: "glob",
					Value:   *pairV1.Request.Method,
				})
			} else {
				methodMatchers = append(methodMatchers, MatcherViewV5{
					Matcher: "exact",
					Value:   *pairV1.Request.Method,
				})
			}
		}

		if pairV1.Request.Destination != nil {
			if isNotRecording {
				destinationMatchers = append(destinationMatchers, MatcherViewV5{
					Matcher: "glob",
					Value:   *pairV1.Request.Destination,
				})
			} else {
				destinationMatchers = append(destinationMatchers, MatcherViewV5{
					Matcher: "exact",
					Value:   *pairV1.Request.Destination,
				})
			}
		}

		if pairV1.Request.Path != nil {
			if isNotRecording {
				pathMatchers = append(pathMatchers, MatcherViewV5{
					Matcher: "glob",
					Value:   *pairV1.Request.Path,
				})
			} else {
				pathMatchers = append(pathMatchers, MatcherViewV5{
					Matcher: "exact",
					Value:   *pairV1.Request.Path,
				})
			}
		}

		if pairV1.Request.Query != nil {
			query, _ := url.QueryUnescape(*pairV1.Request.Query)
			if isNotRecording {
				queryMatchers = append(queryMatchers, MatcherViewV5{
					Matcher: "glob",
					Value:   query,
				})
			} else {
				queryMatchers = append(queryMatchers, MatcherViewV5{
					Matcher: "exact",
					Value:   query,
				})
			}
		}

		if pairV1.Request.Body != nil {
			if isNotRecording {
				bodyMatchers = append(bodyMatchers, MatcherViewV5{
					Matcher: "glob",
					Value:   *pairV1.Request.Body,
				})
			} else {
				bodyMatchers = append(bodyMatchers, MatcherViewV5{
					Matcher: "exact",
					Value:   *pairV1.Request.Body,
				})
			}
		}

		pair := RequestMatcherResponsePairViewV5{
			RequestMatcher: RequestMatcherViewV5{
				Scheme:        schemeMatchers,
				Method:        methodMatchers,
				Destination:   destinationMatchers,
				Path:          pathMatchers,
				Query:         queryMatchers,
				Body:          bodyMatchers,
				Headers:       headers,
				RequiresState: nil,
			},
			Response: ResponseDetailsViewV5{
				Body:        pairV1.Response.Body,
				EncodedBody: pairV1.Response.EncodedBody,
				Headers:     pairV1.Response.Headers,
				Status:      pairV1.Response.Status,
				Templated:   false,
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
					Matcher: "exact",
					Value:   unescapedQuery,
				})
			}
			if requestResponsePairV2.RequestMatcher.Query.GlobMatch != nil {
				unescapedQuery, _ := url.QueryUnescape(*requestResponsePairV2.RequestMatcher.Query.GlobMatch)
				queryMatchers = append(queryMatchers, MatcherViewV5{
					Matcher: "glob",
					Value:   unescapedQuery,
				})
			}
		}
		requestResponsePair := RequestMatcherResponsePairViewV5{
			RequestMatcher: RequestMatcherViewV5{
				Destination:   destinationMatchers,
				Headers:       requestResponsePairV2.RequestMatcher.Headers,
				Method:        methodMatchers,
				Path:          pathMatchers,
				Query:         queryMatchers,
				Scheme:        schemeMatchers,
				Body:          bodyMatchers,
				RequiresState: nil,
			},
			Response: ResponseDetailsViewV5{
				Body:        requestResponsePairV2.Response.Body,
				EncodedBody: requestResponsePairV2.Response.EncodedBody,
				Headers:     requestResponsePairV2.Response.Headers,
				Status:      requestResponsePairV2.Response.Status,
				Templated:   false,
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

func upgradeV3(originalSimulation SimulationViewV3) SimulationViewV5 {
	requestReponsePairs := []RequestMatcherResponsePairViewV5{}

	for _, requestResponsePairV2 := range originalSimulation.DataViewV3.RequestResponsePairs {
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
					Matcher: "exact",
					Value:   unescapedQuery,
				})
			}
			if requestResponsePairV2.RequestMatcher.Query.GlobMatch != nil {
				unescapedQuery, _ := url.QueryUnescape(*requestResponsePairV2.RequestMatcher.Query.GlobMatch)
				queryMatchers = append(queryMatchers, MatcherViewV5{
					Matcher: "glob",
					Value:   unescapedQuery,
				})
			}
		}
		requestResponsePair := RequestMatcherResponsePairViewV5{
			RequestMatcher: RequestMatcherViewV5{
				Destination:   destinationMatchers,
				Headers:       requestResponsePairV2.RequestMatcher.Headers,
				Method:        methodMatchers,
				Path:          pathMatchers,
				Query:         queryMatchers,
				Scheme:        schemeMatchers,
				Body:          bodyMatchers,
				RequiresState: nil,
			},
			Response: ResponseDetailsViewV5{
				Body:        requestResponsePairV2.Response.Body,
				EncodedBody: requestResponsePairV2.Response.EncodedBody,
				Headers:     requestResponsePairV2.Response.Headers,
				Status:      requestResponsePairV2.Response.Status,
				Templated:   requestResponsePairV2.Response.Templated,
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
					Matcher: "exact",
					Value:   unescapedQuery,
				})
			}
			if requestResponsePairV2.RequestMatcher.Query.GlobMatch != nil {
				unescapedQuery, _ := url.QueryUnescape(*requestResponsePairV2.RequestMatcher.Query.GlobMatch)
				queryMatchers = append(queryMatchers, MatcherViewV5{
					Matcher: "glob",
					Value:   unescapedQuery,
				})
			}
		}
		requestResponsePair := RequestMatcherResponsePairViewV5{
			RequestMatcher: RequestMatcherViewV5{
				Destination:   destinationMatchers,
				Headers:       requestResponsePairV2.RequestMatcher.Headers,
				Method:        methodMatchers,
				Path:          pathMatchers,
				Query:         queryMatchers,
				Scheme:        schemeMatchers,
				Body:          bodyMatchers,
				RequiresState: requestResponsePairV2.RequestMatcher.RequiresState,
			},
			Response: ResponseDetailsViewV5{
				Body:        requestResponsePairV2.Response.Body,
				EncodedBody: requestResponsePairV2.Response.EncodedBody,
				Headers:     requestResponsePairV2.Response.Headers,
				Status:      requestResponsePairV2.Response.Status,
				Templated:   requestResponsePairV2.Response.Templated,
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

func v2GetMatchersFromRequestFieldMatchersView(requestFieldMatchers *RequestFieldMatchersView) []MatcherViewV5 {
	matchers := []MatcherViewV5{}
	if requestFieldMatchers != nil {
		if requestFieldMatchers.ExactMatch != nil {
			matchers = append(matchers, MatcherViewV5{
				Matcher: "exact",
				Value:   *requestFieldMatchers.ExactMatch,
			})
		}
		if requestFieldMatchers.GlobMatch != nil {
			matchers = append(matchers, MatcherViewV5{
				Matcher: "glob",
				Value:   *requestFieldMatchers.GlobMatch,
			})
		}
		if requestFieldMatchers.JsonMatch != nil {
			matchers = append(matchers, MatcherViewV5{
				Matcher: "json",
				Value:   *requestFieldMatchers.JsonMatch,
			})
		}
		if requestFieldMatchers.JsonPathMatch != nil {
			matchers = append(matchers, MatcherViewV5{
				Matcher: "jsonpath",
				Value:   *requestFieldMatchers.JsonPathMatch,
			})
		}
		if requestFieldMatchers.RegexMatch != nil {
			fmt.Println("in regex")
			matchers = append(matchers, MatcherViewV5{
				Matcher: "regex",
				Value:   *requestFieldMatchers.RegexMatch,
			})
		}
		if requestFieldMatchers.XmlMatch != nil {
			matchers = append(matchers, MatcherViewV5{
				Matcher: "xml",
				Value:   *requestFieldMatchers.XmlMatch,
			})
		}
		if requestFieldMatchers.XpathMatch != nil {
			matchers = append(matchers, MatcherViewV5{
				Matcher: "xpath",
				Value:   *requestFieldMatchers.XpathMatch,
			})
		}
	}
	return matchers
}

func newMetaView(originalMeta MetaView) MetaView {
	return MetaView{
		SchemaVersion:   "v5",
		HoverflyVersion: originalMeta.HoverflyVersion,
		TimeExported:    originalMeta.TimeExported,
	}
}
