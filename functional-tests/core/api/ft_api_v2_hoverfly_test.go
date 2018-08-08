package api_test

import (
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/api/v2/hoverfly", func() {

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

	Context("GET", func() {

		It("Should get the hoverfly config", func() {
			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly")

			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))

			hoverflyJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(hoverflyJson).To(MatchRegexp(`"destination":"."`))
			Expect(hoverflyJson).To(MatchRegexp(`"middleware":{"binary":"","script":"","remote":""}`))
			Expect(hoverflyJson).To(MatchRegexp(`"usage":{"counters":{"capture":0,"diff":0,"modify":0,"simulate":0,"spy":0,"synthesize":0}}`))
			Expect(hoverflyJson).To(MatchRegexp(`"version":"v\d+.\d+.\d+"`))
			Expect(hoverflyJson).To(MatchRegexp(`"upstreamProxy":""`))
			Expect(hoverflyJson).To(MatchRegexp(`"mode":"simulate","arguments":{"matchingStrategy":"strongest"}`))
		})
	})
})
