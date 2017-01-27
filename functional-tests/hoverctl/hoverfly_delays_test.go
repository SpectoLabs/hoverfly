package hoverctl_end_to_end

import (
	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I use hoverfly-cli", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	Describe("with a running hoverfly", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start()

			WriteConfiguration("localhost", hoverfly.GetAdminPort(), hoverfly.GetProxyPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		Context("I can get the response delay config set on hoverfly", func() {

			It("when no delay is set", func() {
				hoverfly.SetMode("simulate")

				output := functional_tests.Run(hoverctlBinary, "delays")

				Expect(output).To(ContainSubstring("Hoverfly has no delays configured"))
			})

		})

		Context("I can update the response delay config set on hoverfly", func() {

			It("when no delay is set", func() {
				hoverfly.SetMode("simulate")

				output := functional_tests.Run(hoverctlBinary, "delays", "testdata/delays.json")

				Expect(output).To(ContainSubstring("Response delays set in Hoverfly"))
				Expect(output).To(ContainSubstring("host1 - 100ms"))
				Expect(output).To(ContainSubstring("POST | host2 - 110ms"))
			})

		})

	})
})
