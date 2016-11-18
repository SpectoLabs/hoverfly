package hoverfly_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"strings"

	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Running Hoverfly", func() {

	Context("in capture mode", func() {

		var fakeServer *httptest.Server

		BeforeEach(func() {
			hoverflyCmd = startHoverfly(adminPort, proxyPort)
			SetHoverflyMode("capture")

			fakeServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.Header().Set("date", "date")
				w.Write([]byte("Hello world"))
			}))

		})

		AfterEach(func() {
			stopHoverfly()

			fakeServer.Close()
		})

		It("Should not capture if destination does not match", func() {
			SetHoverflyDestination("notlocalhost")

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

		It("Should capture if destination is 127.0.0.1", func() {
			SetHoverflyDestination("127.0.0.1")

			resp := CallFakeServerThroughProxy(fakeServer)

			Expect(resp.StatusCode).To(Equal(200))
			Expect(resp.Header.Get("date")).To(Equal("date"))

			recordsJson, err := ioutil.ReadAll(ExportHoverflyRecords())
			Expect(err).To(BeNil())
			Expect(recordsJson).ToNot(MatchJSON(fmt.Sprintf(`{
				"data": null
			}`)))
		})

		It("Should capture if destination is set to port numbers", func() {
			SetHoverflyDestination(strings.Replace(fakeServer.URL, "http://127.0.0.1", "", 1))

			resp := CallFakeServerThroughProxy(fakeServer)

			Expect(resp.StatusCode).To(Equal(200))
			Expect(resp.Header.Get("date")).To(Equal("date"))

			recordsJson, err := ioutil.ReadAll(ExportHoverflyRecords())
			Expect(err).To(BeNil())
			Expect(recordsJson).ToNot(MatchJSON(fmt.Sprintf(`{
				"data": null
			}`)))
		})

		It("Should capture if destination is set to the path", func() {
			SetHoverflyDestination("/path")

			resp := DoRequestThroughProxy(sling.New().Get(fakeServer.URL + "/path"))

			Expect(resp.StatusCode).To(Equal(200))
			Expect(resp.Header.Get("date")).To(Equal("date"))

			recordsJson, err := ioutil.ReadAll(ExportHoverflyRecords())
			Expect(err).To(BeNil())
			Expect(recordsJson).ToNot(MatchJSON(fmt.Sprintf(`{
				"data": null
			}`)))
		})

		It("Should not capture if destination is set to the wrong path", func() {
			SetHoverflyDestination("/wrongpath")

			resp := DoRequestThroughProxy(sling.New().Get(fakeServer.URL + "/path"))

			Expect(resp.StatusCode).To(Equal(200))
			Expect(resp.Header.Get("date")).To(Equal("date"))

			recordsJson, err := ioutil.ReadAll(ExportHoverflyRecords())
			Expect(err).To(BeNil())
			Expect(recordsJson).To(MatchJSON(fmt.Sprintf(`{
				"data": null
			}`)))
		})
	})
})
