package hoverfly_test

import (
	"bytes"
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Interacting with the API", func() {

	var (
		jsonRequestResponsePair1 *bytes.Buffer
		jsonRequestResponsePair2 *bytes.Buffer
	)

	BeforeEach(func() {
		jsonRequestResponsePair1 = bytes.NewBufferString(`{"data":[{"requestTemplate": {"path": "/path1", "method": "method1", "destination": "destination1", "scheme": "scheme1", "query": "query1", "body": "body1", "headers": {"Header": ["value1"]}}, "response": {"status": 201, "encodedBody": false, "body": "body1", "headers": {"Header": ["value1"]}}}]}`)
		jsonRequestResponsePair2 = bytes.NewBufferString(`{"data":[{"requestTemplate": {"path": "/path2", "method": "method2", "destination": "destination2", "scheme": "scheme2", "query": "query2", "body": "body2", "headers": {"Header": ["value2"]}}, "response": {"status": 202, "encodedBody": false, "body": "body2", "headers": {"Header": ["value2"]}}}]}`)
	})

	Context("GET /api/templates", func() {

		BeforeEach(func() {
			hoverflyCmd = startHoverfly(adminPort, proxyPort)
			ImportHoverflyTemplates(jsonRequestResponsePair1)
			ImportHoverflyTemplates(jsonRequestResponsePair2)
		})

		AfterEach(func() {
			stopHoverfly()
		})

		It("Should retrieve the templates", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/templates")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			templatesJsonBytes, err := ioutil.ReadAll(res.Body)

			Expect(err).To(BeNil())
			Expect(templatesJsonBytes).To(MatchJSON(
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
				      "requestTemplate": {
					"path": "/path1",
					"method": "method1",
					"destination": "destination1",
					"scheme": "scheme1",
					"query": "query1",
					"body": "body1",
					"headers": {
					  "Header": [
					    "value1"
					  ]
					}
				      }
				    },
				    {
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
				      "requestTemplate": {
					"path": "/path2",
					"method": "method2",
					"destination": "destination2",
					"scheme": "scheme2",
					"query": "query2",
					"body": "body2",
					"headers": {
					  "Header": [
					    "value2"
					  ]
					}
				      }
				    }
				  ]
				}`))
		})
	})

	Context("DELETE /api/templates", func() {

		BeforeEach(func() {
			hoverflyCmd = startHoverfly(adminPort, proxyPort)
			ImportHoverflyTemplates(jsonRequestResponsePair1)
			ImportHoverflyTemplates(jsonRequestResponsePair2)
		})

		AfterEach(func() {
			stopHoverfly()
		})

		It("Should delete the templates", func() {
			reqPost := sling.New().Delete(hoverflyAdminUrl + "/api/templates")
			resPost := functional_tests.DoRequest(reqPost)
			Expect(resPost.StatusCode).To(Equal(200))
			responseMessage, err := ioutil.ReadAll(resPost.Body)
			Expect(err).To(BeNil())

			Expect(string(responseMessage)).To(ContainSubstring("Template store wiped successfuly"))

			reqGet := sling.New().Get(hoverflyAdminUrl + "/api/templates")
			resGet := functional_tests.DoRequest(reqGet)
			Expect(resGet.StatusCode).To(Equal(200))
			templatesJson, err := ioutil.ReadAll(resGet.Body)
			Expect(err).To(BeNil())
			Expect(templatesJson).To(MatchJSON(
				`{
				  "data": null
				}`))
		})
	})

	Context("POST /api/templates", func() {

		BeforeEach(func() {
			hoverflyCmd = startHoverfly(adminPort, proxyPort)
		})

		AfterEach(func() {
			stopHoverfly()
		})

		Context("When no templates exist", func() {
			It("Should create the templates", func() {
				res := functional_tests.DoRequest(sling.New().Post(hoverflyAdminUrl + "/api/templates").Body(jsonRequestResponsePair1))
				Expect(res.StatusCode).To(Equal(200))

				reqGet := sling.New().Get(hoverflyAdminUrl + "/api/templates")
				resGet := functional_tests.DoRequest(reqGet)

				Expect(resGet.StatusCode).To(Equal(200))

				templatesJson, err := ioutil.ReadAll(resGet.Body)
				Expect(err).To(BeNil())
				Expect(templatesJson).To(MatchJSON(
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
					      "requestTemplate": {
						"path": "/path1",
						"method": "method1",
						"destination": "destination1",
						"scheme": "scheme1",
						"query": "query1",
						"body": "body1",
						"headers": {
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
				ImportHoverflyTemplates(jsonRequestResponsePair1)
			})

			It("Should append the templates to the existing ones", func() {
				res := functional_tests.DoRequest(sling.New().Post(hoverflyAdminUrl+"/api/templates").Set("Content-Type", "application/json").Body(jsonRequestResponsePair2))
				Expect(res.StatusCode).To(Equal(200))

				reqGet := sling.New().Get(hoverflyAdminUrl + "/api/templates")
				resGet := functional_tests.DoRequest(reqGet)

				Expect(resGet.StatusCode).To(Equal(200))
				templatesJsonBytes, err := ioutil.ReadAll(resGet.Body)
				Expect(err).To(BeNil())
				Expect(templatesJsonBytes).To(MatchJSON(
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
					      "requestTemplate": {
						"path": "/path1",
						"method": "method1",
						"destination": "destination1",
						"scheme": "scheme1",
						"query": "query1",
						"body": "body1",
						"headers": {
						  "Header": [
						    "value1"
						  ]
						}
					      }
					    },
					    {
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
					      "requestTemplate": {
						"path": "/path2",
						"method": "method2",
						"destination": "destination2",
						"scheme": "scheme2",
						"query": "query2",
						"body": "body2",
						"headers": {
						  "Header": [
						    "value2"
						  ]
						}
					      }
					    }
					  ]
					}`))
			})
		})
	})
})
