package hoverfly_test

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/antonholmquist/jason"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Interacting with the API", func() {

	var (
		jsonRequestResponsePair1        *bytes.Buffer
		jsonRequestResponsePair2        *bytes.Buffer
		jsonRequestResponsePairTemplate *bytes.Buffer
	)

	BeforeEach(func() {
		jsonRequestResponsePair1 = bytes.NewBufferString(`{"data":[{"request": {"path": "/path1", "method": "method1", "destination": "destination1", "scheme": "scheme1", "query": "query1", "body": "body1", "headers": {"Header": ["value1"]}}, "response": {"status": 201, "encodedBody": false, "body": "body1", "headers": {"Header": ["value1"]}}}]}`)
		jsonRequestResponsePair2 = bytes.NewBufferString(`{"data":[{"request": {"path": "/path2", "method": "method2", "destination": "destination2", "scheme": "scheme2", "query": "query2", "body": "body2", "headers": {"Header": ["value2"]}}, "response": {"status": 202, "encodedBody": false, "body": "body2", "headers": {"Header": ["value2"]}}}]}`)
		jsonRequestResponsePairTemplate = bytes.NewBufferString(`{"data":[{"request": {"requestType": "template", "path": "/template"}, "response": {"status": 202, "encodedBody": false, "body": "template-body", "headers": {"Header": ["value2"]}}}]}`)
		hoverflyCmd = startHoverfly(adminPort, proxyPort)
	})

	AfterEach(func() {
		stopHoverfly()
	})

	Context("GET /api/v2/simulation", func() {

		BeforeEach(func() {
			ImportHoverflyRecords(jsonRequestResponsePair1)
			ImportHoverflyRecords(jsonRequestResponsePair2)
			ImportHoverflyRecords(jsonRequestResponsePairTemplate)
			SetHoverflyResponseDelays("testdata/delays.json")
		})

		It("Should get all the Hoverfly simulation data in one JSON file", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/simulation")
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

			Expect(pairsArray).To(HaveLen(3))

			Expect(jsonObject.String()).To(ContainSubstring(`{"body":"body1","destination":"destination1","headers":{"Content-Type":["text/plain; charset=utf-8"],"Header":["value1"]},"method":"method1","path":"/path1","query":"query1","requestType":"recording","scheme":"scheme1"}`))
			Expect(jsonObject.String()).To(ContainSubstring(`{"body":"body1","encodedBody":false,"headers":{"Header":["value1"]},"status":201}`))

			Expect(jsonObject.String()).To(ContainSubstring(`{"body":"body2","destination":"destination2","headers":{"Content-Type":["text/plain; charset=utf-8"],"Header":["value2"]},"method":"method2","path":"/path2","query":"query2","requestType":"recording","scheme":"scheme2"}`))
			Expect(jsonObject.String()).To(ContainSubstring(`{"body":"body2","encodedBody":false,"headers":{"Header":["value2"]},"status":202}`))

			Expect(jsonObject.String()).To(ContainSubstring(`{"body":null,"destination":null,"headers":null,"method":null,"path":"/template","query":null,"requestType":"template","scheme":null}`))
			Expect(jsonObject.String()).To(ContainSubstring(`{"body":"template-body","encodedBody":false,"headers":{"Header":["value2"]},"status":202}`))

			globalActionsObject, err := dataObject.GetObject("globalActions")
			Expect(err).To(BeNil())

			delaysArray, err := globalActionsObject.GetObjectArray("delays")
			Expect(err).To(BeNil())

			Expect(delaysArray).To(HaveLen(2))
			Expect(delaysArray[0].String()).To(Equal(`{"delay":100,"httpMethod":"","urlPattern":"virtual\\.com"}`))
			Expect(delaysArray[1].String()).To(Equal(`{"delay":110,"httpMethod":"","urlPattern":"virtual\\.com"}`))
		})
	})

	Context("DELETE /api/v2/simulation", func() {

		It("Should delete all the Hoverfly data", func() {
			req := sling.New().Delete(hoverflyAdminUrl + "/api/v2/simulation")
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

	Context("PUT /api/v2/simulation", func() {
		It("Should import data using a PUT and should be able to get the same data back using a GET", func() {
			req := sling.New().Put(hoverflyAdminUrl + "/api/v2/simulation")
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

			getReq := sling.New().Get(hoverflyAdminUrl + "/api/v2/simulation")

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
			req := sling.New().Put(hoverflyAdminUrl + "/api/v2/simulation")
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
			originalReq := sling.New().Put(hoverflyAdminUrl + "/api/v2/simulation")
			originalPayload := bytes.NewBufferString(`
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
				"meta": {}
			}
			`)

			originalReq.Body(originalPayload)
			functional_tests.DoRequest(originalReq)

			request := sling.New().Put(hoverflyAdminUrl + "/api/v2/simulation")
			payload := bytes.NewBufferString(`
			{
				"data": {
					"pairs": []
				},
				"meta": {}
			}
			`)

			request.Body(payload)
			functional_tests.DoRequest(request)
			getReq := sling.New().Get(hoverflyAdminUrl + "/api/v2/simulation")

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

	Context("GET /api/v2/hoverfly", func() {

		It("Should get the hoverfly config", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly")

			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))

			hoverflyJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(hoverflyJson).To(MatchRegexp(`"destination":"."`))
			Expect(hoverflyJson).To(MatchRegexp(`"middleware":{"binary":"","script":"","remote":""}`))
			Expect(hoverflyJson).To(MatchRegexp(`"usage":{"counters":{"capture":0,"modify":0,"simulate":0,"synthesize":0}}`))
			Expect(hoverflyJson).To(MatchRegexp(`"version":"v\d+.\d+.\d+"`))
			Expect(hoverflyJson).To(MatchRegexp(`"upstream-proxy":""`))
		})
	})

	Context("GET /api/v2/hoverfly/destination", func() {

		It("Should get the mode", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/destination")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"destination":"."}`)))
		})

	})

	Context("PUT /api/v2/hoverfly/destination", func() {

		It("Should put the mode", func() {
			req := sling.New().Put(hoverflyAdminUrl + "/api/v2/hoverfly/destination")
			req.Body(strings.NewReader(`{"destination":"test.com"}`))
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"destination":"test.com"}`)))

			req = sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/destination")
			res = functional_tests.DoRequest(req)
			modeJson, err = ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"destination":"test.com"}`)))
		})

	})

	Context("GET /api/v2/hoverfly/mode", func() {

		It("Should get the mode", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/mode")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"mode":"simulate"}`)))
		})
	})

	Context("PUT /api/v2/hoverfly/mode", func() {

		It("Should put the mode", func() {
			req := sling.New().Put(hoverflyAdminUrl + "/api/v2/hoverfly/mode")
			req.Body(strings.NewReader(`{"mode":"capture"}`))
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"mode":"capture"}`)))

			req = sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/mode")
			res = functional_tests.DoRequest(req)
			modeJson, err = ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"mode":"capture"}`)))
		})

	})

	Context("GET /api/v2/hoverfly/middleware", func() {

		It("Should get the middleware which should be blank", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/middleware")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"binary":"","script":"","remote":""}`)))
		})
	})

	Context("PUT /api/v2/hoverfly/middleware", func() {

		It("Should put the middleware", func() {
			req := sling.New().Put(hoverflyAdminUrl + "/api/v2/hoverfly/middleware")
			req.Body(strings.NewReader(`{"binary":"ruby", "script":"#!/usr/bin/env ruby\n# encoding: utf-8\nwhile payload = STDIN.gets\nnext unless payload\n\nSTDOUT.puts payload\nend"}`))
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"binary":"ruby","script":"#!/usr/bin/env ruby\n# encoding: utf-8\nwhile payload = STDIN.gets\nnext unless payload\n\nSTDOUT.puts payload\nend","remote":""}`)))

			req = sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/middleware")
			res = functional_tests.DoRequest(req)
			modeJson, err = ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"binary":"ruby","script":"#!/usr/bin/env ruby\n# encoding: utf-8\nwhile payload = STDIN.gets\nnext unless payload\n\nSTDOUT.puts payload\nend","remote":""}`)))
		})

	})

	Context("GET /api/v2/hoverfly/usage", func() {

		It("Should get the usage counters", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/usage")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"usage":{"counters":{"capture":0,"modify":0,"simulate":0,"synthesize":0}}}`)))
		})

		It("Should get the usage counters with 1 simulate request when a request has been made", func() {
			proxyReq := sling.New().Get("http://www.google.com")
			DoRequestThroughProxy(proxyReq)
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/usage")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"usage":{"counters":{"capture":0,"modify":0,"simulate":1,"synthesize":0}}}`)))
		})

		It("Should get the usage counters with 1 capture request when a request has been made", func() {
			SetHoverflyMode("capture")

			proxyReq := sling.New().Get("http://www.google.com")
			DoRequestThroughProxy(proxyReq)
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/usage")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"usage":{"counters":{"capture":1,"modify":0,"simulate":0,"synthesize":0}}}`)))
		})

		It("Should get the usage counters with 1 modify request when a request has been made", func() {
			SetHoverflyMode("modify")

			proxyReq := sling.New().Get("http://www.google.com")
			DoRequestThroughProxy(proxyReq)
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/usage")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"usage":{"counters":{"capture":0,"modify":1,"simulate":0,"synthesize":0}}}`)))
		})

		It("Should get the usage counters with 1 modify request when a request has been made", func() {
			SetHoverflyMode("synthesize")

			proxyReq := sling.New().Get("http://www.google.com")
			DoRequestThroughProxy(proxyReq)
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/usage")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"usage":{"counters":{"capture":0,"modify":0,"simulate":0,"synthesize":1}}}`)))
		})
	})

	Context("GET /api/v2/hoverfly/version", func() {

		It("Should get the version", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/version")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(string(modeJson)).To(MatchRegexp(`{"version":"v\d+.\d+.\d+"}`))
		})
	})

	Context("GET /api/v2/hoverfly/upstream-proxy", func() {

		It("Should get the upstream proxy", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/upstream-proxy")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(string(modeJson)).To(MatchRegexp(`{"upstream-proxy":""}`))
		})
	})

	Context("GET /api/records", func() {

		BeforeEach(func() {
			ImportHoverflyRecords(jsonRequestResponsePair1)
			ImportHoverflyRecords(jsonRequestResponsePair2)
		})

		It("Should retrieve the records", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/records")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			recordsJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(recordsJson).To(ContainSubstring(jsonRequestResponsePair1.String()))
			Expect(recordsJson).To(ContainSubstring(jsonRequestResponsePair2.String()))
		})
	})

	Context("DELETE /api/records", func() {

		BeforeEach(func() {
			ImportHoverflyRecords(jsonRequestResponsePair1)
			ImportHoverflyRecords(jsonRequestResponsePair2)
		})

		It("Should delete the records", func() {
			reqPost := sling.New().Delete(hoverflyAdminUrl + "/api/records")
			resPost := functional_tests.DoRequest(reqPost)
			Expect(resPost.StatusCode).To(Equal(200))
			responseMessage, err := ioutil.ReadAll(resPost.Body)
			Expect(err).To(BeNil())

			Expect(string(responseMessage)).To(ContainSubstring("Proxy cache deleted successfuly"))

			reqGet := sling.New().Get(hoverflyAdminUrl + "/api/records")
			resGet := functional_tests.DoRequest(reqGet)
			Expect(resGet.StatusCode).To(Equal(200))
			recordsJson, err := ioutil.ReadAll(resGet.Body)
			Expect(err).To(BeNil())
			Expect(recordsJson).To(MatchJSON(
				`{
				  "data": null
				}`))
		})
	})

	Context("POST /api/records", func() {

		Context("When no records exist", func() {
			It("Should create the records", func() {
				res := functional_tests.DoRequest(sling.New().Post(hoverflyAdminUrl + "/api/records").Body(jsonRequestResponsePair1))
				Expect(res.StatusCode).To(Equal(200))

				reqGet := sling.New().Get(hoverflyAdminUrl + "/api/records")
				resGet := functional_tests.DoRequest(reqGet)

				Expect(resGet.StatusCode).To(Equal(200))

				recordsJson, err := ioutil.ReadAll(resGet.Body)
				Expect(err).To(BeNil())
				Expect(recordsJson).To(MatchJSON(
					`{
					  "data": [
					    {
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
					      	"requestType": "recording",
						"path": "/path1",
						"method": "method1",
						"destination": "destination1",
						"scheme": "scheme1",
						"query": "query1",
						"body": "body1",
						"headers": {
						  "Content-Type": [
						    "text/plain; charset=utf-8"
						  ],
						  "Header": [
						    "value1"
						  ]
						}
					      }
					    }
					  ]
					}`))
			})
		})

		Context("When a record already exists", func() {

			BeforeEach(func() {
				ImportHoverflyRecords(jsonRequestResponsePair1)
			})

			It("Should append the records to the existing ones", func() {
				res := functional_tests.DoRequest(sling.New().Post(hoverflyAdminUrl+"/api/records").Set("Content-Type", "application/json").Body(jsonRequestResponsePair2))
				Expect(res.StatusCode).To(Equal(200))

				reqGet := sling.New().Get(hoverflyAdminUrl + "/api/records")
				resGet := functional_tests.DoRequest(reqGet)

				Expect(resGet.StatusCode).To(Equal(200))

				recordsJson, err := ioutil.ReadAll(resGet.Body)
				Expect(err).To(BeNil())
				Expect(recordsJson).To(ContainSubstring(jsonRequestResponsePair1.String()))
				Expect(recordsJson).To(ContainSubstring(jsonRequestResponsePair2.String()))
			})
		})
	})
})
