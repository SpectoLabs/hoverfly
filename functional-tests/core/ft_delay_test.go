package hoverfly_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/SpectoLabs/hoverfly/v2/functional-tests"
	"github.com/SpectoLabs/hoverfly/v2/functional-tests/testdata"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Running Hoverfly with delays", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("When running in capture mode", func() {

		var fakeServer *httptest.Server

		BeforeEach(func() {
			hoverfly.Start()

			fakeServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.Header().Set("Date", "date")
				w.Write([]byte("Hello world"))
			}))

			url.Parse(fakeServer.URL)
			hoverfly.SetMode("capture")
		})

		AfterEach(func() {
			fakeServer.Close()
		})

		It("Should NOT delay the response", func() {
			start := time.Now()
			resp := hoverfly.Proxy(sling.New().Get(fakeServer.URL))
			end := time.Now()
			reqDuration := end.Sub(start)
			Expect(resp.StatusCode).To(Equal(200))
			_, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			Expect(reqDuration < (100 * time.Millisecond)).To(BeTrue())
		})
	})

	Context("When running in simulate mode", func() {

		BeforeEach(func() {
			hoverfly.Start()
			hoverfly.SetMode("simulate")
		})

		It("should apply global delay to the cached response", func() {
			hoverfly.ImportSimulation(testdata.V3Delays)

			start := time.Now()
			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/path1"))
			end := time.Now()
			reqDuration := end.Sub(start)
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("exact match"))
			Expect(reqDuration > (100 * time.Millisecond)).To(BeTrue())
		})

		It("should apply fixed response delay to the cached response", func() {
			hoverfly.ImportSimulation(testdata.ResponseFixedDelays)

			start := time.Now()
			resp := hoverfly.Proxy(
				sling.New().Set("X-API-Version", "v1").Get("http://test-server.com/api/profile"),
			)
			end := time.Now()

			reqDuration := end.Sub(start)
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("Page is slow"))
			Expect(reqDuration > (100 * time.Millisecond)).To(BeTrue())
		})

		It("should ignore global delay when fixed response delay is specified", func() {
			simulation := `{
"data": {
        "pairs": [
			{
				"request": {"destination": [{"matcher": "exact", "value": "test-server.com"}]},
				"response": {"status": 200, "fixedDelay": 130}
			}
        ],
        "globalActions": {"delays": [{"urlPattern": "test-server\\.com", "delay": 10000}]}
    },
    "meta": {"schemaVersion": "v5", "hoverflyVersion": "v1.2.0"}
}`
			hoverfly.ImportSimulation(simulation)

			start := time.Now()
			hoverfly.Proxy(sling.New().Get("http://test-server.com/path1"))
			end := time.Now()
			reqDuration := end.Sub(start)
			Expect(reqDuration < (1000 * time.Millisecond)).To(BeTrue())
		})

		It("should apply log normal response delay to the cached response", func() {
			hoverfly.ImportSimulation(testdata.ResponseLogNormalDelays)

			start := time.Now()
			resp := hoverfly.Proxy(
				sling.New().Set("X-API-Version", "v1").Get("http://test-server.com/api/profile"),
			)
			end := time.Now()

			reqDuration := end.Sub(start)
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("Page is slow"))
			Expect(reqDuration > (100 * time.Millisecond)).To(BeTrue())
		})

		It("should ignore global delay when log normal response delay is specified", func() {
			simulation := `{
"data": {
        "pairs": [
			{
				"request": {"destination": [{"matcher": "exact", "value": "test-server.com"}]},
				"response": {"status": 200, "logNormalDelay": {"min": 100, "max": 150, "mean": 130, "median": 110}}
			}
        ],
        "globalActions": {"delays": [{"urlPattern": "test-server\\.com", "delay": 10000}]}
    },
    "meta": {"schemaVersion": "v5", "hoverflyVersion": "v1.2.0"}
}`
			hoverfly.ImportSimulation(simulation)

			start := time.Now()
			hoverfly.Proxy(sling.New().Get("http://test-server.com/path1"))
			end := time.Now()
			reqDuration := end.Sub(start)
			Expect(reqDuration < (1000 * time.Millisecond)).To(BeTrue())
		})

		It("should apply both fixed and log normal delays when specified", func() {
			simulation := `{
"data": {
        "pairs": [
			{
				"request": {"destination": [{"matcher": "exact", "value": "test-server.com"}]},
				"response": {
					"status": 200,
					"fixedDelay": 60,
					"logNormalDelay": {"min": 50, "max": 100, "mean": 80, "median": 60}
				}
			}
        ],
        "globalActions": {"delays": []}
    },
    "meta": {"schemaVersion": "v5", "hoverflyVersion": "v1.2.0"}
}`
			hoverfly.ImportSimulation(simulation)

			start := time.Now()
			hoverfly.Proxy(sling.New().Get("http://test-server.com/path1"))
			end := time.Now()
			reqDuration := end.Sub(start)
			Expect(reqDuration > (110 * time.Millisecond)).To(BeTrue())
		})
	})

	Context("When running in synthesize mode (with middleware)", func() {

		It("should apply global delay to the response", func() {
			hoverfly.Start("-middleware", "python testdata/middleware.py")
			hoverfly.SetMode("synthesize")
			hoverfly.ImportSimulation(testdata.V3Delays)

			start := time.Now()
			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/path2"))
			end := time.Now()
			reqDuration := end.Sub(start)
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("CHANGED_RESPONSE_BODY"))
			Expect(reqDuration > (100 * time.Millisecond)).To(BeTrue())
		})

		It("should apply fixed response delay to the response", func() {
			hoverfly.Start("-middleware", "python testdata/response_fixed_delay_middleware.py")
			hoverfly.SetMode("synthesize")
			hoverfly.ImportSimulation(testdata.ResponseFixedDelays)

			start := time.Now()
			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/api/settings"))
			end := time.Now()

			reqDuration := end.Sub(start)
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("CHANGED_RESPONSE_BODY"))
			Expect(reqDuration > (130 * time.Millisecond)).To(BeTrue())
		})

		It("should apply log normal response delay to the response", func() {
			hoverfly.Start("-middleware", "python testdata/response_lognormal_delay_middleware.py")
			hoverfly.SetMode("synthesize")
			hoverfly.ImportSimulation(testdata.ResponseLogNormalDelays)

			start := time.Now()
			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/api/settings"))
			end := time.Now()

			reqDuration := end.Sub(start)
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("CHANGED_RESPONSE_BODY"))
			Expect(reqDuration > (130 * time.Millisecond)).To(BeTrue())
		})
	})

	Context("When running in modify mode", func() {

		BeforeEach(func() {
			hoverfly.Start("-middleware", "python testdata/middleware.py")
			hoverfly.SetMode("modify")
		})

		It("should apply global delay to the response", func() {
			hoverfly.ImportSimulation(testdata.V3Delays)

			start := time.Now()
			resp := hoverfly.Proxy(sling.New().Get("http://localhost:" + hoverfly.GetAdminPort()))
			end := time.Now()
			reqDuration := end.Sub(start)
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("CHANGED_RESPONSE_BODY"))
			Expect(reqDuration > (100 * time.Millisecond)).To(BeTrue())
		})
	})

	Context("When running in modify mode with middleware to set delay", func() {

		BeforeEach(func() {
			hoverfly.Start("-middleware", "python testdata/middleware-set-delay.py")
			hoverfly.SetMode("modify")
		})

		It("should apply response delay to the response", func() {

			start := time.Now()
			resp := hoverfly.Proxy(sling.New().Get("http://localhost:" + hoverfly.GetAdminPort()))
			end := time.Now()

			reqDuration := end.Sub(start)
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("CHANGED_RESPONSE_BODY"))
			Expect(reqDuration > (130 * time.Millisecond)).To(BeTrue())
		})
	})
})
