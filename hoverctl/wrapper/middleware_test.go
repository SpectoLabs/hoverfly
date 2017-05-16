package wrapper

import (
	"fmt"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
	. "github.com/onsi/gomega"
)

func Test_SetMiddleware_ReturnsErrorIfAPIResponsesWithError(t *testing.T) {
	RegisterTestingT(t)

	target := configuration.Target{
		Host:      "localhost",
		AdminPort: 8500,
	}
	hoverfly.PutSimulation(v2.SimulationViewV2{
		v2.DataViewV2{
			RequestResponsePairs: []v2.RequestResponsePairViewV2{
				v2.RequestResponsePairViewV2{
					Request: v2.RequestDetailsViewV2{
						Path: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("/api/v2/hoverfly/middleware"),
						},
					},
					Response: v2.ResponseDetailsView{
						Status: 403,
						Body:   `{"error":"this is a middleware test error"}`,
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	_, err := SetMiddleware(target, "", "", "remote-middleware.com")
	fmt.Println(err.Error())
	Expect(err.Error()).To(ContainSubstring("Hoverfly could not execute this middleware"))
	Expect(err.Error()).To(ContainSubstring("this is a middleware test error"))
}
