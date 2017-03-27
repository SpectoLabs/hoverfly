package hoverfly_test

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Helper function for gzipping strings
func GzipString(s string) string {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	gz.Write([]byte(s))
	return b.String()
}

var _ = Describe("Capture > export > importing > simulate flow", func() {

	Describe("When I import and export", func() {

		Context("A plain text response", func() {

			var afterImportFakeServerResponse *http.Response

			BeforeEach(func() {
				// Start hoverfly
				hoverflyCmd = startHoverfly(adminPort, proxyPort)

				// Spin up a fake server which returns hello world
				fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

					w.Header().Set("Content-Type", "text/plain")
					w.WriteHeader(200)
					fmt.Fprintf(w, "hello_world")
				}))
				defer fakeServer.Close()

				// Switch Hoverfly to capture mode
				SetHoverflyMode("capture")

				// Make a request to the fake server and proxy through Hoverfly
				response := CallFakeServerThroughProxy(fakeServer)
				Expect(response.StatusCode).To(Equal(200))

				// Export the data out of Hoverfly
				exportedRecords := ExportHoverflySimulation()

				// Wipe the records in Hoverfly
				EraseHoverflySimulation()

				// Import the same data into Hoverfly
				ImportHoverflySimulation(exportedRecords)

				// Switch Hoverfly to simulate mode
				SetHoverflyMode("simulate")

				// Make the request to Hoverfly simulate
				afterImportFakeServerResponse = CallFakeServerThroughProxy(fakeServer)
			})

			AfterEach(func() {
				// Stop hoverfly
				stopHoverfly()
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

		Context("A gzipped response", func() {

			var afterImportFakeServerResponse *http.Response

			BeforeEach(func() {
				// Start hoverfly
				hoverflyCmd = startHoverfly(adminPort, proxyPort)

				// Spin up a fake server which returns hello world gzipped
				fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

					w.Header().Set("Content-Encoding", "gzip")
					w.Header().Set("Content-Type", "text/plain")
					gzipWriter := gzip.NewWriter(w)
					gzipWriter.Write([]byte(`hello_world`))
				}))
				defer fakeServer.Close()

				// Switch Hoverfly to capture mode
				SetHoverflyMode("capture")

				// Make a request to the fake server and proxy through Hoverfly
				fakeServerRequest := sling.New().Get(fakeServer.URL).Set("Accept-Encoding", "gzip")

				response := DoRequestThroughProxy(fakeServerRequest)
				Expect(response.StatusCode).To(Equal(200))

				// Export the data out of Hoverfly
				exportedSimulation := ExportHoverflySimulation()

				// Wipe the records in Hoverfly
				EraseHoverflySimulation()

				// Import the same data into Hoverfly
				ImportHoverflySimulation(exportedSimulation)

				// Switch Hoverfly to simulate mode
				SetHoverflyMode("simulate")

				// Make the request to Hoverfly simulate
				afterImportFakeServerRequest := sling.New().Get(fakeServer.URL).Set("Accept-Encoding", "gzip")
				afterImportFakeServerResponse = DoRequestThroughProxy(afterImportFakeServerRequest)
			})

			AfterEach(func() {
				// Stop hoverfly
				stopHoverfly()
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

		Context("An image response", func() {

			var afterImportFakeServerResponse *http.Response

			pwd, _ := os.Getwd()
			imageUri := "/testdata/1x1.png"

			BeforeEach(func() {
				// Start hoverfly
				hoverflyCmd = startHoverfly(adminPort, proxyPort)

				// Spin up a fake server which returns hello world gzipped
				fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "image/jpeg")
					w.WriteHeader(200)
					http.ServeFile(w, r, pwd+imageUri)
				}))

				defer fakeServer.Close()
				//time.Sleep(time.Second * 3)

				// Switch Hoverfly to capture mode
				SetHoverflyMode("capture")

				// Make a request to the fake server and proxy through Hoverfly
				fakeServerRequest := sling.New().Get(fakeServer.URL)
				response := DoRequestThroughProxy(fakeServerRequest)
				Expect(response.StatusCode).To(Equal(200))

				// Export the data out of Hoverfly
				exportedRecords := ExportHoverflySimulation()

				// Wipe the records in Hoverfly
				EraseHoverflySimulation()

				// Import the same data into Hoverfly
				ImportHoverflySimulation(exportedRecords)

				// Switch Hoverfly to simulate mode
				SetHoverflyMode("simulate")

				// Make the request to Hoverfly simulate
				afterImportFakeServerRequest := sling.New().Get(fakeServer.URL)
				afterImportFakeServerResponse = DoRequestThroughProxy(afterImportFakeServerRequest)
			})

			AfterEach(func() {
				// Stop hoverfly
				stopHoverfly()
			})

			It("Returns a status code of 200", func() {
				Expect(afterImportFakeServerResponse.StatusCode).To(Equal(200))
			})

			It("Returns image", func() {
				file, _ := os.Open(pwd + imageUri)
				defer file.Close()
				returnedImageBytes, _ := ioutil.ReadAll(afterImportFakeServerResponse.Body)
				originalImageBytes, _ := ioutil.ReadAll(file)
				Expect(returnedImageBytes).To(Equal(originalImageBytes))
			})

			It("Returns with Hoverfly header", func() {
				Expect(afterImportFakeServerResponse.Header).To(HaveKeyWithValue("Hoverfly", []string{"Was-Here"}))
			})

			It("Returns with image/jpeg Content-Type header", func() {
				Expect(afterImportFakeServerResponse.Header).To(HaveKeyWithValue("Content-Type", []string{"image/jpeg"}))
			})

			It("Returns with a Hoverfly header", func() {
				Expect(afterImportFakeServerResponse.Header).To(HaveKeyWithValue("Hoverfly", []string{"Was-Here"}))
			})
		})
	})
})
