package hoverctl_end_to_end

import (
	"fmt"
	"io/ioutil"
	"os"
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
				output := functional_tests.Run(hoverctlBinary, "delete", "simulations")

				Expect(output).To(ContainSubstring("Simulations have been deleted from Hoverfly"))

				bytes, _ := ioutil.ReadAll(hoverfly.GetSimulation())
				Expect(string(bytes)).To(Equal(`{"data":null}`))
			})
		})

		Context("I can delete the delays in Hoverfly", func() {
			BeforeEach(func() {
				functional_tests.DoRequest(sling.New().Put(fmt.Sprintf("http://localhost:%v/api/delays", hoverfly.GetAdminPort())).Body(strings.NewReader(`
						{
							"data": [
								{
									"hostPattern": "virtual\\.com",
									"delay": 100
								},
								{
									"hostPattern": "virtual\\.com",
									"delay": 110
								}
							]`)))
			})

			It("and they should be removed", func() {
				output := functional_tests.Run(hoverctlBinary, "delete", "delays")

				Expect(output).To(ContainSubstring("Delays have been deleted from Hoverfly"))

				resp := functional_tests.DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/delays", hoverfly.GetAdminPort())))
				bytes, _ := ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).To(Equal(`{"data":[]}`))
			})
		})

		Context("I can delete the middleware in Hoverfly", func() {
			BeforeEach(func() {
				functional_tests.Run(hoverctlBinary, "middleware", "python testdata/add_random_delay.py")
			})

			It("and they should be removed", func() {
				output := functional_tests.Run(hoverctlBinary, "delete", "middleware")

				Expect(output).To(ContainSubstring("Middleware has been deleted from Hoverfly"))

				resp := functional_tests.DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/v2/hoverfly/middleware", hoverfly.GetAdminPort())))
				bytes, _ := ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).To(Equal(`{"binary":"","script":"","remote":""}`))
			})
		})

		Context("I can delete the request templates in Hoverfly", func() {
			BeforeEach(func() {
				fileReader, err := os.Open("testdata/request-template.json")
				defer fileReader.Close()
				if err != nil {
					Fail("Failed to read request template test data")
				}
				resp := functional_tests.DoRequest(sling.New().Post(fmt.Sprintf("http://localhost:%v/api/templates", hoverfly.GetAdminPort())).Body(fileReader))
				bytes, _ := ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).To(ContainSubstring(`{"message":"2 payloads import complete."}`))
			})

			It("and they should be removed", func() {
				output := functional_tests.Run(hoverctlBinary, "delete", "templates")

				Expect(output).To(ContainSubstring("Request templates have been deleted from Hoverfly"))

				resp := functional_tests.DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/templates", hoverfly.GetAdminPort())))
				bytes, _ := ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).To(ContainSubstring(`{"data":null}`))
			})
		})

		Context("I can delete everything in hoverfly", func() {

			BeforeEach(func() {
				functional_tests.DoRequest(sling.New().Put(fmt.Sprintf("http://localhost:%v/api/delays", hoverfly.GetAdminPort())).Body(strings.NewReader(`
					{
						"data": [
							{
								"hostPattern": "virtual\\.com",
								"delay": 100
							},
							{
								"hostPattern": "virtual\\.com",
								"delay": 110
							}
						]`)))

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
				fileReader, err := os.Open("testdata/request-template.json")
				defer fileReader.Close()
				if err != nil {
					Fail("Failed to read request template test data")
				}
				resp := functional_tests.DoRequest(sling.New().Post(fmt.Sprintf("http://localhost:%v/api/templates", hoverfly.GetAdminPort())).Body(fileReader))
				bytes, _ := ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).To(ContainSubstring(`{"message":"2 payloads import complete."}`))

				functional_tests.Run(hoverctlBinary, "middleware", "python testdata/add_random_delay.py")
			})

			It("by calling delete all", func() {
				output := functional_tests.Run(hoverctlBinary, "delete", "all")
				Expect(output).To(ContainSubstring("Delays, middleware, request templates and simulations have all been deleted from Hoverfly"))

				resp := functional_tests.DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/delays", hoverfly.GetAdminPort())))
				bytes, _ := ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).To(Equal(`{"data":[]}`))

				resp = functional_tests.DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/records", hoverfly.GetAdminPort())))
				bytes, _ = ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).To(Equal(`{"data":null}`))

				resp = functional_tests.DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/templates", hoverfly.GetAdminPort())))
				bytes, _ = ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).To(ContainSubstring(`{"data":null}`))

				resp = functional_tests.DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/v2/hoverfly/middleware", hoverfly.GetAdminPort())))
				bytes, _ = ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).To(Equal(`{"binary":"","script":"","remote":""}`))
			})
		})

		Context("I won't delete if I have not specified what to delete", func() {
			It("when I call hoverctl delete with no resource", func() {
				output := functional_tests.Run(hoverctlBinary, "delete")
				Expect(output).To(ContainSubstring("You have not specified a resource to delete from Hoverfly"))
			})

			It("when I call hoverctl delete with an invalid resource", func() {
				output := functional_tests.Run(hoverctlBinary, "delete", "test")
				Expect(output).To(ContainSubstring("You have not specified a valid resource to delete from Hoverfly"))
			})
		})

	})
})
