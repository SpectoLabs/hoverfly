package hoverfly_test

import (
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/SpectoLabs/hoverfly/functional-tests/testdata"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I match with Hoverfly", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("using empty query map request", func() {

		BeforeEach(func() {
			hoverfly.Start()
			hoverfly.SetMode("simulate")
		})

		It("should match", func() {
			hoverfly.ImportSimulation(testdata.EmptyQuery)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(body)).To(Equal("hello"))
		})

		It("should not match with queries", func() {
			hoverfly.ImportSimulation(testdata.EmptyQuery)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com?test=value"))
			Expect(resp.StatusCode).To(Equal(502))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(body)).To(ContainSubstring("There was an error when matching"))
		})

	})

	Context("using no query map request", func() {

		BeforeEach(func() {
			hoverfly.Start()
			hoverfly.SetMode("simulate")
		})

		It("should match", func() {
			hoverfly.ImportSimulation(testdata.NoQuery)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(body)).To(Equal("hello"))
		})

		It("should match with queries", func() {
			hoverfly.ImportSimulation(testdata.NoQuery)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com?test=value"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(body)).To(ContainSubstring("hello"))
		})

	})
})
