package hoverctl_end_to_end

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I use hoverctl", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	Describe("with a running hoverfly", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start()

			WriteConfiguration("localhost", hoverfly.GetAdminPort(), hoverfly.GetProxyPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		Context("I can delete the simulations in Hoverfly", func() {
			BeforeEach(func() {
				functional_tests.DoRequest(sling.New().Post(fmt.Sprintf("http://localhost:%v/api/records", hoverfly.GetAdminPort())).Body(strings.NewReader(`
					{
						"data": [{
							"request": {
								"path": "/api/bookings",
								"method": "POST",
								"destination": "www.my-test.com",
								"scheme": "http",
								"query": "",
								"body": "{\"flightId\": \"1\"}",
								"headers": {
									"Content-Type": [
										"application/json"
									]
								}
							},
							"response": {
								"status": 201,
								"body": "",
								"encodedBody": false,
								"headers": {
									"Location": [
										"http://localhost/api/bookings/1"
									]
								}
							}
						}]
					}`)))
			})

			It("and they should be removed", func() {
				output := functional_tests.Run(hoverctlBinary, "delete", "--force")

				Expect(output).To(ContainSubstring("Simulation data has been deleted from Hoverfly"))

				bytes, _ := ioutil.ReadAll(hoverfly.GetSimulation())
				Expect(string(bytes)).To(Equal(`{"data":null}`))
			})
		})

	})
})
