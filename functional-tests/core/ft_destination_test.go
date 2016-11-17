package hoverfly_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Running Hoverfly", func() {

	Context("in capture mode", func() {

		var fakeServer *httptest.Server
		var fakeServerUrl *url.URL

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
				w.Header().Set("Date", "date")
				w.Write([]byte("Hello world"))
			}))

			defer fakeServer.Close()

			fakeServerUrl, _ = url.Parse(fakeServer.URL)
			resp := CallFakeServerThroughProxy(fakeServer)
			Expect(resp.StatusCode).To(Equal(200))

			recordsJson, err := ioutil.ReadAll(ExportHoverflyRecords())
			Expect(err).To(BeNil())
			Expect(recordsJson).To(MatchJSON(fmt.Sprintf(
				`{
					  "data": null
					}`)))
		})
	})
})
