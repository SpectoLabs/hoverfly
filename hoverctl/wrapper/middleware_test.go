package wrapper

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

func Test_GetMiddleware_GetsMiddlewareFromHoverfly(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV4{
		v2.DataViewV4{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV4{
				v2.RequestMatcherResponsePairViewV4{
					RequestMatcher: v2.RequestMatcherViewV4{
						Method: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("GET"),
						},
						Path: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("/api/v2/hoverfly/middleware"),
						},
					},
					Response: v2.ResponseDetailsViewV4{
						Status: 200,
						Body:   `{"binary": "test-binary", "script": "test.script", "remote": "http://test.com"}`,
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	response, err := GetMiddleware(target)
	Expect(err).To(BeNil())

	Expect(response.Binary).To(Equal("test-binary"))
	Expect(response.Script).To(Equal("test.script"))
	Expect(response.Remote).To(Equal("http://test.com"))
}

func Test_GetMiddleware_ErrorsWhen_HoverflyNotAccessible(t *testing.T) {
	RegisterTestingT(t)

	_, err := GetMiddleware(inaccessibleTarget)

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not connect to Hoverfly at something:1234"))
}

func Test_GetMiddleware_ErrorsWhen_HoverflyReturnsNon200(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV4{
		v2.DataViewV4{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV4{
				v2.RequestMatcherResponsePairViewV4{
					RequestMatcher: v2.RequestMatcherViewV4{
						Method: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("GET"),
						},
						Path: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("/api/v2/hoverfly/middleware"),
						},
					},
					Response: v2.ResponseDetailsViewV4{
						Status: 400,
						Body:   `{"error": "test error"}`,
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	_, err := GetMiddleware(target)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not retrieve middleware\n\ntest error"))
}

func Test_SetMiddleware_ErrorsWhen_HoverflyNotAccessible(t *testing.T) {
	RegisterTestingT(t)

	_, err := SetMiddleware(inaccessibleTarget, "", "", "")

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not connect to Hoverfly at something:1234"))
}

func Test_SetMiddleware_ErrorsWhen_HoverflyReturnsNon200(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV4{
		v2.DataViewV4{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV4{
				v2.RequestMatcherResponsePairViewV4{
					RequestMatcher: v2.RequestMatcherViewV4{
						Method: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("PUT"),
						},
						Path: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("/api/v2/hoverfly/middleware"),
						},
					},
					Response: v2.ResponseDetailsViewV4{
						Status: 400,
						Body:   `{"error": "test error"}`,
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	_, err := SetMiddleware(target, "", "", "")
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not set middleware\n\ntest error"))
}
