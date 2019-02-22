package matching_test

import (
	"sync"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/state"
	. "github.com/onsi/gomega"
)

var testResponse = models.ResponseDetails{
	Body: "request matched",
}

func Test_FirstMatchStrategy_EmptyRequestMatchersShouldMatchOnAnyRequest(t *testing.T) {
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
	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: make(map[string]string)}, &matching.FirstMatchStrategy{})

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchStrategy_RequestMatchersShouldMatchOnBody(t *testing.T) {
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
	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: make(map[string]string)}, &matching.FirstMatchStrategy{})
	Expect(result.Error).To(BeNil())

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchStrategy_ReturnResponseWhenAllHeadersMatch(t *testing.T) {
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

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: make(map[string]string)}, &matching.FirstMatchStrategy{})

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchStrategy_ReturnNilWhenOneHeaderNotPresentInRequest(t *testing.T) {
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

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: make(map[string]string)}, &matching.FirstMatchStrategy{})

	Expect(result.Pair).To(BeNil())
}

func Test_FirstMatchStrategy_ReturnNilWhenOneHeaderValueDifferent(t *testing.T) {
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
	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: make(map[string]string)}, &matching.FirstMatchStrategy{})

	Expect(result.Pair).To(BeNil())
}

func Test_FirstMatchStrategy_ReturnResponseWithMultiValuedHeaderMatch(t *testing.T) {
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
	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: make(map[string]string)}, &matching.FirstMatchStrategy{})

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchStrategy_ReturnNilWithDifferentMultiValuedHeaders(t *testing.T) {
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

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: make(map[string]string)}, &matching.FirstMatchStrategy{})

	Expect(result.Pair).To(BeNil())
}

func Test_FirstMatchStrategy_EndpointMatchWithHeaders(t *testing.T) {
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
	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: make(map[string]string)}, &matching.FirstMatchStrategy{})

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchStrategy_EndpointMismatchWithHeadersReturnsNil(t *testing.T) {
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

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: make(map[string]string)}, &matching.FirstMatchStrategy{})

	Expect(result.Pair).To(BeNil())
}

func Test_FirstMatchStrategy_AbleToMatchAnEmptyPathInAReasonableWay(t *testing.T) {
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
	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: make(map[string]string)}, &matching.FirstMatchStrategy{})

	Expect(result.Pair.Response.Body).To(Equal("request matched"))

	r = models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/a/1",
		Query: map[string][]string{
			"q": {"test"},
		},
	}

	result = matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: make(map[string]string)}, &matching.FirstMatchStrategy{})

	Expect(result.Pair).To(BeNil())
}

func Test_FirstMatchStrategy_RequestMatcherResponsePairCanBeConvertedToARequestResponsePairView_WhileIncomplete(t *testing.T) {
	RegisterTestingT(t)

	requestMatcherResponsePair := models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "POST",
				},
			},
		},
		Response: testResponse,
	}

	pairView := requestMatcherResponsePair.BuildView()

	Expect(pairView.RequestMatcher.Method).To(HaveLen(1))

	Expect(pairView.RequestMatcher.Method[0].Value).To(Equal("POST"))
	Expect(pairView.RequestMatcher.Destination).To(BeNil())
	Expect(pairView.RequestMatcher.Path).To(BeNil())
	Expect(pairView.RequestMatcher.Scheme).To(BeNil())
	Expect(pairView.RequestMatcher.Query).To(BeNil())
	Expect(pairView.RequestMatcher.Headers).To(HaveLen(0))

	Expect(pairView.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchStrategy_RequestMatchersCanUseGlobsAndBeMatched(t *testing.T) {
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

	result := matching.MatchingStrategyRunner(request, false, simulation, &state.State{State: make(map[string]string)}, &matching.FirstMatchStrategy{})
	Expect(result.Error).To(BeNil())

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchStrategy_RequestMatchersCanUseGlobsOnSchemeAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Glob,
					Value:   "H*",
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

	result := matching.MatchingStrategyRunner(request, false, simulation, &state.State{State: make(map[string]string)}, &matching.FirstMatchStrategy{})
	Expect(result.Error).To(BeNil())

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchStrategy_RequestMatchersCanUseGlobsOnHeadersAndBeMatched(t *testing.T) {
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

	result := matching.MatchingStrategyRunner(request, false, simulation, &state.State{State: make(map[string]string)}, &matching.FirstMatchStrategy{})
	Expect(result.Error).To(BeNil())

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchStrategy_RequestMatcherResponsePair_ConvertToRequestResponsePairView_CanBeConvertedToARequestResponsePairView_WhileIncomplete(t *testing.T) {
	RegisterTestingT(t)

	requestMatcherResponsePair := models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "POST",
				},
			},
		},
		Response: testResponse,
	}

	pairView := requestMatcherResponsePair.BuildView()

	Expect(pairView.RequestMatcher.Method[0].Matcher).To(Equal("exact"))
	Expect(pairView.RequestMatcher.Method[0].Value).To(Equal("POST"))
	Expect(pairView.RequestMatcher.Destination).To(BeNil())
	Expect(pairView.RequestMatcher.Path).To(BeNil())
	Expect(pairView.RequestMatcher.Scheme).To(BeNil())
	Expect(pairView.RequestMatcher.Query).To(BeNil())
	Expect(pairView.RequestMatcher.Headers).To(HaveLen(0))

	Expect(pairView.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchShouldNotBeCachableIfMatchedOnEverythingApartFromHeadersAtLeastOnce(t *testing.T) {
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

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: make(map[string]string)}, &matching.FirstMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeFalse())
}

func Test_FirstMatchShouldBeCachableIfMatchedOnEverythingApartFromHeadersZeroTimes(t *testing.T) {
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

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: make(map[string]string)}, &matching.FirstMatchStrategy{})

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

	result = matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: make(map[string]string)}, &matching.FirstMatchStrategy{})

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

	result = matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: make(map[string]string)}, &matching.FirstMatchStrategy{})

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

	result = matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: make(map[string]string)}, &matching.FirstMatchStrategy{})

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

	result = matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: make(map[string]string)}, &matching.FirstMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())
}

func Test_FirstMatchStrategy_RequestMatchersShouldMatchOnStateAndNotBeCachable(t *testing.T) {
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
		&state.State{State: map[string]string{"key1": "value1", "key2": "value2"}},
		&matching.FirstMatchStrategy{})

	Expect(result.Error).To(BeNil())
	Expect(result.Cachable).To(BeFalse())
	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchShouldNotBeCachableIfMatchedOnEverythingApartFromStateAtLeastOnce(t *testing.T) {
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

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{"miss": "me"}}, &matching.FirstMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeFalse())
}

func Test_FirstMatchShouldBeCachableIfMatchedOnEverythingApartFromStateZeroTimes(t *testing.T) {
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

	result := matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{"miss": "me"}}, &matching.FirstMatchStrategy{})

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

	result = matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{"miss": "me"}}, &matching.FirstMatchStrategy{})

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

	result = matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{"miss": "me"}}, &matching.FirstMatchStrategy{})

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

	result = matching.MatchingStrategyRunner(r, false, simulation, &state.State{State: map[string]string{"miss": "me"}}, &matching.FirstMatchStrategy{})

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

	result = matching.MatchingStrategyRunner(r, false, simulation, &state.State{map[string]string{"miss": "me"}, sync.RWMutex{}}, &matching.FirstMatchStrategy{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())
}
