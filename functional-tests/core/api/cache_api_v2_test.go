package api_test

import (
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/SpectoLabs/hoverfly/functional-tests/testdata"
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
		hoverfly.ImportSimulation(testdata.JsonPayload)
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("GET", func() {

		It("should cache matches", func() {
			hoverfly.Proxy(sling.New().Get("http://destination-server.com"))
			cacheView := hoverfly.GetCache()

			Expect(cacheView.Cache).To(HaveLen(1))

			Expect(cacheView.Cache[0].MatchingPair.RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
			Expect(cacheView.Cache[0].MatchingPair.RequestMatcher.Destination[0].Value).To(Equal("destination-server.com"))

			Expect(cacheView.Cache[0].MatchingPair.Response.Status).To(Equal(200))
			Expect(cacheView.Cache[0].MatchingPair.Response.Body).To(Equal("destination matched"))
			Expect(cacheView.Cache[0].MatchingPair.Response.EncodedBody).To(BeFalse())
		})

		It("should cache misses alongside closest miss when strongly matching", func() {
			hoverfly.SetModeWithArgs("simulate", v2.ModeArgumentsView{
				MatchingStrategy: util.StringToPointer("strongest"),
			})

			hoverfly.ImportSimulation(testdata.SingleRequestMatcherToResponse)

			hoverfly.Proxy(sling.New().Get("http://unknown-destination.com"))
			cacheView := hoverfly.GetCache()

			Expect(cacheView.Cache).To(HaveLen(1))

			Expect(cacheView.Cache[0].Key).To(Equal("0dd6716f7e5f5f06067de145a2933b2d"))
			Expect(cacheView.Cache[0].MatchingPair).To(BeNil())
			Expect(cacheView.Cache[0].ClosestMiss).ToNot(BeNil())

			Expect(cacheView.Cache[0].ClosestMiss.RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
			Expect(cacheView.Cache[0].ClosestMiss.RequestMatcher.Destination[0].Value).To(Equal("miss"))
			Expect(cacheView.Cache[0].ClosestMiss.MissedFields).To(ConsistOf("destination"))
			Expect(cacheView.Cache[0].ClosestMiss.Response.Body).To(Equal("body"))
		})

		It("should cache misses without closest miss when firstly matching", func() {
			hoverfly.SetModeWithArgs("simulate", v2.ModeArgumentsView{
				MatchingStrategy: util.StringToPointer("first"),
			})

			hoverfly.ImportSimulation(testdata.SingleRequestMatcherToResponse)

			hoverfly.Proxy(sling.New().Get("http://unknown-destination.com"))
			cacheView := hoverfly.GetCache()

			Expect(cacheView.Cache).To(HaveLen(1))

			Expect(cacheView.Cache[0].Key).To(Equal("0dd6716f7e5f5f06067de145a2933b2d"))
			Expect(cacheView.Cache[0].MatchingPair).To(BeNil())
			Expect(cacheView.Cache[0].ClosestMiss).To(BeNil())
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

		It("should be able to set cache size", func() {
			hoverfly.Stop()
			hoverfly.Start("-cache-size=1")

			hoverfly.Proxy(sling.New().Get("http://destination-server.com"))
			hoverfly.Proxy(sling.New().Get("http://unknown-destination.com"))
			cacheView := hoverfly.GetCache()

			Expect(cacheView.Cache).To(HaveLen(1))
			Expect(cacheView.Cache[0].Key).To(Equal("0dd6716f7e5f5f06067de145a2933b2d"))
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
