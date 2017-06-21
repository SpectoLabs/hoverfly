package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

var testResponse = models.ResponseDetails{
	Body: "request matched",
}

func Test_FirstMatchRequestMatcher_EmptyRequestMatchersShouldMatchOnAnyRequest(t *testing.T) {
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
	result, _ := matching.FirstMatchRequestMatcher(r, false, simulation)

	Expect(result.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchRequestMatcher_RequestMatchersShouldMatchOnBody(t *testing.T) {
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
	result, err := matching.FirstMatchRequestMatcher(r, false, simulation)
	Expect(err).To(BeNil())

	Expect(result.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchRequestMatcher_ReturnResponseWhenAllHeadersMatch(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
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
			"header1": []string{"val1"},
			"header2": []string{"val2"},
		},
	}

	result, _ := matching.FirstMatchRequestMatcher(r, false, simulation)

	Expect(result.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchRequestMatcher_ReturnNilWhenOneHeaderNotPresentInRequest(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
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
			"header1": []string{"val1"},
		},
	}

	result, _ := matching.FirstMatchRequestMatcher(r, false, simulation)

	Expect(result).To(BeNil())
}

func Test_FirstMatchRequestMatcher_ReturnNilWhenOneHeaderValueDifferent(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
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
			"header1": []string{"val1"},
			"header2": []string{"different"},
		},
	}
	result, _ := matching.FirstMatchRequestMatcher(r, false, simulation)

	Expect(result).To(BeNil())
}

func Test_FirstMatchRequestMatcher_ReturnResponseWithMultiValuedHeaderMatch(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
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
			"header1": []string{"val1-a", "val1-b"},
			"header2": []string{"val2"},
		},
	}
	result, _ := matching.FirstMatchRequestMatcher(r, false, simulation)

	Expect(result.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchRequestMatcher_ReturnNilWithDifferentMultiValuedHeaders(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
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
			"header1": []string{"val1-a", "val1-differnet"},
			"header2": []string{"val2"},
		},
	}

	result, _ := matching.FirstMatchRequestMatcher(r, false, simulation)

	Expect(result).To(BeNil())
}

func Test_FirstMatchRequestMatcher_EndpointMatchWithHeaders(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
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
			"header1": []string{"val1-a", "val1-b"},
			"header2": []string{"val2"},
		},
	}
	result, _ := matching.FirstMatchRequestMatcher(r, false, simulation)

	Expect(result.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchRequestMatcher_EndpointMismatchWithHeadersReturnsNil(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
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
			"header1": []string{"val1-a", "val1-b"},
			"header2": []string{"val2"},
		},
	}

	result, _ := matching.FirstMatchRequestMatcher(r, false, simulation)

	Expect(result).To(BeNil())
}

func Test_FirstMatchRequestMatcher_AbleToMatchAnEmptyPathInAReasonableWay(t *testing.T) {
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
	result, _ := matching.FirstMatchRequestMatcher(r, false, simulation)

	Expect(result.Response.Body).To(Equal("request matched"))

	r = models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/a/1",
		Query:       "q=test",
	}

	result, _ = matching.FirstMatchRequestMatcher(r, false, simulation)

	Expect(result).To(BeNil())
}

func Test_FirstMatchRequestMatcher_RequestMatcherResponsePairCanBeConvertedToARequestResponsePairView_WhileIncomplete(t *testing.T) {
	RegisterTestingT(t)

	method := "POST"

	requestMatcherResponsePair := models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				ExactMatch: &method,
			},
		},
		Response: testResponse,
	}

	pairView := requestMatcherResponsePair.BuildView()

	Expect(pairView.RequestMatcher.Method.ExactMatch).To(Equal(StringToPointer("POST")))
	Expect(pairView.RequestMatcher.Destination).To(BeNil())
	Expect(pairView.RequestMatcher.Path).To(BeNil())
	Expect(pairView.RequestMatcher.Scheme).To(BeNil())
	Expect(pairView.RequestMatcher.Query).To(BeNil())

	Expect(pairView.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchRequestMatcher_RequestMatchersCanUseGlobsAndBeMatched(t *testing.T) {
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

	response, err := matching.FirstMatchRequestMatcher(request, false, simulation)
	Expect(err).To(BeNil())

	Expect(response.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchRequestMatcher_RequestMatchersCanUseGlobsOnSchemeAndBeMatched(t *testing.T) {
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

	response, err := matching.FirstMatchRequestMatcher(request, false, simulation)
	Expect(err).To(BeNil())

	Expect(response.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchRequestMatcher_RequestMatchersCanUseGlobsOnHeadersAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: map[string][]string{
				"unique-header": []string{"*"},
			},
		},
		Response: testResponse,
	})

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
		Headers: map[string][]string{
			"unique-header": []string{"totally-unique"},
		},
	}

	response, err := matching.FirstMatchRequestMatcher(request, false, simulation)
	Expect(err).To(BeNil())

	Expect(response.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchRequestMatcher_RequestMatcherResponsePair_ConvertToRequestResponsePairView_CanBeConvertedToARequestResponsePairView_WhileIncomplete(t *testing.T) {
	RegisterTestingT(t)

	method := "POST"

	requestMatcherResponsePair := models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				ExactMatch: &method,
			},
		},
		Response: testResponse,
	}

	pairView := requestMatcherResponsePair.BuildView()

	Expect(pairView.RequestMatcher.Method.ExactMatch).To(Equal(StringToPointer("POST")))
	Expect(pairView.RequestMatcher.Destination).To(BeNil())
	Expect(pairView.RequestMatcher.Path).To(BeNil())
	Expect(pairView.RequestMatcher.Scheme).To(BeNil())
	Expect(pairView.RequestMatcher.Query).To(BeNil())

	Expect(pairView.Response.Body).To(Equal("request matched"))
}

func Test_ShouldStoreIfMatchedOnEverythingApartFromHeadersAtLeastOnce(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("POST"),
			},
			Body: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("body"),
			},
			Scheme: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("http"),
			},
			Query: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("?foo=bar"),
			},
			Path: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("/foo"),
			},
			Destination: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("www.test.com"),
			},
			Headers: map[string][]string {
				"foo" : {"bar"},
			},
		},
		Response: testResponse,
	})

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query:       "?foo=bar",
		Scheme: "http",
		Body: "body",
		Path: "/foo",
		Headers: map[string][]string {
			"miss" : {"me"},
		},
	}

	_, err := matching.FirstMatchRequestMatcher(r, false, simulation)

	Expect(err).ToNot(BeNil())
	Expect(err.MatchedOnAllButHeadersAtLeastOnce).To(BeTrue())
}

func Test_ShouldStoreIfMatchedOnEverythingApartFromHeadersZeroTimes(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("POST"),
			},
			Body: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("body"),
			},
			Scheme: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("http"),
			},
			Query: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("?foo=bar"),
			},
			Path: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("/foo"),
			},
			Destination: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("www.test.com"),
			},
			Headers: map[string][]string {
				"foo" : {"bar"},
			},
		},
		Response: testResponse,
	})

	simulation.MatchingPairs = append(simulation.MatchingPairs, models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "MISS",
		Destination: "www.test.com",
		Query:       "?foo=bar",
		Scheme: "http",
		Body: "body",
		Path: "/foo",
		Headers: map[string][]string {
			"miss" : {"me"},
		},
	}

	_, err := matching.FirstMatchRequestMatcher(r, false, simulation)

	Expect(err).ToNot(BeNil())
	Expect(err.MatchedOnAllButHeadersAtLeastOnce).To(BeFalse())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "miss",
		Query:       "?foo=bar",
		Scheme: "http",
		Body: "body",
		Path: "/foo",
		Headers: map[string][]string {
			"miss" : {"me"},
		},
	}

	_, err = matching.FirstMatchRequestMatcher(r, false, simulation)

	Expect(err).ToNot(BeNil())
	Expect(err.MatchedOnAllButHeadersAtLeastOnce).To(BeFalse())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query:       "miss",
		Scheme: "http",
		Body: "body",
		Path: "/foo",
		Headers: map[string][]string {
			"miss" : {"me"},
		},
	}

	_, err = matching.FirstMatchRequestMatcher(r, false, simulation)

	Expect(err).ToNot(BeNil())
	Expect(err.MatchedOnAllButHeadersAtLeastOnce).To(BeFalse())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query:       "?foo=bar",
		Scheme: "http",
		Body: "miss",
		Path: "/foo",
		Headers: map[string][]string {
			"miss" : {"me"},
		},
	}

	_, err = matching.FirstMatchRequestMatcher(r, false, simulation)

	Expect(err).ToNot(BeNil())
	Expect(err.MatchedOnAllButHeadersAtLeastOnce).To(BeFalse())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query:       "?foo=bar",
		Scheme: "http",
		Body: "body",
		Path: "miss",
		Headers: map[string][]string {
			"miss" : {"me"},
		},
	}

	_, err = matching.FirstMatchRequestMatcher(r, false, simulation)

	Expect(err).ToNot(BeNil())
	Expect(err.MatchedOnAllButHeadersAtLeastOnce).To(BeFalse())
}