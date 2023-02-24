package matching_test

import (
	"github.com/SpectoLabs/hoverfly/v2/core/cache"
	"github.com/SpectoLabs/hoverfly/v2/core/matching"
	"github.com/SpectoLabs/hoverfly/v2/core/matching/matchers"
	"github.com/SpectoLabs/hoverfly/v2/core/models"
	. "github.com/onsi/gomega"
	"testing"
)

func Test_CacheMatcher_GetCachedResponse_WillReturnErrorIfCacheIsNil(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{}

	_, err := unit.GetCachedResponse(&models.RequestDetails{})
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("No cache set"))
}

func Test_CacheMatcher_GetAllResponses_WillReturnErrorIfCacheIsNil(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{}

	_, err := unit.GetAllResponses()
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("No cache set"))
}

func Test_CacheMatcher_SaveRequestMatcherResponsePair_WillReturnErrorIfCacheIsNil(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{}

	cachedResponse, err := unit.SaveRequestMatcherResponsePair(models.RequestDetails{}, nil, nil)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("No cache set"))
	Expect(cachedResponse).To(BeNil())
}

func Test_CacheMatcher_SaveRequestMatcherResponsePair_CanSaveNilPairs(t *testing.T) {
	RegisterTestingT(t)

	unit := matching.CacheMatcher{
		RequestCache: cache.NewDefaultLRUCache(),
	}

	cachedResponse, err := unit.SaveRequestMatcherResponsePair(models.RequestDetails{}, nil, nil)
	Expect(err).To(BeNil())

	Expect(cachedResponse.MatchingPair).To(BeNil())
}

func Test_CacheMatcher_FlushCache_WillReturnErrorIfCacheIsNil(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{}

	err := unit.FlushCache()
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("No cache set"))
}

func Test_CacheMatcher_PreloadCache_WillReturnErrorIfCacheIsNil(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{}

	simulation := models.Simulation{}
	err := unit.PreloadCache(&simulation)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("No cache set"))
}

func Test_CacheMatcher_PreloadCache_WillNotCacheIncompleteRequestMatchers(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{
		RequestCache: cache.NewDefaultLRUCache(),
	}

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Regex,
					Value:   "loose",
				},
			},
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "body",
		},
	})

	err := unit.PreloadCache(simulation)

	Expect(err).To(BeNil())
	Expect(unit.RequestCache.RecordsCount()).To(Equal(0))
}

func Test_CacheMatcher_PreloadCache_WillPreemptivelyCacheFullExactMatchRequestMatchers(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{
		RequestCache: cache.NewDefaultLRUCache(),
	}

	simulation := models.NewSimulation()

	pair1 := &models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
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
					Value:   "query",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "scheme",
				},
			},
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "body 1",
		},
	}

	pair2 := &models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
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
			Query: &models.QueryRequestFieldMatchers{
				"queryKey": {
					{
						Matcher: matchers.Exact,
						Value:   "queryValue",
					},
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "scheme",
				},
			},
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "body 2",
		},
	}
	simulation.AddPair(pair1)
	simulation.AddPair(pair2)

	err := unit.PreloadCache(simulation)

	Expect(err).To(BeNil())
	Expect(unit.RequestCache.RecordsCount()).To(Equal(2))

	cacheable1 := *pair1.RequestMatcher.ToEagerlyCacheable()
	cached1, _ := unit.RequestCache.Get(cacheable1.Hash())
	var cachedResponse1 *models.CachedResponse
	cachedResponse1 = cached1.(*models.CachedResponse)
	Expect(cachedResponse1.MatchingPair.Response.Body).To(Equal("body 1"))
	Expect(cachedResponse1.MatchingPair.RequestMatcher.Query).To(BeNil())

	cacheable2 := *pair2.RequestMatcher.ToEagerlyCacheable()
	cached2, _ := unit.RequestCache.Get(cacheable2.Hash())
	var cachedResponse2 *models.CachedResponse
	cachedResponse2 = cached2.(*models.CachedResponse)
	Expect(cachedResponse2.MatchingPair.Response.Body).To(Equal("body 2"))
	Expect(cachedResponse2.MatchingPair.RequestMatcher.Query.Get("queryKey")[0].Matcher).To(Equal(matchers.Exact))
	Expect(cachedResponse2.MatchingPair.RequestMatcher.Query.Get("queryKey")[0].Value).To(Equal("queryValue"))
}

func Test_CacheMatcher_PreloadCache_WillNotPreemptivelyCacheRequestMatchersWithoutExactMatches(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{
		RequestCache: cache.NewDefaultLRUCache(),
	}

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Regex,
					Value:   "destination",
				},
			},
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "body",
		},
	})

	err := unit.PreloadCache(simulation)

	Expect(err).To(BeNil())
	Expect(unit.RequestCache.RecordsCount()).To(Equal(0))
}

func Test_CacheMatcher_PreloadCache_WillCheckAllRequestMatchersInSimulation(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{
		RequestCache: cache.NewDefaultLRUCache(),
	}

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Regex,
					Value:   "destination",
				},
			},
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "body",
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
					Value:   "scheme",
				},
			},
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "body",
		},
	})

	err := unit.PreloadCache(simulation)

	Expect(err).To(BeNil())
	Expect(unit.RequestCache.RecordsCount()).To(Equal(1))
}

func Test_CacheMatcher_PreloadCache_WillNotCacheMatchersWithHeaders(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{
		RequestCache: cache.NewDefaultLRUCache(),
	}

	simulation := models.NewSimulation()

	simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Regex,
					Value:   "destination",
				},
			},
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "body",
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
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "destination",
				},
			},
			Headers: map[string][]models.RequestFieldMatchers{
				"Headers": {
					{
						Matcher: matchers.Exact,
						Value:   "value",
					},
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
					Value:   "scheme",
				},
			},
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "body",
		},
	})

	err := unit.PreloadCache(simulation)

	Expect(err).To(BeNil())
	Expect(unit.RequestCache.RecordsCount()).To(Equal(0))
}
