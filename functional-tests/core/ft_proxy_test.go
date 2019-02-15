package hoverfly_test

import (
	"io/ioutil"
	"net/http"

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

	Context("using standard configuration", func() {

		BeforeEach(func() {
			hoverfly.Start()
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("should response OK with a query with escaped query parameters", func() {

			hoverfly.ImportSimulation(`{
				"data": {
					"pairs": [
						{
							"request": {
								"query": {
									"exactMatch": "query=something with a space"
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
			response := hoverfly.Proxy(sling.New().Get("http://hoverfly.io?query=something%20with%20a%20space"))
			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})

		It("should set Content-Length if empty", func() {

			hoverfly.ImportSimulation(`{
				"data": {
					"pairs": [
						{
							"request": {
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
			response := hoverfly.Proxy(sling.New().Get("http://hoverfly.io/path"))
			Expect(response.StatusCode).To(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(response.Body)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("OK"))

			Expect(response.Header.Get("Content-Length")).To(Equal("2"))
			Expect(response.Header.Get("Transfer-Encoding")).To(Equal(""))

			Expect(response.ContentLength).To(Equal(int64(2)))
			Expect(response.TransferEncoding).To(BeNil())
		})

		It("should error if Content-Length is incorrect", func() {

			hoverfly.ImportSimulation(`{
				"data": {
					"pairs": [
						{
							"request": {
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

			response := hoverfly.Proxy(sling.New().Get("http://hoverfly.io/path"))
			Expect(response.StatusCode).To(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(response.Body)
			Expect(err).To(Not(BeNil()))
			Expect(err.Error()).To(Equal("unexpected EOF"))
			Expect(string(body)).To(Equal("OK"))

			Expect(response.Header.Get("Content-length")).To(Equal("5555"))  // Just to make sure we don't delete the Content-Length Header
			Expect(response.Header.Get("Transfer-Encoding")).To(Equal(""))

			Expect(response.ContentLength).To(Equal(int64(5555)))
			Expect(response.TransferEncoding).To(BeNil())
		})

		It("should not set Content-Length if Transfer-Encoding set", func() {

			hoverfly.ImportSimulation(`{
				"data": {
					"pairs": [
						{
							"request": {
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

			response := hoverfly.Proxy(sling.New().Get("http://hoverfly.io/path"))
			Expect(response.StatusCode).To(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(response.Body)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("OK"))

			// Should be empty as Transfer-Encoding was set
			Expect(response.Header.Get("Content-length")).To(Equal(""))
			// Will always be empty as they are excluded by net/http
			Expect(response.Header.Get("Transfer-Encoding")).To(Equal(""))

			Expect(response.ContentLength).To(Equal(int64(-1)))
			Expect(response.TransferEncoding).To(Equal([]string{"chunked"}))
		})
	})

	Context("using plain http tunneling", func() {

		BeforeEach(func() {
			hoverfly.Start("-plain-http-tunneling")
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("should response OK on CONNECT request", func() {

			hoverfly.ImportSimulation(`{
				"data": {
					"pairs": [
						{
							"request": {
								"destination": {
									"exactMatch": "hoverfly.io"
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
			req, _ := http.NewRequest(http.MethodConnect, "http://hoverfly.io", nil)
			response := hoverfly.ProxyRequest(req)
			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})
	})
})
