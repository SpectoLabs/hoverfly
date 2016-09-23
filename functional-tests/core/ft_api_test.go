package hoverfly_test

import (
	"bytes"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"strings"
)

var _ = Describe("Interacting with the API", func() {

	var (
		jsonRequestResponsePair1 *bytes.Buffer
		jsonRequestResponsePair2 *bytes.Buffer
	)

	BeforeEach(func() {
		jsonRequestResponsePair1 = bytes.NewBufferString(`{"data":[{"request": {"path": "/path1", "method": "method1", "destination": "destination1", "scheme": "scheme1", "query": "query1", "body": "body1", "headers": {"Header": ["value1"]}}, "response": {"status": 201, "encodedBody": false, "body": "body1", "headers": {"Header": ["value1"]}}}]}`)
		jsonRequestResponsePair2 = bytes.NewBufferString(`{"data":[{"request": {"path": "/path2", "method": "method2", "destination": "destination2", "scheme": "scheme2", "query": "query2", "body": "body2", "headers": {"Header": ["value2"]}}, "response": {"status": 202, "encodedBody": false, "body": "body2", "headers": {"Header": ["value2"]}}}]}`)

		hoverflyCmd = startHoverfly(adminPort, proxyPort)
	})

	AfterEach(func() {
		stopHoverfly()
	})

	Context("GET /api/v2/hoverfly/destination", func() {

		It("Should get the mode", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/destination")
			res := DoRequest(req)
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
			res := DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"destination":"test.com"}`)))

			req = sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/destination")
			res = DoRequest(req)
			modeJson, err = ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"destination":"test.com"}`)))
		})

	})

	Context("GET /api/v2/hoverfly/mode", func() {

		It("Should get the mode", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/mode")
			res := DoRequest(req)
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
			res := DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"mode":"capture"}`)))

			req = sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/mode")
			res = DoRequest(req)
			modeJson, err = ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"mode":"capture"}`)))
		})

	})

	Context("GET /api/v2/hoverfly/middleware", func() {

		It("Should get the middleware which should be blank", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/middleware")
			res := DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"middleware":""}`)))
		})
	})

	Context("PUT /api/v2/hoverfly/middleware", func() {

		It("Should put the middleware", func() {
			req := sling.New().Put(hoverflyAdminUrl + "/api/v2/hoverfly/middleware")
			req.Body(strings.NewReader(`{"middleware":"cat"}`))
			res := DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"middleware":"cat"}`)))

			req = sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/middleware")
			res = DoRequest(req)
			modeJson, err = ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"middleware":"cat"}`)))
		})

	})

	Context("GET /api/v2/hoverfly/usage", func() {

		It("Should get the usage counters", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/usage")
			res := DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"usage":{"counters":{"capture":0,"modify":0,"simulate":0,"synthesize":0}}}`)))
		})

		It("Should get the usage counters with 1 simulate request when a request has been made", func() {
			proxyReq := sling.New().Get("http://www.google.com")
			DoRequestThroughProxy(proxyReq)
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/usage")
			res := DoRequest(req)
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
			res := DoRequest(req)
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
			res := DoRequest(req)
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
			res := DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"usage":{"counters":{"capture":0,"modify":0,"simulate":0,"synthesize":1}}}`)))
		})
	})

	Context("GET /api/records", func() {

		BeforeEach(func() {
			ImportHoverflyRecords(jsonRequestResponsePair1)
			ImportHoverflyRecords(jsonRequestResponsePair2)
		})

		It("Should retrieve the records", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/records")
			res := DoRequest(req)
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
			resPost := DoRequest(reqPost)
			Expect(resPost.StatusCode).To(Equal(200))
			responseMessage, err := ioutil.ReadAll(resPost.Body)
			Expect(err).To(BeNil())

			Expect(string(responseMessage)).To(ContainSubstring("Proxy cache deleted successfuly"))

			reqGet := sling.New().Get(hoverflyAdminUrl + "/api/records")
			resGet := DoRequest(reqGet)
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
				res := DoRequest(sling.New().Post(hoverflyAdminUrl + "/api/records").Body(jsonRequestResponsePair1))
				Expect(res.StatusCode).To(Equal(200))

				reqGet := sling.New().Get(hoverflyAdminUrl + "/api/records")
				resGet := DoRequest(reqGet)

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
				res := DoRequest(sling.New().Post(hoverflyAdminUrl+"/api/records").Set("Content-Type", "application/json").Body(jsonRequestResponsePair2))
				Expect(res.StatusCode).To(Equal(200))

				reqGet := sling.New().Get(hoverflyAdminUrl + "/api/records")
				resGet := DoRequest(reqGet)

				Expect(resGet.StatusCode).To(Equal(200))

				recordsJson, err := ioutil.ReadAll(resGet.Body)
				Expect(err).To(BeNil())
				Expect(recordsJson).To(ContainSubstring(jsonRequestResponsePair1.String()))
				Expect(recordsJson).To(ContainSubstring(jsonRequestResponsePair2.String()))
			})
		})
	})
})
