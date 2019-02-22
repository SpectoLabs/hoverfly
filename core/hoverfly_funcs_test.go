package hoverfly

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
)

func Test_Hoverfly_DoRequest_DoesNotPanicWhenCannotMakeRequest(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	ioutil.NopCloser(bytes.NewBuffer([]byte("")))
	request, err := http.NewRequest("GET", "w.specto.fake", ioutil.NopCloser(bytes.NewBuffer([]byte(""))))
	Expect(err).To(BeNil())

	response, err := unit.DoRequest(request)
	Expect(response).To(BeNil())
	Expect(err).ToNot(BeNil())
}

func Test_Hoverfly_DoRequest_FailedHTTP(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	requestBody := []byte("fizz=buzz")

	body := ioutil.NopCloser(bytes.NewBuffer(requestBody))

	req, err := http.NewRequest("POST", "http://capture_body.com", body)
	Expect(err).To(BeNil())

	_, err = unit.DoRequest(req)
	Expect(err).ToNot(BeNil())
}

func Test_Hoverfly_GetResponse_CanReturnResponseFromCache(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.CacheMatcher.SaveRequestMatcherResponsePair(models.RequestDetails{
		Destination: "somehost.com",
		Method:      "POST",
		Scheme:      "http",
	}, &models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "somehost.com",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "POST",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "http",
				},
			},
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "cached response",
		},
	}, nil)

	response, err := unit.GetResponse(models.RequestDetails{
		Destination: "somehost.com",
		Method:      "POST",
		Scheme:      "http",
	})

	Expect(err).To(BeNil())
	Expect(response).ToNot(BeNil())

	Expect(response.Status).To(Equal(http.StatusOK))
	Expect(response.Body).To(Equal("cached response"))
}

func Test_Hoverfly_GetResponse_CanReturnResponseFromSimulationAndNotCache(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "somehost.com",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "POST",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "http",
				},
			},
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "response body",
		},
	})

	response, err := unit.GetResponse(models.RequestDetails{
		Destination: "somehost.com",
		Method:      "POST",
		Scheme:      "http",
	})

	Expect(err).To(BeNil())
	Expect(response).ToNot(BeNil())

	Expect(response.Status).To(Equal(http.StatusOK))
	Expect(response.Body).To(Equal("response body"))
}

func Test_Hoverfly_GetResponse_WillCacheResponseIfNotInCache(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "somehost.com",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "POST",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "http",
				},
			},
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "response body",
		},
	})

	unit.GetResponse(models.RequestDetails{
		Destination: "somehost.com",
		Method:      "POST",
		Scheme:      "http",
	})

	Expect(unit.CacheMatcher.RequestCache.RecordsCount()).Should(Equal(1))

	cachedRequestResponsePair, found := unit.CacheMatcher.RequestCache.Get("75b4ae6efa2a3f6d3ee6b9fed4d8c8c5")
	Expect(found).To(BeTrue())

	Expect(cachedRequestResponsePair.(*models.CachedResponse).MatchingPair.Response.Body).To(Equal("response body"))

	unit.Simulation = models.NewSimulation()
	response, err := unit.GetResponse(models.RequestDetails{
		Destination: "somehost.com",
		Method:      "POST",
		Scheme:      "http",
	})

	Expect(err).To(BeNil())
	Expect(response).ToNot(BeNil())

	Expect(response.Status).To(Equal(http.StatusOK))
	Expect(response.Body).To(Equal("response body"))
}

func Test_Hoverfly_GetResponse_WillReturnCachedResponseIfHeaderMatchIsFalse(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	requestDetails := models.RequestDetails{
		Destination: "somehost.com",
		Method:      "POST",
		Scheme:      "http",
	}

	unit.CacheMatcher.SaveRequestMatcherResponsePair(requestDetails, &models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{},
		Response: models.ResponseDetails{
			Body: "cached response",
		},
	}, nil)

	response, err := unit.GetResponse(requestDetails)
	Expect(err).To(BeNil())

	Expect(response.Body).To(Equal("cached response"))
}

func Test_Hoverfly_GetResponse_WillCacheMisses(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	requestDetails := models.RequestDetails{
		Destination: "somehost.com",
		Method:      "POST",
		Scheme:      "http",
	}

	_, err := unit.GetResponse(requestDetails)
	Expect(err.Error()).To(Equal("Could not find a match for request, create or record a valid matcher first!"))

	cachedResponse, matchingErr := unit.CacheMatcher.GetCachedResponse(&requestDetails)
	Expect(matchingErr).To(BeNil())

	Expect(cachedResponse.MatchingPair).To(BeNil())
	Expect(cachedResponse.ClosestMiss).To(BeNil())
}

func Test_Hoverfly_GetResponse_WillCacheClosestMiss(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	unit.PutSimulation(v2.SimulationViewV5{
		v2.DataViewV5{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV5{
				{
					RequestMatcher: v2.RequestMatcherViewV5{
						Method: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "closest",
							},
						},
					},
					Response: v2.ResponseDetailsViewV5{
						Body: "closest",
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v3",
		},
	})

	requestDetails := models.RequestDetails{
		Destination: "somehost.com",
		Method:      "POST",
		Scheme:      "http",
	}

	_, err := unit.GetResponse(requestDetails)
	Expect(err.Error()).ToNot(BeNil())

	cachedResponse, matchingErr := unit.CacheMatcher.GetCachedResponse(&requestDetails)
	Expect(matchingErr).To(BeNil())

	Expect(cachedResponse.MatchingPair).To(BeNil())
	Expect(cachedResponse.ClosestMiss.RequestMatcher.Method[0].Matcher).To(Equal("exact"))
	Expect(cachedResponse.ClosestMiss.RequestMatcher.Method[0].Value).To(Equal("closest"))

	Expect(cachedResponse.ClosestMiss.Response.Body).To(Equal("closest"))
	Expect(cachedResponse.ClosestMiss.MissedFields).To(ConsistOf("method"))
}

func Test_Hoverfly_GetResponse_WillCacheTemplateIfNotInCache(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "somehost.com",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "POST",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "http",
				},
			},
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "{{ randomUuid }}",
			Templated: true,
		},
	})

	unit.GetResponse(models.RequestDetails{
		Destination: "somehost.com",
		Method:      "POST",
		Scheme:      "http",
	})

	Expect(unit.CacheMatcher.RequestCache.RecordsCount()).Should(Equal(1))

	cachedRequestResponsePair, found := unit.CacheMatcher.RequestCache.Get("75b4ae6efa2a3f6d3ee6b9fed4d8c8c5")
	Expect(found).To(BeTrue())

	Expect(cachedRequestResponsePair.(*models.CachedResponse).MatchingPair.Response.Body).To(Equal("{{ randomUuid }}"))
	Expect(cachedRequestResponsePair.(*models.CachedResponse).ResponseTemplate).NotTo(BeNil())
}

func Test_Hoverfly_GetResponse_ShouldReturnEmptyTextIfResponseTemplateIsNotRenderable(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "somehost.com",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "POST",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "http",
				},
			},
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "hello {{ unknownFunc }}",
			Templated: true,
		},
	})

	response, err := unit.GetResponse(models.RequestDetails{
		Destination: "somehost.com",
		Method:      "POST",
		Scheme:      "http",
	})

	Expect(err).To(BeNil())
	Expect(response.Body).To(Equal("hello "))

	Expect(unit.CacheMatcher.RequestCache.RecordsCount()).Should(Equal(1))

	cachedRequestResponsePair, found := unit.CacheMatcher.RequestCache.Get("75b4ae6efa2a3f6d3ee6b9fed4d8c8c5")
	Expect(found).To(BeTrue())

	Expect(cachedRequestResponsePair.(*models.CachedResponse).MatchingPair.Response.Body).To(Equal("hello {{ unknownFunc }}"))
}

func Test_Hoverfly_GetResponse_TransitioningBetweenStatesWhenSimulating(t *testing.T) {
	RegisterTestingT(t)

	simulation := `{
		"data": {
			"pairs": [{
					"request": {
						"path": [
							{
								"matcher": "exact",
								"value": "/basket"
							}
						]
					},
					"response": {
						"status": 200,
						"body": "empty"
					}
				},
				{
					"request": {
						"path": [
							{
								"matcher": "exact",
								"value": "/basket"
							}
						],
						"requiresState": {
							"eggs": "present"
						}
					},
					"response": {
						"status": 200,
						"body": "eggs"
					}
				},
				{
					"request": {
						"path": [
							{
								"matcher": "exact",
								"value": "/basket"
							}
						],
						"requiresState": {
							"bacon": "present"
						}
					},
					"response": {
						"status": 200,
						"body": "bacon"
					}
				},
				{
					"request": {
						"path": [
							{
								"matcher": "exact",
								"value": "/basket"
							}
						],
						"requiresState": {
							"eggs": "present",
							"bacon": "present"
						}
					},
					"response": {
						"status": 200,
						"body": "eggs, bacon"
					}
				},
				{
					"request": {
						"path": [
							{
								"matcher": "exact",
								"value": "/add-eggs"
							}
						]
					},
					"response": {
						"status": 200,
						"body": "added eggs",
						"transitionsState": {
							"eggs": "present"
						}
					}
				},
				{
					"request": {
						"path": [
							{
								"matcher": "exact",
								"value": "/add-bacon"
							}
						]
					},
					"response": {
						"status": 200,
						"body": "added bacon",
						"transitionsState": {
							"bacon": "present"
						}
					}
				},
				{
					"request": {
						"path": [
							{
								"matcher": "exact",
								"value": "/remove-eggs"
							}
						]
					},
					"response": {
						"status": 200,
						"body": "removed eggs",
						"removesState": ["eggs"]
					}
				},
				{
					"request": {
						"path": [
							{
								"matcher": "exact",
								"value": "/remove-bacon"
							}
						]
					},
					"response": {
						"status": 200,
						"body": "removed bacon",
						"removesState": ["bacon"]
					}
				}
			],
			"globalActions": {
				"delays": []
			}
		},
		"meta": {
			"schemaVersion": "v5",
			"hoverflyVersion": "v0.10.2",
			"timeExported": "2017-02-23T12:43:48Z"
		}
	}`

	v5 := &v2.SimulationViewV5{}

	json.Unmarshal([]byte(simulation), v5)

	hoverfly := NewHoverfly()
	hoverfly.CacheMatcher = matching.CacheMatcher{
		RequestCache: cache.NewDefaultLRUCache(),
	}
	hoverfly.PutSimulation(*v5)

	hoverfly.SetModeWithArguments(v2.ModeView{Mode: "simulate"})

	response, _ := hoverfly.GetResponse(models.RequestDetails{
		Path: "/basket",
	})
	Expect(string(response.Body)).To(Equal(`empty`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/add-eggs",
	})
	Expect(string(response.Body)).To(Equal(`added eggs`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/basket",
	})
	Expect(string(response.Body)).To(Equal(`eggs`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/add-bacon",
	})
	Expect(string(response.Body)).To(Equal(`added bacon`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/basket",
	})
	Expect(string(response.Body)).To(Equal(`eggs, bacon`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/remove-eggs",
	})
	Expect(string(response.Body)).To(Equal(`removed eggs`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/basket",
	})
	Expect(string(response.Body)).To(Equal(`bacon`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/remove-bacon",
	})
	Expect(string(response.Body)).To(Equal(`removed bacon`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/basket",
	})
	Expect(string(response.Body)).To(Equal(`empty`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/basket",
	})
	Expect(string(response.Body)).To(Equal(`empty`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/add-eggs",
	})
	Expect(string(response.Body)).To(Equal(`added eggs`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/basket",
	})
	Expect(string(response.Body)).To(Equal(`eggs`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/add-bacon",
	})
	Expect(string(response.Body)).To(Equal(`added bacon`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/basket",
	})

	Expect(string(response.Body)).To(Equal(`eggs, bacon`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/remove-eggs",
	})
	Expect(string(response.Body)).To(Equal(`removed eggs`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/basket",
	})
	Expect(string(response.Body)).To(Equal(`bacon`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/remove-bacon",
	})
	Expect(string(response.Body)).To(Equal(`removed bacon`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/basket",
	})
	Expect(string(response.Body)).To(Equal(`empty`))
}

func Test_Hoverfly_GetResponse_GetNotRecordedRequest(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	request, err := http.NewRequest("POST", "http://capture_body.com", nil)
	Expect(err).To(BeNil())

	requestDetails, err := models.NewRequestDetailsFromHttpRequest(request)
	Expect(err).To(BeNil())

	response, err := unit.GetResponse(requestDetails)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not find a match for request, create or record a valid matcher first!"))

	Expect(response).To(BeNil())
}

func Test_Hoverfly_Save_SavesRequestAndResponseToSimulation(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Save(&models.RequestDetails{
		Body:        "testbody",
		Destination: "testdestination",
		Headers:     map[string][]string{"testheader": {"testvalue"}},
		Method:      "testmethod",
		Path:        "/testpath",
		Query: map[string][]string{
			"query": {"test"},
		},
		Scheme: "http",
	}, &models.ResponseDetails{
		Body:    "testresponsebody",
		Headers: map[string][]string{"testheader": {"testvalue"}},
		Status:  200,
	}, nil, false)

	Expect(unit.Simulation.GetMatchingPairs()).To(HaveLen(1))

	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Body).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Body[0].Matcher).To(Equal("exact"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Body[0].Value).To(Equal("testbody"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Destination).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Destination[0].Value).To(Equal("testdestination"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Method).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Method[0].Matcher).To(Equal("exact"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Method[0].Value).To(Equal("testmethod"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Path).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Path[0].Matcher).To(Equal("exact"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Path[0].Value).To(Equal("/testpath"))
	Expect(*unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Query).To(HaveLen(1))
	Expect(*unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Query).To(HaveKeyWithValue("query", []models.RequestFieldMatchers{
		{
			Matcher: matchers.Exact,
			Value:   "test",
		},
	}))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Scheme).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Scheme[0].Matcher).To(Equal("exact"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Scheme[0].Value).To(Equal("http"))

	Expect(unit.Simulation.GetMatchingPairs()[0].Response.Body).To(Equal("testresponsebody"))
	Expect(unit.Simulation.GetMatchingPairs()[0].Response.Headers).To(HaveKeyWithValue("testheader", []string{"testvalue"}))
	Expect(unit.Simulation.GetMatchingPairs()[0].Response.Status).To(Equal(200))
}


func Test_Hoverfly_Save_SavesRequestContainsMultiValueQuery(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	_ = unit.Save(&models.RequestDetails{
		Query: map[string][]string{
			"query": {"value1", "value2"},
		},
	}, &models.ResponseDetails{
		Body:    "testresponsebody",
		Headers: map[string][]string{"testheader": {"testvalue"}},
		Status:  200,
	}, nil, false)

	Expect(unit.Simulation.GetMatchingPairs()).To(HaveLen(1))

	Expect(*unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Query).To(HaveLen(1))
	Expect(*unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Query).To(HaveKeyWithValue("query", []models.RequestFieldMatchers{
		{
			Matcher: matchers.Exact,
			Value:   "value1;value2",
		},
	}))
}

func Test_Hoverfly_Save_DoesNotSaveRequestHeadersWhenGivenHeadersArrayIsNil(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Save(&models.RequestDetails{
		Headers: map[string][]string{"testheader": {"testvalue"}},
	}, &models.ResponseDetails{
		Body:    "testresponsebody",
		Headers: map[string][]string{"testheader": {"testvalue"}},
		Status:  200,
	}, nil, false)

	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers).To(BeEmpty())
}

func Test_Hoverfly_Save_SavesAllRequestHeadersWhenGivenAnAsterisk(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Save(&models.RequestDetails{
		Headers: map[string][]string{
			"testheader":  {"testvalue"},
			"testheader2": {"testvalue2"},
		},
	}, &models.ResponseDetails{
		Body:    "testresponsebody",
		Headers: map[string][]string{"testheader": {"testvalue"}},
		Status:  200,
	}, []string{"*"}, false)

	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers).To(HaveLen(2))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers["testheader"]).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers["testheader"][0]).To(Equal(models.RequestFieldMatchers{
		Matcher: "exact",
		Value:   "testvalue",
	}))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers["testheader2"]).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers["testheader2"][0]).To(Equal(models.RequestFieldMatchers{
		Matcher: "exact",
		Value:   "testvalue2",
	}))
}

func Test_Hoverfly_Save_SavesSpecificRequestHeadersWhenSpecifiedInHeadersArray(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Save(&models.RequestDetails{
		Headers: map[string][]string{
			"testheader":  {"testvalue"},
			"testheader2": {"testvalue2"},
		},
	}, &models.ResponseDetails{
		Body:    "testresponsebody",
		Headers: map[string][]string{"testheader": {"testvalue"}},
		Status:  200,
	}, []string{"testheader"}, false)

	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers["testheader"]).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers["testheader"][0]).To(Equal(models.RequestFieldMatchers{
		Matcher: "exact",
		Value:   "testvalue",
	}))
}

func Test_Hoverfly_Save_DoesNotSaveAnyRequestHeaderIfItDoesNotMatchEntryInHeadersArray(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Save(&models.RequestDetails{
		Headers: map[string][]string{
			"testheader":  {"testvalue"},
			"testheader2": {"testvalue2"},
		},
	}, &models.ResponseDetails{
		Body:    "testresponsebody",
		Headers: map[string][]string{"testheader": {"testvalue"}},
		Status:  200,
	}, []string{"nonmatch"}, false)

	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers).To(BeEmpty())
}

func Test_Hoverfly_Save_SavesMultipleRequestHeadersWhenMultiplesSpecifiedInHeadersArray(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Save(&models.RequestDetails{
		Headers: map[string][]string{
			"testheader":  {"testvalue"},
			"testheader2": {"testvalue2"},
			"nonmatch":    {"nonmatchvalue"},
		},
	}, &models.ResponseDetails{
		Body:    "testresponsebody",
		Headers: map[string][]string{"testheader": {"testvalue"}},
		Status:  200,
	}, []string{"testheader", "nonmatch"}, false)

	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers).To(HaveLen(2))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers["testheader"]).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers["testheader"][0]).To(Equal(models.RequestFieldMatchers{
		Matcher: "exact",
		Value:   "testvalue",
	}))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers["nonmatch"]).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers["nonmatch"][0]).To(Equal(models.RequestFieldMatchers{
		Matcher: "exact",
		Value:   "nonmatchvalue",
	}))
}

func Test_Hoverfly_Save_SavesIncompleteRequestAndResponseToSimulation(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Save(&models.RequestDetails{
		Destination: "testdestination",
	}, &models.ResponseDetails{
		Body:    "testresponsebody",
		Headers: map[string][]string{"testheader": {"testvalue"}},
		Status:  200,
	}, nil, false)

	Expect(unit.Simulation.GetMatchingPairs()).To(HaveLen(1))

	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Method).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Method[0]).To(Equal(models.RequestFieldMatchers{
		Matcher: "exact",
		Value:   "",
	}))

	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Scheme).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Scheme[0]).To(Equal(models.RequestFieldMatchers{
		Matcher: "exact",
		Value:   "",
	}))

	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Destination).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Destination[0]).To(Equal(models.RequestFieldMatchers{
		Matcher: "exact",
		Value:   "testdestination",
	}))

	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Path).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Path[0]).To(Equal(models.RequestFieldMatchers{
		Matcher: "exact",
		Value:   "",
	}))

	Expect(*unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Query).To(HaveLen(0))

	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers).To(HaveLen(0))

	Expect(unit.Simulation.GetMatchingPairs()[0].Response.Body).To(Equal("testresponsebody"))
	Expect(unit.Simulation.GetMatchingPairs()[0].Response.Headers).To(HaveKeyWithValue("testheader", []string{"testvalue"}))
	Expect(unit.Simulation.GetMatchingPairs()[0].Response.Status).To(Equal(200))
}

func Test_Hoverfly_Save_SavesRequestBodyAsJsonPathIfContentTypeIsJson(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Save(&models.RequestDetails{
		Body: `{"test": []}`,
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
	}, &models.ResponseDetails{}, nil, false)

	Expect(unit.Simulation.GetMatchingPairs()).To(HaveLen(1))

	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Body).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Body[0].Matcher).To(Equal("json"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Body[0].Value).To(Equal(`{"test": []}`))
}

func Test_Hoverfly_Save_SavesRequestBodyAsXmlPathIfContentTypeIsXml(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Save(&models.RequestDetails{
		Body: `<xml>`,
		Headers: map[string][]string{
			"Content-Type": {"application/xml"},
		},
	}, &models.ResponseDetails{}, nil, false)

	Expect(unit.Simulation.GetMatchingPairs()).To(HaveLen(1))

	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Body).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Body[0].Matcher).To(Equal("xml"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Body[0].Value).To(Equal(`<xml>`))
}

func Test_Hoverfly_Save_CanAddPairStatefully(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Save(&models.RequestDetails{
		Body: `body`,
	}, &models.ResponseDetails{}, nil, true)

	unit.Save(&models.RequestDetails{
		Body: `body`,
	}, &models.ResponseDetails{}, nil, true)

	Expect(unit.Simulation.GetMatchingPairs()).To(HaveLen(2))

	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.RequiresState).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.RequiresState["sequence:1"]).To(Equal("1"))

	Expect(unit.Simulation.GetMatchingPairs()[1].RequestMatcher.RequiresState).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[1].RequestMatcher.RequiresState["sequence:1"]).To(Equal("2"))
}
