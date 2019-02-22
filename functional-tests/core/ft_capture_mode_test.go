package hoverfly_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
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

				payload := hoverfly.ExportSimulation()

				Expect(payload.RequestResponsePairs).To(HaveLen(1))

				Expect(payload.RequestResponsePairs[0].RequestMatcher).To(Equal(v2.RequestMatcherViewV5{
					Path: []v2.MatcherViewV5{
						{
							Matcher: matchers.Exact,
							Value:   "/",
						},
					},
					Method: []v2.MatcherViewV5{
						{
							Matcher: matchers.Exact,
							Value:   "GET",
						},
					},
					Destination: []v2.MatcherViewV5{
						{
							Matcher: matchers.Exact,
							Value:   expectedDestination,
						},
					},
					Scheme: []v2.MatcherViewV5{
						{
							Matcher: matchers.Exact,
							Value:   "http",
						},
					},
					Body: []v2.MatcherViewV5{
						{
							Matcher: matchers.Exact,
							Value:   "",
						},
					},
					Query: &v2.QueryMatcherViewV5{},
				}))

				Expect(payload.RequestResponsePairs[0].Response).To(Equal(v2.ResponseDetailsViewV5{
					Status:      200,
					Body:        "Hello world",
					EncodedBody: false,
					Headers: map[string][]string{
						"Content-Length": {"11"},
						"Content-Type":   {"text/plain"},
						"Date":           {"date"},
						"Hoverfly":       {"Was-Here"},
					},
					Templated: false,
				}))
			})

			It("Should capture all request headers if argument is set to *", func() {
				hoverfly.SetModeWithArgs("capture", v2.ModeArgumentsView{
					Headers: []string{"*"},
				})

				fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/plain")
					w.Header().Set("Date", "date")
					w.Write([]byte("Hello world"))
				}))

				defer fakeServer.Close()

				resp := hoverfly.Proxy(sling.New().Get(fakeServer.URL))
				Expect(resp.StatusCode).To(Equal(200))

				payload := hoverfly.ExportSimulation()
				Expect(payload.RequestResponsePairs).To(HaveLen(1))

				Expect(payload.RequestResponsePairs[0].RequestMatcher.Headers).To(Equal(
					map[string][]v2.MatcherViewV5{
						"Accept-Encoding": {
							{
								Matcher: "exact",
								Value:   "gzip",
							},
						},
						"User-Agent": {
							{
								Matcher: "exact",
								Value:   "Go-http-client/1.1",
							},
						},
					}))
			})

			It("Should capture User-Agent request headers if argument is set to User-Agent", func() {
				hoverfly.SetModeWithArgs("capture", v2.ModeArgumentsView{
					Headers: []string{"User-Agent"},
				})

				fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/plain")
					w.Header().Set("Date", "date")
					w.Write([]byte("Hello world"))
				}))

				defer fakeServer.Close()

				resp := hoverfly.Proxy(sling.New().Get(fakeServer.URL))
				Expect(resp.StatusCode).To(Equal(200))

				recordsJson, err := ioutil.ReadAll(hoverfly.GetSimulation())
				Expect(err).To(BeNil())

				payload := v2.SimulationViewV5{}

				Expect(json.Unmarshal(recordsJson, &payload)).To(Succeed())
				Expect(payload.RequestResponsePairs).To(HaveLen(1))

				Expect(payload.RequestResponsePairs[0].RequestMatcher.Headers).To(Equal(
					map[string][]v2.MatcherViewV5{
						"User-Agent": {
							{
								Matcher: "exact",
								Value:   "Go-http-client/1.1",
							},
						},
					}))
			})

			It("Should capture User-Agent and Test request headers if argument is set to User-Agent,Test", func() {
				hoverfly.SetModeWithArgs("capture", v2.ModeArgumentsView{
					Headers: []string{"User-Agent", "Test"},
				})

				fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/plain")
					w.Header().Set("Date", "date")
					w.Write([]byte("Hello world"))
				}))

				defer fakeServer.Close()

				resp := hoverfly.Proxy(sling.New().Get(fakeServer.URL).Add("Test", "value"))
				Expect(resp.StatusCode).To(Equal(200))

				recordsJson, err := ioutil.ReadAll(hoverfly.GetSimulation())
				Expect(err).To(BeNil())

				payload := v2.SimulationViewV5{}

				Expect(json.Unmarshal(recordsJson, &payload)).To(Succeed())
				Expect(payload.RequestResponsePairs).To(HaveLen(1))

				Expect(payload.RequestResponsePairs[0].RequestMatcher.Headers).To(Equal(
					map[string][]v2.MatcherViewV5{
						"User-Agent": {
							{
								Matcher: "exact",
								Value:   "Go-http-client/1.1",
							},
						},

						"Test": {
							{
								Matcher: "exact",
								Value:   "value",
							},
						},
					},
				))
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

				payload := v2.SimulationViewV5{}

				functional_tests.Unmarshal(recordsJson, &payload)
				Expect(payload.RequestResponsePairs).To(HaveLen(1))

				Expect(payload.RequestResponsePairs[0].RequestMatcher.Destination[0].Matcher).To(Equal(`exact`))
				Expect(payload.RequestResponsePairs[0].RequestMatcher.Destination[0].Value).To(Equal(expectedRedirectDestination))

				Expect(payload.RequestResponsePairs[0].Response.Status).To(Equal(301))
				Expect(payload.RequestResponsePairs[0].Response.Headers["Location"][0]).To(Equal(fakeServerUrl.String()))
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

			It("Should capture a JSON request body as a jsonMatch", func() {

				fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("okay"))
				}))

				defer fakeServer.Close()

				resp := hoverfly.Proxy(sling.New().Post(fakeServer.URL).Add("Content-Type", "application/json").Body(bytes.NewBuffer([]byte(`{"title": "a todo"}`))))
				Expect(resp.StatusCode).To(Equal(200))

				recordsJson, err := ioutil.ReadAll(hoverfly.GetSimulation())
				Expect(err).To(BeNil())

				payload := v2.SimulationViewV5{}

				functional_tests.Unmarshal(recordsJson, &payload)
				Expect(payload.RequestResponsePairs).To(HaveLen(1))

				Expect(payload.RequestResponsePairs[0].RequestMatcher.Body[0].Matcher).To(Equal(`json`))
				Expect(payload.RequestResponsePairs[0].RequestMatcher.Body[0].Value).To(Equal(`{"title": "a todo"}`))
			})

			It("Should capture a XML request body as a xmlMatch", func() {

				fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("okay"))
				}))

				defer fakeServer.Close()

				resp := hoverfly.Proxy(sling.New().Post(fakeServer.URL).Add("Content-Type", "application/xml").Body(bytes.NewBuffer([]byte(`<document/>`))))
				Expect(resp.StatusCode).To(Equal(200))

				recordsJson, err := ioutil.ReadAll(hoverfly.GetSimulation())
				Expect(err).To(BeNil())

				payload := v2.SimulationViewV5{}

				functional_tests.Unmarshal(recordsJson, &payload)
				Expect(payload.RequestResponsePairs).To(HaveLen(1))

				Expect(payload.RequestResponsePairs[0].RequestMatcher.Body[0].Matcher).To(Equal(`xml`))
				Expect(payload.RequestResponsePairs[0].RequestMatcher.Body[0].Value).To(Equal(`<document/>`))
			})

			It("Should pass through the original query", func() {

				var capturedRequestQuery string

				fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					capturedRequestQuery = r.URL.RawQuery

					w.Write([]byte("okay"))
				}))

				defer fakeServer.Close()

				request, _ := sling.New().Post(fakeServer.URL + "?z=1&y=2&x=3").Request()
				request.URL.RawQuery = "z=1&y=2&x=3"
				hoverfly.ProxyRequest(request)

				Expect(capturedRequestQuery).To(Equal("z=1&y=2&x=3"))
			})
		})
	})

	Context("When running in capture mode with stateful capturing enabled", func() {

		BeforeEach(func() {
			hoverfly.SetModeWithArgs("capture", v2.ModeArgumentsView{
				Stateful: true,
			})
		})

		It("Should capture duplicate pair", func() {

			statefulServerResponse := 0

			fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				statefulServerResponse = statefulServerResponse + 1
				w.Header().Set("Content-Type", "text/plain")
				w.Header().Set("Date", "date")
				w.Write([]byte(strconv.Itoa(statefulServerResponse)))
			}))

			defer fakeServer.Close()

			resp := hoverfly.Proxy(sling.New().Get(fakeServer.URL))
			Expect(resp.StatusCode).To(Equal(200))

			resp = hoverfly.Proxy(sling.New().Get(fakeServer.URL))
			Expect(resp.StatusCode).To(Equal(200))

			resp = hoverfly.Proxy(sling.New().Get(fakeServer.URL))
			Expect(resp.StatusCode).To(Equal(200))

			payload := hoverfly.ExportSimulation()

			Expect(payload.RequestResponsePairs).To(HaveLen(3))

			Expect(payload.RequestResponsePairs[0].RequestMatcher.RequiresState).To(Equal(map[string]string{"sequence:1": "1"}))
			Expect(payload.RequestResponsePairs[0].Response.Body).To(Equal("1"))
			Expect(payload.RequestResponsePairs[0].Response.TransitionsState).To(Equal(map[string]string{"sequence:1": "2"}))

			Expect(payload.RequestResponsePairs[1].RequestMatcher.RequiresState).To(Equal(map[string]string{"sequence:1": "2"}))
			Expect(payload.RequestResponsePairs[1].Response.Body).To(Equal("2"))
			Expect(payload.RequestResponsePairs[1].Response.TransitionsState).To(Equal(map[string]string{"sequence:1": "3"}))

			Expect(payload.RequestResponsePairs[2].RequestMatcher.RequiresState).To(Equal(map[string]string{"sequence:1": "3"}))
			Expect(payload.RequestResponsePairs[2].Response.Body).To(Equal("3"))
			Expect(payload.RequestResponsePairs[2].Response.TransitionsState).To(BeNil())
		})
	})
})
