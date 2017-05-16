package wrapper

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

func Test_GetLogs_GetsLogsWithCorrect_Text_Plain_AcceptHeader(t *testing.T) {
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
							ExactMatch: util.StringToPointer("/api/v2/logs"),
						},
						Headers: map[string][]string{
							"Accept": []string{
								"text/plain",
							},
						},
					},
					Response: v2.ResponseDetailsView{
						Status: 200,
						Body:   `logs line 1`,
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	logs, err := GetLogs(target, "plain")
	Expect(err).To(BeNil())
	Expect(logs[0]).To(Equal("logs line 1"))
}

func Test_GetLogs_GetsLogsWithCorrect_JSON_AcceptHeader(t *testing.T) {
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
							ExactMatch: util.StringToPointer("/api/v2/logs"),
						},
						Headers: map[string][]string{
							"Accept": []string{
								"application/json",
							},
						},
					},
					Response: v2.ResponseDetailsView{
						Status: 200,
						Body:   `{"logs":[{"msg": "logs line 1"}]}`,
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	logs, err := GetLogs(target, "json")
	Expect(err).To(BeNil())
	Expect(logs[0]).To(Equal(`{"msg":"logs line 1"}`))
}

func Test_GetLogs_ErrorsWhen_HoverflyNotAccessible(t *testing.T) {
	RegisterTestingT(t)

	_, err := GetLogs(inaccessibleTarget, "plain")

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not connect to Hoverfly at something:1234"))
}

func Test_GetLogs_ErrorsWhen_HoverflyReturnsNon200(t *testing.T) {
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
							ExactMatch: util.StringToPointer("/api/v2/logs"),
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

	_, err := GetLogs(target, "plain")
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not retrieve logs\n\ntest error"))
}
