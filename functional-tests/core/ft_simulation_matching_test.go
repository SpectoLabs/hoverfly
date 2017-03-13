package hoverfly_test

import (
	"encoding/xml"
	"io/ioutil"

	"bytes"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When using different matchers", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
		hoverfly.Start()
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("Using `xpathMatch`", func() {

		BeforeEach(func() {
			hoverfly.ImportSimulation(functional_tests.XpathSimulation)
		})

		It("should match on the body", func() {
			req := sling.New().Get("http://test.com")
			req.Body(bytes.NewBufferString(xml.Header + "<item></item><item></item><item></item><item></item><item></item>"))

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(ioutil.ReadAll(response.Body)).Should(Equal([]byte("xpath match")))
		})

		It("should not match on no body", func() {
			req := sling.New().Get("http://test.com")

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(502))

			Expect(ioutil.ReadAll(response.Body)).Should(ContainSubstring("There was an error when matching"))
		})
	})

	Context("Using `jsonPathMatch`", func() {

		BeforeEach(func() {
			hoverfly.ImportSimulation(functional_tests.JsonPathMatchSimulation)
		})

		It("should match on the body", func() {
			req := sling.New().Get("http://test.com")
			req.Body(bytes.NewBufferString(`{"items": [{}, {}, {}, {}, {}]}`))

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(ioutil.ReadAll(response.Body)).Should(Equal([]byte("json match")))
		})

		It("should not match on no body", func() {
			req := sling.New().Get("http://test.com")

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(502))

			Expect(ioutil.ReadAll(response.Body)).Should(ContainSubstring("There was an error when matching"))
		})
	})

	Context("Using `regexMatch`", func() {

		BeforeEach(func() {
			hoverfly.ImportSimulation(functional_tests.RegexMatchSimulation)
		})

		It("should match on the body", func() {
			req := sling.New().Get("http://test.com")
			req.Body(bytes.NewBufferString(xml.Header + "<items><item field=something></item></items>"))

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(ioutil.ReadAll(response.Body)).Should(Equal([]byte("regex match")))
		})

		It("should not match on no body", func() {
			req := sling.New().Get("http://test.com")

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(502))

			Expect(ioutil.ReadAll(response.Body)).Should(ContainSubstring("There was an error when matching"))
		})
	})

	Context("Using `globMatch`", func() {

		BeforeEach(func() {
			hoverfly.ImportSimulation(functional_tests.GlobMatchSimulation)
		})

		It("should match on the body", func() {
			req := sling.New().Get("http://test.com")
			req.Body(bytes.NewBufferString(xml.Header + "<items><item field=something></item></items>"))

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(ioutil.ReadAll(response.Body)).Should(Equal([]byte("glob match")))
		})

		It("should not match on no body", func() {
			req := sling.New().Get("http://test.com")

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(502))

			Expect(ioutil.ReadAll(response.Body)).Should(ContainSubstring("There was an error when matching"))
		})
	})
})
