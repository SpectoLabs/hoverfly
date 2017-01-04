package hoverfly_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Running Hoverfly with delays", func() {

	Context("When running in capture mode", func() {

		var fakeServer *httptest.Server
		var fakeServerUrl *url.URL

		BeforeEach(func() {
			hoverflyCmd = startHoverfly(adminPort, proxyPort)

			fakeServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.Header().Set("Date", "date")
				w.Write([]byte("Hello world"))
			}))

			fakeServerUrl, _ = url.Parse(fakeServer.URL)
			SetHoverflyMode("capture")
		})

		AfterEach(func() {
			stopHoverfly()
			fakeServer.Close()
		})

		It("Should NOT delay the response", func() {
			start := time.Now()
			resp := CallFakeServerThroughProxy(fakeServer)
			end := time.Now()
			reqDuration := end.Sub(start)
			Expect(resp.StatusCode).To(Equal(200))
			_, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			Expect(reqDuration < (100 * time.Millisecond)).To(BeTrue())
		})
	})

	Context("When running in simulate mode", func() {

		var (
			jsonRequestResponsePair *bytes.Buffer
		)

		BeforeEach(func() {
			jsonRequestResponsePair = bytes.NewBufferString(`{"data":[{"request": {"path": "/path1", "method": "GET", "destination": "www.virtual.com", "scheme": "http", "query": "", "body": "", "headers": {"Header": ["value1"]}}, "response": {"status": 201, "encodedBody": false, "body": "body1", "headers": {"Header": ["value1"]}}}, {"request": {"path": "/path2", "method": "GET", "destination": "www.virtual.com", "scheme": "http", "query": "", "body": "", "headers": {"Header": ["value2"]}}, "response": {"status": 202, "body": "body2", "headers": {"Header": ["value2"]}}}]}`)
			hoverflyCmd = startHoverfly(adminPort, proxyPort)
			SetHoverflyResponseDelays("testdata/delays.json")
			SetHoverflyMode("simulate")
			ImportHoverflyRecords(jsonRequestResponsePair)
		})

		It("should delay returning the cached response", func() {
			start := time.Now()
			resp := DoRequestThroughProxy(sling.New().Get("http://www.virtual.com/path2"))
			end := time.Now()
			reqDuration := end.Sub(start)
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("body2"))
			Expect(reqDuration > (100 * time.Millisecond)).To(BeTrue())
		})

		AfterEach(func() {
			stopHoverfly()
		})
	})

	Context("When running in synthesise mode (with middleware)", func() {

		BeforeEach(func() {
			hoverflyCmd = startHoverflyWithMiddleware(adminPort, proxyPort, "python testdata/middleware.py")
			SetHoverflyResponseDelays("testdata/delays.json")
			SetHoverflyMode("synthesize")
		})

		It("should delay returning the response", func() {
			start := time.Now()
			resp := DoRequestThroughProxy(sling.New().Get("http://www.virtual.com/path2"))
			end := time.Now()
			reqDuration := end.Sub(start)
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("CHANGED_RESPONSE_BODY"))
			Expect(reqDuration > (100 * time.Millisecond)).To(BeTrue())

		})

		AfterEach(func() {
			stopHoverfly()
		})

	})

	Context("When running in modify mode", func() {

		BeforeEach(func() {
			hoverflyCmd = startHoverflyWithMiddleware(adminPort, proxyPort, "python testdata/middleware.py")
			SetHoverflyResponseDelays("testdata/delays.json")
			SetHoverflyMode("modify")
		})

		It("should delay returning the response", func() {
			start := time.Now()
			resp := DoRequestThroughProxy(sling.New().Get("http://www.virtual.com/path2"))
			end := time.Now()
			reqDuration := end.Sub(start)
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("CHANGED_RESPONSE_BODY"))
			Expect(reqDuration > (100 * time.Millisecond)).To(BeTrue())
		})

		AfterEach(func() {
			stopHoverfly()
		})

	})
})
