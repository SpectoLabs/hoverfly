package wrapper

import (
	"encoding/json"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_ExportSimulation_GetsModeFromHoverfly(t *testing.T) {
	RegisterTestingT(t)

	responseBody := `{"simulation": true}`
	simulationList := v2.SimulationViewV5{
		DataViewV5: v2.DataViewV5{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV5{
				{
					RequestMatcher: v2.RequestMatcherViewV5{
						Method: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "GET",
							},
						},
						Path: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "/api/v2/simulation",
							},
						},
					},
					Response: v2.ResponseDetailsViewV5{
						Status: 200,
						Body:   responseBody,
					},
				},
			},
		},
		MetaView: v2.MetaView{
			SchemaVersion: "v2",
		},
	}

	simulationListBytes, err := json.Marshal(simulationList)
	Expect(err).To(BeNil())

	simulationList.RequestResponsePairs[0].Response.Body = string(simulationListBytes[:])
	hoverfly.ReplaceSimulation(simulationList)
	simulationList.RequestResponsePairs[0].Response.Body = responseBody

	view, err := ExportSimulation(target, "")
	Expect(err).To(BeNil())
	Expect(view).To(Equal(simulationList))
}

func Test_ExportSimulation_WithUrlPattern(t *testing.T) {
	RegisterTestingT(t)

	responseBody := `{"simulation": true}`
	simulationList := v2.SimulationViewV5{
		DataViewV5: v2.DataViewV5{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV5{
				{
					RequestMatcher: v2.RequestMatcherViewV5{
						Method: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "GET",
							},
						},
						Path: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "/api/v2/simulation",
							},
						},
						Query: &v2.QueryMatcherViewV5{
							"urlPattern": []v2.MatcherViewV5{
								{
									Matcher: matchers.Exact,
									Value:   "test-(.+).com",
								},
							},
						},
					},
					Response: v2.ResponseDetailsViewV5{
						Status: 200,
						Body:   responseBody,
					},
				},
			},
		},
		MetaView: v2.MetaView{
			SchemaVersion: "v2",
		},
	}

	simulationListBytes, err := json.Marshal(simulationList)
	Expect(err).To(BeNil())

	simulationList.RequestResponsePairs[0].Response.Body = string(simulationListBytes[:])
	hoverfly.ReplaceSimulation(simulationList)
	simulationList.RequestResponsePairs[0].Response.Body = responseBody

	view, err := ExportSimulation(target, "test-(.+).com")
	Expect(err).To(BeNil())
	Expect(view).To(Equal(simulationList))
}

func Test_ExportSimulation_ErrorsWhen_HoverflyNotAccessible(t *testing.T) {
	RegisterTestingT(t)

	_, err := ExportSimulation(inaccessibleTarget, "")

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not connect to Hoverfly at something:1234"))
}

func Test_ExportSimulation_ErrorsWhen_HoverflyReturnsNon200(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV5{
		DataViewV5: v2.DataViewV5{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV5{
				{
					RequestMatcher: v2.RequestMatcherViewV5{
						Method: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "GET",
							},
						},
						Path: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "/api/v2/simulation",
							},
						},
					},
					Response: v2.ResponseDetailsViewV5{
						Status: 400,
						Body:   "{\"error\":\"test error\"}",
					},
				},
			},
		},
		MetaView: v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	_, err := ExportSimulation(target, "")
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not retrieve simulation\n\ntest error"))
}

func Test_ImportSimulation_SendsCorrectHTTPRequest(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV5{
		DataViewV5: v2.DataViewV5{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV5{
				{
					RequestMatcher: v2.RequestMatcherViewV5{
						Method: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "PUT",
							},
						},
						Path: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "/api/v2/simulation",
							},
						},
						Body: []v2.MatcherViewV5{
							{
								Matcher: "json",
								Value:   `{"simulation": true}`,
							},
						},
					},
					Response: v2.ResponseDetailsViewV5{
						Status: 200,
						Body:   `{"simulation": true}`,
					},
				},
			},
		},
		MetaView: v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	err := ImportSimulation(target, `{"simulation": true}`)
	Expect(err).To(BeNil())
}

func Test_ImportSimulation_ErrorsWhen_HoverflyNotAccessible(t *testing.T) {
	RegisterTestingT(t)

	err := ImportSimulation(inaccessibleTarget, "")

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not connect to Hoverfly at something:1234"))
}

func Test_ImportSimulation_ErrorsWhen_HoverflyReturnsNon200(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV5{
		DataViewV5: v2.DataViewV5{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV5{
				{
					RequestMatcher: v2.RequestMatcherViewV5{
						Method: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "PUT",
							},
						},
						Path: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "/api/v2/simulation",
							},
						},
					},
					Response: v2.ResponseDetailsViewV5{
						Status: 400,
						Body:   "{\"error\":\"test error\"}",
					},
				},
			},
		},
		MetaView: v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	err := ImportSimulation(target, "")
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not import simulation\n\ntest error"))
}

func Test_AddSimulation_SendsCorrectHTTPRequest(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV5{
		DataViewV5: v2.DataViewV5{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV5{
				{
					RequestMatcher: v2.RequestMatcherViewV5{
						Method: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "POST",
							},
						},
						Path: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "/api/v2/simulation",
							},
						},
						Body: []v2.MatcherViewV5{
							{
								Matcher: "json",
								Value:   `{"simulation": true}`,
							},
						},
					},
					Response: v2.ResponseDetailsViewV5{
						Status: 200,
						Body:   `{"simulation": true}`,
					},
				},
			},
		},
		MetaView: v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	err := AddSimulation(target, `{"simulation": true}`)
	Expect(err).To(BeNil())
}

func Test_DeleteSimulations_SendsCorrectHTTPRequest(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV5{
		DataViewV5: v2.DataViewV5{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV5{
				{
					RequestMatcher: v2.RequestMatcherViewV5{
						Method: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "DELETE",
							},
						},
						Path: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "/api/v2/simulation",
							},
						},
					},
					Response: v2.ResponseDetailsViewV5{
						Status: 200,
						Body:   `{"simulation": true}`,
					},
				},
			},
		},
		MetaView: v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	err := DeleteSimulations(target)
	Expect(err).To(BeNil())
}

func Test_DeleteSimulations_ErrorsWhen_HoverflyNotAccessible(t *testing.T) {
	RegisterTestingT(t)

	err := DeleteSimulations(inaccessibleTarget)

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not connect to Hoverfly at something:1234"))
}

func Test_DeleteSimulations_ErrorsWhen_HoverflyReturnsNon200(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV5{
		DataViewV5: v2.DataViewV5{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV5{
				{
					RequestMatcher: v2.RequestMatcherViewV5{
						Method: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "DELETE",
							},
						},
						Path: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "/api/v2/simulation",
							},
						},
					},
					Response: v2.ResponseDetailsViewV5{
						Status: 400,
						Body:   "{\"error\":\"test error\"}",
					},
				},
			},
		},
		MetaView: v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	err := DeleteSimulations(target)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not delete simulation\n\ntest error"))
}
