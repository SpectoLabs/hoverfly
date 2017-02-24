package hoverfly_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I run Hoverfly", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
		hoverfly.Start()
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("When running in capture mode", func() {

		BeforeEach(func() {
			hoverfly.SetMode("capture")
		})

		Context("without middleware", func() {

			It("Should capture the request and response", func() {

				fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/plain")
					w.Header().Set("Date", "date")
					w.Write([]byte("Hello world"))
				}))

				defer fakeServer.Close()

				resp := hoverfly.Proxy(sling.New().Get(fakeServer.URL))
				Expect(resp.StatusCode).To(Equal(200))

				expectedDestination := strings.Replace(fakeServer.URL, "http://", "", 1)

				recordsJson, err := ioutil.ReadAll(hoverfly.GetSimulation())
				Expect(err).To(BeNil())

				payload := v2.SimulationView{}

				json.Unmarshal(recordsJson, &payload)
				Expect(payload.DataView.RequestResponsePairs).To(HaveLen(1))

				Expect(payload.DataView.RequestResponsePairs[0].Request).To(Equal(v2.RequestDetailsView{
					Path:        util.StringToPointer("/"),
					Method:      util.StringToPointer("GET"),
					Destination: util.StringToPointer(expectedDestination),
					Scheme:      util.StringToPointer("http"),
					Query:       util.StringToPointer(""),
					Body:        util.StringToPointer(""),
					Headers: map[string][]string{
						"Accept-Encoding": []string{"gzip"},
						"User-Agent":      []string{"Go-http-client/1.1"},
					},
				}))

				Expect(payload.DataView.RequestResponsePairs[0].Response).To(Equal(v2.ResponseDetailsView{
					Status:      200,
					Body:        "Hello world",
					EncodedBody: false,
					Headers: map[string][]string{
						"Content-Length": []string{"11"},
						"Content-Type":   []string{"text/plain"},
						"Date":           []string{"date"},
						"Hoverfly":       []string{"Was-Here"},
					},
				}))
			})

			It("Should capture a redirect response", func() {

				fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/plain")
					w.Header().Set("Date", "date")
					w.Write([]byte("redirection got you here"))
				}))
				fakeServerUrl, _ := url.Parse(fakeServer.URL)
				fakeServerRedirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Location", fakeServer.URL)
					w.Header().Set("Content-Type", "text/plain")
					w.WriteHeader(301)
				}))
				fakeServerRedirectUrl, _ := url.Parse(fakeServerRedirect.URL)

				defer fakeServer.Close()
				defer fakeServerRedirect.Close()

				resp := hoverfly.Proxy(sling.New().Get(fakeServerRedirect.URL))
				Expect(resp.StatusCode).To(Equal(301))

				expectedRedirectDestination := strings.Replace(fakeServerRedirectUrl.String(), "http://", "", 1)

				recordsJson, err := ioutil.ReadAll(hoverfly.GetSimulation())
				Expect(err).To(BeNil())

				payload := v2.SimulationView{}

				json.Unmarshal(recordsJson, &payload)
				Expect(payload.DataView.RequestResponsePairs).To(HaveLen(1))

				Expect(payload.DataView.RequestResponsePairs[0].Request.Destination).To(Equal(&expectedRedirectDestination))

				Expect(payload.DataView.RequestResponsePairs[0].Response.Status).To(Equal(301))
				Expect(payload.DataView.RequestResponsePairs[0].Response.Headers["Location"][0]).To(Equal(fakeServerUrl.String()))
			})

			It("Should capture a request body from POST", func() {

				var capturedRequestBody string

				fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					requestBody, err := util.GetRequestBody(r)
					Expect(err).To(BeNil())
					capturedRequestBody = requestBody

					w.Write([]byte("okay"))
				}))

				defer fakeServer.Close()

				resp := hoverfly.Proxy(sling.New().Post(fakeServer.URL).Body(bytes.NewBuffer([]byte(`{"title": "a todo"}`))))
				Expect(resp.StatusCode).To(Equal(200))

				Expect(capturedRequestBody).To(Equal(`{"title": "a todo"}`))
			})
		})
	})
})
