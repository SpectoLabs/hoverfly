package hoverfly_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Running Hoverfly in various modes", func() {

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
		var fakeServerUrl *url.URL

		Context("with middleware", func() {

			BeforeEach(func() {
				hoverfly.Start("-middleware", "python testdata/middleware.py")

				fakeServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/plain")
					w.Header().Set("Date", "date")
					w.Write([]byte("Hello world"))
				}))

				fakeServerUrl, _ = url.Parse(fakeServer.URL)

				hoverfly.SetMode("capture")
			})

			AfterEach(func() {
				fakeServer.Close()
			})

			It("Should modify the request but not the response", func() {
				hoverfly.Proxy(sling.New().Get(fakeServer.URL))
				expectedDestination := strings.Replace(fakeServerUrl.String(), "http://", "", 1)
				simulation := hoverfly.ExportSimulation()

				Expect(simulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
				Expect(simulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Value).To(Equal(expectedDestination))
				Expect(simulation.RequestResponsePairs[0].RequestMatcher.Body[0].Matcher).To(Equal("exact"))
				Expect(simulation.RequestResponsePairs[0].RequestMatcher.Body[0].Value).To(Equal("CHANGED_REQUEST_BODY"))
			})
		})
	})

	Context("When running in synthesise mode", func() {

		Context("With middleware", func() {

			BeforeEach(func() {
				hoverfly.Start("-middleware", "python testdata/middleware.py")
				hoverfly.SetMode("synthesize")
			})

			It("Should generate responses using middleware", func() {
				resp := hoverfly.Proxy(sling.New().Get("http://www.virtual.com/path2"))
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(string(body)).To(Equal("CHANGED_RESPONSE_BODY"))
			})
		})

		Context("Without middleware", func() {

			BeforeEach(func() {
				hoverfly.Start()
				hoverfly.SetMode("synthesize")
			})

			It("Should fail to generate responses using middleware", func() {
				resp := hoverfly.Proxy(sling.New().Get("http://www.virtual.com/path2"))
				Expect(resp.StatusCode).To(Equal(http.StatusBadGateway))
			})
		})
	})

	Context("When running in modify mode", func() {

		var fakeServer *httptest.Server
		var requestBody string
		var requestQuery string

		Context("With middleware", func() {

			BeforeEach(func() {
				hoverfly.Start("-middleware", "python testdata/middleware.py")
				hoverfly.SetMode("modify")

				fakeServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					body, _ := ioutil.ReadAll(r.Body)
					requestBody = string(body)
					requestQuery = r.URL.RawQuery
					w.Header().Set("Content-Type", "text/plain")
					w.Header().Set("Date", "date")
					w.Write([]byte("Hello world"))
				}))
			})

			It("Should modify the request using middleware", func() {
				resp := hoverfly.Proxy(sling.New().Get(fakeServer.URL + "?test=a&test=b c,d"))
				Expect(resp.StatusCode).To(Equal(200))
				Expect(requestBody).To(Equal("CHANGED_REQUEST_BODY"))
				Expect(requestQuery).To(Equal("test=a&test=b%20c,d"))
			})

			It("Should modify the response using middleware", func() {
				resp := hoverfly.Proxy(sling.New().Get(fakeServer.URL))
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(string(body)).To(Equal("CHANGED_RESPONSE_BODY"))
			})

			AfterEach(func() {
				fakeServer.Close()
			})
		})

		Context("Without middleware", func() {

			BeforeEach(func() {
				hoverfly.Start()
				hoverfly.SetMode("modify")
			})

			It("Should fail to generate responses using middleware", func() {
				resp := hoverfly.Proxy(sling.New().Get(fakeServer.URL))
				Expect(resp.StatusCode).To(Equal(http.StatusBadGateway))
			})
		})
	})

	Context("Using middleware with binary data", func() {

		var expectedImage []byte

		BeforeEach(func() {
			hoverfly.Start("-middleware", "python testdata/binary_middleware.py")
			hoverfly.SetMode("synthesize")
			pwd, _ := os.Getwd()
			expectedFile := "/testdata/1x1.png"
			expectedImage, _ = ioutil.ReadFile(pwd + expectedFile)
		})

		It("Should render an image correctly after base64 encoding it using middleware", func() {
			resp := hoverfly.Proxy(sling.New().Get("http://www.foo.com"))
			responseBytes, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			Expect(responseBytes).To(Equal(expectedImage))
		})
	})
})
