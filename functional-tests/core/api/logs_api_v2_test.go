package api_test

import (
	"io/ioutil"
	"net/http"

	"strconv"
	"time"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/SpectoLabs/hoverfly/functional-tests/testdata"
	"github.com/antonholmquist/jason"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func MakeLogsArray(logsJson []*jason.Object) []string {
	var logs []string
	for _, log := range logsJson {
		logMessage, _ := log.GetString("msg")
		logs = append(logs, logMessage)
	}

	return logs
}

var _ = Describe("/api/v2/logs", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("with log-size=100", func() {

		BeforeEach(func() {
			hoverfly.Start("-logs-size=100")
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		Context("GET", func() {

			It("should not exceed size", func() {
				for i := 0; i < 111; i++ {
					hoverfly.Proxy(sling.New().Get("http://hoverfly.io"))
				}

				req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/logs")
				res := functional_tests.DoRequest(req)

				Expect(res.StatusCode).To(Equal(200))

				var logs v2.LogsView

				functional_tests.UnmarshalFromResponse(res, &logs)

				Expect(logs.Logs).To(HaveLen(100))
			})
		})
	})

	Context("with log-size=0", func() {

		BeforeEach(func() {
			hoverfly.Start("-logs-size=0")
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		Context("GET", func() {

			It("should be disabled", func() {
				for i := 0; i < 111; i++ {
					hoverfly.Proxy(sling.New().Get("http://hoverfly.io"))
				}

				req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/logs")
				res := functional_tests.DoRequest(req)

				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))

				var errorView handlers.ErrorView

				functional_tests.UnmarshalFromResponse(res, &errorView)

				Expect(errorView.Error).To(Equal("Logs disabled"))
			})
		})
	})

	Context("with standard configuration", func() {

		BeforeEach(func() {
			hoverfly.Start()
			hoverfly.SetMode("simulate")
			hoverfly.ImportSimulation(testdata.JsonPayload)
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

				logs := MakeLogsArray(logsArray)

				Expect(logs).Should(ContainElement("Proxy prepared..."))
				Expect(logs).Should(ContainElement("current proxy configuration"))
				Expect(logs).Should(ContainElement("serving proxy"))
				Expect(logs).Should(ContainElement("Admin interface is starting..."))
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

				Expect(logsArray[0].GetString("msg")).Should(Equal("payloads imported"))
			})

			It("should query the logs by from time", func() {
				time.Sleep(time.Second)
				now := time.Now()

				hoverfly.Proxy(sling.New().Get("https://hoverfly.io"))

				req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/logs?from=" + strconv.FormatInt(now.Unix(), 10))
				res := functional_tests.DoRequest(req)
				Expect(res.StatusCode).To(Equal(200))
				responseJson, err := ioutil.ReadAll(res.Body)
				Expect(err).To(BeNil())

				jsonObject, err := jason.NewObjectFromBytes(responseJson)
				Expect(err).To(BeNil())

				logsArray, err := jsonObject.GetObjectArray("logs")
				Expect(err).To(BeNil())

				Expect(len(logsArray)).To(Equal(2))

				for _, log := range logsArray {
					timeStr, _ := log.GetString("time")
					logTime, _ := time.Parse(time.RFC3339, timeStr)
					Expect(logTime.Unix()).Should(BeNumerically(">=", now.Unix()))
				}
			})
		})

		Context("GET Accept text/plain", func() {

			It("should get logs", func() {
				req := sling.New().Get("http://localhost:"+hoverfly.GetAdminPort()+"/api/v2/logs").Add("Accept", "text/plain")
				res := functional_tests.DoRequest(req)
				Expect(res.StatusCode).To(Equal(200))
				responseBody, err := ioutil.ReadAll(res.Body)
				Expect(err).To(BeNil())

				Expect(responseBody).To(ContainSubstring("Proxy prepared..."))
				Expect(responseBody).To(ContainSubstring("=."))
				Expect(responseBody).To(ContainSubstring("=simulate"))
				Expect(responseBody).To(ContainSubstring("=" + hoverfly.GetProxyPort()))
			})

			It("should limit the logs it returns", func() {
				req := sling.New().Get("http://localhost:"+hoverfly.GetAdminPort()+"/api/v2/logs?limit=1").Add("Accept", "text/plain")
				res := functional_tests.DoRequest(req)
				Expect(res.StatusCode).To(Equal(200))
				responseBody, err := ioutil.ReadAll(res.Body)
				Expect(err).To(BeNil())

				Expect(responseBody).To(ContainSubstring("payloads imported"))

				Expect(responseBody).ToNot(ContainSubstring("Using memory backend"))
				Expect(responseBody).ToNot(ContainSubstring("Proxy prepared"))
				Expect(responseBody).ToNot(ContainSubstring("current proxy configuration"))
				Expect(responseBody).ToNot(ContainSubstring("Admin interface is starting..."))
			})
		})

		Context("GET Content-Type text/plain", func() {

			It("should get logs", func() {
				req := sling.New().Get("http://localhost:"+hoverfly.GetAdminPort()+"/api/v2/logs").Add("Content-Type", "text/plain")
				res := functional_tests.DoRequest(req)
				Expect(res.StatusCode).To(Equal(200))
				responseBody, err := ioutil.ReadAll(res.Body)
				Expect(err).To(BeNil())

				Expect(responseBody).To(ContainSubstring("Proxy prepared..."))
				Expect(responseBody).To(ContainSubstring("=."))
				Expect(responseBody).To(ContainSubstring("=simulate"))
				Expect(responseBody).To(ContainSubstring("=" + hoverfly.GetProxyPort()))
			})

			It("should limit the logs it returns", func() {
				req := sling.New().Get("http://localhost:"+hoverfly.GetAdminPort()+"/api/v2/logs?limit=1").Add("Content-Type", "text/plain")
				res := functional_tests.DoRequest(req)
				Expect(res.StatusCode).To(Equal(200))
				responseBody, err := ioutil.ReadAll(res.Body)
				Expect(err).To(BeNil())

				Expect(responseBody).To(ContainSubstring("payloads imported"))

				Expect(responseBody).ToNot(ContainSubstring("Using memory backend"))
				Expect(responseBody).ToNot(ContainSubstring("Proxy prepared"))
				Expect(responseBody).ToNot(ContainSubstring("current proxy configuration"))
				Expect(responseBody).ToNot(ContainSubstring("Admin interface is starting..."))
			})
		})
	})
})
