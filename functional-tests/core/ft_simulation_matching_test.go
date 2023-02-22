package hoverfly_test

import (
	"encoding/xml"
	"io"

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

			Expect(io.ReadAll(response.Body)).Should(Equal([]byte("xpath match")))
		})

		It("should not match on no body", func() {
			req := sling.New().Get("http://test.com")

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(502))

			Expect(io.ReadAll(response.Body)).Should(ContainSubstring("There was an error when matching"))
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

			Expect(io.ReadAll(response.Body)).Should(Equal([]byte("xml match")))
		})

		It("should not match on wrong body", func() {
			req := sling.New().Get("http://test.com")
			req.Body(bytes.NewBufferString("<items><item>one</item><item>two</item></items>"))

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(502))

			Expect(io.ReadAll(response.Body)).Should(ContainSubstring("There was an error when matching"))
		})
	})

	Context("Using `xmlTemplatedMatch`", func() {

		BeforeEach(func() {
			hoverfly.ImportSimulation(testdata.XmlTemplatedMatch)
		})

		It("should match on the body", func() {
			req := sling.New().Get("http://test.com")
			req.Body(bytes.NewBufferString("<items><item>A12345</item><item>here can be any string</item><item>123</item></items>"))

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(io.ReadAll(response.Body)).Should(Equal([]byte("xml match")))
		})

		It("should not match on wrong body", func() {
			req := sling.New().Get("http://test.com")
			req.Body(bytes.NewBufferString("<items><item>A1234</item><item>here can be any string</item><item>123</item></items>"))

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(502))

			Expect(io.ReadAll(response.Body)).Should(ContainSubstring("There was an error when matching"))
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

			Expect(io.ReadAll(response.Body)).Should(Equal([]byte("json match")))
		})

		It("should not match on no body", func() {
			req := sling.New().Get("http://test.com")
			req.Body(bytes.NewBufferString(`{"test": [ ] }`))

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(502))

			Expect(io.ReadAll(response.Body)).Should(ContainSubstring("There was an error when matching"))
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

			Expect(io.ReadAll(response.Body)).Should(Equal([]byte("json match")))
		})

		It("should not match on no body", func() {
			req := sling.New().Get("http://test.com")

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(502))

			Expect(io.ReadAll(response.Body)).Should(ContainSubstring("There was an error when matching"))
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

			Expect(io.ReadAll(response.Body)).Should(Equal([]byte("regex match")))
		})

		It("should not match on no body", func() {
			req := sling.New().Get("http://test.com")

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(502))

			Expect(io.ReadAll(response.Body)).Should(ContainSubstring("There was an error when matching"))
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

			Expect(io.ReadAll(response.Body)).Should(Equal([]byte("glob match")))
		})

		It("should not match on no body", func() {
			req := sling.New().Get("http://test.com")

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(502))

			Expect(io.ReadAll(response.Body)).Should(ContainSubstring("There was an error when matching"))
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

			Expect(io.ReadAll(response.Body)).Should(Equal([]byte("multiple matches")))
		})

		It("should not match on wrong body", func() {
			req := sling.New().Get("http://test.com")
			req.Body(bytes.NewBufferString(xml.Header + "<items><item field=nothing></item></items>"))

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(502))

			Expect(io.ReadAll(response.Body)).Should(ContainSubstring("There was an error when matching"))
		})

		It("should match on the destination", func() {
			req := sling.New().Get("http://destination.com")

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(io.ReadAll(response.Body)).Should(Equal([]byte("multiple matches 2")))
		})

		It("should not match on wrong destination", func() {
			req := sling.New().Get("http://destination.io")

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(502))

			Expect(io.ReadAll(response.Body)).Should(ContainSubstring("There was an error when matching"))
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

			Expect(io.ReadAll(response.Body)).Should(Equal([]byte("header matchers matches")))
		})

		It("should match on the headers", func() {
			req := sling.New().Get("http://test.com")
			req.Set("test2", "one;two;three")

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(io.ReadAll(response.Body)).Should(Equal([]byte("header matchers matches")))
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

			Expect(io.ReadAll(response.Body)).Should(Equal([]byte("query matchers matches")))
		})

		It("should match on the queries", func() {
			req := sling.New().Get("http://test.com?test=test1&test=test2")

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(io.ReadAll(response.Body)).Should(Equal([]byte("query matchers matches")))
		})
	})

	Context("Using array matchers", func() {

		BeforeEach(func() {
			hoverfly.ImportSimulation(testdata.ArrayMatcher)
		})

		It("should match multiple header values with array matcher", func() {
			req := sling.New().Get("http://test.com")
			req.Set("test1", "a;b;c")

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(io.ReadAll(response.Body)).Should(Equal([]byte("array matchers matches")))
		})

		It("should match multiple query values with array matcher ignoring orders", func() {
			req := sling.New().Get("http://test.com?test=value3&test=value1&test=value2")

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(io.ReadAll(response.Body)).Should(Equal([]byte("array matchers matches query")))
		})
	})

	Context("Using JWT matchers", func() {

		BeforeEach(func() {
			hoverfly.ImportSimulation(testdata.JwtMatcher)
		})

		It("should match JWT token in header with JWT matcher", func() {
			req := sling.New().Get("http://test.com")
			req.Set("Authorisation", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(io.ReadAll(response.Body)).Should(Equal([]byte("jwt matchers matches")))
		})
	})

	Context("Using matcher chaining", func() {

		BeforeEach(func() {
			hoverfly.ImportSimulation(testdata.MatcherChaining)
		})

		It("should match on the body", func() {
			req := sling.New().Get("http://test.com")
			req.Body(bytes.NewBufferString(`{"items": [{}, {}, {}, {}, {"name": "pineapple", "price": 1.99}]}`))

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(io.ReadAll(response.Body)).Should(Equal([]byte("matcher chaining")))
		})

	})

	Context("Using form matcher for request body", func() {

		type PseudoOauthParams struct {
			ClientAssertion string `url:"client_assertion,omitempty"`
			GrantType       string `url:"grant_type,omitempty"`
			Code            string `url:"code,omitempty"`
		}

		BeforeEach(func() {
			hoverfly.ImportSimulation(testdata.FormDataMatch)
		})

		It("should match some form data", func() {
			req := sling.New().Post("http://test.com/test").BodyForm(&PseudoOauthParams{
				ClientAssertion: "some-client-assertion",
				GrantType:       "authorization_code",
				Code:            "some-auth-code",
			})

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(io.ReadAll(response.Body)).Should(Equal([]byte("form data matches")))
		})

		It("should match all form data", func() {
			req := sling.New().Post("http://test.com/test").BodyForm(&PseudoOauthParams{
				ClientAssertion: "fake-client-assertion",
				GrantType:       "authorization_code",
				Code:            "fake-auth-code-1",
			})

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(io.ReadAll(response.Body)).Should(Equal([]byte("all form data matches")))
		})

		It("should match jwt in form data", func() {
			req := sling.New().Post("http://test.com/test").BodyForm(&PseudoOauthParams{
				ClientAssertion: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
				GrantType:       "authorization_code",
				Code:            "fake-auth-code-2",
			})

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(io.ReadAll(response.Body)).Should(Equal([]byte("jwt in form data matches")))
		})

		It("should match jwt in form data using matcher chaining", func() {
			req := sling.New().Post("http://test.com/test").BodyForm(&PseudoOauthParams{
				ClientAssertion: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
				GrantType:       "authorization_code",
				Code:            "fake-auth-code-3",
			})

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(io.ReadAll(response.Body)).Should(Equal([]byte("jwt in form data matches with chaining")))
		})

	})
})
