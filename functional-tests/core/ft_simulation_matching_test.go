package hoverfly_test

import (
	"encoding/xml"
	"io/ioutil"

	"bytes"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/SpectoLabs/hoverfly/functional-tests/testdata"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("	When using different matchers", func() {

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
			hoverfly.ImportSimulation(testdata.XpathMatch)
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

	Context("Using `xmlMatch`", func() {

		BeforeEach(func() {
			hoverfly.ImportSimulation(testdata.XmlMatch)
		})

		It("should match on the body", func() {
			req := sling.New().Get("http://test.com")
			req.Body(bytes.NewBufferString("<items><item>one</item></items>"))

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(ioutil.ReadAll(response.Body)).Should(Equal([]byte("xml match")))
		})

		It("should not match on wrong body", func() {
			req := sling.New().Get("http://test.com")
			req.Body(bytes.NewBufferString("<items><item>one</item><item>two</item></items>"))

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(502))

			Expect(ioutil.ReadAll(response.Body)).Should(ContainSubstring("There was an error when matching"))
		})
	})

	Context("Using `jsonMatch`", func() {

		BeforeEach(func() {
			hoverfly.ImportSimulation(testdata.JsonMatch)
		})

		It("should match on the body", func() {
			req := sling.New().Get("http://test.com")
			req.Body(bytes.NewBufferString(`{
				"test": "data"
			}`))

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(ioutil.ReadAll(response.Body)).Should(Equal([]byte("json match")))
		})

		It("should not match on no body", func() {
			req := sling.New().Get("http://test.com")
			req.Body(bytes.NewBufferString(`{"test": [ ] }`))

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(502))

			Expect(ioutil.ReadAll(response.Body)).Should(ContainSubstring("There was an error when matching"))
		})
	})

	Context("Using `jsonPathMatch`", func() {

		BeforeEach(func() {
			hoverfly.ImportSimulation(testdata.JsonPathMatch)
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
			hoverfly.ImportSimulation(testdata.RegexMatch)
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
			hoverfly.ImportSimulation(testdata.GlobMatch)
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

	Context("Using multiple matchers", func() {

		BeforeEach(func() {
			hoverfly.ImportSimulation(testdata.MultipleMatch)
		})

		It("should match on the body", func() {
			req := sling.New().Get("http://test.com")
			req.Body(bytes.NewBufferString(xml.Header + "<items><item field=something></item></items>"))

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(ioutil.ReadAll(response.Body)).Should(Equal([]byte("multiple matches")))
		})

		It("should not match on wrong body", func() {
			req := sling.New().Get("http://test.com")
			req.Body(bytes.NewBufferString(xml.Header + "<items><item field=nothing></item></items>"))

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(502))

			Expect(ioutil.ReadAll(response.Body)).Should(ContainSubstring("There was an error when matching"))
		})

		It("should match on the destination", func() {
			req := sling.New().Get("http://destination.com")

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(ioutil.ReadAll(response.Body)).Should(Equal([]byte("multiple matches 2")))
		})

		It("should not match on wrong destination", func() {
			req := sling.New().Get("http://destination.io")

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(502))

			Expect(ioutil.ReadAll(response.Body)).Should(ContainSubstring("There was an error when matching"))
		})
	})

	Context("Using header matchers", func() {

		BeforeEach(func() {
			hoverfly.ImportSimulation(testdata.HeaderMatchers)
		})

		It("should match on the headers", func() {
			req := sling.New().Get("http://test.com")
			req.Set("test", "test")

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(ioutil.ReadAll(response.Body)).Should(Equal([]byte("header matchers matches")))
		})

		It("should match on the headers", func() {
			req := sling.New().Get("http://test.com")
			req.Set("test2", "one;two;three")

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(ioutil.ReadAll(response.Body)).Should(Equal([]byte("header matchers matches")))
		})
	})

	Context("Using query matchers", func() {

		BeforeEach(func() {
			hoverfly.ImportSimulation(testdata.QueryMatchers)
		})

		It("should match on the queries", func() {
			req := sling.New().Get("http://test.com/?test=test")

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(ioutil.ReadAll(response.Body)).Should(Equal([]byte("query matchers matches")))
		})

		It("should match on the queries", func() {
			req := sling.New().Get("http://test.com?test=test1&test=test2")

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(ioutil.ReadAll(response.Body)).Should(Equal([]byte("query matchers matches")))
		})
	})
})
