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

		hoverfly.ImportSimulation(functional_tests.XpathSimulation)
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("Using `xpathMatch`", func() {

		It("should match on the body", func() {
			req := sling.New().Get("http://test.com")
			req.Body(bytes.NewBufferString(xml.Header + "<item></item><item></item><item></item><item></item><item></item>"))

			response := hoverfly.Proxy(req)
			Expect(response.StatusCode).To(Equal(200))

			Expect(ioutil.ReadAll(response.Body)).Should(Equal([]byte("xpath match")))
		})
	})
})
