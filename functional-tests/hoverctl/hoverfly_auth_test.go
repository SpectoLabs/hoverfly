package hoverctl_suite

import (
	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I use hoverctl with a running an authenticated hoverfly", func() {
	var (
		hoverfly *functional_tests.Hoverfly

		username = "ft_user"
		password = "ft_password"
	)

	Describe("and the credentials are in the hoverctl config", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start("-db", "boltdb", "-auth", "-username", username, "-password", password)

			WriteConfigurationWithAuth("localhost", hoverfly.GetAdminPort(), hoverfly.GetProxyPort(), false, username, password)
		})

		AfterEach(func() {
			hoverfly.Stop()
			hoverfly.DeleteBoltDb()
		})

		Context("you can get the mode", func() {

			It("and it returns the correct mode", func() {
				output := functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to simulate mode"))
			})
		})

		Context("you can set the mode", func() {

			It("and it correctly sets it", func() {
				output := functional_tests.Run(hoverctlBinary, "mode", "capture")

				Expect(output).To(ContainSubstring("Hoverfly has been set to capture mode"))

				output = functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to capture mode"))
			})
		})

	})

	Describe("and the credentials are not in the hoverctl config", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start("-db", "boltdb", "-auth", "-username", username, "-password", password)

			WriteConfiguration("localhost", hoverfly.GetAdminPort(), hoverfly.GetProxyPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
			hoverfly.DeleteBoltDb()
		})

		Context("you cannot get the mode", func() {

			It("and it returns an error", func() {
				output := functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly requires authentication"))
			})
		})

		Context("you cannot set the mode", func() {

			It("and it returns an error", func() {
				output := functional_tests.Run(hoverctlBinary, "mode", "capture")

				Expect(output).To(ContainSubstring("Hoverfly requires authentication"))
			})
		})
	})
})
