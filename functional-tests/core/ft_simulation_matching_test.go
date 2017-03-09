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

		It("should match on the body", func() {
			hoverfly.ImportSimulation(functional_tests.XpathSimulation)
			req := sling.New().Get("http://test.com")
			req.Body(bytes.NewBufferString(xml.Header + "<item></item><item></item><item></item><item></item><item></item>"))

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(ioutil.ReadAll(response.Body)).Should(Equal([]byte("xpath match")))
		})
	})

	Context("Using `jsonMatch`", func() {

		It("should match on the body", func() {
			hoverfly.ImportSimulation(functional_tests.JsonMatchSimulation)
			req := sling.New().Get("http://test.com")
			req.Body(bytes.NewBufferString(`{"items": [{}, {}, {}, {}, {}]}`))

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(ioutil.ReadAll(response.Body)).Should(Equal([]byte("json match")))
		})
	})

	Context("Using `regexMatch`", func() {

		It("should match on the body", func() {
			hoverfly.ImportSimulation(functional_tests.RegexMatchSimulation)
			req := sling.New().Get("http://test.com")
			req.Body(bytes.NewBufferString(xml.Header + "<items><item field=\"something\"></item></items>"))

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(ioutil.ReadAll(response.Body)).Should(Equal([]byte("regex match")))
		})
	})
})
