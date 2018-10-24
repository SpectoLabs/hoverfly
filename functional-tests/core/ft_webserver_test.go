package hoverfly_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/SpectoLabs/hoverfly/functional-tests/testdata"
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
		hoverfly.ImportSimulation(testdata.JsonGetAndPost)
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

				It("and it should increment the usage counter", func() {
					request := sling.New().Get("http://localhost:" + hoverfly.GetProxyPort() + "/path1")
					functional_tests.DoRequest(request)

					req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/usage")
					res := functional_tests.DoRequest(req)
					Expect(res.StatusCode).To(Equal(200))
					modeJson, err := ioutil.ReadAll(res.Body)
					Expect(err).To(BeNil())
					Expect(modeJson).To(Equal([]byte(`{"usage":{"counters":{"capture":0,"diff":0,"modify":0,"simulate":1,"spy":0,"synthesize":0}}}`)))
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
	})
})
