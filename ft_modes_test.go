package hoverfly_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	//"compress/gzip"
	"io/ioutil"
	"github.com/SpectoLabs/hoverfly"
	//"compress/gzip"
	"fmt"
	"strings"
	"net/url"
	"os"
)

var _ = Describe("Running Hoverfly in various modes", func() {

	Context("When running in capture mode", func() {

		var fakeServer * httptest.Server
		var fakeServerUrl * url.URL

		Context("without middleware", func() {

			BeforeEach(func() {
				requestCache.DeleteData()
				fakeServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/plain")
					w.Header().Set("Date", "date")
					w.Write([]byte("Hello world"))
				}))

				defer fakeServer.Close()

				fakeServerUrl, _ = url.Parse(fakeServer.URL)
				SetHoverflyMode(hoverfly.CaptureMode)
				resp := CallFakeServerThroughProxy(fakeServer)
				Expect(resp.StatusCode).To(Equal(200))
			})

			It("Should capture the request and response", func() {
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
		})

		Context("with middleware", func() {
			BeforeEach(func() {
				requestCache.DeleteData()
				fakeServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/plain")
					w.Header().Set("Date", "date")
					w.Write([]byte("Hello world"))
				}))

				fakeServerUrl, _ = url.Parse(fakeServer.URL)
				SetHoverflyMode(hoverfly.CaptureMode)

				wd, err := os.Getwd()
				Expect(err).To(BeNil())
				hf.Cfg.Middleware = wd + "/testdata/middleware.py"
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
						"path": "/",
						"method": "GET",
						"destination": "%v",
						"scheme": "http",
						"query": "",
						"body": "CHANGED",
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

			AfterEach(func() {
				hf.Cfg.Middleware = ""
				fakeServer.Close()
			})
		})
	})
})
