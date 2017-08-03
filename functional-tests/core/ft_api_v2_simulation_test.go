package hoverfly_test

import (
	"bytes"
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/util"
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
			Expect(schemaVersion).To(Equal("v4"))
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

			bodyMatchers, err := pairOneRequest.GetObject("body")
			Expect(err).To(BeNil())

			Expect(bodyMatchers.GetString("exactMatch")).Should(Equal(""))

			destinationMatchers, err := pairOneRequest.GetObject("destination")
			Expect(err).To(BeNil())

			Expect(destinationMatchers.GetString("exactMatch")).Should(Equal("test-server.com"))

			methodMatchers, err := pairOneRequest.GetObject("method")
			Expect(err).To(BeNil())

			Expect(methodMatchers.GetString("exactMatch")).Should(Equal("GET"))

			pathMatchers, err := pairOneRequest.GetObject("path")
			Expect(err).To(BeNil())

			Expect(pathMatchers.GetString("exactMatch")).Should(Equal("/path1"))

			queryMatchers, err := pairOneRequest.GetObject("query")
			Expect(err).To(BeNil())

			Expect(queryMatchers.GetString("exactMatch")).Should(Equal(""))

			schemeMatchers, err := pairOneRequest.GetObject("scheme")
			Expect(err).To(BeNil())

			Expect(schemeMatchers.GetString("exactMatch")).Should(Equal("http"))

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

			destinationMatchers, err = pairTwoRequest.GetObject("destination")
			Expect(err).To(BeNil())

			Expect(destinationMatchers.GetString("exactMatch")).Should(Equal("destination-server.com"))
			Expect(pairTwoRequest.GetNull("method")).ShouldNot(Succeed())
			Expect(pairTwoRequest.GetNull("path")).ShouldNot(Succeed())
			Expect(pairTwoRequest.GetNull("destination")).ShouldNot(Succeed())
			Expect(pairTwoRequest.GetNull("query")).ShouldNot(Succeed())
			Expect(pairTwoRequest.GetNull("scheme")).ShouldNot(Succeed())
			Expect(pairTwoRequest.GetNull("body")).ShouldNot(Succeed())

			Expect(pairTwoRequest.GetNull("headers")).ShouldNot(Succeed())

			pairTwoResponse, err := pairsArray[1].GetObject("response")

			Expect(pairTwoResponse.GetInt64("status")).Should(Equal(int64(200)))
			Expect(pairTwoResponse.GetString("body")).Should(Equal("destination matched"))
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
			Expect(schemaVersion).To(Equal("v4"))
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
							"destination": {
								"exactMatch": "destination-server.com"
							},
							"requiresState" : {
								"burger" : "present"
							}
						},
						"response": {
							"status": 200,
							"body": "destination matched",
							"encodedBody": false,
							"headers": {},
							"templated" : false,
							"transitionsState" : {
								"foo" : "bar"
							},
							"removesState" : ["ham"]
						}
					}]
				},
				"meta": {
					"schemaVersion": "v4"
				}
			}
			`)

			req.Body(payload)
			res := functional_tests.DoRequest(req)
			bytes, _ := ioutil.ReadAll(res.Body)
			GinkgoWriter.Write(bytes)
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

			destinationMatchers, err := requestObject.GetObject("destination")
			Expect(err).To(BeNil())

			destination, err := destinationMatchers.GetString("exactMatch")
			Expect(err).To(BeNil())
			Expect(destination).To(Equal("destination-server.com"))

			responseObject, err := pairsArray[0].GetObject("response")
			Expect(err).To(BeNil())

			status, err := responseObject.GetNumber("status")
			Expect(err).To(BeNil())
			Expect(status.String()).To(Equal("200"))

			body, err := responseObject.GetString("body")
			Expect(err).To(BeNil())
			Expect(body).To(Equal("destination matched"))

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
							"destination": {
								"exactMatch": "destination-server.com"
							}
						},
						"response": {
							"status": 200,
							"body": "destination matched",
							"encodedBody": false,
							"headers": {},
							"templated": true
						}
					}]
				},
				"meta": {
					"schemaVersion": "v3"
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

			destinationMatchers, err := requestObject.GetObject("destination")
			Expect(err).To(BeNil())

			destination, err := destinationMatchers.GetString("exactMatch")
			Expect(err).To(BeNil())
			Expect(destination).To(Equal("destination-server.com"))

			responseObject, err := pairsArray[0].GetObject("response")
			Expect(err).To(BeNil())

			status, err := responseObject.GetNumber("status")
			Expect(err).To(BeNil())
			Expect(status.String()).To(Equal("200"))

			body, err := responseObject.GetString("body")
			Expect(err).To(BeNil())
			Expect(body).To(Equal("destination matched"))

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
                    "schemaVersion": "v2",
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

		It("should import old v1 simulations and upgrade them to v2 simulations", func() {
			hoverfly.ImportSimulation(functional_tests.JsonPayloadV1)

			simulation := hoverfly.ExportSimulation()

			Expect(simulation.DataViewV4.RequestResponsePairs[0].RequestMatcher).To(Equal(v2.RequestMatcherViewV4{
				Destination: &v2.RequestFieldMatchersView{
					ExactMatch: util.StringToPointer("v1-simulation.com"),
				}}))

			Expect(simulation.DataViewV4.RequestResponsePairs[1].RequestMatcher).To(Equal(v2.RequestMatcherViewV4{
				Scheme: &v2.RequestFieldMatchersView{
					ExactMatch: util.StringToPointer("http"),
				},
				Method: &v2.RequestFieldMatchersView{
					ExactMatch: util.StringToPointer("GET"),
				},
				Destination: &v2.RequestFieldMatchersView{
					ExactMatch: util.StringToPointer("v1-simulation.com"),
				},
				Path: &v2.RequestFieldMatchersView{
					ExactMatch: util.StringToPointer("/path"),
				},
				Query: &v2.RequestFieldMatchersView{
					ExactMatch: util.StringToPointer(""),
				},
				Body: &v2.RequestFieldMatchersView{
					ExactMatch: util.StringToPointer(""),
				}}))

		})

		It("should error when importing unknown version", func() {
			request := sling.New().Put("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/simulation")
			payload := bytes.NewBufferString(`{
				"data": {},
				"meta": {
					"schemaVersion": "r3"
				}
			}`)

			request.Body(payload)
			response := functional_tests.DoRequest(request)
			Expect(response.StatusCode).To(Equal(400))

			responseBody, _ := ioutil.ReadAll(response.Body)
			Expect(string(responseBody)).To(Equal(`{"error":"Invalid simulation: schema version r3 is not supported by this version of Hoverfly, you may need to update Hoverfly"}`))
		})
	})
})
