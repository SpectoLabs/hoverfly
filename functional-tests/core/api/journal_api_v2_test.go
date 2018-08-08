package api_test

import (
	"bytes"
	"net/http"

	"time"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/api/v2/journal", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	Context("With journal enabled", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start()
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		Context("GET", func() {

			It("should display an empty journal", func() {
				req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
				res := functional_tests.DoRequest(req)

				Expect(res.StatusCode).To(Equal(200))

				var journalView v2.JournalView

				functional_tests.UnmarshalFromResponse(res, &journalView)

				Expect(journalView.Journal).To(HaveLen(0))
			})

			It("should display one item in the journal", func() {
				hoverfly.Proxy(sling.New().Get("http://hoverfly.io"))

				req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
				res := functional_tests.DoRequest(req)

				Expect(res.StatusCode).To(Equal(200))

				var journalView v2.JournalView

				functional_tests.UnmarshalFromResponse(res, &journalView)

				Expect(journalView.Journal).To(HaveLen(1))

				Expect(*journalView.Journal[0].Request.Scheme).To(Equal("http"))
				Expect(*journalView.Journal[0].Request.Method).To(Equal("GET"))
				Expect(*journalView.Journal[0].Request.Destination).To(Equal("hoverfly.io"))
				Expect(*journalView.Journal[0].Request.Path).To(Equal("/"))
				Expect(*journalView.Journal[0].Request.Query).To(Equal(""))
				Expect(journalView.Journal[0].Request.Headers["Accept-Encoding"]).To(ContainElement("gzip"))
				Expect(journalView.Journal[0].Request.Headers["User-Agent"]).To(ContainElement("Go-http-client/1.1"))

				Expect(journalView.Journal[0].Response.Status).To(Equal(502))
				Expect(journalView.Journal[0].Response.Body).To(Equal("Hoverfly Error!\n\nThere was an error when matching\n\nGot error: Could not find a match for request, create or record a valid matcher first!"))
				Expect(journalView.Journal[0].Response.Headers["Content-Type"]).To(ContainElement("text/plain"))

				Expect(journalView.Journal[0].Latency).To(BeNumerically("<", time.Millisecond))
				Expect(journalView.Journal[0].Mode).To(Equal("simulate"))
			})

			It("should display multiple items in the journal", func() {
				hoverfly.Proxy(sling.New().Get("http://hoverfly.io"))
				hoverfly.Proxy(sling.New().Get("http://github.com/SpectoLabs/hoverfly"))
				hoverfly.Proxy(sling.New().Get("http://specto.io"))

				req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
				res := functional_tests.DoRequest(req)

				Expect(res.StatusCode).To(Equal(200))

				var journalView v2.JournalView

				functional_tests.UnmarshalFromResponse(res, &journalView)

				Expect(journalView.Journal).To(HaveLen(3))

				Expect(*journalView.Journal[0].Request.Destination).To(Equal("hoverfly.io"))
				Expect(*journalView.Journal[0].Request.Path).To(Equal("/"))

				Expect(*journalView.Journal[1].Request.Destination).To(Equal("github.com"))
				Expect(*journalView.Journal[1].Request.Path).To(Equal("/SpectoLabs/hoverfly"))

				Expect(*journalView.Journal[2].Request.Destination).To(Equal("specto.io"))
				Expect(*journalView.Journal[2].Request.Path).To(Equal("/"))
			})

			It("should display the mode each request was in", func() {
				hoverfly.SetMode("simulate")
				hoverfly.Proxy(sling.New().Get("http://localhost:" + hoverfly.GetAdminPort()))

				hoverfly.SetMode("capture")
				hoverfly.Proxy(sling.New().Get("http://localhost:" + hoverfly.GetAdminPort()))

				req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
				res := functional_tests.DoRequest(req)

				Expect(res.StatusCode).To(Equal(200))

				var journalView v2.JournalView

				functional_tests.UnmarshalFromResponse(res, &journalView)

				Expect(journalView.Journal).To(HaveLen(2))

				Expect(journalView.Journal[0].Mode).To(Equal("simulate"))
				Expect(journalView.Journal[1].Mode).To(Equal("capture"))
			})
		})

		Context("POST", func() {

			BeforeEach(func() {
				hoverfly.Proxy(sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/first"))
				hoverfly.Proxy(sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/second"))
			})

			It("should filter", func() {
				req := sling.New().Post("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
				req.Body(bytes.NewBufferString(`{
					"request": {
						"path": [
							{
								"matcher": "exact",
								"value": "/first"
							}
						]
					}
				}`))
				res := functional_tests.DoRequest(req)

				Expect(res.StatusCode).To(Equal(http.StatusOK))

				var journalView v2.JournalView

				functional_tests.UnmarshalFromResponse(res, &journalView)

				Expect(journalView.Journal).To(HaveLen(1))
				Expect(*journalView.Journal[0].Request.Path).To(Equal("/first"))
			})

			It("should error when body is malformed JSON", func() {
				req := sling.New().Post("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
				req.Body(bytes.NewBufferString(`not json`))
				res := functional_tests.DoRequest(req)

				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))

				var errorView handlers.ErrorView

				functional_tests.UnmarshalFromResponse(res, &errorView)

				Expect(errorView.Error).To(Equal("Malformed JSON"))
			})

			It("should error when body has no request", func() {
				req := sling.New().Post("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
				req.Body(bytes.NewBufferString(`{"norequest": true}`))
				res := functional_tests.DoRequest(req)

				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))

				var errorView handlers.ErrorView

				functional_tests.UnmarshalFromResponse(res, &errorView)

				Expect(errorView.Error).To(Equal("No \"request\" object in search parameters"))
			})
		})

		Context("DELETE", func() {
			It("should delete journal entries", func() {
				hoverfly.Proxy(sling.New().Get("http://localhost:" + hoverfly.GetAdminPort()))
				hoverfly.Proxy(sling.New().Get("http://localhost:" + hoverfly.GetAdminPort()))

				req := sling.New().Delete("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
				functional_tests.DoRequest(req)

				req = sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
				res := functional_tests.DoRequest(req)

				Expect(res.StatusCode).To(Equal(200))

				var journalView v2.JournalView

				functional_tests.UnmarshalFromResponse(res, &journalView)

				Expect(journalView.Journal).To(HaveLen(0))
			})
		})
	})

	Context("With journal enabled and Hoverfly as a webserver", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start("-webserver")
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		Context("GET", func() {

			It("should display an empty journal", func() {
				req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
				res := functional_tests.DoRequest(req)

				Expect(res.StatusCode).To(Equal(200))

				var journalView v2.JournalView

				functional_tests.UnmarshalFromResponse(res, &journalView)

				Expect(journalView.Journal).To(HaveLen(0))
			})

			It("should display one item in the journal", func() {
				functional_tests.DoRequest(sling.New().Get("http://localhost:" + hoverfly.GetProxyPort()))

				req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
				res := functional_tests.DoRequest(req)

				Expect(res.StatusCode).To(Equal(200))

				var journalView v2.JournalView

				functional_tests.UnmarshalFromResponse(res, &journalView)

				Expect(journalView.Journal).To(HaveLen(1))

				Expect(*journalView.Journal[0].Request.Scheme).To(Equal("http"))
				Expect(*journalView.Journal[0].Request.Method).To(Equal("GET"))
				Expect(*journalView.Journal[0].Request.Destination).To(Equal("localhost:" + hoverfly.GetProxyPort()))
				Expect(*journalView.Journal[0].Request.Path).To(Equal("/"))
				Expect(*journalView.Journal[0].Request.Query).To(Equal(""))
				Expect(journalView.Journal[0].Request.Headers["Accept-Encoding"]).To(ContainElement("gzip"))
				Expect(journalView.Journal[0].Request.Headers["User-Agent"]).To(ContainElement("Go-http-client/1.1"))

				Expect(journalView.Journal[0].Response.Status).To(Equal(502))
				Expect(journalView.Journal[0].Response.Body).To(Equal("Hoverfly Error!\n\nThere was an error when matching\n\nGot error: Could not find a match for request, create or record a valid matcher first!"))
				Expect(journalView.Journal[0].Response.Headers["Content-Type"]).To(ContainElement("text/plain"))

				Expect(journalView.Journal[0].Latency).To(BeNumerically("<", time.Millisecond))
				Expect(journalView.Journal[0].Mode).To(Equal("simulate"))
			})

			It("should display multiple items in the journal", func() {
				functional_tests.DoRequest(sling.New().Get("http://localhost:" + hoverfly.GetProxyPort() + "/first"))
				functional_tests.DoRequest(sling.New().Get("http://localhost:" + hoverfly.GetProxyPort() + "/second"))
				functional_tests.DoRequest(sling.New().Get("http://localhost:" + hoverfly.GetProxyPort() + "/third"))

				req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
				res := functional_tests.DoRequest(req)

				Expect(res.StatusCode).To(Equal(200))

				var journalView v2.JournalView

				functional_tests.UnmarshalFromResponse(res, &journalView)

				Expect(journalView.Journal).To(HaveLen(3))

				Expect(*journalView.Journal[0].Request.Path).To(ContainSubstring("/first"))

				Expect(*journalView.Journal[1].Request.Path).To(ContainSubstring("/second"))

				Expect(*journalView.Journal[2].Request.Path).To(ContainSubstring("/third"))
			})
		})

		Context("POST", func() {

			BeforeEach(func() {
				functional_tests.DoRequest(sling.New().Get("http://localhost:" + hoverfly.GetProxyPort() + "/first"))
				functional_tests.DoRequest(sling.New().Get("http://localhost:" + hoverfly.GetProxyPort() + "/second"))
			})

			It("should filter", func() {
				req := sling.New().Post("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
				req.Body(bytes.NewBufferString(`{
					"request": {
						"path": [
							{
								"matcher": "exact",
								"value": "/first"
							}
						]
					}
				}`))
				res := functional_tests.DoRequest(req)

				Expect(res.StatusCode).To(Equal(http.StatusOK))

				var journalView v2.JournalView

				functional_tests.UnmarshalFromResponse(res, &journalView)

				Expect(journalView.Journal).To(HaveLen(1))
				Expect(*journalView.Journal[0].Request.Path).To(Equal("/first"))
			})

			It("should error when body is malformed JSON", func() {
				req := sling.New().Post("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
				req.Body(bytes.NewBufferString(`not json`))
				res := functional_tests.DoRequest(req)

				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))

				var errorView handlers.ErrorView

				functional_tests.UnmarshalFromResponse(res, &errorView)

				Expect(errorView.Error).To(Equal("Malformed JSON"))
			})

			It("should error when body has no request", func() {
				req := sling.New().Post("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
				req.Body(bytes.NewBufferString(`{"norequest": true}`))
				res := functional_tests.DoRequest(req)

				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))

				var errorView handlers.ErrorView

				functional_tests.UnmarshalFromResponse(res, &errorView)

				Expect(errorView.Error).To(Equal("No \"request\" object in search parameters"))
			})
		})

		Context("DELETE", func() {
			It("should delete journal entries", func() {
				hoverfly.Proxy(sling.New().Get("http://localhost:" + hoverfly.GetAdminPort()))
				hoverfly.Proxy(sling.New().Get("http://localhost:" + hoverfly.GetAdminPort()))

				req := sling.New().Delete("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
				functional_tests.DoRequest(req)

				req = sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
				res := functional_tests.DoRequest(req)

				Expect(res.StatusCode).To(Equal(200))

				var journalView v2.JournalView

				functional_tests.UnmarshalFromResponse(res, &journalView)

				Expect(journalView.Journal).To(HaveLen(0))
			})
		})

	})

	Context("With journal disabled", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start("-journal-size=0")
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		Context("GET", func() {

			It("should return an error", func() {
				req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
				res := functional_tests.DoRequest(req)

				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))

				var errorView handlers.ErrorView

				functional_tests.UnmarshalFromResponse(res, &errorView)

				Expect(errorView.Error).To(Equal("Journal disabled"))
			})
		})

		Context("POST", func() {

			It("should return an error", func() {
				req := sling.New().Post("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
				req.Body(bytes.NewBufferString(`{
					"request": {
						"path": [
							{
								"matcher": "exact",
								"value": "/first"
							}
						]
					}
				}`))
				res := functional_tests.DoRequest(req)

				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))

				var errorView handlers.ErrorView

				functional_tests.UnmarshalFromResponse(res, &errorView)

				Expect(errorView.Error).To(Equal("Journal disabled"))
			})
		})

		Context("DELETE", func() {

			It("should return an error", func() {
				req := sling.New().Delete("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
				res := functional_tests.DoRequest(req)

				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))

				var errorView handlers.ErrorView

				functional_tests.UnmarshalFromResponse(res, &errorView)

				Expect(errorView.Error).To(Equal("Journal disabled"))
			})
		})
	})

	Context("with -journal-size=10", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start("-journal-size=10")
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		Context("GET", func() {

			It("should not exceed size", func() {
				for i := 0; i < 10; i++ {
					hoverfly.Proxy(sling.New().Get("http://hoverfly.io"))
				}

				req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
				res := functional_tests.DoRequest(req)

				Expect(res.StatusCode).To(Equal(200))

				var journalView v2.JournalView

				functional_tests.UnmarshalFromResponse(res, &journalView)

				Expect(journalView.Journal).To(HaveLen(10))
			})
		})
	})
})
