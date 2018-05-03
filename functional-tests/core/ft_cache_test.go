package hoverfly_test

import (
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/SpectoLabs/hoverfly/functional-tests/testdata"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hoverfly cache", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
		hoverfly.Start()
		hoverfly.ImportSimulation(testdata.JsonPayload)
		hoverfly.Proxy(sling.New().Get("http://destination-server.com"))
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	It("should be flushed when changing to capture mode", func() {
		hoverfly.SetMode("capture")

		cacheView := hoverfly.GetCache()

		Expect(cacheView.Cache).To(HaveLen(0))
	})

	It("should be flushed when importing via the API", func() {
		hoverfly.ImportSimulation(testdata.JsonPayload)
		cacheView := hoverfly.GetCache()

		Expect(cacheView.Cache).To(HaveLen(0))
	})

	It("should preload the cache with exact match request matcher when put into simulate mode", func() {
		hoverfly.ImportSimulation(testdata.PreloadCache)

		hoverfly.SetMode("simulate")

		cacheView := hoverfly.GetCache()

		Expect(cacheView.Cache).To(HaveLen(1))
	})

	It("should not preload the cache if simulation does not contain exact match request matcher", func() {
		hoverfly.ImportSimulation(testdata.JsonMatch)

		hoverfly.SetMode("simulate")

		cacheView := hoverfly.GetCache()

		Expect(cacheView.Cache).To(HaveLen(0))
	})

	It("should not cache hits when matching on headers", func() {
		hoverfly.ImportSimulation(testdata.ExactMatch)

		hoverfly.SetMode("simulate")

		req := sling.New().Get("http://test-server.com/path1").Add("Header", "value1")

		res := hoverfly.Proxy(req)
		Expect(res.StatusCode).To(Equal(200))

		responseBytes, err := ioutil.ReadAll(res.Body)
		Expect(err).To(BeNil())

		Expect(responseBytes).To(Equal([]byte("exact match 1")))

		Expect(hoverfly.GetCache().Cache).To(BeEmpty())
	})

	It("should not cache misses when matched on all fields but headers", func() {
		hoverfly.ImportSimulation(testdata.ExactMatch)

		hoverfly.SetMode("simulate")

		req := sling.New().Get("http://test-server.com/path1").Add("Header", "miss")

		res := hoverfly.Proxy(req)
		Expect(res.StatusCode).To(Equal(502))

		Expect(hoverfly.GetCache().Cache).To(BeEmpty())
	})
})
