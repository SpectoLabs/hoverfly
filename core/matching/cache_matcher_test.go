package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
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

	err := unit.SaveRequestMatcherResponsePair(models.RequestDetails{}, nil, nil)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("No cache set"))
}

func Test_CacheMatcher_SaveRequestMatcherResponsePair_CanSaveNilPairs(t *testing.T) {
	RegisterTestingT(t)

	unit := matching.CacheMatcher{
		RequestCache: cache.NewInMemoryCache(),
	}

	err := unit.SaveRequestMatcherResponsePair(models.RequestDetails{}, nil, nil)
	Expect(err).To(BeNil())

	cacheValues, err := unit.RequestCache.Get([]byte("d41d8cd98f00b204e9800998ecf8427e"))
	Expect(err).To(BeNil())

	cachedResponse, err := models.NewCachedResponseFromBytes(cacheValues)
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

// func Test_CacheMatcher_PreloadCache_WillNotCacheIncompleteRequestMatchers(t *testing.T) {
// 	RegisterTestingT(t)
// 	unit := matching.CacheMatcher{
// 		RequestCache: cache.NewInMemoryCache(),
// 	}

// 	simulation := models.NewSimulation()

// 	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
// 		RequestMatcher: models.RequestMatcher{
// 			Body: []models.RequestFieldMatchers{
// 				{
// 					Matcher: "regex",
// 					Value:   "loose",
// 				},
// 			},
// 		},
// 		Response: models.ResponseDetails{
// 			Status: 200,
// 			Body:   "body",
// 		},
// 	})

// 	err := unit.PreloadCache(*simulation)

// 	Expect(err).To(BeNil())
// 	Expect(unit.RequestCache.GetAllKeys()).To(HaveLen(0))
// }

func Test_CacheMatcher_PreloadCache_WillPreemptivelyCacheFullExactMatchRequestMatchers(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{
		RequestCache: cache.NewInMemoryCache(),
	}

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "body",
				},
			},

			Destination: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "destination",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "method",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "path",
				},
			},
			Query: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "query",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
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
	Expect(unit.RequestCache.GetAllKeys()).To(HaveLen(1))
}

// func Test_CacheMatcher_PreloadCache_WillNotPreemptivelyCacheRequestMatchersWithoutExactMatches(t *testing.T) {
// 	RegisterTestingT(t)
// 	unit := matching.CacheMatcher{
// 		RequestCache: cache.NewInMemoryCache(),
// 	}

// 	simulation := models.NewSimulation()

// 	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
// 		RequestMatcher: models.RequestMatcher{
// 			Destination: []models.RequestFieldMatchers{
// 				{
// 					Matcher: "regex",
// 					Value:   "destination",
// 				},
// 			},
// 		},
// 		Response: models.ResponseDetails{
// 			Status: 200,
// 			Body:   "body",
// 		},
// 	})

// 	err := unit.PreloadCache(*simulation)

// 	Expect(err).To(BeNil())
// 	Expect(unit.RequestCache.GetAllKeys()).To(HaveLen(0))
// }

func Test_CacheMatcher_PreloadCache_WillCheckAllRequestMatchersInSimulation(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{
		RequestCache: cache.NewInMemoryCache(),
	}

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: "regex",
					Value:   "destination",
				},
			},
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "body",
		},
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "body",
				},
			},
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "destination",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "method",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "path",
				},
			},
			Query: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "query",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
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
	Expect(unit.RequestCache.GetAllKeys()).To(HaveLen(1))
}

// func Test_CacheMatcher_PreloadCache_WillNotCacheMatchersWithHeaders(t *testing.T) {
// 	RegisterTestingT(t)
// 	unit := matching.CacheMatcher{
// 		RequestCache: cache.NewInMemoryCache(),
// 	}

// 	simulation := models.NewSimulation()

// 	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
// 		RequestMatcher: models.RequestMatcher{
// 			Destination: []models.RequestFieldMatchers{
// 				{
// 					Matcher: "regex",
// 					Value:   "destination",
// 				},
// 			},
// 		},
// 		Response: models.ResponseDetails{
// 			Status: 200,
// 			Body:   "body",
// 		},
// 	})

// 	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
// 		RequestMatcher: models.RequestMatcher{
// 			Body: []models.RequestFieldMatchers{
// 				{
// 					Matcher: "exact",
// 					Value:   "body",
// 				},
// 			},
// 			Destination: []models.RequestFieldMatchers{
// 				{
// 					Matcher: "exact",
// 					Value:   "destination",
// 				},
// 			},
// 			Headers: map[string][]string{
// 				"Headers": {"value"},
// 			},
// 			Method: []models.RequestFieldMatchers{
// 				{
// 					Matcher: "exact",
// 					Value:   "method",
// 				},
// 			},
// 			Path: []models.RequestFieldMatchers{
// 				{
// 					Matcher: "exact",
// 					Value:   "path",
// 				},
// 			},
// 			Query: []models.RequestFieldMatchers{
// 				{
// 					Matcher: "exact",
// 					Value:   "query",
// 				},
// 			},
// 			Scheme: []models.RequestFieldMatchers{
// 				{
// 					Matcher: "exact",
// 					Value:   "scheme",
// 				},
// 			},
// 		},
// 		Response: models.ResponseDetails{
// 			Status: 200,
// 			Body:   "body",
// 		},
// 	})

// 	err := unit.PreloadCache(*simulation)

// 	Expect(err).To(BeNil())
// 	Expect(unit.RequestCache.GetAllKeys()).To(HaveLen(0))
// }
