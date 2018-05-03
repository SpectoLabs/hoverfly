package hoverctl_suite

import (
	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/SpectoLabs/hoverfly/functional-tests/testdata"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("hoverctl flush cache", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
		hoverfly.Start()
		hoverfly.SetMode("simulate")
		hoverfly.ImportSimulation(testdata.JsonPayload)
		hoverfly.Proxy(sling.New().Get("http://destination-server.com"))

		functional_tests.Run(hoverctlBinary, "targets", "update", "local", "--admin-port", hoverfly.GetAdminPort())
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	It("should flush cache", func() {
		output := functional_tests.Run(hoverctlBinary, "flush", "--force")

		Expect(output).To(ContainSubstring("Successfully flushed cache"))

		cacheView := hoverfly.GetCache()

		Expect(cacheView.Cache).To(HaveLen(0))
	})

	It("should error nicely when trying to flush but cache is disabled", func() {
		hoverfly.Stop()
		hoverfly.Start("-disable-cache")
		output := functional_tests.Run(hoverctlBinary, "flush", "--force")

		Expect(output).To(ContainSubstring("Could not flush cache"))
		Expect(output).To(ContainSubstring("No cache set"))
	})

	It("should error nicely when there is no hoverfly", func() {
		functional_tests.Run(hoverctlBinary, "targets", "create", "alt-port", "--admin-port", "12345")
		output := functional_tests.Run(hoverctlBinary, "flush", "-t", "alt-port", "--force")

		Expect(output).To(ContainSubstring("Could not connect to Hoverfly"))
	})

	Context("with a target that doesn't exist", func() {
		It("should error", func() {
			output := functional_tests.Run(hoverctlBinary, "flush", "--target", "test-target")

			Expect(output).To(ContainSubstring("test-target is not a target"))
			Expect(output).To(ContainSubstring("Run `hoverctl targets create test-target`"))
		})
	})

})
