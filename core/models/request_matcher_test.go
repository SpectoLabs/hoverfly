package models_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
)

func Test_NewRequestFieldMatchersFromView_ReturnsNewStruct(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewRequestFieldMatchersFromView([]v2.MatcherViewV5{
		{
			Matcher: matchers.Exact,
			Value:   "exactly",
		},
	})

	Expect(unit).ToNot(BeNil())
	Expect(unit).To(HaveLen(1))
	Expect(unit[0].Matcher).To(Equal("exact"))
	Expect(unit[0].Value).To(Equal("exactly"))
}

func Test_NewRequestFieldMatchersFromView_WillReturnNilIfGivenNil(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewRequestFieldMatchersFromView(nil)

	Expect(unit).To(BeNil())
}

func Test_NewRequestFieldMatchers_BuildView(t *testing.T) {
	RegisterTestingT(t)

	unit := models.RequestFieldMatchers{
		Matcher: matchers.Exact,
		Value:   "exactly",
	}

	view := unit.BuildView()
	Expect(view.Matcher).To(Equal("exact"))
	Expect(view.Value).To(Equal("exactly"))
}

func Test_NewRequestMatcherResponsePairFromView_BuildsPair(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewRequestMatcherResponsePairFromView(&v2.RequestMatcherResponsePairViewV5{
		RequestMatcher: v2.RequestMatcherViewV5{
			Path: []v2.MatcherViewV5{
				{
					Matcher: matchers.Exact,
					Value:   "/",
				},
			},
			Headers: map[string][]v2.MatcherViewV5{
				"Header": {
					{
						Matcher: matchers.Exact,
						Value:   "header value",
					},
				},
			},
			Query: &v2.QueryMatcherViewV5{
				"Query": {
					{
						Matcher: matchers.Exact,
						Value:   "query value",
					},
				},
			},
		},
		Response: v2.ResponseDetailsViewV5{
			Body: "body",
		},
	})

	Expect(unit.RequestMatcher.Path).To(HaveLen(1))
	Expect(unit.RequestMatcher.Path[0].Matcher).To(Equal("exact"))
	Expect(unit.RequestMatcher.Path[0].Value).To(Equal("/"))
	Expect(unit.RequestMatcher.Headers["Header"][0].Matcher).To(Equal("exact"))
	Expect(unit.RequestMatcher.Headers["Header"][0].Value).To(Equal("header value"))
	Expect((*unit.RequestMatcher.Query)["Query"][0].Matcher).To(Equal("exact"))
	Expect((*unit.RequestMatcher.Query)["Query"][0].Value).To(Equal("query value"))
	Expect(unit.RequestMatcher.Destination).To(BeNil())

	Expect(unit.Response.Body).To(Equal("body"))
}

func Test_NewRequestMatcherResponsePairFromView_LeavesHeadersWithMatchersNil(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewRequestMatcherResponsePairFromView(&v2.RequestMatcherResponsePairViewV5{
		RequestMatcher: v2.RequestMatcherViewV5{
			Path: []v2.MatcherViewV5{
				{
					Matcher: matchers.Exact,
					Value:   "/",
				},
			},
		},
		Response: v2.ResponseDetailsViewV5{},
	})

	Expect(unit.RequestMatcher.Headers).To(BeNil())
}

func Test_NewRequestMatcherResponsePairFromView_LeavesQueriesWithMatchersNil(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewRequestMatcherResponsePairFromView(&v2.RequestMatcherResponsePairViewV5{
		RequestMatcher: v2.RequestMatcherViewV5{
			Path: []v2.MatcherViewV5{
				{
					Matcher: matchers.Exact,
					Value:   "/",
				},
			},
		},
		Response: v2.ResponseDetailsViewV5{},
	})

	Expect(unit.RequestMatcher.Query).To(BeNil())
}

func Test_NewRequestMatcherResponsePairFromView_SortsDeprecatedQuery(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewRequestMatcherResponsePairFromView(&v2.RequestMatcherResponsePairViewV5{
		RequestMatcher: v2.RequestMatcherViewV5{
			DeprecatedQuery: []v2.MatcherViewV5{
				{
					Matcher: matchers.Exact,
					Value:   "b=b&a=a",
				},
			},
		},
		Response: v2.ResponseDetailsViewV5{
			Body: "body",
		},
	})

	Expect(unit.RequestMatcher.DeprecatedQuery[0].Value).To(Equal("a=a&b=b"))
}

func Test_NewRequestMatcherResponsePairFromView_StoresTemplated(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewRequestMatcherResponsePairFromView(&v2.RequestMatcherResponsePairViewV5{
		RequestMatcher: v2.RequestMatcherViewV5{
			Path: []v2.MatcherViewV5{
				{
					Matcher: matchers.Exact,
					Value:   "/",
				},
			},
		},
		Response: v2.ResponseDetailsViewV5{
			Body:      "body",
			Templated: true,
		},
	})

	Expect(unit.Response.Templated).To(BeTrue())
}

func Test_RequestMatcher_BuildRequestDetailsFromExactMatches_GeneratesARequestDetails(t *testing.T) {
	RegisterTestingT(t)

	unit := models.RequestMatcher{
		Body: []models.RequestFieldMatchers{
			{
				Matcher: matchers.Exact,
				Value:   "body",
			},
		},
		Destination: []models.RequestFieldMatchers{
			{
				Matcher: matchers.Exact,
				Value:   "destination",
			},
		},
		Method: []models.RequestFieldMatchers{
			{
				Matcher: matchers.Exact,
				Value:   "method",
			},
		},
		Path: []models.RequestFieldMatchers{
			{
				Matcher: matchers.Exact,
				Value:   "path",
			},
		},
		DeprecatedQuery: []models.RequestFieldMatchers{
			{
				Matcher: matchers.Exact,
				Value:   "query=two",
			},
		},
		Scheme: []models.RequestFieldMatchers{
			{
				Matcher: matchers.Exact,
				Value:   "scheme",
			},
		},
	}

	Expect(unit.ToEagerlyCachable()).ToNot(BeNil())
	Expect(unit.ToEagerlyCachable()).To(Equal(&models.RequestDetails{
		Body:        "body",
		Destination: "destination",
		Method:      "method",
		Path:        "path",
		Query:       map[string][]string{"query": {"two"}},
		Scheme:      "scheme",
	}))
}

func Test_RequestMatcher_BuildRequestDetailsFromExactMatches_ReturnsNilIfEmpty(t *testing.T) {
	RegisterTestingT(t)

	unit := models.RequestMatcher{}

	Expect(unit.ToEagerlyCachable()).To(BeNil())
}

func Test_RequestMatcher_BuildRequestDetailsFromExactMatches_ReturnsNilIfMissingAnExactMatch(t *testing.T) {
	RegisterTestingT(t)

	unit := models.RequestMatcher{
		Destination: []models.RequestFieldMatchers{
			{
				Matcher: matchers.Exact,
				Value:   "destination",
			},
		},
		Method: []models.RequestFieldMatchers{
			{
				Matcher: matchers.Exact,
				Value:   "method",
			},
		},
		Path: []models.RequestFieldMatchers{
			{
				Matcher: matchers.Exact,
				Value:   "path",
			},
		},
		DeprecatedQuery: []models.RequestFieldMatchers{
			{
				Matcher: matchers.Exact,
				Value:   "query",
			},
		},
		Scheme: []models.RequestFieldMatchers{
			{
				Matcher: matchers.Exact,
				Value:   "query",
			},
		},
	}

	Expect(unit.ToEagerlyCachable()).To(BeNil())
}
