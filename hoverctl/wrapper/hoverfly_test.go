package wrapper

import (
	"fmt"
	"testing"

	"github.com/SpectoLabs/hoverfly/core"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
	. "github.com/onsi/gomega"
)

func Test_isLocal_WhenLocalhost(t *testing.T) {
	RegisterTestingT(t)

	Expect(IsLocal("localhost")).To(BeTrue())
}

func Test_isLocal_WhenLocalhost_WithHttp(t *testing.T) {
	RegisterTestingT(t)

	Expect(IsLocal("http://localhost")).To(BeTrue())
}

func Test_isLocal_WhenLocalhostIP(t *testing.T) {
	RegisterTestingT(t)

	Expect(IsLocal("127.0.0.1")).To(BeTrue())
}

func Test_isLocal_WhenLocalhostIP_WithHttp(t *testing.T) {
	RegisterTestingT(t)

	Expect(IsLocal("http://127.0.0.1")).To(BeTrue())
}

func Test_isLocal_WhenAnotherDNS(t *testing.T) {
	RegisterTestingT(t)

	Expect(IsLocal("specto.io")).To(BeFalse())
}

func Test_BuildUrl_AddsHostAdminPortAndPath(t *testing.T) {
	RegisterTestingT(t)

	target := configuration.Target{
		Host:      "http://localhost",
		AdminPort: 1234,
	}

	Expect(BuildURL(target, "/something")).To(Equal("http://localhost:1234/something"))
}

func Test_BuildUrl_AddsHostAdminPortAndPath_Https(t *testing.T) {
	RegisterTestingT(t)

	target := configuration.Target{
		Host:      "https://localhost",
		AdminPort: 1234,
	}

	Expect(BuildURL(target, "/something")).To(Equal("https://localhost:1234/something"))
}

func Test_BuildUrl_AddsHttpIfHostIsLocalhost(t *testing.T) {
	RegisterTestingT(t)

	target := configuration.Target{
		Host:      "localhost",
		AdminPort: 1234,
	}

	Expect(BuildURL(target, "/something")).To(Equal("http://localhost:1234/something"))
}

func Test_BuildUrl_AddsHttpIfHostIsExternal(t *testing.T) {
	RegisterTestingT(t)

	target := configuration.Target{
		Host:      "test-instance.hoverfly.io",
		AdminPort: 1234,
	}

	Expect(BuildURL(target, "/something")).To(Equal("https://test-instance.hoverfly.io:1234/something"))
}

func Test_SetMiddleware_ReturnsErrorIfAPIResponsesWithError(t *testing.T) {
	RegisterTestingT(t)

	hf := hoverfly.NewHoverfly()
	hf.Cfg.Webserver = true
	hf.StartProxy()

	target := configuration.Target{
		Host:      "localhost",
		AdminPort: 8500,
	}
	hf.PutSimulation(v2.SimulationViewV2{
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

	hf.StopProxy()
}
