package wrapper

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_SetPACFile_CanSetPACFile(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV6{
		v2.DataViewV6{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV6{
				{
					RequestMatcher: v2.RequestMatcherViewV6{
						Method: []v2.MatcherViewV6{
							{
								Matcher: matchers.Exact,
								Value:   "PUT",
							},
						},
						Path: []v2.MatcherViewV6{
							{
								Matcher: matchers.Exact,
								Value:   "/api/v2/hoverfly/pac",
							},
						},
					},
					Response: v2.ResponseDetailsViewV6{
						Status: 200,
						Body:   `PACFILE`,
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	err := SetPACFile(target)
	Expect(err).To(BeNil())
}

func Test_SetPACFile_ServerError(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV6{
		v2.DataViewV6{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV6{
				{
					RequestMatcher: v2.RequestMatcherViewV6{
						Method: []v2.MatcherViewV6{
							{
								Matcher: matchers.Exact,
								Value:   "PUT",
							},
						},
						Path: []v2.MatcherViewV6{
							{
								Matcher: matchers.Exact,
								Value:   "/api/v2/hoverfly/pac",
							},
						},
					},
					Response: v2.ResponseDetailsViewV6{
						Status: 400,
						Body:   `PACFILE`,
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	err := SetPACFile(target)
	Expect(err).To(Not(BeNil()))
}
