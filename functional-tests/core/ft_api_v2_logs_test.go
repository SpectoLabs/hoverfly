package hoverfly_test

import (
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/antonholmquist/jason"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/api/v2/logs", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
		hoverfly.Start()
		hoverfly.SetMode("simulate")
		hoverfly.ImportSimulation(functional_tests.JsonPayload)
		hoverfly.Proxy(sling.New().Get("http://destination-server.com"))
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("GET", func() {

		It("should get logs", func() {
			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/logs")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			responseJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())

			jsonObject, err := jason.NewObjectFromBytes(responseJson)
			Expect(err).To(BeNil())

			logsArray, err := jsonObject.GetObjectArray("logs")
			Expect(err).To(BeNil())

			Expect(len(logsArray)).To(BeNumerically(">=", 4))

			Expect(logsArray[0].GetString("msg")).Should(Equal("Proxy prepared..."))
			Expect(logsArray[1].GetString("msg")).Should(Equal("current proxy configuration"))
			Expect(logsArray[2].GetString("msg")).Should(Equal("serving proxy"))
			Expect(logsArray[3].GetString("msg")).Should(Equal("Admin interface is starting..."))
		})

		It("should limit the logs it returns", func() {
			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/logs?limit=1")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			responseJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())

			jsonObject, err := jason.NewObjectFromBytes(responseJson)
			Expect(err).To(BeNil())

			logsArray, err := jsonObject.GetObjectArray("logs")
			Expect(err).To(BeNil())

			Expect(len(logsArray)).To(Equal(1))

			Expect(logsArray[0].GetString("msg")).Should(Equal("started handling request"))
		})
	})

	Context("GET Content-Type text/plain", func() {

		It("should get logs", func() {
			req := sling.New().Get("http://localhost:"+hoverfly.GetAdminPort()+"/api/v2/logs").Add("Content-Type", "text/plain")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			responseBody, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())

			Expect(responseBody).To(ContainSubstring(`msg="Proxy prepared..."`))
			Expect(responseBody).To(ContainSubstring(`Destination=.`))
			Expect(responseBody).To(ContainSubstring(`Mode=simulate`))
			Expect(responseBody).To(ContainSubstring(`ProxyPort=` + hoverfly.GetProxyPort()))
		})

		It("should limit the logs it returns", func() {
			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/logs?limit=1")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			responseJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())

			jsonObject, err := jason.NewObjectFromBytes(responseJson)
			Expect(err).To(BeNil())

			logsArray, err := jsonObject.GetObjectArray("logs")
			Expect(err).To(BeNil())

			Expect(len(logsArray)).To(Equal(1))

			Expect(logsArray[0].GetString("msg")).Should(Equal("started handling request"))
		})
	})
})
