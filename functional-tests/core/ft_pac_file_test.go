package hoverfly_test

import (
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I run Hoverfly with a PAC file", func() {

	var (
		hoverflyPassThrough, hoverflyUpstream1, hoverflyUpstream2 *functional_tests.Hoverfly
		hoverflyUpstream1URL, hoverflyUpstream2URL                string
	)

	BeforeEach(func() {
		hoverflyPassThrough = functional_tests.NewHoverfly()
		hoverflyUpstream1 = functional_tests.NewHoverfly()
		hoverflyUpstream2 = functional_tests.NewHoverfly()
	})

	AfterEach(func() {
		hoverflyPassThrough.Stop()
		hoverflyUpstream1.Stop()
		hoverflyUpstream2.Stop()
	})

	Context("and configure the upstream proxy", func() {

		BeforeEach(func() {
			hoverflyUpstream1.Start()
			hoverflyUpstream1.SetMode("simulate")
			hoverflyUpstream1.ImportSimulation(`{
				"data": {
					"pairs": [
						{
							"request": {
								"destination": [
									{
										"matcher": "glob",
										"value": "*"
									}
								]
							},
							"response": {
								"status": 200,
								"body": "hoverfly upstream proxy 1",
								"encodedBody": false,
								"templated": false
							}
						}
					],
					"globalActions": {
						"delays": []
					}
				},
				"meta": {
					"schemaVersion": "v5",
					"hoverflyVersion": "v0.17.0",
					"timeExported": "2018-05-03T12:08:35+01:00"
				}
			}`)

			hoverflyUpstream2.Start()
			hoverflyUpstream2.SetMode("simulate")
			hoverflyUpstream2.ImportSimulation(`{
				"data": {
					"pairs": [
						{
							"request": {
								"destination": [
									{
										"matcher": "glob",
										"value": "*"
									}
								]
							},
							"response": {
								"status": 200,
								"body": "hoverfly upstream proxy 2",
								"encodedBody": false,
								"templated": false
							}
						}
					],
					"globalActions": {
						"delays": []
					}
				},
				"meta": {
					"schemaVersion": "v5",
					"hoverflyVersion": "v0.17.0",
					"timeExported": "2018-05-03T12:08:35+01:00"
				}
			}`)

			hoverflyUpstream1URL = "localhost:" + hoverflyUpstream1.GetProxyPort()
			hoverflyUpstream2URL = "localhost:" + hoverflyUpstream2.GetProxyPort()

			hoverflyPassThrough.Start()
			hoverflyPassThrough.SetMode("capture")
		})

		It("Should use PAC file to determine upstream proxy", func() {
			hoverflyPassThrough.SetPACFile(`function FindProxyForURL(url, host) {
				return "PROXY ` + hoverflyUpstream1URL + `";
			}`)

			resp := hoverflyPassThrough.Proxy(sling.New().Get("http://example.com"))
			Expect(resp.StatusCode).To(Equal(200))

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(bodyBytes)).To(Equal("hoverfly upstream proxy 1"))
		})

		It("Should use PAC file to dynamically set upstream proxy per request", func() {
			hoverflyPassThrough.SetPACFile(`function FindProxyForURL(url, host) {
				if (shExpMatch(host, "*.com"))
				{
					return "PROXY ` + hoverflyUpstream1URL + `";
				}

				if (shExpMatch(host, "*.org"))
				{
					return "PROXY ` + hoverflyUpstream2URL + `";
				}
			
				return "PROXY ` + hoverflyUpstream1URL + `";
			}`)

			resp := hoverflyPassThrough.Proxy(sling.New().Get("http://example.com"))
			Expect(resp.StatusCode).To(Equal(200))

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(bodyBytes)).To(Equal("hoverfly upstream proxy 1"))

			resp = hoverflyPassThrough.Proxy(sling.New().Get("http://example.org"))
			Expect(resp.StatusCode).To(Equal(200))

			bodyBytes, err = ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(bodyBytes)).To(Equal("hoverfly upstream proxy 2"))

			resp = hoverflyPassThrough.Proxy(sling.New().Get("http://example.io"))
			Expect(resp.StatusCode).To(Equal(200))

			bodyBytes, err = ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(bodyBytes)).To(Equal("hoverfly upstream proxy 1"))
		})

		It("Should error appropriately if PAC file is invalid", func() {
			hoverflyPassThrough.SetPACFile(`BADPACFILE`)

			resp := hoverflyPassThrough.Proxy(sling.New().Get("http://example.com"))
			Expect(resp.StatusCode).To(Equal(502))

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(bodyBytes)).To(ContainSubstring("Got error: Unable to parse PAC file"))
			Expect(string(bodyBytes)).To(ContainSubstring("ReferenceError: 'BADPACFILE' is not defined"))

			hoverflyPassThrough.SetPACFile(`function FindProxyForURL(url, host) {
				&BADPACFILE&
			}`)

			resp = hoverflyPassThrough.Proxy(sling.New().Get("http://example.com"))
			Expect(resp.StatusCode).To(Equal(502))

			bodyBytes, err = ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(bodyBytes)).To(ContainSubstring("Got error: Unable to parse PAC file"))
			Expect(string(bodyBytes)).To(ContainSubstring("(anonymous): Line 2:5 Unexpected token & (and 1 more errors)"))
		})
	})
})
