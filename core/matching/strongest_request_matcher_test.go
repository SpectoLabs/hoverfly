package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

func Test_ClosestRequestMatcherRequestMatcher_EmptyRequestMatchersShouldMatchOnAnyRequest(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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
	result, _, _ := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(result).ToNot(BeNil())
	Expect(result.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_RequestMatchersShouldMatchOnBody(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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
	result, err, _ := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})
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

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	result, _, _ := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(result.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_ReturnNilWhenOneHeaderNotPresentInRequest(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": {"val1"},
		"header2": {"val2"},
	}

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	result, _, _ := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(result).To(BeNil())
}

func Test_ClosestRequestMatcherRequestMatcher_ReturnNilWhenOneHeaderValueDifferent(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": {"val1"},
		"header2": {"val2"},
	}

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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
	result, _, _ := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(result).To(BeNil())
}

func Test_ClosestRequestMatcherRequestMatcher_ReturnResponseWithMultiValuedHeaderMatch(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": {"val1-a", "val1-b"},
		"header2": {"val2"},
	}

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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
	result, _, _ := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(result.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_ReturnNilWithDifferentMultiValuedHeaders(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": {"val1-a", "val1-b"},
		"header2": {"val2"},
	}

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	result, _, _ := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

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

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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
		Query: map[string][]string{
			"q": []string{"test"},
		},
		Headers: map[string][]string{
			"header1": {"val1-a", "val1-b"},
			"header2": {"val2"},
		},
	}
	result, _, _ := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

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

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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
		Query: map[string][]string{
			"q": []string{"different"},
		},
		Headers: map[string][]string{
			"header1": {"val1-a", "val1-b"},
			"header2": {"val2"},
		},
	}

	result, _, _ := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(result).To(BeNil())
}

func Test_ClosestRequestMatcherRequestMatcher_AbleToMatchAnEmptyPathInAReasonableWay(t *testing.T) {
	RegisterTestingT(t)

	destination := "testhost.com"
	method := "GET"
	path := ""
	query := "q=test"
	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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
		Query: map[string][]string{
			"q": []string{"test"},
		},
	}
	result, _, _ := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(result.Response.Body).To(Equal("request matched"))

	r = models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/a/1",
		Query: map[string][]string{
			"q": []string{"test"},
		},
	}

	result, _, _ = matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(result).To(BeNil())
}

func Test_ClosestRequestMatcherRequestMatcher_RequestMatchersCanUseGlobsAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	response, err, _ := matching.StrongestMatchRequestMatcher(request, false, simulation, map[string]string{})
	Expect(err).To(BeNil())

	Expect(response.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_RequestMatchersCanUseGlobsOnSchemeAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	response, err, _ := matching.StrongestMatchRequestMatcher(request, false, simulation, map[string]string{})
	Expect(err).To(BeNil())

	Expect(response.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_RequestMatchersCanUseGlobsOnHeadersAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	response, err, _ := matching.StrongestMatchRequestMatcher(request, false, simulation, map[string]string{})
	Expect(err).To(BeNil())

	Expect(response.Response.Body).To(Equal("request matched"))
}

func Test_ShouldReturnClosestMissIfMatchIsNotFound(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	result, err, _ := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(err).ToNot(BeNil())
	Expect(result).To(BeNil())
	Expect(err.ClosestMiss).ToNot(BeNil())
	Expect(*err.ClosestMiss.RequestMatcher.Body.ExactMatch).To(Equal(`body`))
	Expect(*err.ClosestMiss.RequestMatcher.Body.GlobMatch).To(Equal(`bod*`))
	Expect(*err.ClosestMiss.RequestMatcher.Path.ExactMatch).To(Equal(`path`))
	Expect(err.ClosestMiss.Response.Body).To(Equal(`two`))
	Expect(err.ClosestMiss.RequestDetails.Body).To(Equal(`body`))
}

func Test_ShouldReturnClosestMissIfMatchIsNotFoundAgain(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	result, err, _ := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(err).ToNot(BeNil())
	Expect(result).To(BeNil())
	Expect(err.ClosestMiss).ToNot(BeNil())
	Expect(*err.ClosestMiss.RequestMatcher.Body.RegexMatch).To(Equal(`.*`))
	Expect(*err.ClosestMiss.RequestMatcher.Path.ExactMatch).To(Equal(`miss`))
	Expect(*err.ClosestMiss.RequestMatcher.Method.ExactMatch).To(Equal(`GET`))
	Expect(err.ClosestMiss.Response.Body).To(Equal(`one`))
}

func Test_ShouldNotReturnClosestMissWhenThereIsAMatch(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	result, err, _ := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(err).To(BeNil())
	Expect(result).ToNot(BeNil())
}

func Test__NotBeCachableIfMatchedOnEverythingApartFromHeadersAtLeastOnce(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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
				ExactMatch: StringToPointer("foo=bar"),
			},
			Path: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("/foo"),
			},
			Destination: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("www.test.com"),
			},
			Headers: map[string][]string{
				"foo": {"bar"},
			},
		},
		Response: testResponse,
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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
		Query: map[string][]string{
			"foo": []string{"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	_, err, cachable := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(err).ToNot(BeNil())
	Expect(cachable).To(BeFalse())
}

func Test__ShouldBeCachableIfMatchedOnEverythingApartFromHeadersZeroTimes(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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
			Headers: map[string][]string{
				"foo": {"bar"},
			},
		},
		Response: testResponse,
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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
		Query: map[string][]string{
			"foo": []string{"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	_, err, cachable := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(err).ToNot(BeNil())
	Expect(cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "miss",
		Query: map[string][]string{
			"foo": []string{"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	_, err, cachable = matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(err).ToNot(BeNil())
	Expect(cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"miss": []string{""},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	_, err, cachable = matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(err).ToNot(BeNil())
	Expect(cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": []string{"bar"},
		},
		Scheme: "http",
		Body:   "miss",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	_, err, cachable = matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(err).ToNot(BeNil())
	Expect(cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": []string{"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "miss",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	_, err, _ = matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(err).ToNot(BeNil())
	Expect(cachable).To(BeTrue())
}

func Test_ShouldReturnStrongestMatchWhenThereAreMultipleMatches(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	result, err, _ := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(err).To(BeNil())
	Expect(result).ToNot(BeNil())
	Expect(result.Response.Body).To(Equal("two"))
}

func Test_ShouldReturnStrongestMatchWhenThereAreMultipleMatchesAgain(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	result, err, _ := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(err).To(BeNil())
	Expect(result).ToNot(BeNil())
	Expect(result.Response.Body).To(Equal("one"))
}

func Test_ShouldSetClosestMissBackToNilIfThereIsAMatchLaterOn(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	_, err, _ := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(err).To(BeNil())
}

func Test_ShouldIncludeHeadersInCalculationForStrongestMatch(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	result, err, _ := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(err).To(BeNil())
	Expect(result).ToNot(BeNil())
	Expect(result.Response.Body).To(Equal("one"))
}

func Test_ShouldIncludeHeadersInCalculationForClosestMiss(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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

	result, err, _ := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(err).ToNot(BeNil())
	Expect(result).To(BeNil())
	Expect(err.ClosestMiss).ToNot(BeNil())
	Expect(err.ClosestMiss.Response.Body).To(Equal("one"))
}

func Test_ShouldIncludeStateInCalculationForClosestMiss(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				RegexMatch: StringToPointer(".*"),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("GET"),
			},
			RequiresState: map[string]string{
				"one":   "one",
				"two":   "one",
				"three": "one",
			},
		},
		Response: models.ResponseDetails{
			Body: "one",
		},
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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
	}

	result, err, _ := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{
		"one":   "one",
		"two":   "one",
		"three": "one",
	})

	Expect(err).ToNot(BeNil())
	Expect(result).To(BeNil())
	Expect(err.ClosestMiss).ToNot(BeNil())
	Expect(err.ClosestMiss.Response.Body).To(Equal("one"))
}

func Test_ShouldReturnFieldsMissedInClosestMiss(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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
		Method:      "hit",
		Destination: "hit",
		Headers: map[string][]string{
			"hitKey": {"hitValue"},
		},
	}

	result, err, _ := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(err).ToNot(BeNil())
	Expect(result).To(BeNil())
	Expect(err.ClosestMiss).ToNot(BeNil())
	//TODO: Scheme matching?
	Expect(err.ClosestMiss.MissedFields).To(ConsistOf(`body`, `path`, `query`))
}

func Test_ShouldReturnFieldsMissedInClosestMissAgain(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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
				ExactMatch: StringToPointer("hit="),
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
		Query: map[string][]string{
			"hit": []string{""},
		},
	}

	result, err, _ := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{})

	Expect(err).ToNot(BeNil())
	Expect(result).To(BeNil())
	Expect(err.ClosestMiss).ToNot(BeNil())
	//TODO: Scheme matching?
	Expect(err.ClosestMiss.MissedFields).To(ConsistOf(`method`, `destination`, `headers`))
}

func Test_ShouldReturnMessageForClosestMiss(t *testing.T) {
	RegisterTestingT(t)

	miss := &models.ClosestMiss{
		RequestDetails: models.RequestDetails{
			Path:        "path",
			Method:      "method",
			Destination: "destination",
			Scheme:      "scheme",
			Query: map[string][]string{
				"query": {""},
			},
			Body: "body",
			Headers: map[string][]string{
				"miss": {"miss"},
			},
		},
		State: map[string]string{
			"key1": "value2",
			"key3": "value4",
		},
		Response: v2.ResponseDetailsViewV4{
			Body: "hello world",
			Headers: map[string][]string{
				"hello": {"world"},
			},
			Status: 200,
		},
		RequestMatcher: v2.RequestMatcherViewV4{
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

	message := miss.GetMessage()
	Expect(message).To(Equal(
		`

The following request was made, but was not matched by Hoverfly:

{
    "Path": "path",
    "Method": "method",
    "Destination": "destination",
    "Scheme": "scheme",
    "Query": {
        "query": [
            ""
        ]
    },
    "Body": "body",
    "Headers": {
        "miss": [
            "miss"
        ]
    }
}

Whilst Hoverfly has the following state:

{
    "key1": "value2",
    "key3": "value4"
}

The matcher which came closest was:

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
    },
    "templated": false
}`))
}

func Test_StrongestMatch_ShouldNotBeCachableIfMatchedOnEverythingApartFromHeadersAtLeastOnce(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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
				ExactMatch: StringToPointer("foo=bar"),
			},
			Path: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("/foo"),
			},
			Destination: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("www.test.com"),
			},
			Headers: map[string][]string{
				"foo": {"bar"},
			},
		},
		Response: testResponse,
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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
		Query: map[string][]string{
			"foo": []string{"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	_, err, cachable := matching.StrongestMatchRequestMatcher(r, false, simulation, make(map[string]string))

	Expect(err).ToNot(BeNil())
	Expect(cachable).To(BeFalse())
}

func Test_StrongestMatch__ShouldBeCachableIfMatchedOnEverythingApartFromHeadersZeroTimes(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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
			Headers: map[string][]string{
				"foo": {"bar"},
			},
		},
		Response: testResponse,
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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
		Query: map[string][]string{
			"foo": []string{"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	_, err, cachable := matching.StrongestMatchRequestMatcher(r, false, simulation, make(map[string]string))

	Expect(err).ToNot(BeNil())
	Expect(cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "miss",
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	_, err, cachable = matching.StrongestMatchRequestMatcher(r, false, simulation, make(map[string]string))

	Expect(err).ToNot(BeNil())
	Expect(cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"miss": []string{""},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	_, err, cachable = matching.StrongestMatchRequestMatcher(r, false, simulation, make(map[string]string))

	Expect(err).ToNot(BeNil())
	Expect(cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "miss",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	_, err, cachable = matching.StrongestMatchRequestMatcher(r, false, simulation, make(map[string]string))

	Expect(err).ToNot(BeNil())
	Expect(cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "miss",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	_, err, cachable = matching.StrongestMatchRequestMatcher(r, false, simulation, make(map[string]string))

	Expect(err).ToNot(BeNil())
	Expect(cachable).To(BeTrue())
}

func Test_StrongestMatchRequestMatcher_RequestMatchersShouldMatchOnStateAndNotBeCachable(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			RequiresState: map[string]string{"key1": "value1", "key2": "value2"},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Body: "body",
	}

	result, err, cachable := matching.StrongestMatchRequestMatcher(
		r,
		false,
		simulation,
		map[string]string{"key1": "value1", "key2": "value2"})

	Expect(err).To(BeNil())
	Expect(cachable).To(BeFalse())
	Expect(result.Response.Body).To(Equal("request matched"))
}

func Test_StrongestMatch_ShouldNotBeCachableIfMatchedOnEverythingApartFromStateAtLeastOnce(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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
				ExactMatch: StringToPointer("foo=bar"),
			},
			Path: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("/foo"),
			},
			Destination: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("www.test.com"),
			},
			RequiresState: map[string]string{
				"foo": "bar",
			},
		},
		Response: testResponse,
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
	}

	_, err, cachable := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{"miss": "me"})

	Expect(err).ToNot(BeNil())
	Expect(cachable).To(BeFalse())
}

func Test_StrongestMatch__ShouldBeCachableIfMatchedOnEverythingApartFromStateZeroTimes(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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
			RequiresState: map[string]string{
				"foo": "bar",
			},
		},
		Response: testResponse,
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
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
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
	}

	_, err, cachable := matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{"miss": "me"})

	Expect(err).ToNot(BeNil())
	Expect(cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "miss",
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
	}

	_, err, cachable = matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{"miss": "me"})

	Expect(err).ToNot(BeNil())
	Expect(cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"miss": []string{""},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
	}

	_, err, cachable = matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{"miss": "me"})

	Expect(err).ToNot(BeNil())
	Expect(cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "miss",
		Path:   "/foo",
	}

	_, err, cachable = matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{"miss": "me"})

	Expect(err).ToNot(BeNil())
	Expect(cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "miss",
	}

	_, err, cachable = matching.StrongestMatchRequestMatcher(r, false, simulation, map[string]string{"miss": "me"})

	Expect(err).ToNot(BeNil())
	Expect(cachable).To(BeTrue())
}
