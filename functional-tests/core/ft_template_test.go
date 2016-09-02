package hoverfly_test

import (
	"bytes"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
)

var _ = Describe("Using Hoverfly to return responses by request templates", func() {

	Context("With a request template loaded for matching on URL + headers", func() {

		var (
			jsonRequestResponsePair *bytes.Buffer
		)

		BeforeEach(func() {
			jsonRequestResponsePair = bytes.NewBufferString(`
			{
				"data": [{
					"request": {
						"requestType": "template",
						"path": "/path1",
						"method": "GET",
						"destination": "www.virtual.com"
					},
					"response": {
						"status": 201,
						"encodedBody": false,
						"body": "body1",
						"headers": {
							"Header": ["value1"]
						}
					}
				}, {
					"request": {
						"requestType": "template",
						"path": "/path2",
						"method": "GET",
						"destination": "www.virtual.com",
						"headers": {
							"Header": ["value2"]
						}
					},
					"response": {
						"status": 202,
						"body": "body2",
						"headers": {
							"Header": ["value2"]
						}
					}
				}, {
					"request": {
						"requestType": "template",
						"method": "GET",
						"destination": "www.randomheader.com",
						"headers": {
							"Random": ["*"]
						}
					},
					"response": {
						"status": 203,
						"body": "body3",
						"headers": {}
					}
				}]
			}
			`)
		})

		Context("When running in proxy mode", func() {

			BeforeEach(func() {
				hoverflyCmd = startHoverfly(adminPort, proxyPort)
				SetHoverflyMode("simulate")
				ImportHoverflyRecords(jsonRequestResponsePair)
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
		})

		Context("When running in webserver mode", func() {

			BeforeEach(func() {
				hoverflyCmd = startHoverflyWebServer(adminPort, proxyPort)
				ImportHoverflyRecords(jsonRequestResponsePair)
			})

			AfterEach(func() {
				stopHoverfly()
			})

			It("Should find a match", func() {
				request := sling.New().Get("http://localhost:"+proxyPortAsString+"/path2").Add("Header", "value2")

				resp := DoRequest(request)
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(202))
				Expect(string(body)).To(Equal("body2"))
			})

			It("Should find a match using wildcards", func() {
				request := sling.New().Get("http://localhost:"+proxyPortAsString+"/unmatched_path").Add("Random", "whatever-you-like")

				resp := DoRequest(request)
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(203))
				Expect(string(body)).To(Equal("body3"))
			})

		})

	})

})
