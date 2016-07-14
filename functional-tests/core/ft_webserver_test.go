package hoverfly_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//"io/ioutil"
	//"bytes"
	//"github.com/dghubble/sling"
	//"net/http/httptest"
	//"net/url"
	//"strings"
	//"fmt"
	//"net/http"
	//"os"
	//"time"
	"bytes"
	"io/ioutil"
	"github.com/dghubble/sling"
	"strings"
)

var _ = Describe("When running Hoverfly as a webserver", func() {

	Context("and its in simulate mode", func() {

		BeforeEach(func() {
			hoverflyCmd = startHoverflyWebServer(adminPort, proxyPort)
			getPayload1 := bytes.NewBufferString(`{"data":[{"request": {"path": "/path1", "method": "GET", "destination": "destination1", "scheme": "", "query": "", "body": "", "headers": {"Header": ["value1"]}}, "response": {"status": 201, "encodedBody": false, "body": "body1", "headers": {"Header": ["value1"]}}}]}`)
			getPayload2 := bytes.NewBufferString(`{"data":[{"request": {"path": "/path1/resource", "method": "GET", "destination": "another-destination.com", "scheme": "", "query": "", "body": "", "headers": {"Header": ["value1"]}}, "response": {"status": 201, "encodedBody": false, "body": "another-host.com body1", "headers": {"Header": ["value1"]}}}]}`)
			postPayload1 := bytes.NewBufferString(`{"data":[{"request": {"path": "/path2", "method": "POST", "destination": "destination1", "scheme": "", "query": "", "body": "", "headers": {"Header": ["value1"]}}, "response": {"status": 201, "encodedBody": false, "body": "body2", "headers": {"Header": ["value1"]}}}]}`)
			postPayload2 := bytes.NewBufferString(`{"data":[{"request": {"path": "/path2/resource", "method": "POST", "destination": "another-destination.com", "scheme": "", "query": "", "body": "", "headers": {"Header": ["value1"]}}, "response": {"status": 201, "encodedBody": false, "body": "another-host.com body2", "headers": {"Header": ["value1"]}}}]}`)

			ImportHoverflyRecords(getPayload1)
			ImportHoverflyRecords(postPayload1)
			ImportHoverflyRecords(getPayload2)
			ImportHoverflyRecords(postPayload2)
		})

		AfterEach(func() {
			hoverflyCmd.Process.Kill()
		})

		Context("I can request an endpoint", func() {
			Context("using GET", func() {
				It("and it should return the response", func() {
					request := sling.New().Get("http://localhost:" + proxyPortAsString + "/path1")

					response := DoRequest(request)

					responseBody, err := ioutil.ReadAll(response.Body)
					Expect(err).To(BeNil())

					Expect(string(responseBody)).To(Equal("body1"))
				})
			})

			Context("using POST", func() {
				It("and it should return the response", func() {
					request := sling.New().Post("http://localhost:" + proxyPortAsString + "/path2")

					response := DoRequest(request)

					responseBody, err := ioutil.ReadAll(response.Body)
					Expect(err).To(BeNil())

					Expect(string(responseBody)).To(Equal("body2"))
				})
			})
		})

		Context("I can request an endpoint on another host", func() {
			Context("using GET", func() {
				It("and it should still return the response", func() {
					request := sling.New().Get("http://localhost:" + proxyPortAsString + "/path1/resource")

					response := DoRequest(request)

					responseBody, err := ioutil.ReadAll(response.Body)
					Expect(err).To(BeNil())

					Expect(string(responseBody)).To(Equal("another-host.com body1"))
				})
			})

			Context("using POST", func() {
				It("and it should still return the response", func() {
					request := sling.New().Post("http://localhost:" + proxyPortAsString + "/path2/resource")

					response := DoRequest(request)

					responseBody, err := ioutil.ReadAll(response.Body)
					Expect(err).To(BeNil())

					Expect(string(responseBody)).To(Equal("another-host.com body2"))
				})
			})
		})

		Context("I cannot change the mode", func() {

			It("it should start in simulate mode", func() {
				request := sling.New().Get(hoverflyAdminUrl + "/api/state")
				response := DoRequest(request)

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).To(BeNil())

				Expect(string(responseBody)).To(ContainSubstring("simulate"))
			})

			It("it should not be switchable", func() {
				request := sling.New().Post(hoverflyAdminUrl + "/api/state").Body(strings.NewReader(`{"mode":"capture"}`))
				response := DoRequest(request)

				Expect(response.StatusCode).To(Equal(403))

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).To(BeNil())

				Expect(string(responseBody)).To(ContainSubstring("Hoverfly is currently configured to act as webserver, which can only operate in simulate mode"))
			})
		})
	})

	Context("using a request.db from the disk", func() {

		Context("and the requests.db was created by Hoverfly proxy", func() {
			BeforeEach(func() {
				hoverflyCmd = startHoverflyWithDatabase(adminPort, proxyPort)
				getPayload1 := bytes.NewBufferString(`{"data":[{"request": {"path": "/path1", "method": "GET", "destination": "destination1", "scheme": "", "query": "", "body": "", "headers": {"Header": ["value1"]}}, "response": {"status": 200, "encodedBody": false, "body": "body1", "headers": {"Header": ["value1"]}}}]}`)
				ImportHoverflyRecords(getPayload1)

				hoverflyCmd.Process.Kill()

				hoverflyCmd = startHoverflyWebServerWithDatabase(adminPort, proxyPort)
			})

			AfterEach(func() {
				hoverflyCmd.Process.Kill()
			})

			It("should recache the requests so that the destination is not included in the cache", func() {
				request := sling.New().Get("http://localhost:" + proxyPortAsString + "/path1")

				response := DoRequest(request)

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).To(BeNil())

				Expect(string(responseBody)).To(Equal("body1"))
			})

		})
	})
})
