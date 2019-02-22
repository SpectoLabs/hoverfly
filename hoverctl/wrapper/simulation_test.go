package wrapper

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_ExportSimulation_GetsModeFromHoverfly(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV5{
		v2.DataViewV5{
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
						Body:   `{"simulation": true}`,
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	simulation, err := ExportSimulation(target, "")
	Expect(err).To(BeNil())

	Expect(string(simulation)).To(Equal("{\n\t\"simulation\": true\n}"))
}

func Test_ExportSimulation_WithUrlPattern(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV5{
		v2.DataViewV5{
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
						Body:   `{"simulation": true}`,
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	simulation, err := ExportSimulation(target, "test-(.+).com")
	Expect(err).To(BeNil())

	Expect(string(simulation)).To(Equal("{\n\t\"simulation\": true\n}"))
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
		v2.DataViewV5{
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
		v2.MetaView{
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
		v2.DataViewV5{
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
		v2.MetaView{
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
		v2.DataViewV5{
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
		v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	err := ImportSimulation(target, "")
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not import simulation\n\ntest error"))
}

func Test_DeleteSimulations_SendsCorrectHTTPRequest(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV5{
		v2.DataViewV5{
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
		v2.MetaView{
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
		v2.DataViewV5{
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
		v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	err := DeleteSimulations(target)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not delete simulation\n\ntest error"))
}
