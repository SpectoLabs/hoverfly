package hoverfly_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	functional_tests "github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/SpectoLabs/hoverfly/functional-tests/testdata"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manage journal indexing in hoverfly", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("get templated journal response", func() {

		Context("hoverfly with journal indexing with query params", func() {

			BeforeEach(func() {
				hoverfly.Start("-journal-indexing-key", "Request.QueryParam.id")
			})

			It("Should return templated journal response", func() {
				hoverfly.SetMode("capture")

				fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.Header().Set("Date", "date")
					w.Write([]byte("{\"name\":\"Application Testing\"}"))
				}))

				resp := hoverfly.Proxy(sling.New().Get(fakeServer.URL + "?id=123"))
				Expect(resp.StatusCode).To(Equal(200))

				hoverfly.SetMode("simulate")
				hoverfly.ImportSimulation(testdata.JournalTemplatingWithQueryParamIndexEnabled)

				simulationResponse := hoverfly.Proxy(sling.New().Get("http://test-server.com/journaltest"))
				Expect(resp.StatusCode).To(Equal(200))

				body, err := io.ReadAll(simulationResponse.Body)
				Expect(err).To(BeNil())

				Expect(string(body)).To(Equal("Application Testing"))

				// hasJournalKey function should return true
				simulationResponse = hoverfly.Proxy(sling.New().Get("http://test-server.com/checkJournalKey"))
				Expect(simulationResponse.StatusCode).To(Equal(200))

				body, err = io.ReadAll(simulationResponse.Body)
				Expect(err).To(BeNil())

				Expect(string(body)).To(Equal("123: true 345: false"))

			})
		})

		Context("hoverfly with journal indexing with body", func() {

			BeforeEach(func() {
				hoverfly.Start("-journal-indexing-key", "Request.Body 'jsonpath' '$.id'")
			})

			It("Should return templated journal response", func() {
				hoverfly.SetMode("capture")

				fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.Header().Set("Date", "date")
					w.Write([]byte("{\"name\":\"Application Testing\"}"))
				}))

				defer fakeServer.Close()

				postData := "{\"id\":\"1234\"}"

				resp := hoverfly.Proxy(sling.New().Post(fakeServer.URL).Body(strings.NewReader(postData)))
				Expect(resp.StatusCode).To(Equal(200))

				hoverfly.SetMode("simulate")
				hoverfly.ImportSimulation(testdata.JournalTemplatingWithBodyIndexEnabled)

				simulationResponse := hoverfly.Proxy(sling.New().Get("http://test-server.com/journaltest"))
				Expect(resp.StatusCode).To(Equal(200))

				body, err := io.ReadAll(simulationResponse.Body)
				Expect(err).To(BeNil())

				Expect(string(body)).To(Equal("Application Testing"))

			})
		})
	})
})
