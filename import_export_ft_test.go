package hoverfly_test

import (
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"fmt"
	"io/ioutil"
	"github.com/SpectoLabs/hoverfly"
	"compress/gzip"
	"bytes"
	"os"
)

// Helper function for gzipping strings
func GzipString(s string) (string) {
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
				// Spin up a fake server which returns hello world
				fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

					w.Header().Set("Content-Type", "text/plain")
					w.WriteHeader(200)
					fmt.Fprintf(w, "hello_world")
				}))
				defer fakeServer.Close()

				// Switch Hoverfly to capture mode
				SetHoverflyMode(hoverfly.CaptureMode)

				// Make a request to the fake server and proxy through Hoverfly
				fakeServerRequest := sling.New().Get(fakeServer.URL)

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
				afterImportFakeServerRequest := sling.New().Get(fakeServer.URL)
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

		Context("A gzipped response", func() {

			var afterImportFakeServerResponse *http.Response

			BeforeEach(func() {
				// Spin up a fake server which returns hello world gzipped
				fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

					w.Header().Set("Content-Encoding", "gzip")
					w.Header().Set("Content-Type", "text/plain")
					gzipWriter := gzip.NewWriter(w)
					gzipWriter.Write([]byte(`hello_world`))
				}))
				defer fakeServer.Close()

				// Switch Hoverfly to capture mode
				SetHoverflyMode(hoverfly.CaptureMode)

				// Make a request to the fake server and proxy through Hoverfly
				fakeServerRequest := sling.New().Get(fakeServer.URL).Set("Accept-Encoding", "gzip")

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
				afterImportFakeServerRequest := sling.New().Get(fakeServer.URL).Set("Accept-Encoding", "gzip")
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

		Context("An image response", func() {

			var afterImportFakeServerResponse *http.Response


			pwd, _ := os.Getwd()
			imageUri := "/testdata/1x1.png"

			BeforeEach(func() {
				// Spin up a fake server which returns hello world gzipped
				fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "image/jpeg")
					w.WriteHeader(200)
					http.ServeFile(w, r, pwd + imageUri)
				}))

				defer fakeServer.Close()
				fmt.Println(fakeServer.URL)
				//time.Sleep(time.Second * 3)

				// Switch Hoverfly to capture mode
				SetHoverflyMode(hoverfly.CaptureMode)

				// Make a request to the fake server and proxy through Hoverfly
				fakeServerRequest := sling.New().Get(fakeServer.URL)
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
				afterImportFakeServerRequest := sling.New().Get(fakeServer.URL)
				afterImportFakeServerResponse = DoRequestThroughProxy(afterImportFakeServerRequest)
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
