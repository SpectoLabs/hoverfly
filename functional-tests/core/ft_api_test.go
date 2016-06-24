package hoverfly_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"github.com/dghubble/sling"
	"bytes"
)

var _ = Describe("Interacting with the API", func() {

	var (
		jsonPayload1 *bytes.Buffer
		jsonPayload2 *bytes.Buffer
	)

	BeforeEach(func() {
		jsonPayload1 = bytes.NewBufferString(`{"data":[{"request": {"path": "/path1", "method": "method1", "destination": "destination1", "scheme": "scheme1", "query": "query1", "body": "body1", "headers": {"Header": ["value1"]}}, "response": {"status": 201, "encodedBody": false, "body": "body1", "headers": {"Header": ["value1"]}}}]}`)
		jsonPayload2 = bytes.NewBufferString(`{"data":[{"request": {"path": "/path2", "method": "method2", "destination": "destination2", "scheme": "scheme2", "query": "query2", "body": "body2", "headers": {"Header": ["value2"]}}, "response": {"status": 202, "encodedBody": false, "body": "body2", "headers": {"Header": ["value2"]}}}]}`)
	})

	Context("GET /api/records", func() {

		BeforeEach(func() {
			hoverflyCmd = startHoverfly(adminPort, proxyPort)
			ImportHoverflyRecords(jsonPayload1)
			ImportHoverflyRecords(jsonPayload2)
		})

		AfterEach(func() {
			stopHoverfly()
		})

		It("Should retrieve the records", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/records")
			res := DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			recordsJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(recordsJson).To(ContainSubstring(jsonPayload1.String()))
			Expect(recordsJson).To(ContainSubstring(jsonPayload2.String()))
		})
	})

	Context("DELETE /api/records", func() {

		BeforeEach(func() {
			hoverflyCmd = startHoverfly(adminPort, proxyPort)
			ImportHoverflyRecords(jsonPayload1)
			ImportHoverflyRecords(jsonPayload2)
		})

		AfterEach(func() {
			stopHoverfly()
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

		BeforeEach(func() {
			hoverflyCmd = startHoverfly(adminPort, proxyPort)
		})

		AfterEach(func() {
			stopHoverfly()
		})

		Context("When no records exist", func() {
			It("Should create the records", func() {
				res := DoRequest(sling.New().Post(hoverflyAdminUrl + "/api/records").Body(jsonPayload1))
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
				ImportHoverflyRecords(jsonPayload1)
			})

			It("Should append the records to the existing ones", func() {
				res := DoRequest(sling.New().Post(hoverflyAdminUrl + "/api/records").Set("Content-Type", "application/json").Body(jsonPayload2))
				Expect(res.StatusCode).To(Equal(200))

				reqGet := sling.New().Get(hoverflyAdminUrl + "/api/records")
				resGet := DoRequest(reqGet)

				Expect(resGet.StatusCode).To(Equal(200))

				recordsJson, err := ioutil.ReadAll(resGet.Body)
				Expect(err).To(BeNil())
				Expect(recordsJson).To(ContainSubstring(jsonPayload1.String()))
				Expect(recordsJson).To(ContainSubstring(jsonPayload2.String()))
			})
		})
	})
})
