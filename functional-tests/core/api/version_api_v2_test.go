package api_test

import (
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/api/v2/hoverfly/version", func() {

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

		It("Should get the version", func() {
			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/version")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(string(modeJson)).To(MatchRegexp(`{"version":"v\d+.\d+.\d+(-rc.\d)*"}`))
		})
	})
})
