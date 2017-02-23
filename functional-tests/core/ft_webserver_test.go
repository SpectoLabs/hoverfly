package hoverfly_test

import (
	"bytes"
	"io/ioutil"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
)

var _ = Describe("When running Hoverfly as a webserver", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
		hoverfly.Start("-webserver")
		hoverfly.SetMode("simulate")
		hoverfly.ImportSimulation(functional_tests.JsonSimulationGetAndPost)
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("and its in simulate mode", func() {

		// BeforeEach(func() {
		// hoverflyCmd = startHoverflyWebServer(adminPort, proxyPort)
		// getPair1 := bytes.NewBufferString(`{"data":[{"request": {"path": "/path1", "method": "GET", "destination": "destination1", "scheme": "", "query": "", "body": "", "headers": {"Header": ["value1"]}}, "response": {"status": 201, "encodedBody": false, "body": "body1", "headers": {"Header": ["value1", "value2"]}}}]}`)
		// getPair2 := bytes.NewBufferString(`{"data":[{"request": {"path": "/path1/resource", "method": "GET", "destination": "another-destination.com", "scheme": "", "query": "", "body": "", "headers": {"Header": ["value1"]}}, "response": {"status": 201, "encodedBody": false, "body": "another-host.com body1", "headers": {"Header": ["value1"]}}}]}`)
		// postPair1 := bytes.NewBufferString(`{"data":[{"request": {"path": "/path2", "method": "POST", "destination": "destination1", "scheme": "", "query": "", "body": "", "headers": {"Header": ["value1"]}}, "response": {"status": 201, "encodedBody": false, "body": "body2", "headers": {"Header": ["value1"]}}}]}`)
		// postPair2 := bytes.NewBufferString(`{"data":[{"request": {"path": "/path2/resource", "method": "POST", "destination": "another-destination.com", "scheme": "", "query": "", "body": "", "headers": {"Header": ["value1"]}}, "response": {"status": 201, "encodedBody": false, "body": "another-host.com body2", "headers": {"Header": ["value1"]}}}]}`)

		// ImportHoverflySimulation()
		// })

		// AfterEach(func() {
		// hoverflyCmd.Process.Kill()
		// })

		Context("I can request an endpoint", func() {
			Context("using GET", func() {
				It("and it should return the response", func() {
					request := sling.New().Get("http://localhost:" + hoverfly.GetProxyPort() + "/path1")

					response := functional_tests.DoRequest(request)

					responseBody, err := ioutil.ReadAll(response.Body)
					Expect(err).To(BeNil())

					Expect(string(responseBody)).To(Equal("body1"))
				})

				It("and it should return the correct headers on the response", func() {
					request := sling.New().Get("http://localhost:" + hoverfly.GetProxyPort() + "/path1")

					response := functional_tests.DoRequest(request)

					Expect(response.Header).To(HaveKeyWithValue("Header", []string{"value1", "value2"}))
				})
			})

			Context("using POST", func() {
				It("and it should return the response", func() {
					request := sling.New().Post("http://localhost:" + hoverfly.GetProxyPort() + "/path2/resource")

					response := functional_tests.DoRequest(request)

					responseBody, err := ioutil.ReadAll(response.Body)
					Expect(err).To(BeNil())

					Expect(string(responseBody)).To(Equal("POST body response"))
				})
			})
		})

		Context("I cannot change the mode", func() {

			It("it should start in simulate mode", func() {
				request := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/state")
				response := functional_tests.DoRequest(request)

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).To(BeNil())

				Expect(string(responseBody)).To(ContainSubstring("simulate"))
			})

			It("it should not be switchable", func() {
				request := sling.New().Post("http://localhost:" + hoverfly.GetAdminPort() + "/api/state").Body(strings.NewReader(`{"mode":"capture"}`))
				response := functional_tests.DoRequest(request)

				Expect(response.StatusCode).To(Equal(403))

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).To(BeNil())

				Expect(string(responseBody)).To(ContainSubstring("Hoverfly is currently configured to act as webserver, which can only operate in simulate mode"))
			})
		})
	})

	Context("using a request.db from the disk", func() {

		Context("and the requests.db was created by Hoverfly proxy", func() {
			BeforeEach(func() {
				hoverflyCmd = startHoverflyWithDatabase(adminPort, proxyPort)
				ImportHoverflySimulation(bytes.NewBufferString(functional_tests.JsonSimulationGetAndPost))

				hoverflyCmd.Process.Kill()

				hoverflyCmd = startHoverflyWebServerWithDatabase(adminPort, proxyPort)
			})

			AfterEach(func() {
				hoverflyCmd.Process.Kill()
			})

			It("should recache the requests so that the destination is not included in the cache", func() {
				request := sling.New().Get("http://localhost:" + proxyPortAsString + "/path1")

				response := functional_tests.DoRequest(request)

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).To(BeNil())

				Expect(string(responseBody)).To(Equal("body1"))
			})

		})
	})
})
