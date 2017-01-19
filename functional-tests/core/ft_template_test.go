package hoverfly_test

import (
	"bytes"
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Using Hoverfly to return responses by request templates", func() {

	Context("With a request template loaded for matching on URL + headers", func() {

		var (
			jsonRequestResponsePair *bytes.Buffer
		)

		BeforeEach(func() {
			jsonRequestResponsePair = bytes.NewBufferString(`
			{
	"data": {
		"pairs": [{
			"response": {
				"status": 201,
				"body": "body1",
				"encodedBody": false,
				"headers": {
					"Header": [
						"value1"
					]
				}
			},
			"request": {
				"requestType": "template",
				"path": "/path1",
				"method": "GET",
				"destination": "www.virtual.com"
			}
		}, {
			"response": {
				"status": 202,
				"body": "body2",
				"encodedBody": false,
				"headers": {
					"Header": [
						"value2"
					]
				}
			},
			"request": {
				"requestType": "template",
				"path": "/path2",
				"method": "GET",
				"destination": "www.virtual.com",
				"headers": {
					"Header": ["value2"]
				}
			}
		}, {
			"response": {
				"status": 203,
				"body": "body3",
				"encodedBody": false,
				"headers": {}
			},
			"request": {
				"requestType": "template",
				"method": "GET",
				"destination": "www.randomheader.com",
				"headers": {
					"Random": ["*"]
				}
			}
		}, {
			"response": {
				"status": 200,
				"body": "body",
				"encodedBody": false,
				"headers": {}
			},
			"request": {
				"requestType": "template",
				"method": "GET",
				"destination": "www.query.com",
				"query": "query1=one&query2=two&query2=three"
			}
		}, {
			"response": {
				"status": 200,
				"body": "body",
				"encodedBody": false,
				"headers": {}
			},
			"request": {
				"requestType": "template",
				"method": "GET",
				"destination": "www.query.com",
				"query": "query2=two&query1=one"
			}
		}],
		"globalActions": {
			"delays": []
		}
	},
	"meta": {
		"schemaVersion": "v1",
		"hoverflyVersion": "v0.10.0",
		"timeExported": "2017-01-20T14:03:04Z"
	}
}
			`)
		})

		Context("When running in proxy mode", func() {

			BeforeEach(func() {
				hoverflyCmd = startHoverfly(adminPort, proxyPort)
				SetHoverflyMode("simulate")
				ImportHoverflySimulation(jsonRequestResponsePair)
			})

			AfterEach(func() {
				stopHoverfly()
			})

			It("Should find a match", func() {
				resp := DoRequestThroughProxy(sling.New().Get("http://www.virtual.com/path2").Add("Header", "value2"))
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(202))
				Expect(string(body)).To(Equal("body2"))
			})

			It("Should find a match using wildcards", func() {
				resp := DoRequestThroughProxy(sling.New().Get("http://www.randomheader.com/unmatched_path").Add("Random", "value2"))
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(203))
				Expect(string(body)).To(Equal("body3"))
			})

			It("Should find a match using a different order set of query parameters", func() {
				resp := DoRequestThroughProxy(sling.New().Get("http://www.query.com/?query1=one&query2=two").Add("Random", "value2"))
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(200))
				Expect(string(body)).To(Equal("body"))
			})

			It("Should find a match with two query parameter keys", func() {
				resp := DoRequestThroughProxy(sling.New().Get("http://www.query.com/?query2=two&query1=one&query2=three").Add("Random", "value2"))
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(200))
				Expect(string(body)).To(Equal("body"))
			})
		})

		Context("When running in webserver mode", func() {

			BeforeEach(func() {
				hoverflyCmd = startHoverflyWebServer(adminPort, proxyPort)
				ImportHoverflySimulation(jsonRequestResponsePair)
			})

			AfterEach(func() {
				stopHoverfly()
			})

			It("Should find a match", func() {
				request := sling.New().Get("http://localhost:"+proxyPortAsString+"/path2").Add("Header", "value2")

				resp := functional_tests.DoRequest(request)
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(202))
				Expect(string(body)).To(Equal("body2"))
			})

			It("Should find a match using wildcards", func() {
				request := sling.New().Get("http://localhost:"+proxyPortAsString+"/unmatched_path").Add("Random", "whatever-you-like")

				resp := functional_tests.DoRequest(request)
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(203))
				Expect(string(body)).To(Equal("body3"))
			})

		})

	})

})
