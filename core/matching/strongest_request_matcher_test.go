package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
)

func Test_ClosestRequestMatcherRequestMatcher_EmptyRequestMatchersShouldMatchOnAnyRequest(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{},
		Response:       testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "somehost.com",
		Headers: map[string][]string{
			"sdv": {"ascd"},
		},
	}
	result, _, _ := matching.StrongestMatchRequestMatcher(r, false, simulation)

	Expect(result).ToNot(BeNil())
	Expect(result.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_RequestMatchersShouldMatchOnBody(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("body"),
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Body: "body",
	}
	result, _, err := matching.StrongestMatchRequestMatcher(r, false, simulation)
	Expect(err).To(BeNil())

	Expect(result.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_ReturnResponseWhenAllHeadersMatch(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": {"val1"},
		"header2": {"val2"},
	}

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: headers,
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Headers: map[string][]string{
			"header1": {"val1"},
			"header2": {"val2"},
		},
	}

	result, _, _ := matching.StrongestMatchRequestMatcher(r, false, simulation)

	Expect(result.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_ReturnNilWhenOneHeaderNotPresentInRequest(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": {"val1"},
		"header2": {"val2"},
	}

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: headers,
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Headers: map[string][]string{
			"header1": {"val1"},
		},
	}

	result, _, _ := matching.StrongestMatchRequestMatcher(r, false, simulation)

	Expect(result).To(BeNil())
}

func Test_ClosestRequestMatcherRequestMatcher_ReturnNilWhenOneHeaderValueDifferent(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": {"val1"},
		"header2": {"val2"},
	}

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: headers,
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "somehost.com",
		Headers: map[string][]string{
			"header1": {"val1"},
			"header2": {"different"},
		},
	}
	result, _, _ := matching.StrongestMatchRequestMatcher(r, false, simulation)

	Expect(result).To(BeNil())
}

func Test_ClosestRequestMatcherRequestMatcher_ReturnResponseWithMultiValuedHeaderMatch(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": {"val1-a", "val1-b"},
		"header2": {"val2"},
	}

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: headers,
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Body:        "test-body",
		Headers: map[string][]string{
			"header1": {"val1-a", "val1-b"},
			"header2": {"val2"},
		},
	}
	result, _, _ := matching.StrongestMatchRequestMatcher(r, false, simulation)

	Expect(result.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_ReturnNilWithDifferentMultiValuedHeaders(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": {"val1-a", "val1-b"},
		"header2": {"val2"},
	}

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: headers,
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Headers: map[string][]string{
			"header1": {"val1-a", "val1-differnet"},
			"header2": {"val2"},
		},
	}

	result, _, _ := matching.StrongestMatchRequestMatcher(r, false, simulation)

	Expect(result).To(BeNil())
}

func Test_ClosestRequestMatcherRequestMatcher_EndpointMatchWithHeaders(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": {"val1-a", "val1-b"},
		"header2": {"val2"},
	}

	destination := "testhost.com"
	method := "GET"
	path := "/a/1"
	query := "q=test"

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: headers,
			Destination: &models.RequestFieldMatchers{
				ExactMatch: &destination,
			},
			Path: &models.RequestFieldMatchers{
				ExactMatch: &path,
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: &method,
			},
			Query: &models.RequestFieldMatchers{
				ExactMatch: &query,
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/a/1",
		Query:       "q=test",
		Headers: map[string][]string{
			"header1": {"val1-a", "val1-b"},
			"header2": {"val2"},
		},
	}
	result, _, _ := matching.StrongestMatchRequestMatcher(r, false, simulation)

	Expect(result.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_EndpointMismatchWithHeadersReturnsNil(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": {"val1-a", "val1-b"},
		"header2": {"val2"},
	}

	destination := "testhost.com"
	method := "GET"
	path := "/a/1"
	query := "q=test"

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: headers,
			Destination: &models.RequestFieldMatchers{
				ExactMatch: &destination,
			},
			Path: &models.RequestFieldMatchers{
				ExactMatch: &path,
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: &method,
			},
			Query: &models.RequestFieldMatchers{
				ExactMatch: &query,
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://testhost.com",
		Path:        "/a/1",
		Query:       "q=different",
		Headers: map[string][]string{
			"header1": {"val1-a", "val1-b"},
			"header2": {"val2"},
		},
	}

	result, _, _ := matching.StrongestMatchRequestMatcher(r, false, simulation)

	Expect(result).To(BeNil())
}

func Test_ClosestRequestMatcherRequestMatcher_AbleToMatchAnEmptyPathInAReasonableWay(t *testing.T) {
	RegisterTestingT(t)

	destination := "testhost.com"
	method := "GET"
	path := ""
	query := "q=test"
	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: &models.RequestFieldMatchers{
				ExactMatch: &destination,
			},
			Path: &models.RequestFieldMatchers{
				ExactMatch: &path,
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: &method,
			},
			Query: &models.RequestFieldMatchers{
				ExactMatch: &query,
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Query:       "q=test",
	}
	result, _, _ := matching.StrongestMatchRequestMatcher(r, false, simulation)

	Expect(result.Response.Body).To(Equal("request matched"))

	r = models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/a/1",
		Query:       "q=test",
	}

	result, _, _ = matching.StrongestMatchRequestMatcher(r, false, simulation)

	Expect(result).To(BeNil())
}

func Test_ClosestRequestMatcherRequestMatcher_RequestMatchersCanUseGlobsAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: &models.RequestFieldMatchers{
				GlobMatch: StringToPointer("*.com"),
			},
		},
		Response: testResponse,
	})

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
	}

	response, _, err := matching.StrongestMatchRequestMatcher(request, false, simulation)
	Expect(err).To(BeNil())

	Expect(response.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_RequestMatchersCanUseGlobsOnSchemeAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Scheme: &models.RequestFieldMatchers{
				GlobMatch: StringToPointer("H*"),
			},
		},
		Response: testResponse,
	})

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Scheme:      "http",
		Path:        "/api/1",
	}

	response, _, err := matching.StrongestMatchRequestMatcher(request, false, simulation)
	Expect(err).To(BeNil())

	Expect(response.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_RequestMatchersCanUseGlobsOnHeadersAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: map[string][]string{
				"unique-header": {"*"},
			},
		},
		Response: testResponse,
	})

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
		Headers: map[string][]string{
			"unique-header": {"totally-unique"},
		},
	}

	response, _, err := matching.StrongestMatchRequestMatcher(request, false, simulation)
	Expect(err).To(BeNil())

	Expect(response.Response.Body).To(Equal("request matched"))
}

func Test_ShouldReturnClosestMissIfMatchIsNotFound(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("completemiss"),
			},
			Path: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("completemiss"),
			},
		},
		Response: models.ResponseDetails{
			Body: "one",
		},
	})

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("body"),
				GlobMatch:  StringToPointer("bod*"),
			},
			Path: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("path"),
			},
		},
		Response: models.ResponseDetails{
			Body: "two",
		},
	})

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("body"),
			},
			Path: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("path"),
			},
		},
		Response: models.ResponseDetails{
			Body: "three",
		},
	})

	r := models.RequestDetails{
		Body: "body",
		Path: "nomatch",
	}

	result, closestMiss, err := matching.StrongestMatchRequestMatcher(r, false, simulation)

	Expect(err).ToNot(BeNil())
	Expect(result).To(BeNil())
	Expect(closestMiss).ToNot(BeNil())
	Expect(*closestMiss.RequestMatcher.Body.ExactMatch).To(Equal(`body`))
	Expect(*closestMiss.RequestMatcher.Body.GlobMatch).To(Equal(`bod*`))
	Expect(*closestMiss.RequestMatcher.Path.ExactMatch).To(Equal(`path`))
	Expect(closestMiss.Response.Body).To(Equal(`two`))
}

func Test_ShouldReturnClosestMissIfMatchIsNotFoundAgain(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				RegexMatch: StringToPointer(".*"),
			},
			Path: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("miss"),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: models.ResponseDetails{
			Body: "one",
		},
	})

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer(".*"),
				GlobMatch:  StringToPointer("miss"),
			},
			Path: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("miss"),
			},
		},
		Response: models.ResponseDetails{
			Body: "two",
		},
	})

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("miss"),
			},
			Path: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("miss"),
			},
		},
		Response: models.ResponseDetails{
			Body: "three",
		},
	})

	r := models.RequestDetails{
		Body:   "foo",
		Method: "GET",
	}

	result, closestMiss, err := matching.StrongestMatchRequestMatcher(r, false, simulation)

	Expect(err).ToNot(BeNil())
	Expect(result).To(BeNil())
	Expect(closestMiss).ToNot(BeNil())
	Expect(*closestMiss.RequestMatcher.Body.RegexMatch).To(Equal(`.*`))
	Expect(*closestMiss.RequestMatcher.Path.ExactMatch).To(Equal(`miss`))
	Expect(*closestMiss.RequestMatcher.Method.ExactMatch).To(Equal(`GET`))
	Expect(closestMiss.Response.Body).To(Equal(`one`))
}

func Test_ShouldNotReturnClosestMissWhenThereIsAMatch(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				RegexMatch: StringToPointer(".*"),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: models.ResponseDetails{
			Body: "one",
		},
	})

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("miss"),
			},
			Path: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: models.ResponseDetails{
			Body: "two",
		},
	})

	r := models.RequestDetails{
		Body:   "foo",
		Method: "GET",
	}

	result, closestMiss, err := matching.StrongestMatchRequestMatcher(r, false, simulation)

	Expect(err).To(BeNil())
	Expect(result).ToNot(BeNil())
	Expect(closestMiss).To(BeNil())
	Expect(result).ToNot(BeNil())
}

func Test_ShouldReturnStrongestMatchWhenThereAreMultipleMatches(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				RegexMatch: StringToPointer(".*"),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: models.ResponseDetails{
			Body: "one",
		},
	})

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("foo"),
				RegexMatch: StringToPointer(".*"),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("GET"),
				RegexMatch: StringToPointer(".*"),
			},
		},
		Response: models.ResponseDetails{
			Body: "two",
		},
	})

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("foo"),
				RegexMatch: StringToPointer(".*"),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: models.ResponseDetails{
			Body: "three",
		},
	})

	r := models.RequestDetails{
		Body:   "foo",
		Method: "GET",
	}

	result, closestMiss, err := matching.StrongestMatchRequestMatcher(r, false, simulation)

	Expect(err).To(BeNil())
	Expect(closestMiss).To(BeNil())
	Expect(result).ToNot(BeNil())
	Expect(result.Response.Body).To(Equal("two"))
}

func Test_ShouldReturnStrongestMatchWhenThereAreMultipleMatchesAgain(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer(`{"foo": "bar"}`),
				JsonMatch:  StringToPointer(`{"foo": "bar"}`),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("GET"),
				RegexMatch: StringToPointer(".*"),
			},
		},
		Response: models.ResponseDetails{
			Body: "one",
		},
	})

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer(`{"foo": "bar"}`),
				JsonMatch:  StringToPointer(`{"foo": "bar"}`),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: models.ResponseDetails{
			Body: "two",
		},
	})

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				RegexMatch: StringToPointer(`.*`),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("GET"),
				RegexMatch: StringToPointer(".*"),
			},
		},
		Response: models.ResponseDetails{
			Body: "three",
		},
	})

	r := models.RequestDetails{
		Body:   `{"foo": "bar"}`,
		Method: "GET",
	}

	result, closestMiss, err := matching.StrongestMatchRequestMatcher(r, false, simulation)

	Expect(err).To(BeNil())
	Expect(closestMiss).To(BeNil())
	Expect(result).ToNot(BeNil())
	Expect(result.Response.Body).To(Equal("one"))
}

func Test_ShouldSetClosestMissBackToNilIfThereIsAMatchLaterOn(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer(`body`),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: models.ResponseDetails{
			Body: "one",
		},
	})

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer(`body`),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("POST"),
			},
		},
		Response: models.ResponseDetails{
			Body: "two",
		},
	})

	r := models.RequestDetails{
		Body:   `body`,
		Method: "POST",
	}

	_, closestMiss, _ := matching.StrongestMatchRequestMatcher(r, false, simulation)

	Expect(closestMiss).To(BeNil())
}

func Test_ShouldIncludeHeadersInCalculationForStrongestMatch(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				RegexMatch: StringToPointer(".*"),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("GET"),
			},
			Headers: map[string][]string{
				"one":   {"one"},
				"two":   {"one"},
				"three": {"one"},
			},
		},
		Response: models.ResponseDetails{
			Body: "one",
		},
	})

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("foo"),
				RegexMatch: StringToPointer(".*"),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("GET"),
				RegexMatch: StringToPointer(".*"),
			},
		},
		Response: models.ResponseDetails{
			Body: "two",
		},
	})

	r := models.RequestDetails{
		Body:   "foo",
		Method: "GET",
		Headers: map[string][]string{
			"one":   {"one"},
			"two":   {"one"},
			"three": {"one"},
		},
	}

	result, closestMiss, err := matching.StrongestMatchRequestMatcher(r, false, simulation)

	Expect(err).To(BeNil())
	Expect(closestMiss).To(BeNil())
	Expect(result).ToNot(BeNil())
	Expect(result.Response.Body).To(Equal("one"))
}

func Test_ShouldIncludeHeadersInCalculationForClosestMiss(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				RegexMatch: StringToPointer(".*"),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("GET"),
			},
			Headers: map[string][]string{
				"one":   {"one"},
				"two":   {"one"},
				"three": {"one"},
			},
		},
		Response: models.ResponseDetails{
			Body: "one",
		},
	})

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("foo"),
				RegexMatch: StringToPointer(".*"),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("GET"),
				RegexMatch: StringToPointer(".*"),
			},
		},
		Response: models.ResponseDetails{
			Body: "two",
		},
	})

	r := models.RequestDetails{
		Body:   "foo",
		Method: "MISS",
		Headers: map[string][]string{
			"one":   {"one"},
			"two":   {"one"},
			"three": {"one"},
		},
	}

	result, closestMiss, err := matching.StrongestMatchRequestMatcher(r, false, simulation)

	Expect(err).ToNot(BeNil())
	Expect(result).To(BeNil())
	Expect(closestMiss).ToNot(BeNil())
	Expect(closestMiss.Response.Body).To(Equal("one"))
}

func Test_ShouldReturnFieldsMissedInClosestMiss(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				GlobMatch: StringToPointer("miss"),
			},
			Path: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("miss"),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("hit"),
			},
			Destination: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("hit"),
			},
			Query: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("miss"),
			},

			Scheme: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("miss"),
			},
			Headers: map[string][]string{
				"hitKey": {"hitValue"},
			},

		},
		Response: models.ResponseDetails{
			Body: "two",
		},
	})

	r := models.RequestDetails{
		Method: "hit",
		Destination: "hit",
		Headers: map[string][]string{
			"hitKey": {"hitValue"},
		},
	}

	result, closestMiss, err := matching.StrongestMatchRequestMatcher(r, false, simulation)

	Expect(err).ToNot(BeNil())
	Expect(result).To(BeNil())
	Expect(closestMiss).ToNot(BeNil())
	//TODO: Scheme matching?
	Expect(closestMiss.MissedFields).To(ConsistOf(`body`, `path`, `query`))
}

func Test_ShouldReturnFieldsMissedInClosestMissAgain(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				GlobMatch: StringToPointer("hit"),
			},
			Path: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("hit"),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("miss"),
			},
			Destination: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("miss"),
			},
			Query: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("hit"),
			},
			Scheme: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("hit"),
			},
			Headers: map[string][]string{
				"miss": {"miss"},
			},

		},
		Response: models.ResponseDetails{
			Body: "two",
		},
	})

	r := models.RequestDetails{
		Body: "hit",
		Path: "hit",
		Query: "hit",
	}

	result, closestMiss, err := matching.StrongestMatchRequestMatcher(r, false, simulation)

	Expect(err).ToNot(BeNil())
	Expect(result).To(BeNil())
	Expect(closestMiss).ToNot(BeNil())
	//TODO: Scheme matching?
	Expect(closestMiss.MissedFields).To(ConsistOf(`method`, `destination`, `headers`))
}

func Test_ShouldReturnMessageForClosestMiss(t *testing.T) {
	RegisterTestingT(t)

	miss := &matching.ClosestMiss{
		RequestDetails: &models.RequestDetails{
			Path:        "path",
			Method:      "method",
			Destination: "destination",
			Scheme:      "scheme",
			Query:       "query",
			Body:        "body",
			Headers: map[string][]string{
				"miss": {"miss"},
			},
		},
		Response: &v2.ResponseDetailsView{
			Body: "hello world",
			Headers: map[string][]string{
				"hello": {"world"},
			},
			Status: 200,
		},
		RequestMatcher: &v2.RequestMatcherViewV2{
			Body: &v2.RequestFieldMatchersView{
				GlobMatch: StringToPointer("hit"),
			},
			Path: &v2.RequestFieldMatchersView{
				ExactMatch: StringToPointer("hit"),
			},
			Method: &v2.RequestFieldMatchersView{
				ExactMatch: StringToPointer("miss"),
			},
			Destination: &v2.RequestFieldMatchersView{
				ExactMatch: StringToPointer("miss"),
			},
			Query: &v2.RequestFieldMatchersView{
				ExactMatch: StringToPointer("hit"),
			},
			Scheme: &v2.RequestFieldMatchersView{
				ExactMatch: StringToPointer("hit"),
			},
			Headers: map[string][]string{
				"miss": {"miss"},
			},
		},
		MissedFields: []string{"body", "path", "method"},
	}

	Expect(miss.GetMessage()).To(Equal(
		`

The following request was made, but was not matched by Hoverfly:

{
    "Path": "path",
    "Method": "method",
    "Destination": "destination",
    "Scheme": "scheme",
    "Query": "query",
    "Body": "body",
    "Headers": {
        "miss": [
            "miss"
        ]
    }
}

The closest miss was the following matcher:

{
    "path": {
        "exactMatch": "hit"
    },
    "method": {
        "exactMatch": "miss"
    },
    "destination": {
        "exactMatch": "miss"
    },
    "scheme": {
        "exactMatch": "hit"
    },
    "query": {
        "exactMatch": "hit"
    },
    "body": {
        "globMatch": "hit"
    },
    "headers": {
        "miss": [
            "miss"
        ]
    }
}

But it did not match on the following fields:

[body, path, method]

Which if hit would have given the following response:

{
    "status": 200,
    "body": "hello world",
    "encodedBody": false,
    "headers": {
        "hello": [
            "world"
        ]
    }
}`))}
