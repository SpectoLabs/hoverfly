package matching_test

import (
	"sync"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/state"
	. "github.com/onsi/gomega"
)

func Test_ClosestRequestMatcherRequestMatcher_EmptyRequestMatchersShouldMatchOnAnyRequest(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
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
	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Pair).ToNot(BeNil())
	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_RequestMatchersShouldMatchOnBody(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "body",
				},
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Body: "body",
	}
	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})
	Expect(result.Error).To(BeNil())

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_ReturnResponseWhenAllHeadersMatch(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]models.RequestFieldMatchers{
		"header1": {
			{
				Matcher: matchers.Exact,
				Value:   "val1",
			},
		},
		"header2": {
			{
				Matcher: matchers.Exact,
				Value:   "val2",
			},
		},
	}

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
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

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_ReturnNilWhenOneHeaderNotPresentInRequest(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]models.RequestFieldMatchers{
		"header1": {
			{
				Matcher: matchers.Exact,
				Value:   "val1",
			},
		},
		"header2": {
			{
				Matcher: matchers.Exact,
				Value:   "val2",
			},
		},
	}

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
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

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Pair).To(BeNil())
}

func Test_ClosestRequestMatcherRequestMatcher_ReturnNilWhenOneHeaderValueDifferent(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]models.RequestFieldMatchers{
		"header1": {
			{
				Matcher: matchers.Exact,
				Value:   "val1",
			},
		},
		"header2": {
			{
				Matcher: matchers.Exact,
				Value:   "val2",
			},
		},
	}

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
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
	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Pair).To(BeNil())
}

func Test_ClosestRequestMatcherRequestMatcher_ReturnResponseWithMultiValuedHeaderMatch(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]models.RequestFieldMatchers{
		"header1": {
			{
				Matcher: matchers.Exact,
				Value:   "val1-a;val1-b",
			},
		},
		"header2": {
			{
				Matcher: matchers.Exact,
				Value:   "val2",
			},
		},
	}

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
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
	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_ReturnNilWithDifferentMultiValuedHeaders(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]models.RequestFieldMatchers{
		"header1": {
			{
				Matcher: matchers.Exact,
				Value:   "val1-a;val1-b",
			},
		},
		"header2": {
			{
				Matcher: matchers.Exact,
				Value:   "val2",
			},
		},
	}
	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
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

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Pair).To(BeNil())
}

func Test_ClosestRequestMatcherRequestMatcher_EndpointMatchWithHeaders(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]models.RequestFieldMatchers{
		"header1": {
			{
				Matcher: matchers.Exact,
				Value:   "val1-a;val1-b",
			},
		},
		"header2": {
			{
				Matcher: matchers.Exact,
				Value:   "val2",
			},
		},
	}

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: headers,
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "testhost.com",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "/a/1",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "GET",
				},
			},
			DeprecatedQuery: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "q=test",
				},
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/a/1",
		Query: map[string][]string{
			"q": {"test"},
		},
		Headers: map[string][]string{
			"header1": {"val1-a", "val1-b"},
			"header2": {"val2"},
		},
	}
	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_EndpointMismatchWithHeadersReturnsNil(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]models.RequestFieldMatchers{
		"header1": {
			{
				Matcher: matchers.Exact,
				Value:   "val1-a;val1-b",
			},
		},
		"header2": {
			{
				Matcher: matchers.Exact,
				Value:   "val2",
			},
		},
	}

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: headers,
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "testhost.com",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "/a/1",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "GET",
				},
			},
			DeprecatedQuery: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "q=test",
				},
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://testhost.com",
		Path:        "/a/1",
		Query: map[string][]string{
			"q": {"different"},
		},
		Headers: map[string][]string{
			"header1": {"val1-a", "val1-b"},
			"header2": {"val2"},
		},
	}

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Pair).To(BeNil())
}

func Test_ClosestRequestMatcherRequestMatcher_AbleToMatchAnEmptyPathInAReasonableWay(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "testhost.com",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "GET",
				},
			},
			DeprecatedQuery: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "q=test",
				},
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Query: map[string][]string{
			"q": {"test"},
		},
	}
	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Pair.Response.Body).To(Equal("request matched"))

	r = models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/a/1",
		Query: map[string][]string{
			"q": {"test"},
		},
	}

	result = matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Pair).To(BeNil())
}

func Test_ClosestRequestMatcherRequestMatcher_RequestMatchersCanUseGlobsAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Glob,
					Value:   "*.com",
				},
			},
		},
		Response: testResponse,
	})

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
	}

	result := matching.MatchingStrategyRunner(request, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})
	Expect(result.Error).To(BeNil())

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_RequestMatchersCanUseGlobsOnSchemeAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Glob,
					Value:   "*.com",
				},
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

	result := matching.MatchingStrategyRunner(request, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})
	Expect(result.Error).To(BeNil())

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_RequestMatchersCanUseGlobsOnHeadersAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: map[string][]models.RequestFieldMatchers{
				"unique-header": {
					{
						Matcher: matchers.Glob,
						Value:   "*",
					},
				},
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

	result := matching.MatchingStrategyRunner(request, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})
	Expect(result.Error).To(BeNil())

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_ShouldReturnClosestMissIfMatchIsNotFound(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "completemiss",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "completemiss",
				},
			},
		},
		Response: models.ResponseDetails{
			Body: "one",
		},
	})

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "body",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "path",
				},
			},
		},
		Response: models.ResponseDetails{
			Body: "two",
		},
	})

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "body",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "path",
				},
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

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Pair).To(BeNil())
	Expect(result.Error.ClosestMiss).ToNot(BeNil())
	Expect(result.Error.ClosestMiss.RequestMatcher.Body[0].Matcher).To(Equal(`exact`))
	Expect(result.Error.ClosestMiss.RequestMatcher.Body[0].Value).To(Equal(`body`))
	Expect(result.Error.ClosestMiss.RequestMatcher.Path[0].Matcher).To(Equal(`exact`))
	Expect(result.Error.ClosestMiss.RequestMatcher.Path[0].Value).To(Equal(`path`))
	Expect(result.Error.ClosestMiss.Response.Body).To(Equal(`two`))
	Expect(result.Error.ClosestMiss.RequestDetails.Body).To(Equal(`body`))
}

func Test_ShouldReturnClosestMissIfMatchIsNotFoundAgain(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Regex,
					Value:   ".*",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "miss",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "GET",
				},
			},
		},
		Response: models.ResponseDetails{
			Body: "one",
		},
	})

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   ".*",
				},
				// GlobMatch:  StringToPointer("miss"),
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "miss",
				},
			},
		},
		Response: models.ResponseDetails{
			Body: "two",
		},
	})

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "miss",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "miss",
				},
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

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Pair).To(BeNil())
	Expect(result.Error.ClosestMiss).ToNot(BeNil())
	Expect(result.Error.ClosestMiss.RequestMatcher.Body[0].Matcher).To(Equal(`regex`))
	Expect(result.Error.ClosestMiss.RequestMatcher.Body[0].Value).To(Equal(`.*`))
	Expect(result.Error.ClosestMiss.RequestMatcher.Path[0].Matcher).To(Equal(`exact`))
	Expect(result.Error.ClosestMiss.RequestMatcher.Path[0].Value).To(Equal(`miss`))
	Expect(result.Error.ClosestMiss.RequestMatcher.Method[0].Matcher).To(Equal(`exact`))
	Expect(result.Error.ClosestMiss.RequestMatcher.Method[0].Value).To(Equal(`GET`))
	Expect(result.Error.ClosestMiss.Response.Body).To(Equal(`one`))
}

func Test_ShouldNotReturnClosestMissWhenThereIsAMatch(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Regex,
					Value:   ".*",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "GET",
				},
			},
		},
		Response: models.ResponseDetails{
			Body: "one",
		},
	})

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "miss",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "GET",
				},
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

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).To(BeNil())
	Expect(result.Pair).ToNot(BeNil())
}

func Test__NotBeCachableIfMatchedOnEverythingApartFromHeadersAtLeastOnce(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "POST",
				},
			},
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "body",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "http",
				},
			},
			DeprecatedQuery: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "foo=bar",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "/foo",
				},
			},
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "www.test.com",
				},
			},
			Headers: map[string][]models.RequestFieldMatchers{
				"foo": {
					{
						Matcher: matchers.Exact,
						Value:   "bar",
					},
				},
			},
		},
		Response: testResponse,
	})

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "GET",
				},
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
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeFalse())
}

func Test__ShouldBeCachableIfMatchedOnEverythingApartFromHeadersZeroTimes(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "POST",
				},
			},
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "body",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "http",
				},
			},
			DeprecatedQuery: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "?foo=bar",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "/foo",
				},
			},
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "www.test.com",
				},
			},
			Headers: map[string][]models.RequestFieldMatchers{
				"foo": {
					{
						Matcher: matchers.Exact,
						Value:   "bar",
					},
				},
			},
		},
		Response: testResponse,
	})

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "GET",
				},
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
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

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

	result = matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"miss": {""},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result = matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

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

	result = matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

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

	result = matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())
}

func Test_ShouldSetClosestMissBackToNilIfThereIsAMatchLaterOn(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "body",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "GET",
				},
			},
		},
		Response: models.ResponseDetails{
			Body: "one",
		},
	})

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "body",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "POST",
				},
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

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).To(BeNil())
}

func Test_ShouldIncludeHeadersInCalculationForStrongestMatch(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Regex,
					Value:   ".*",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "GET",
				},
			},
			Headers: map[string][]models.RequestFieldMatchers{
				"one": {
					{
						Matcher: matchers.Exact,
						Value:   "one",
					},
				},
				"two": {
					{
						Matcher: matchers.Exact,
						Value:   "one",
					},
				},
				"three": {
					{
						Matcher: matchers.Exact,
						Value:   "one",
					},
				},
			},
		},
		Response: models.ResponseDetails{
			Body: "one",
		},
	})

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Regex,
					Value:   ".*",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "GET",
				},
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

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).To(BeNil())
	Expect(result.Pair).ToNot(BeNil())
	Expect(result.Pair.Response.Body).To(Equal("one"))
}

func Test_ShouldIncludeHeadersInCalculationForClosestMiss(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "GET",
				},
			},
			Headers: map[string][]models.RequestFieldMatchers{
				"one": {
					{
						Matcher: matchers.Exact,
						Value:   "one",
					},
				},
				"two": {
					{
						Matcher: matchers.Exact,
						Value:   "one",
					},
				},
				"three": {
					{
						Matcher: matchers.Exact,
						Value:   "one",
					},
				},
			},
		},
		Response: models.ResponseDetails{
			Body: "one",
		},
	})

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Regex,
					Value:   "GET",
				},
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

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Pair).To(BeNil())
	Expect(result.Error.ClosestMiss).ToNot(BeNil())
	Expect(result.Error.ClosestMiss.Response.Body).To(Equal("one"))
}

func Test_ShouldReturnFieldsMissedInClosestMiss(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Glob,
					Value:   "miss",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "miss",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "hit",
				},
			},
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "hit",
				},
			},
			DeprecatedQuery: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "miss",
				},
			},

			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "miss",
				},
			},
			Headers: map[string][]models.RequestFieldMatchers{
				"hitKey": {
					{
						Matcher: matchers.Exact,
						Value:   "hitValue",
					},
				},
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

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Pair).To(BeNil())
	Expect(result.Error.ClosestMiss).ToNot(BeNil())
	//TODO: Scheme matching?
	Expect(result.Error.ClosestMiss.MissedFields).To(ConsistOf(`body`, `path`, `query`))
}

func Test_ShouldReturnFieldsMissedInClosestMissAgain(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Glob,
					Value:   "hit",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "hit",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "miss",
				},
			},
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "miss",
				},
			},
			DeprecatedQuery: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "hit=",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "hit",
				},
			},
			Headers: map[string][]models.RequestFieldMatchers{
				"miss": {
					{
						Matcher: matchers.Exact,
						Value:   "miss",
					},
				},
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
			"hit": {""},
		},
	}

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Pair).To(BeNil())
	Expect(result.Error.ClosestMiss).ToNot(BeNil())
	//TODO: Scheme matching?
	Expect(result.Error.ClosestMiss.MissedFields).To(ConsistOf(`method`, `destination`, `headers`))
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
		Response: v2.ResponseDetailsViewV5{
			Body: "hello world",
			Headers: map[string][]string{
				"hello": {"world"},
			},
			Status: 200,
		},
		RequestMatcher: v2.RequestMatcherViewV5{
			Body: []v2.MatcherViewV5{
				{
					Matcher: matchers.Glob,
					Value:   "hit",
				},
			},
			Path: []v2.MatcherViewV5{
				{
					Matcher: matchers.Exact,
					Value:   "hit",
				},
			},
			Method: []v2.MatcherViewV5{
				{
					Matcher: matchers.Exact,
					Value:   "miss",
				},
			},
			Destination: []v2.MatcherViewV5{
				{
					Matcher: matchers.Exact,
					Value:   "miss",
				},
			},
			Query: &v2.QueryMatcherViewV5{
				"query": []v2.MatcherViewV5{
					{
						Matcher: matchers.Exact,
						Value:   "hit",
					},
				},
			},
			Scheme: []v2.MatcherViewV5{
				{
					Matcher: matchers.Exact,
					Value:   "hit",
				},
			},
			Headers: map[string][]v2.MatcherViewV5{
				"miss": {
					{
						Matcher: matchers.Exact,
						Value:   "miss",
					},
				},
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
    "path": [
        {
            "matcher": "exact",
            "value": "hit"
        }
    ],
    "method": [
        {
            "matcher": "exact",
            "value": "miss"
        }
    ],
    "destination": [
        {
            "matcher": "exact",
            "value": "miss"
        }
    ],
    "scheme": [
        {
            "matcher": "exact",
            "value": "hit"
        }
    ],
    "body": [
        {
            "matcher": "glob",
            "value": "hit"
        }
    ],
    "headers": {
        "miss": [
            {
                "matcher": "exact",
                "value": "miss"
            }
        ]
    },
    "query": {
        "query": [
            {
                "matcher": "exact",
                "value": "hit"
            }
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

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "POST",
				},
			},
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "body",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "http",
				},
			},
			DeprecatedQuery: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "foo=bar",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "/foo",
				},
			},
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "www.test.com",
				},
			},
			Headers: map[string][]models.RequestFieldMatchers{
				"foo": {
					{
						Matcher: matchers.Exact,
						Value:   "bar",
					},
				},
			},
		},
		Response: testResponse,
	})

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "GET",
				},
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
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeFalse())
}

func Test_StrongestMatch__ShouldBeCachableIfMatchedOnEverythingApartFromHeadersZeroTimes(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "POST",
				},
			},
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "body",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "http",
				},
			},
			DeprecatedQuery: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "?foo=bar",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "/foo",
				},
			},
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "www.test.com",
				},
			},
			Headers: map[string][]models.RequestFieldMatchers{
				"foo": {
					{
						Matcher: matchers.Exact,
						Value:   "bar",
					},
				},
			},
		},
		Response: testResponse,
	})

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "GET",
				},
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
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

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

	result = matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"miss": {""},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result = matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

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

	result = matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

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

	result = matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())
}

func Test_MatchingStrategyRunner_RequestMatchersShouldMatchOnStateAndNotBeCachable(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			RequiresState: map[string]string{"key1": "value1", "key2": "value2"},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Body: "body",
	}

	result := matching.MatchingStrategyRunner(
		r,
		false,
		simulation,
		&state.State{map[string]string{"key1": "value1", "key2": "value2"}, sync.RWMutex{}},
		&matching.StrongestMatchStrategy{})

	Expect(result.Error).To(BeNil())
	Expect(result.Cachable).To(BeFalse())
	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_StrongestMatch_ShouldNotBeCachableIfMatchedOnEverythingApartFromStateAtLeastOnce(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "POST",
				},
			},
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "body",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "http",
				},
			},
			DeprecatedQuery: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "foo=bar",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "/foo",
				},
			},
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "www.test.com",
				},
			},
			RequiresState: map[string]string{
				"foo": "bar",
			},
		},
		Response: testResponse,
	})

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "GET",
				},
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

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{map[string]string{"miss": "me"}, sync.RWMutex{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeFalse())
}

func Test_StrongestMatch__ShouldBeCachableIfMatchedOnEverythingApartFromStateZeroTimes(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "POST",
				},
			},
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "body",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "http",
				},
			},
			DeprecatedQuery: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "?foo=bar",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "/foo",
				},
			},
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "www.test.com",
				},
			},
			RequiresState: map[string]string{
				"foo": "bar",
			},
		},
		Response: testResponse,
	})

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "GET",
				},
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

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{map[string]string{"miss": "me"}, sync.RWMutex{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

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

	result = matching.MatchingStrategyRunner(r, false, simulation, &state.State{map[string]string{"miss": "me"}, sync.RWMutex{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"miss": {""},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
	}

	result = matching.MatchingStrategyRunner(r, false, simulation, &state.State{map[string]string{"miss": "me"}, sync.RWMutex{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

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

	result = matching.MatchingStrategyRunner(r, false, simulation, &state.State{map[string]string{"miss": "me"}, sync.RWMutex{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

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

	result = matching.MatchingStrategyRunner(r, false, simulation, &state.State{map[string]string{"miss": "me"}, sync.RWMutex{}}, &matching.StrongestMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())
}
