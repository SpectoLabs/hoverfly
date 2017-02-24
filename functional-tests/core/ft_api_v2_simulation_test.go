package hoverfly_test

import (
	"bytes"
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/antonholmquist/jason"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/api/v2/simulation", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
		hoverfly.Start()
		hoverfly.SetMode("simulate")
		hoverfly.ImportSimulation(functional_tests.JsonPayload)
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("GET", func() {

		It("Should get all the Hoverfly simulation data in one JSON file", func() {
			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/simulation")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			responseJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())

			jsonObject, err := jason.NewObjectFromBytes(responseJson)
			Expect(err).To(BeNil())

			metaObject, err := jsonObject.GetObject("meta")
			Expect(err).To(BeNil())
			schemaVersion, err := metaObject.GetString("schemaVersion")
			Expect(err).To(BeNil())
			Expect(schemaVersion).To(Equal("v1"))
			hoverflyVersion, err := metaObject.GetString("hoverflyVersion")
			Expect(err).To(BeNil())
			Expect(hoverflyVersion).ToNot(BeNil())
			timeExported, err := metaObject.GetString("timeExported")
			Expect(err).To(BeNil())
			Expect(timeExported).ToNot(BeNil())

			dataObject, err := jsonObject.GetObject("data")
			Expect(err).To(BeNil())

			pairsArray, err := dataObject.GetObjectArray("pairs")
			Expect(err).To(BeNil())

			Expect(pairsArray).To(HaveLen(2))

			pairOneRequest, err := pairsArray[0].GetObject("request")

			Expect(pairOneRequest.GetString("body")).Should(Equal(""))
			Expect(pairOneRequest.GetString("destination")).Should(Equal("test-server.com"))
			Expect(pairOneRequest.GetString("method")).Should(Equal("GET"))
			Expect(pairOneRequest.GetString("path")).Should(Equal("/path1"))
			Expect(pairOneRequest.GetString("query")).Should(Equal(""))
			Expect(pairOneRequest.GetString("scheme")).Should(Equal("http"))

			pairOneRequestHeaders, _ := pairOneRequest.GetObject("headers")
			Expect(pairOneRequestHeaders.GetStringArray("Accept-Encoding")).Should(ContainElement("gzip"))
			Expect(pairOneRequestHeaders.GetStringArray("User-Agent")).Should(ContainElement("Go-http-client/1.1"))

			pairOneResponse, err := pairsArray[0].GetObject("response")

			Expect(pairOneResponse.GetInt64("status")).Should(Equal(int64(200)))
			Expect(pairOneResponse.GetString("body")).Should(Equal("exact match"))
			Expect(pairOneResponse.GetBoolean("encodedBody")).Should(BeFalse())

			pairOneResponseHeaders, _ := pairOneResponse.GetObject("headers")
			Expect(pairOneResponseHeaders.GetStringArray("Header")).Should(ContainElement("value1"))
			Expect(pairOneResponseHeaders.GetStringArray("Header")).Should(ContainElement("value2"))

			pairTwoRequest, err := pairsArray[1].GetObject("request")

			Expect(pairTwoRequest.GetNull("body")).Should(BeNil())
			Expect(pairTwoRequest.GetString("destination")).Should(Equal("template-server.com"))
			Expect(pairTwoRequest.GetNull("method")).Should(BeNil())
			Expect(pairTwoRequest.GetNull("path")).Should(BeNil())
			Expect(pairTwoRequest.GetNull("query")).Should(BeNil())
			Expect(pairTwoRequest.GetNull("scheme")).Should(BeNil())

			Expect(pairTwoRequest.GetNull("headers")).Should(BeNil())

			pairTwoResponse, err := pairsArray[1].GetObject("response")

			Expect(pairTwoResponse.GetInt64("status")).Should(Equal(int64(200)))
			Expect(pairTwoResponse.GetString("body")).Should(Equal("template match"))
			Expect(pairTwoResponse.GetBoolean("encodedBody")).Should(BeFalse())

			globalActionsObject, err := dataObject.GetObject("globalActions")
			Expect(err).To(BeNil())

			delaysArray, err := globalActionsObject.GetObjectArray("delays")
			Expect(err).To(BeNil())

			Expect(delaysArray).To(HaveLen(0))
		})
	})

	Context("DELETE", func() {

		It("Should delete all the Hoverfly data", func() {
			req := sling.New().Delete("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/simulation")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			responseJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())

			jsonObject, err := jason.NewObjectFromBytes(responseJson)
			Expect(err).To(BeNil())

			dataObject, err := jsonObject.GetObject("data")
			Expect(err).To(BeNil())

			pairsArray, err := dataObject.GetObjectArray("pairs")
			Expect(err).To(BeNil())

			Expect(pairsArray).To(HaveLen(0))

			metaObject, err := jsonObject.GetObject("meta")
			Expect(err).To(BeNil())
			schemaVersion, err := metaObject.GetString("schemaVersion")
			Expect(err).To(BeNil())
			Expect(schemaVersion).To(Equal("v1"))
			hoverflyVersion, err := metaObject.GetString("hoverflyVersion")
			Expect(err).To(BeNil())
			Expect(hoverflyVersion).ToNot(BeNil())
			timeExported, err := metaObject.GetString("timeExported")
			Expect(err).To(BeNil())
			Expect(timeExported).ToNot(BeNil())

			globalActionsObject, err := dataObject.GetObject("globalActions")
			Expect(err).To(BeNil())

			delaysArray, err := globalActionsObject.GetObjectArray("delays")
			Expect(err).To(BeNil())

			Expect(delaysArray).To(HaveLen(0))
		})
	})

	Context("PUT", func() {

		It("Should import data using a PUT and should be able to get the same data back using a GET", func() {
			req := sling.New().Put("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/simulation")
			payload := bytes.NewBufferString(`
			{
				"data": {
					"pairs": [{
						"request": {
							"requestType": "template",
							"destination": "templatedurl.com"
						},
						"response": {
							"status": 200,
							"body": "This is the body for the template",
							"encodedBody": false,
							"headers": {}
						}
					}]
				},
				"meta": {
					"schemaVersion": "v1"
				}
			}
			`)

			req.Body(payload)
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))

			getReq := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/simulation")

			getRes := functional_tests.DoRequest(getReq)
			Expect(getRes.StatusCode).To(Equal(200))

			defer getRes.Body.Close()

			schemaObject, err := jason.NewObjectFromReader(getRes.Body)
			Expect(err).To(BeNil())

			dataObject, err := schemaObject.GetObject("data")
			Expect(err).To(BeNil())

			pairsArray, err := dataObject.GetObjectArray("pairs")
			Expect(err).To(BeNil())

			Expect(pairsArray).To(HaveLen(1))

			requestObject, err := pairsArray[0].GetObject("request")
			Expect(err).To(BeNil())

			requestType, err := requestObject.GetString("requestType")
			Expect(err).To(BeNil())
			Expect(requestType).To(Equal("template"))

			destination, err := requestObject.GetString("destination")
			Expect(err).To(BeNil())
			Expect(destination).To(Equal("templatedurl.com"))

			responseObject, err := pairsArray[0].GetObject("response")
			Expect(err).To(BeNil())

			status, err := responseObject.GetNumber("status")
			Expect(err).To(BeNil())
			Expect(status.String()).To(Equal("200"))

			body, err := responseObject.GetString("body")
			Expect(err).To(BeNil())
			Expect(body).To(Equal("This is the body for the template"))

			encodedBody, err := responseObject.GetBoolean("encodedBody")
			Expect(err).To(BeNil())
			Expect(encodedBody).To(BeFalse())
		})

		It("Should import data using a PUT and should return the new state", func() {
			req := sling.New().Put("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/simulation")
			payload := bytes.NewBufferString(`
			{
				"data": {
					"pairs": [{
						"request": {
							"requestType": "template",
							"destination": "templatedurl.com"
						},
						"response": {
							"status": 200,
							"body": "This is the body for the template",
							"encodedBody": false,
							"headers": {}
						}
					}]
				},
				"meta": {
					"schemaVersion": "v1"
				}
			}
			`)

			req.Body(payload)
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))

			defer res.Body.Close()

			schemaObject, err := jason.NewObjectFromReader(res.Body)
			Expect(err).To(BeNil())

			dataObject, err := schemaObject.GetObject("data")
			Expect(err).To(BeNil())

			pairsArray, err := dataObject.GetObjectArray("pairs")
			Expect(err).To(BeNil())

			Expect(pairsArray).To(HaveLen(1))

			requestObject, err := pairsArray[0].GetObject("request")
			Expect(err).To(BeNil())

			requestType, err := requestObject.GetString("requestType")
			Expect(err).To(BeNil())
			Expect(requestType).To(Equal("template"))

			destination, err := requestObject.GetString("destination")
			Expect(err).To(BeNil())
			Expect(destination).To(Equal("templatedurl.com"))

			responseObject, err := pairsArray[0].GetObject("response")
			Expect(err).To(BeNil())

			status, err := responseObject.GetNumber("status")
			Expect(err).To(BeNil())
			Expect(status.String()).To(Equal("200"))

			body, err := responseObject.GetString("body")
			Expect(err).To(BeNil())
			Expect(body).To(Equal("This is the body for the template"))

			encodedBody, err := responseObject.GetBoolean("encodedBody")
			Expect(err).To(BeNil())
			Expect(encodedBody).To(BeFalse())
		})

		It("should delete previous data when putting new data in", func() {
			request := sling.New().Put("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/simulation")
			payload := bytes.NewBufferString(`
			{
				"data": {
					"pairs": []
				},
				"meta": {
                    "schemaVersion": "v1",
                    "hoverflyVersion": "v0.10.2",
                    "timeExported": "2017-02-23T12:43:48Z"
                }
			}
			`)

			request.Body(payload)
			functional_tests.DoRequest(request)
			getReq := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/simulation")

			getRes := functional_tests.DoRequest(getReq)
			Expect(getRes.StatusCode).To(Equal(200))

			defer getRes.Body.Close()

			schemaObject, err := jason.NewObjectFromReader(getRes.Body)
			Expect(err).To(BeNil())

			dataObject, err := schemaObject.GetObject("data")
			Expect(err).To(BeNil())

			pairsArray, err := dataObject.GetObjectArray("pairs")
			Expect(err).To(BeNil())

			Expect(pairsArray).To(HaveLen(0))
		})
	})
})
