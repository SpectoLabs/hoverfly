package v2_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

func Test_SimulationViewV2_Upgrade_UnescapesExactMatchRequestQueryParameters(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.SimulationViewV2{
		v2.DataViewV2{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV2{
				{
					RequestMatcher: v2.RequestMatcherViewV2{
						Query: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("q=10%20Downing%20Street%20London"),
						},
					},
					Response: v2.ResponseDetailsView{
						Status:      200,
						Body:        "body",
						EncodedBody: false,
						Headers: map[string][]string{
							"Test": []string{"headers"},
						},
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion:   "v1",
			HoverflyVersion: "test",
			TimeExported:    "today",
		},
	}

	simulationViewV3 := unit.Upgrade()

	Expect(simulationViewV3.RequestResponsePairs).To(HaveLen(1))
	Expect(*simulationViewV3.RequestResponsePairs[0].RequestMatcher.Query.ExactMatch).To(Equal("q=10 Downing Street London"))
}

func Test_SimulationViewV2_Upgrade_UnescapesGlobMatchRequestQueryParameters(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.SimulationViewV2{
		v2.DataViewV2{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV2{
				{
					RequestMatcher: v2.RequestMatcherViewV2{
						Query: &v2.RequestFieldMatchersView{
							GlobMatch: util.StringToPointer("q=*%20London"),
						},
					},
					Response: v2.ResponseDetailsView{
						Status:      200,
						Body:        "body",
						EncodedBody: false,
						Headers: map[string][]string{
							"Test": []string{"headers"},
						},
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion:   "v2",
			HoverflyVersion: "test",
			TimeExported:    "today",
		},
	}

	simulationViewV3 := unit.Upgrade()

	Expect(simulationViewV3.RequestResponsePairs).To(HaveLen(1))
	Expect(*simulationViewV3.RequestResponsePairs[0].RequestMatcher.Query.GlobMatch).To(Equal("q=* London"))
}

func Test_SimulationViewV2_Upgrade_KeepsEncodedResponsesEncoded(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.SimulationViewV2{
		v2.DataViewV2{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV2{
				{
					RequestMatcher: v2.RequestMatcherViewV2{
						Query: &v2.RequestFieldMatchersView{
							GlobMatch: util.StringToPointer("q=*%20London"),
						},
					},
					Response: v2.ResponseDetailsView{
						Status:      200,
						Body:        "YmFzZTY0IGVuY29kZWQ=",
						EncodedBody: true,
						Headers: map[string][]string{
							"Test": []string{"headers"},
						},
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion:   "v2",
			HoverflyVersion: "test",
			TimeExported:    "today",
		},
	}

	simulationViewV3 := unit.Upgrade()

	Expect(simulationViewV3.RequestResponsePairs).To(HaveLen(1))
	Expect(simulationViewV3.RequestResponsePairs[0].Response.EncodedBody).To(BeTrue())
	Expect(simulationViewV3.RequestResponsePairs[0].Response.Body).To(Equal("YmFzZTY0IGVuY29kZWQ="))
}
