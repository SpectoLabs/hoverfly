package hoverfly_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"

	"github.com/SpectoLabs/hoverfly/core/views"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Running Hoverfly in various modes", func() {

	Context("When running in capture mode", func() {

		var fakeServer *httptest.Server
		var fakeServerUrl *url.URL

		Context("without middleware", func() {

			BeforeEach(func() {
				hoverflyCmd = startHoverfly(adminPort, proxyPort)

				SetHoverflyMode("capture")
			})

			AfterEach(func() {
				stopHoverfly()
			})

			It("Should capture the request and response", func() {

				fakeServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/plain")
					w.Header().Set("Date", "date")
					w.Write([]byte("Hello world"))
				}))

				defer fakeServer.Close()

				fakeServerUrl, _ = url.Parse(fakeServer.URL)
				resp := CallFakeServerThroughProxy(fakeServer)
				Expect(resp.StatusCode).To(Equal(200))

				expectedDestination := strings.Replace(fakeServerUrl.String(), "http://", "", 1)

				recordsJson, err := ioutil.ReadAll(ExportHoverflyRecords())
				Expect(err).To(BeNil())
				Expect(recordsJson).To(MatchJSON(fmt.Sprintf(
					`{
					  "data": [
					    {
					      "response": {
						"status": 200,
						"body": "Hello world",
						"encodedBody": false,
						"headers": {
						  "Content-Length": [
						    "11"
						  ],
						  "Content-Type": [
						    "text/plain"
						  ],
						  "Date": [
						    "date"
						  ],
						  "Hoverfly": [
						    "Was-Here"
						  ]
						}
					      },
					      "request": {
					      	"requestType": "recording",
						"path": "/",
						"method": "GET",
						"destination": "%v",
						"scheme": "http",
						"query": "",
						"body": "",
						"headers": {
						  "Accept-Encoding": [
						    "gzip"
						  ],
						  "User-Agent": [
						    "Go-http-client/1.1"
						  ]
						}
					      }
					    }
					  ]
					}`, expectedDestination)))
			})

			It("Should capture the request and response in capture mode after the same request is made in simulate mode", func() {

				fakeServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/plain")
					w.Header().Set("Date", "date")
					w.Write([]byte("Hello world"))
				}))

				defer fakeServer.Close()

				fakeServerUrl, _ = url.Parse(fakeServer.URL)

				SetHoverflyMode("simulate")
				resp := CallFakeServerThroughProxy(fakeServer)
				Expect(resp.StatusCode).To(Equal(412))
				SetHoverflyMode("capture")

				resp = CallFakeServerThroughProxy(fakeServer)
				Expect(resp.StatusCode).To(Equal(200))

				expectedDestination := strings.Replace(fakeServerUrl.String(), "http://", "", 1)

				recordsJson, err := ioutil.ReadAll(ExportHoverflyRecords())
				Expect(err).To(BeNil())
				Expect(recordsJson).To(MatchJSON(fmt.Sprintf(
					`{
					  "data": [
					    {
					      "response": {
						"status": 200,
						"body": "Hello world",
						"encodedBody": false,
						"headers": {
						  "Content-Length": [
						    "11"
						  ],
						  "Content-Type": [
						    "text/plain"
						  ],
						  "Date": [
						    "date"
						  ],
						  "Hoverfly": [
						    "Was-Here"
						  ]
						}
					      },
					      "request": {
					      	"requestType": "recording",
						"path": "/",
						"method": "GET",
						"destination": "%v",
						"scheme": "http",
						"query": "",
						"body": "",
						"headers": {
						  "Accept-Encoding": [
						    "gzip"
						  ],
						  "User-Agent": [
						    "Go-http-client/1.1"
						  ]
						}
					      }
					    }
					  ]
					}`, expectedDestination)))
			})

			It("Should capture a redirect response", func() {

				fakeServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/plain")
					w.Header().Set("Date", "date")
					w.Write([]byte("redirection got you here"))
				}))
				fakeServerUrl, _ := url.Parse(fakeServer.URL)
				fakeServerRedirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Location", fakeServer.URL)
					w.Header().Set("Content-Type", "text/plain")
					w.WriteHeader(301)
				}))
				fakeServerRedirectUrl, _ := url.Parse(fakeServerRedirect.URL)

				defer fakeServer.Close()
				defer fakeServerRedirect.Close()

				resp := CallFakeServerThroughProxy(fakeServerRedirect)
				Expect(resp.StatusCode).To(Equal(301))

				expectedRedirectDestination := strings.Replace(fakeServerRedirectUrl.String(), "http://", "", 1)

				recordsJson, err := ioutil.ReadAll(ExportHoverflyRecords())
				Expect(err).To(BeNil())

				payload := views.RequestResponsePairPayload{}

				json.Unmarshal(recordsJson, &payload)
				Expect(payload.Data).To(HaveLen(1))

				Expect(payload.Data[0].Request.Destination).To(Equal(expectedRedirectDestination))

				Expect(payload.Data[0].Response.Status).To(Equal(301))
				Expect(payload.Data[0].Response.Headers["Location"][0]).To(Equal(fakeServerUrl.String()))
			})
		})

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
				recordsJson, err := ioutil.ReadAll(ExportHoverflyRecords())
				Expect(err).To(BeNil())
				Expect(recordsJson).To(MatchJSON(fmt.Sprintf(
					`{
					  "data": [
					    {
					      "response": {
						"status": 200,
						"body": "Hello world",
						"encodedBody": false,
						"headers": {
						  "Content-Length": [
						    "11"
						  ],
						  "Content-Type": [
						    "text/plain"
						  ],
						  "Date": [
						    "date"
						  ],
						  "Hoverfly": [
						    "Was-Here"
						  ]
						}
					      },
					      "request": {
					      	"requestType": "recording",
						"path": "/",
						"method": "GET",
						"destination": "%v",
						"scheme": "http",
						"query": "",
						"body": "CHANGED_REQUEST_BODY",
						"headers": {
						  "Accept-Encoding": [
						    "gzip"
						  ],
						  "User-Agent": [
						    "Go-http-client/1.1"
						  ]
						}
					      }
					    }
					  ]
					}`, expectedDestination)))
			})
		})
	})

	Context("When running in simulate mode", func() {

		var (
			jsonRequestResponsePair *bytes.Buffer
		)

		BeforeEach(func() {
			jsonRequestResponsePair = bytes.NewBufferString(`{"data":[{"request": {"path": "/path1", "method": "GET", "destination": "www.virtual.com", "scheme": "http", "query": "", "body": "", "headers": {"Header": ["value1"]}}, "response": {"status": 201, "encodedBody": false, "body": "body1", "headers": {"Header": ["value1", "value2"]}}}, {"request": {"path": "/path2", "method": "GET", "destination": "www.virtual.com", "scheme": "http", "query": "", "body": "", "headers": {"Header": ["value2"]}}, "response": {"status": 202, "body": "body2", "headers": {"Header": ["value2"]}}}]}`)
		})

		Context("without middleware", func() {

			BeforeEach(func() {
				hoverflyCmd = startHoverfly(adminPort, proxyPort)
				SetHoverflyMode("simulate")
				ImportHoverflyRecords(jsonRequestResponsePair)
			})

			AfterEach(func() {
				stopHoverfly()
			})

			It("should return the cached response", func() {
				resp := DoRequestThroughProxy(sling.New().Get("http://www.virtual.com/path1"))
				Expect(resp.StatusCode).To(Equal(201))
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(string(body)).To(Equal("body1"))
				Expect(resp.Header).To(HaveKeyWithValue("Header", []string{"value1", "value2"}))
			})
		})

		Context("with middleware", func() {

			BeforeEach(func() {
				hoverflyCmd = startHoverflyWithMiddleware(adminPort, proxyPort, "python testdata/middleware.py")
				SetHoverflyMode("simulate")
				ImportHoverflyRecords(jsonRequestResponsePair)
			})

			It("should apply middleware to the cached response", func() {
				resp := DoRequestThroughProxy(sling.New().Get("http://www.virtual.com/path2"))
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(string(body)).To(Equal("CHANGED_RESPONSE_BODY"))
			})

			AfterEach(func() {
				stopHoverfly()
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
				Expect(resp.StatusCode).To(Equal(503))
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
				Expect(resp.StatusCode).To(Equal(503))
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
