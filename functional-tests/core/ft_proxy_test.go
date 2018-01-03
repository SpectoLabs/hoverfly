package hoverfly_test

import (
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
