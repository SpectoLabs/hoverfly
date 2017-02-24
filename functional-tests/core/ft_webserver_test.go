package hoverfly_test

import (
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
})
