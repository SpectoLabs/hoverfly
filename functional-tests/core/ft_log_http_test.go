package hoverfly_test

import (
	"io/ioutil"
	"regexp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	functional_tests "github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/SpectoLabs/hoverfly/functional-tests/testdata"
	"github.com/dghubble/sling"
)

// remove ANSI escaped color codes
func removeEscapedChars(str string) string {
	re := regexp.MustCompile(`\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])`)
	s := re.ReplaceAllString(str, "")
	return s
}

var _ = Describe("When running Hoverfly as a webserver", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("and its in simulate mode with flag -log-http", func() {

		BeforeEach(func() {
			hoverfly.Start("-webserver", "-log-http")
			hoverfly.SetMode("simulate")
			hoverfly.ImportSimulation(testdata.JsonGetAndPost)
		})

		Context("I can request an endpoint", func() {
			Context("using GET", func() {
				It("and it should write HTTP request/response messages into log", func() {
					request := sling.New().Get("http://localhost:" + hoverfly.GetProxyPort() + "/path1")
					response := functional_tests.DoRequest(request)

					responseBody, err := ioutil.ReadAll(response.Body)
					Expect(err).To(BeNil())

					Expect(string(responseBody)).To(Equal("body1"))

					req := sling.New().Get("http://localhost:"+hoverfly.GetAdminPort()+"/api/v2/logs").Add("Accept", "text/plain")
					res := functional_tests.DoRequest(req)
					Expect(res.StatusCode).To(Equal(200))

					logs, err := ioutil.ReadAll(res.Body)
					Expect(err).To(BeNil())
					log := removeEscapedChars(string(logs))
					Expect(log).To(ContainSubstring("HTTP Message"))
					Expect(log).To(ContainSubstring("Path=/path1"))
					Expect(log).To(ContainSubstring("Scheme=http"))
					Expect(log).To(ContainSubstring("Method=GET"))
					Expect(log).To(ContainSubstring("ResponseBody=" + string(responseBody)))
					Expect(log).To(ContainSubstring("Status=201"))
				})
			})

			Context("using POST", func() {
				It("and it should write HTTP request/response messages into log", func() {
					request := sling.New().Post("http://localhost:" + hoverfly.GetProxyPort() + "/path2/resource")

					response := functional_tests.DoRequest(request)

					responseBody, err := ioutil.ReadAll(response.Body)
					Expect(err).To(BeNil())

					Expect(string(responseBody)).To(Equal("POST body response"))

					req := sling.New().Get("http://localhost:"+hoverfly.GetAdminPort()+"/api/v2/logs").Add("Accept", "text/plain")
					res := functional_tests.DoRequest(req)
					Expect(res.StatusCode).To(Equal(200))

					logs, err := ioutil.ReadAll(res.Body)
					Expect(err).To(BeNil())
					log := removeEscapedChars((string(logs)))
					Expect(log).To(ContainSubstring("HTTP Message"))
					Expect(log).To(ContainSubstring("Path=/path2/resource"))
					Expect(log).To(ContainSubstring("Scheme=http"))
					Expect(log).To(ContainSubstring("Method=POST"))
					Expect(log).To(ContainSubstring("ResponseBody=\"" + string(responseBody) + "\""))
					Expect(log).To(ContainSubstring("Status=200"))
				})
			})
		})
	})
})
