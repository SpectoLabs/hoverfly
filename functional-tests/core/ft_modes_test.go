package hoverfly_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"

	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Running Hoverfly in various modes", func() {

	Context("When running in capture mode", func() {

		var fakeServer *httptest.Server
		var fakeServerUrl *url.URL

		Context("with middleware", func() {
			BeforeEach(func() {
				hoverflyCmd = startHoverflyWithMiddleware(adminPort, proxyPort, "python testdata/middleware.py")

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

			It("Should modify the request but not the response", func() {
				CallFakeServerThroughProxy(fakeServer)
				expectedDestination := strings.Replace(fakeServerUrl.String(), "http://", "", 1)
				recordsJson, err := ioutil.ReadAll(ExportHoverflySimulation())
				Expect(err).To(BeNil())
				Expect(recordsJson).To(ContainSubstring(`"destination":{"exactMatch":"` + expectedDestination + `"}`))
				Expect(recordsJson).To(ContainSubstring(`"body":{"exactMatch":"CHANGED_REQUEST_BODY"}`))
			})
		})
	})

	Context("When running in synthesise mode", func() {

		Context("With middleware", func() {

			BeforeEach(func() {
				hoverflyCmd = startHoverflyWithMiddleware(adminPort, proxyPort, "python testdata/middleware.py")
				SetHoverflyMode("synthesize")
			})

			It("Should generate responses using middleware", func() {
				resp := DoRequestThroughProxy(sling.New().Get("http://www.virtual.com/path2"))
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(string(body)).To(Equal("CHANGED_RESPONSE_BODY"))
			})

			AfterEach(func() {
				stopHoverfly()
			})

		})

		Context("Without middleware", func() {

			BeforeEach(func() {
				hoverflyCmd = startHoverfly(adminPort, proxyPort)
				SetHoverflyMode("synthesize")
			})

			It("Should fail to generate responses using middleware", func() {
				resp := DoRequestThroughProxy(sling.New().Get("http://www.virtual.com/path2"))
				Expect(resp.StatusCode).To(Equal(http.StatusBadGateway))
			})

			AfterEach(func() {
				stopHoverfly()
			})
		})
	})

	Context("When running in modify mode", func() {

		var fakeServer *httptest.Server
		var requestBody string

		Context("With middleware", func() {

			BeforeEach(func() {
				hoverflyCmd = startHoverflyWithMiddleware(adminPort, proxyPort, "python testdata/middleware.py")
				SetHoverflyMode("modify")
				fakeServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					body, _ := ioutil.ReadAll(r.Body)
					requestBody = string(body)
					w.Header().Set("Content-Type", "text/plain")
					w.Header().Set("Date", "date")
					w.Write([]byte("Hello world"))
				}))
			})

			It("Should modify the request using middleware", func() {
				resp := DoRequestThroughProxy(sling.New().Get(fakeServer.URL))
				Expect(resp.StatusCode).To(Equal(200))
				Expect(requestBody).To(Equal("CHANGED_REQUEST_BODY"))
			})

			It("Should modify the response using middleware", func() {
				resp := DoRequestThroughProxy(sling.New().Get(fakeServer.URL))
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(string(body)).To(Equal("CHANGED_RESPONSE_BODY"))
			})

			AfterEach(func() {
				stopHoverfly()
				fakeServer.Close()
			})
		})

		Context("Without middleware", func() {

			BeforeEach(func() {
				hoverflyCmd = startHoverfly(adminPort, proxyPort)
				SetHoverflyMode("modify")
			})

			It("Should fail to generate responses using middleware", func() {
				resp := DoRequestThroughProxy(sling.New().Get(fakeServer.URL))
				Expect(resp.StatusCode).To(Equal(http.StatusBadGateway))
			})

			AfterEach(func() {
				stopHoverfly()
				fakeServer.Close()
			})
		})

	})

	Context("Using middleware with binary data", func() {

		var expectedImage []byte

		BeforeEach(func() {
			hoverflyCmd = startHoverflyWithMiddleware(adminPort, proxyPort, "python testdata/binary_middleware.py")
			SetHoverflyMode("synthesize")
			pwd, _ := os.Getwd()
			expectedFile := "/testdata/1x1.png"
			expectedImage, _ = ioutil.ReadFile(pwd + expectedFile)
		})

		It("Should render an image correctly after base64 encoding it using middleware", func() {
			resp := DoRequestThroughProxy(sling.New().Get("http://www.foo.com"))
			responseBytes, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			Expect(responseBytes).To(Equal(expectedImage))
		})

		AfterEach(func() {
			stopHoverfly()
		})
	})
})
