package hoverfly_test

import (
	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

	Context("and proxy HTTPS requests", func() {

		BeforeEach(func() {
			hoverfly.Start()
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("should uses chunked transfer encoding by default", func() {

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

			// TODO we need to find a way to keep the response header for assertion
			// These will always be empty as they are excluded  by net/http
			Expect(response.Header.Get("Content-Length")).To(Equal(""))
			Expect(response.Header.Get("Transfer-Encoding")).To(Equal(""))

			Expect(response.ContentLength).To(Equal(int64(-1)))
			Expect(response.TransferEncoding).To(Equal([]string{"chunked"}))
		})

		It("should ignore content length and use chunked transfer encoding", func() {

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
			Expect(response.ContentLength).To(Equal(int64(-1)))
			Expect(response.TransferEncoding).To(Equal([]string{"chunked"}))
		})

		It("should not set content-length if chunked transfer encoding is already set", func() {

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
									"Transfer-Encoding": ["chunked"]
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
			Expect(response.ContentLength).To(Equal(int64(-1)))
			Expect(response.TransferEncoding).To(Equal([]string{"chunked"}))
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
