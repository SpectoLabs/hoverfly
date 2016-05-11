package hoverfly_test

import (
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	//"compress/gzip"
	"fmt"
	"io/ioutil"
	"github.com/SpectoLabs/hoverfly"
	//"compress/gzip"
	"compress/gzip"
	"bytes"
)

// Helper function for gzipping strings
func GzipString(s string) (string) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	gz.Write([]byte(s))
	return b.String()
}

var _ = Describe("Capture > export > importing > simulate flow", func() {


	Describe("Import, Export", func() {
		Context("The captured response should be returned after exporting and importing", func() {

			var afterImportFakeServerResponse *http.Response

			BeforeEach(func() {
				// Spin up a fake server which returns hello world
				fakeGzipServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

					w.Header().Set("Content-Type", "text/plain")
					w.WriteHeader(200)
					fmt.Fprintf(w, "hello_world")
				}))

				// Switch Hoverfly to capture mode
				SetHoverflyMode(hoverfly.CaptureMode)

				// Make a request to the fake server and proxy through Hoverfly
				fakeServerUrl := fakeGzipServer.URL

				fakeServerRequest := sling.New().Get(fakeServerUrl)

				response := DoRequestThroughProxy(fakeServerRequest)
				Expect(response.StatusCode).To(Equal(200))

				// Export the data out of Hoverfly
				exportedRecords := ExportHoverflyRecords()

				// Wipe the records in Hoverfly
				EraseHoverflyRecords()

				// Import the same data into Hoverfly
				ImportHoverflyRecords(exportedRecords)

				// Switch Hoverfly to simulate mode
				SetHoverflyMode(hoverfly.SimulateMode)

				// Make the request to Hoverfly simulate
				afterImportFakeServerRequest := sling.New().Get(fakeServerUrl)
				afterImportFakeServerResponse = DoRequestThroughProxy(afterImportFakeServerRequest)
			})


			It("Returns a status code of 200", func() {
				Expect(afterImportFakeServerResponse.StatusCode).To(Equal(200))
			})

			It("Returns hello world", func() {
				responseBody, _ := ioutil.ReadAll(afterImportFakeServerResponse.Body)
				Expect(string(responseBody)).To(Equal(`hello_world`))
			})

			It("Returns with Hoverfly header", func() {
				Expect(afterImportFakeServerResponse.Header).To(HaveKeyWithValue("Hoverfly", []string{"Was-Here"}))
			})

			It("Returns with text/plain Content-Type header", func() {
				Expect(afterImportFakeServerResponse.Header).To(HaveKeyWithValue("Content-Type", []string{"text/plain"}))
			})

			It("Returns with a Hoverfly header", func() {
				Expect(afterImportFakeServerResponse.Header).To(HaveKeyWithValue("Hoverfly", []string{"Was-Here"}))
			})
		})

		Context("The captured response should be returned after exporting and importing when gzipped", func() {

			var afterImportFakeServerResponse *http.Response

			BeforeEach(func() {
				// Spin up a fake server which returns hello world gzipped
				fakeGzipServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

					w.Header().Set("Content-Encoding", "gzip")
					w.Header().Set("Content-Type", "text/plain")
					gzipWriter := gzip.NewWriter(w)
					gzipWriter.Write([]byte(`hello_world`))
				}))

				// Switch Hoverfly to capture mode
				SetHoverflyMode(hoverfly.CaptureMode)

				// Make a request to the fake server and proxy through Hoverfly
				fakeServerUrl := fakeGzipServer.URL

				fakeServerRequest := sling.New().Get(fakeServerUrl).Set("Accept-Encoding", "gzip")

				response := DoRequestThroughProxy(fakeServerRequest)
				Expect(response.StatusCode).To(Equal(200))

				// Export the data out of Hoverfly
				exportedRecords := ExportHoverflyRecords()

				// Wipe the records in Hoverfly
				EraseHoverflyRecords()

				// Import the same data into Hoverfly
				ImportHoverflyRecords(exportedRecords)

				// Switch Hoverfly to simulate mode
				SetHoverflyMode(hoverfly.SimulateMode)

				// Make the request to Hoverfly simulate
				afterImportFakeServerRequest := sling.New().Get(fakeServerUrl).Set("Accept-Encoding", "gzip")
				afterImportFakeServerResponse = DoRequestThroughProxy(afterImportFakeServerRequest)
			})


			It("Returns a status code of 200", func() {
				Expect(afterImportFakeServerResponse.StatusCode).To(Equal(200))
			})

			It("Returns hello world", func() {
				responseBody, _ := ioutil.ReadAll(afterImportFakeServerResponse.Body)
				Expect(string(responseBody)).To(Equal(GzipString(`hello_world`)))
			})

			It("Returns with Hoverfly header", func() {
				Expect(afterImportFakeServerResponse.Header).To(HaveKeyWithValue("Hoverfly", []string{"Was-Here"}))
			})

			It("Returns with gzip Content-Encoding header", func() {
				Expect(afterImportFakeServerResponse.Header).To(HaveKeyWithValue("Content-Encoding", []string{"gzip"}))
			})

			It("Returns with text/plain Content-Type header", func() {
				Expect(afterImportFakeServerResponse.Header).To(HaveKeyWithValue("Content-Type", []string{"text/plain"}))
			})

			It("Returns with a Hoverfly header", func() {
				Expect(afterImportFakeServerResponse.Header).To(HaveKeyWithValue("Hoverfly", []string{"Was-Here"}))
			})
		})
	})
})
