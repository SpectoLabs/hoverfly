package models_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
)

func Test_NewRequestFieldMatchersFromView_ReturnsNewStruct(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewRequestFieldMatchersFromView([]v2.MatcherViewV5{
		{
			Matcher: "exact",
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
		Matcher: "exact",
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
					Matcher: "exact",
					Value:   "/",
				},
			},
			HeadersWithMatchers: map[string][]v2.MatcherViewV5{
				"Header": {
					{
						Matcher: "exact",
						Value:   "header value",
					},
				},
			},
			QueriesWithMatchers: map[string][]v2.MatcherViewV5{
				"Query": {
					{
						Matcher: "exact",
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
	Expect(unit.RequestMatcher.HeadersWithMatchers["Header"][0].Matcher).To(Equal("exact"))
	Expect(unit.RequestMatcher.HeadersWithMatchers["Header"][0].Value).To(Equal("header value"))
	Expect(unit.RequestMatcher.QueriesWithMatchers["Query"][0].Matcher).To(Equal("exact"))
	Expect(unit.RequestMatcher.QueriesWithMatchers["Query"][0].Value).To(Equal("query value"))
	Expect(unit.RequestMatcher.Destination).To(BeNil())

	Expect(unit.Response.Body).To(Equal("body"))
}

func Test_NewRequestMatcherResponsePairFromView_LeavesHeadersWithMatchersNil(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewRequestMatcherResponsePairFromView(&v2.RequestMatcherResponsePairViewV5{
		RequestMatcher: v2.RequestMatcherViewV5{
			Path: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "/",
				},
			},
		},
		Response: v2.ResponseDetailsViewV5{},
	})

	Expect(unit.RequestMatcher.HeadersWithMatchers).To(BeNil())
}

func Test_NewRequestMatcherResponsePairFromView_LeavesQueriesWithMatchersNil(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewRequestMatcherResponsePairFromView(&v2.RequestMatcherResponsePairViewV5{
		RequestMatcher: v2.RequestMatcherViewV5{
			Path: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "/",
				},
			},
		},
		Response: v2.ResponseDetailsViewV5{},
	})

	Expect(unit.RequestMatcher.QueriesWithMatchers).To(BeNil())
}

func Test_NewRequestMatcherResponsePairFromView_SortsQuery(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewRequestMatcherResponsePairFromView(&v2.RequestMatcherResponsePairViewV5{
		RequestMatcher: v2.RequestMatcherViewV5{
			Query: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "b=b&a=a",
				},
			},
		},
		Response: v2.ResponseDetailsViewV5{
			Body: "body",
		},
	})

	Expect(unit.RequestMatcher.Query[0].Value).To(Equal("a=a&b=b"))
}

func Test_NewRequestMatcherResponsePairFromView_StoresTemplated(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewRequestMatcherResponsePairFromView(&v2.RequestMatcherResponsePairViewV5{
		RequestMatcher: v2.RequestMatcherViewV5{
			Path: []v2.MatcherViewV5{
				{
					Matcher: "exact",
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

// func Test_RequestMatcher_BuildRequestDetailsFromExactMatches_GeneratesARequestDetails(t *testing.T) {
// 	RegisterTestingT(t)

// 	unit := models.RequestMatcher{
// 		Body: []models.RequestFieldMatchers{
// 			{
// 				Matcher: "exact",
// 				Value:   "body",
// 			},
// 		},
// 		Destination: []models.RequestFieldMatchers{
// 			{
// 				Matcher: "exact",
// 				Value:   "destination",
// 			},
// 		},
// 		Method: []models.RequestFieldMatchers{
// 			{
// 				Matcher: "exact",
// 				Value:   "method",
// 			},
// 		},
// 		Path: []models.RequestFieldMatchers{
// 			{
// 				Matcher: "exact",
// 				Value:   "path",
// 			},
// 		},
// 		Query: []models.RequestFieldMatchers{
// 			{
// 				Matcher: "exact",
// 				Value:   "query=two",
// 			},
// 		},
// 		Scheme: []models.RequestFieldMatchers{
// 			{
// 				Matcher: "exact",
// 				Value:   "scheme",
// 			},
// 		},
// 	}

// 	Expect(unit.ToEagerlyCachable()).ToNot(BeNil())
// 	Expect(unit.ToEagerlyCachable()).To(Equal(&models.RequestDetails{
// 		Body:        "body",
// 		Destination: "destination",
// 		Method:      "method",
// 		Path:        "path",
// 		Query:       map[string][]string{"query": []string{"two"}},
// 		Scheme:      "scheme",
// 	}))
// }

// func Test_RequestMatcher_BuildRequestDetailsFromExactMatches_ReturnsNilIfEmpty(t *testing.T) {
// 	RegisterTestingT(t)

// 	unit := models.RequestMatcher{}

// 	Expect(unit.ToEagerlyCachable()).To(BeNil())
// }

// func Test_RequestMatcher_BuildRequestDetailsFromExactMatches_ReturnsNilIfMissingAnExactMatch(t *testing.T) {
// 	RegisterTestingT(t)

// 	unit := models.RequestMatcher{
// 		Destination: []models.RequestFieldMatchers{
// 			{
// 				Matcher: "exact",
// 				Value:   "destination",
// 			},
// 		},
// 		Method: []models.RequestFieldMatchers{
// 			{
// 				Matcher: "exact",
// 				Value:   "method",
// 			},
// 		},
// 		Path: []models.RequestFieldMatchers{
// 			{
// 				Matcher: "exact",
// 				Value:   "path",
// 			},
// 		},
// 		Query: []models.RequestFieldMatchers{
// 			{
// 				Matcher: "exact",
// 				Value:   "query",
// 			},
// 		},
// 		Scheme: []models.RequestFieldMatchers{
// 			{
// 				Matcher: "exact",
// 				Value:   "query",
// 			},
// 		},
// 	}

// 	Expect(unit.ToEagerlyCachable()).To(BeNil())
// }
