package hoverfly_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Running Hoverfly", func() {

	Context("in capture mode", func() {

		var fakeServer *httptest.Server

		BeforeEach(func() {
			hoverflyCmd = startHoverfly(adminPort, proxyPort)
			SetHoverflyMode("capture")
		})

		AfterEach(func() {
			stopHoverfly()
		})

		It("Should not capture capture if destination does not match", func() {
			SetHoverflyDestination("notlocalhost")

			fakeServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.Header().Set("date", "date")
				w.Write([]byte("Hello world"))
			}))

			defer fakeServer.Close()

			resp := CallFakeServerThroughProxy(fakeServer)

			Expect(resp.StatusCode).To(Equal(200))
			Expect(resp.Header.Get("date")).To(Equal("date"))

			recordsJson, err := ioutil.ReadAll(ExportHoverflyRecords())
			Expect(err).To(BeNil())
			Expect(recordsJson).To(MatchJSON(fmt.Sprintf(
				`{
					  "data": null
					}`)))
		})

		It("Should capture capture if destination is 127.0.0.1", func() {
			fakeServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.Header().Set("date", "date")
				w.Write([]byte("Hello world"))
			}))

			defer fakeServer.Close()

			SetHoverflyDestination("127.0.0.1")

			resp := CallFakeServerThroughProxy(fakeServer)

			Expect(resp.StatusCode).To(Equal(200))
			Expect(resp.Header.Get("date")).To(Equal("date"))

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
					}`, strings.Replace(fakeServer.URL, "http://", "", 1))))
		})
	})
})
