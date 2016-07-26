package hoverfly_end_to_end_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os/exec"
	"strings"
	"strconv"
	"github.com/phayes/freeport"
	"fmt"
	"io/ioutil"
	"github.com/dghubble/sling"
	"os"
)

var _ = Describe("When I use hoverctl", func() {
	var (
		hoverflyCmd *exec.Cmd

		adminPort = freeport.GetPort()
		adminPortAsString = strconv.Itoa(adminPort)

		proxyPort = freeport.GetPort()
		proxyPortAsString = strconv.Itoa(proxyPort)
	)

	Describe("with a running hoverfly", func() {

		BeforeEach(func() {
			hoverflyCmd = startHoverfly(adminPort, proxyPort, workingDirectory)
			WriteConfiguration("localhost", adminPortAsString, proxyPortAsString)
		})

		AfterEach(func() {
			hoverflyCmd.Process.Kill()
		})

		Context("I can delete the simulations in Hoverfly", func() {
			BeforeEach(func() {
				DoRequest(sling.New().Post(fmt.Sprintf("http://localhost:%v/api/records", adminPort)).Body(strings.NewReader(`
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
				out, _ := exec.Command(hoverctlBinary, "delete", "simulations").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Simulations have been deleted from Hoverfly"))

				resp := DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/records", adminPort)))
				bytes, _ := ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).To(Equal(`{"data":null}`))
			})
		})

		Context("I can delete the delays in Hoverfly", func() {
			BeforeEach(func() {
				DoRequest(sling.New().Put(fmt.Sprintf("http://localhost:%v/api/delays", adminPort)).Body(strings.NewReader(`
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
				out, _ := exec.Command(hoverctlBinary, "delete", "delays").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Delays have been deleted from Hoverfly"))

				resp := DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/delays", adminPort)))
				bytes, _ := ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).To(Equal(`{"data":[]}`))
			})
		})

		Context("I can delete the middleware in Hoverfly", func() {
			BeforeEach(func() {
				exec.Command(hoverctlBinary, "middleware", "python testdata/add_random_delay.py").Output()
			})

			It("and they should be removed", func() {
				out, _ := exec.Command(hoverctlBinary, "delete", "middleware").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Middleware has been deleted from Hoverfly"))

				resp := DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/middleware", adminPort)))
				bytes, _ := ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).To(Equal(`{"middleware":""}`))
			})
		})


		Context("I can delete the request templates in Hoverfly", func() {
			BeforeEach(func() {
				fileReader, err := os.Open("testdata/request-template.json")
				defer fileReader.Close()
				if err != nil {
					Fail("Failed to read request template test data")
				}
				resp := DoRequest(sling.New().Post(fmt.Sprintf("http://localhost:%v/api/templates", adminPort)).Body(fileReader))
				bytes, _ := ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).To(ContainSubstring(`{"message":"2 payloads import complete."}`))
			})

			It("and they should be removed", func() {
				out, _ := exec.Command(hoverctlBinary, "delete", "templates").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Request templates have been deleted from Hoverfly"))

				resp := DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/templates", adminPort)))
				bytes, _ := ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).To(ContainSubstring(`{"data":null}`))
			})
		})

		Context("I can delete everything in hoverfly", func() {

			BeforeEach(func() {
				DoRequest(sling.New().Put(fmt.Sprintf("http://localhost:%v/api/delays", adminPort)).Body(strings.NewReader(`
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

				DoRequest(sling.New().Post(fmt.Sprintf("http://localhost:%v/api/records", adminPort)).Body(strings.NewReader(`
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
				resp := DoRequest(sling.New().Post(fmt.Sprintf("http://localhost:%v/api/templates", adminPort)).Body(fileReader))
				bytes, _ := ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).To(ContainSubstring(`{"message":"2 payloads import complete."}`))
			})

			Context("I can delete all data in Hoverfly", func() {

				It("by calling delete all", func() {
					out, _ := exec.Command(hoverctlBinary, "delete", "all").Output()
					output := strings.TrimSpace(string(out))
					Expect(output).To(ContainSubstring("Delays, request templates and simulations have been deleted from Hoverfly"))

					resp := DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/delays", adminPort)))
					bytes, _ := ioutil.ReadAll(resp.Body)
					Expect(string(bytes)).To(Equal(`{"data":[]}`))

					resp = DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/records", adminPort)))
					bytes, _ = ioutil.ReadAll(resp.Body)
					Expect(string(bytes)).To(Equal(`{"data":null}`))

					resp = DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/templates", adminPort)))
					bytes, _ = ioutil.ReadAll(resp.Body)
					Expect(string(bytes)).To(ContainSubstring(`{"data":null}`))
				})
			})
		})

		Context("I won't delete if I have not specified what to delete", func() {
			It("when I call hoverctl delete", func() {
				out, _ := exec.Command(hoverctlBinary, "delete").Output()
				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("You have not specified what to delete from Hoverfly"))
			})
		})

	})
})