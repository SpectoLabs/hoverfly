package hoverfly_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"

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
	})

	Context("and it uses default certificate and key configuration", func() {

		BeforeEach(func() {
			hoverfly.Start()
			hoverfly.SetMode("capture")
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("should respond with HTTPS responses with the default Hoverfly certificate", func() {

			fakeServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.Header().Set("Date", "date")
				w.Write([]byte("Hello world"))
			}))

			defer fakeServer.Close()

			response := hoverfly.Proxy(sling.New().Get(fakeServer.URL))

			Expect(response.TLS.PeerCertificates[0].Issuer.CommonName).To(Equal("hoverfly.proxy"))
			Expect(response.TLS.PeerCertificates[0].Issuer.Organization).To(ContainElement("Hoverfly Authority"))
			Expect(response.TLS.PeerCertificates[0].Subject.Names[0].Value).To(Equal("GoProxy untrusted MITM proxy Inc"))
		})
	})

	Context("it sends the correct headers", func() {

		BeforeEach(func() {
			hoverfly.Start()
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("should set Content-Length if empty", func() {

			hoverfly.ImportSimulation(`{
				"data": {
					"pairs": [
						{
							"request": {
								"scheme": {
									"exactMatch": "https"
								},
								"path": {
									"exactMatch": "/path"
								}
							},
							"response": {
								"status": 200,
								"body": "OK"
							}
						}
					]
				},
				"meta": {
					"schemaVersion": "v3"
				}
			}`)
			response := hoverfly.Proxy(sling.New().Get("https://hoverfly.io/path"))
			Expect(response.StatusCode).To(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(response.Body)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("OK"))

			// These will always be empty as they are excluded  by net/http
			Expect(response.Header.Get("Content-Length")).To(Equal(""))
			Expect(response.Header.Get("Transfer-Encoding")).To(Equal(""))
		})

		It("should not set Content-Length if not empty", func() {

			hoverfly.ImportSimulation(`{
				"data": {
					"pairs": [
						{
							"request": {
								"scheme": {
									"exactMatch": "https"
								},
								"path": {
									"exactMatch": "/path"
								}
							},
							"response": {
								"status": 200,
								"body": "OK",
								"headers": {
									"Content-Length": ["5555"]
								}
							}
						}
					]
				},
				"meta": {
					"schemaVersion": "v3"
				}
			}`)
			response := hoverfly.Proxy(sling.New().Get("https://hoverfly.io/path"))
			Expect(response.StatusCode).To(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(response.Body)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("OK"))

			// These will always be empty as they are excluded by net/http
			Expect(response.Header.Get("Content-Length")).To(Equal(""))
			Expect(response.Header.Get("Transfer-Encoding")).To(Equal(""))
		})
	})

	Context("and it uses default certificate and key configuration", func() {

		BeforeEach(func() {
			hoverfly.Start("-cert", "testdata/cert.pem", "-key", "testdata/key.pem")
			hoverfly.SetMode("capture")
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("should respond with HTTPS responses with the default Hoverfly certificate", func() {

			fakeServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.Header().Set("Date", "date")
				w.Write([]byte("Hello world"))
			}))

			defer fakeServer.Close()

			response := hoverfly.Proxy(sling.New().Get(fakeServer.URL))

			Expect(response.TLS.PeerCertificates[0].Issuer.CommonName).To(Equal("test.cert"))
			Expect(response.TLS.PeerCertificates[0].Issuer.Organization).To(ContainElement("Testdata Certificate Authority"))
			Expect(response.TLS.PeerCertificates[0].Subject.Names[0].Value).To(Equal("GoProxy untrusted MITM proxy Inc"))
		})
	})
})
