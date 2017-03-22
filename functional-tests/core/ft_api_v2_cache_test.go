package hoverfly_test

import (
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/antonholmquist/jason"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/api/v2/cache", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
		hoverfly.Start()
		hoverfly.ImportSimulation(functional_tests.JsonPayload)
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("GET", func() {

		It("should cache matches", func() {
			hoverfly.Proxy(sling.New().Get("http://template-server.com"))
			cacheView := hoverfly.GetCache()

			Expect(cacheView.Cache).To(HaveLen(1))

			Expect(*cacheView.Cache[0].MatchingPair.Request.Destination.ExactMatch).To(Equal("template-server.com"))

			Expect(cacheView.Cache[0].MatchingPair.Response.Status).To(Equal(200))
			Expect(cacheView.Cache[0].MatchingPair.Response.Body).To(Equal("template match"))
			Expect(cacheView.Cache[0].MatchingPair.Response.EncodedBody).To(BeFalse())
		})

		It("should cache failures", func() {
			hoverfly.Proxy(sling.New().Get("http://unknown-destination.com"))
			cacheView := hoverfly.GetCache()

			Expect(cacheView.Cache).To(HaveLen(1))

			Expect(cacheView.Cache[0].Key).To(Equal("0dd6716f7e5f5f06067de145a2933b2d"))
			Expect(cacheView.Cache[0].MatchingPair).To(BeNil())
		})

		It("should get error when cache is disabled", func() {
			hoverfly.Stop()
			hoverfly.Start("-disable-cache")

			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/cache")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(500))
			responseJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())

			jsonObject, err := jason.NewObjectFromBytes(responseJson)
			Expect(err).To(BeNil())

			Expect(jsonObject.GetString("error")).Should(Equal("No cache set"))
		})
	})

	Context("DELETE", func() {

		It("should flush cache", func() {

			cacheView := hoverfly.FlushCache()

			Expect(cacheView.Cache).To(HaveLen(0))
		})

		It("should get error when cache is disabled", func() {
			hoverfly.Stop()
			hoverfly.Start("-disable-cache")

			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/cache")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(500))
			responseJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())

			jsonObject, err := jason.NewObjectFromBytes(responseJson)
			Expect(err).To(BeNil())

			Expect(jsonObject.GetString("error")).Should(Equal("No cache set"))
		})
	})
})
