package hoverfly_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/SpectoLabs/hoverfly/functional-tests/testdata"
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
			hoverfly.ImportSimulation(testdata.V3Delays)
			hoverfly.SetMode("simulate")
		})

		It("should delay returning the cached response", func() {
			start := time.Now()
			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/path1"))
			end := time.Now()
			reqDuration := end.Sub(start)
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("exact match"))
			Expect(reqDuration > (100 * time.Millisecond)).To(BeTrue())
		})
	})

	Context("When running in synthesise mode (with middleware)", func() {

		BeforeEach(func() {
			hoverfly.Start("-middleware", "python testdata/middleware.py")
			hoverfly.ImportSimulation(testdata.V3Delays)
			hoverfly.SetMode("synthesize")
		})

		It("should delay returning the response", func() {
			start := time.Now()
			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/path2"))
			end := time.Now()
			reqDuration := end.Sub(start)
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("CHANGED_RESPONSE_BODY"))
			Expect(reqDuration > (100 * time.Millisecond)).To(BeTrue())

		})
	})

	Context("When running in modify mode", func() {

		BeforeEach(func() {
			hoverfly.Start("-middleware", "python testdata/middleware.py")
			hoverfly.ImportSimulation(testdata.V3Delays)
			hoverfly.SetMode("modify")
		})

		It("should delay returning the response", func() {
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
})
