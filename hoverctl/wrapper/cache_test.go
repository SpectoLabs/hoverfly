package wrapper

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_FlushCache_GetsMiddlewareFromHoverfly(t *testing.T) {
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
								Value:   "DELETE",
							},
						},
						Path: []v2.MatcherViewV6{
							{
								Matcher: matchers.Exact,
								Value:   "/api/v2/cache",
							},
						},
					},
					Response: v2.ResponseDetailsViewV6{
						Status: 200,
						Body:   `{"binary": "test-binary", "script": "test.script", "remote": "http://test.com"}`,
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	}, false)

	err := FlushCache(target)
	Expect(err).To(BeNil())
}

func Test_FlushCache_ErrorsWhen_HoverflyNotAccessible(t *testing.T) {
	RegisterTestingT(t)

	err := FlushCache(inaccessibleTarget)

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not connect to Hoverfly at something:1234"))
}

func Test_FlushCache_ErrorsWhen_HoverflyReturnsNon200(t *testing.T) {
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
								Value:   "DELETE",
							},
						},
						Path: []v2.MatcherViewV6{
							{
								Matcher: matchers.Exact,
								Value:   "/api/v2/cache",
							},
						},
					},
					Response: v2.ResponseDetailsViewV6{
						Status: 400,
						Body:   "{\"error\":\"test error\"}",
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	}, false)

	err := FlushCache(target)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not flush cache\n\ntest error"))
}
