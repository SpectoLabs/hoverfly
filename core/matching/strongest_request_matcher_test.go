package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
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

func Test_ShouldReturnClosestMatchIfMatchIsNotFound(t *testing.T) {
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
		Response: testResponse,
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
		Response: testResponse,
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
		Response: testResponse,
	})

	r := models.RequestDetails{
		Body: "body",
		Path: "nomatch",
	}

	result, closestMatch, err := matching.StrongestMatchRequestMatcher(r, false, simulation)

	Expect(err).ToNot(BeNil())
	Expect(result).To(BeNil())
	Expect(closestMatch).ToNot(BeNil())
	Expect(*closestMatch.RequestMatcher.Body.ExactMatch).To(Equal(`body`))
	Expect(*closestMatch.RequestMatcher.Body.GlobMatch).To(Equal(`bod*`))
	Expect(*closestMatch.RequestMatcher.Path.ExactMatch).To(Equal(`path`))
}
