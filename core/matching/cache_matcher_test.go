package matching_test

import (
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	"github.com/SpectoLabs/hoverfly/core/models"
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

	err := unit.PreloadCache(models.Simulation{})
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

	err := unit.PreloadCache(*simulation)

	Expect(err).To(BeNil())
	Expect(unit.RequestCache.RecordsCount()).To(Equal(0))
}

func Test_CacheMatcher_PreloadCache_WillPreemptivelyCacheFullExactMatchRequestMatchers(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{
		RequestCache: cache.NewDefaultLRUCache(),
	}

	simulation := models.NewSimulation()

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

	err := unit.PreloadCache(*simulation)

	Expect(err).To(BeNil())
	Expect(unit.RequestCache.RecordsCount()).To(Equal(1))
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

	err := unit.PreloadCache(*simulation)

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

	err := unit.PreloadCache(*simulation)

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

	err := unit.PreloadCache(*simulation)

	Expect(err).To(BeNil())
	Expect(unit.RequestCache.RecordsCount()).To(Equal(0))
}
