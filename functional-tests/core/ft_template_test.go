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
			jsonPayload *bytes.Buffer
		)

		BeforeEach(func() {
			jsonPayload = bytes.NewBufferString(`{"data":[{"requestTemplate": {"path": "/path1", "method": "GET", "destination": "www.virtual.com"}, "response": {"status": 201, "encodedBody": false, "body": "body1", "headers": {"Header": ["value1"]}}}, {"requestTemplate": {"path": "/path2", "method": "GET", "destination": "www.virtual.com", "headers": {"Header": ["value2"]}}, "response": {"status": 202, "body": "body2", "headers": {"Header": ["value2"]}}}]}`)

		})

		Context("When running in proxy mode", func() {

			BeforeEach(func() {
				hoverflyCmd = startHoverfly(adminPort, proxyPort)
				SetHoverflyMode("simulate")
				ImportHoverflyTemplates(jsonPayload)
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
		})

		Context("When running in webserver mode", func() {

			BeforeEach(func() {
				hoverflyCmd = startHoverflyWebServer(adminPort, proxyPort)
				ImportHoverflyTemplates(jsonPayload)
			})

			It("Should find a match", func() {
				request := sling.New().Get("http://localhost:"+proxyPortAsString+"/path2").Add("Header", "value2")

				resp := DoRequest(request)
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(202))
				Expect(string(body)).To(Equal("body2"))
			})

			AfterEach(func() {
				stopHoverfly()
			})
		})

	})

})
