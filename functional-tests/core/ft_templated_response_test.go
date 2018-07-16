package hoverfly_test

import (
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/SpectoLabs/hoverfly/functional-tests/testdata"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I run Hoverfly", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("in simulate mode", func() {

		BeforeEach(func() {
			hoverfly.Start()
			hoverfly.SetMode("simulate")
		})

		It("should not template response if templating is disabled explicitely", func() {
			hoverfly.ImportSimulation(testdata.TemplatingDisabled)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com?one=foo"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(body)).To(Equal("{{ Request.QueryParam.singular }}"))
		})

		It("should not template response if templating is not explcitely enabled or disabled", func() {
			hoverfly.ImportSimulation(testdata.TemplatingDisabledByDefault)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com?one=foo"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(body)).To(Equal("{{ Request.QueryParam.one }}"))
		})

		It("should template response if templating is enabled and cache template not response", func() {
			hoverfly.ImportSimulation(testdata.TemplatingEnabled)

			hoverfly.WriteLogsIfError()

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com?one=foo"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(body)).To(Equal("foo"))

			resp = hoverfly.Proxy(sling.New().Get("http://test-server.com?one=bar"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err = ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(body)).To(Equal("bar"))
		})

		It("should be able to use state in templating", func() {
			hoverfly.ImportSimulation(testdata.TemplatingEnabledWithStateInBody)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/one"))
			Expect(resp.StatusCode).To(Equal(200))

			resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/two"))
			Expect(resp.StatusCode).To(Equal(200))
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("state for eggs"))
		})

		It("should not crash when templating a response if templating variable does not exist", func() {
			hoverfly.ImportSimulation(testdata.TemplatingEnabled)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com?wrong=foo"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(body)).To(Equal(""))
		})
	})
})
