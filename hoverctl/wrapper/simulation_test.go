package wrapper

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

func Test_ExportSimulation_GetsModeFromHoverfly(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV2{
		v2.DataViewV2{
			RequestResponsePairs: []v2.RequestResponsePairViewV2{
				v2.RequestResponsePairViewV2{
					Request: v2.RequestDetailsViewV2{
						Method: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("GET"),
						},
						Path: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("/api/v2/simulation"),
						},
					},
					Response: v2.ResponseDetailsView{
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

	simulation, err := ExportSimulation(target)
	Expect(err).To(BeNil())

	Expect(string(simulation)).To(Equal("{\n\t\"simulation\": true\n}"))
}

func Test_ExportSimulation_ErrorsWhen_HoverflyNotAccessible(t *testing.T) {
	RegisterTestingT(t)

	_, err := ExportSimulation(inaccessibleTarget)

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not connect to Hoverfly at something:1234"))
}

func Test_ExportSimulation_ErrorsWhen_HoverflyReturnsNon200(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV2{
		v2.DataViewV2{
			RequestResponsePairs: []v2.RequestResponsePairViewV2{
				v2.RequestResponsePairViewV2{
					Request: v2.RequestDetailsViewV2{
						Method: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("GET"),
						},
						Path: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("/api/v2/simulation"),
						},
					},
					Response: v2.ResponseDetailsView{
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

	_, err := ExportSimulation(target)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not retrieve simulation\n\ntest error"))
}

func Test_ImportSimulation_SendsCorrectHTTPRequest(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV2{
		v2.DataViewV2{
			RequestResponsePairs: []v2.RequestResponsePairViewV2{
				v2.RequestResponsePairViewV2{
					Request: v2.RequestDetailsViewV2{
						Method: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("PUT"),
						},
						Path: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("/api/v2/simulation"),
						},
						Body: &v2.RequestFieldMatchersView{
							JsonMatch: util.StringToPointer(`{"simulation": true}`),
						},
					},
					Response: v2.ResponseDetailsView{
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
	hoverfly.PutSimulation(v2.SimulationViewV2{
		v2.DataViewV2{
			RequestResponsePairs: []v2.RequestResponsePairViewV2{
				v2.RequestResponsePairViewV2{
					Request: v2.RequestDetailsViewV2{
						Method: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("PUT"),
						},
						Path: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("/api/v2/simulation"),
						},
					},
					Response: v2.ResponseDetailsView{
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
	hoverfly.PutSimulation(v2.SimulationViewV2{
		v2.DataViewV2{
			RequestResponsePairs: []v2.RequestResponsePairViewV2{
				v2.RequestResponsePairViewV2{
					Request: v2.RequestDetailsViewV2{
						Method: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("DELETE"),
						},
						Path: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("/api/v2/simulation"),
						},
					},
					Response: v2.ResponseDetailsView{
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
	hoverfly.PutSimulation(v2.SimulationViewV2{
		v2.DataViewV2{
			RequestResponsePairs: []v2.RequestResponsePairViewV2{
				v2.RequestResponsePairViewV2{
					Request: v2.RequestDetailsViewV2{
						Method: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("DELETE"),
						},
						Path: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("/api/v2/simulation"),
						},
					},
					Response: v2.ResponseDetailsView{
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
