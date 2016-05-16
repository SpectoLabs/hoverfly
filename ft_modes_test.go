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
)

var _ = Describe("Running Hoverfly in capture mode", func() {

	Context("When capturing http traffic", func() {

		var fakeServerUrl * url.URL

		BeforeEach(func() {
			requestCache.DeleteData()

			fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.Header().Set("Date", "date")
				w.Write([]byte("Hello world"))
			}))

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
})
