package hoverfly_test

import (
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	//"compress/gzip"
	"strings"
	"fmt"
	"io/ioutil"
)

var _ = Describe("Capturing, exporting, importing and simulating", func() {


	Describe("Import, Export", func() {
		Context("The captured response should be returned after exporting and importing", func() {

			var afterImportFakeServerResponse *http.Response

			BeforeEach(func() {
				// Spin up a fake server which returns hello world
				fakeGzipServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

					w.WriteHeader(200)
					w.Header().Set("Content-Type", "text/plain")
					fmt.Fprintf(w, "hello_world")
					//w.Header().Set("Content-Encoding", "gzip")
					//gzipWriter := gzip.NewWriter(w)
					//gzipWriter.Write([]byte(`hello_world`))
				}))

				// Switch Hoverfly to capture mode
				hoverflyCaptureModeRequest := sling.New().Post(hoverflyAdminUrl + "/api/state").Body(strings.NewReader(`{"mode":"capture"}`))
				hoverflyCaptureModeResponse := DoRequest(hoverflyCaptureModeRequest)
				Expect(hoverflyCaptureModeResponse.StatusCode).To(Equal(200))

				// Make a request to the fake server and proxy through Hoverfly
				fakeServerUrl := fakeGzipServer.URL

				fakeServerRequest := sling.New().Get(fakeServerUrl)

				response := DoRequestThroughProxy(fakeServerRequest)
				Expect(response.StatusCode).To(Equal(200))

				// Export the data out of Hoverfly
				hoverflyRecordsRequest := sling.New().Get(hoverflyAdminUrl + "/api/records")

				hoverflyRecordsResponse := DoRequest(hoverflyRecordsRequest)
				Expect(hoverflyRecordsResponse.StatusCode).To(Equal(200))

				// Wipe the records in Hoverfly
				hoverflyDeleteRequest := sling.New().Delete(hoverflyAdminUrl + "/api/records")

				hoverflyDeleteResponse := DoRequest(hoverflyDeleteRequest)
				Expect(hoverflyDeleteResponse.StatusCode).To(Equal(200))

				// Import the same data into Hoverfly
				hoverflyImportRequest := sling.New().Post(hoverflyAdminUrl + "/api/records").Body(hoverflyRecordsResponse.Body)

				hoverflyImportResponse := DoRequest(hoverflyImportRequest)
				Expect(hoverflyImportResponse.StatusCode).To(Equal(200))

				// Switch Hoverfly to simulate mode

				hoverflySimulateModeRequest := sling.New().Post(hoverflyAdminUrl + "/api/state").Body(strings.NewReader(`{"mode":"simulate"}`))

				hoverflySimulateModeResponse := DoRequest(hoverflySimulateModeRequest)
				Expect(hoverflySimulateModeResponse.StatusCode).To(Equal(200))

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

			//It("Returns with gzip Content-Encoding header", func() {
			//	Expect(afterImportFakeServerResponse.Header).To(HaveKeyWithValue("Content-Encoding", "gzip"))
			//})

			It("Returns with text/plain Content-Type header", func() {
				Expect(afterImportFakeServerResponse.Header).To(HaveKeyWithValue("Content-Type", []string{"text/plain; charset=utf-8"}))
			})

			It("Returns with a Hoverfly header", func() {
				Expect(afterImportFakeServerResponse.Header).To(HaveKeyWithValue("Hoverfly", []string{"Was-Here"}))
			})
		})
	})
})
