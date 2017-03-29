package hoverctl_suite

import (
	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When using the `login` command", func() {

	var (
		hoverfly *functional_tests.Hoverfly

		username = "ft_user"
		password = "ft_password"
	)

	Context("logging into Hoverfly", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start("-auth", "-username", username, "-password", password)

			functional_tests.Run(hoverctlBinary, "targets", "create", "--target", "default", "--admin-port", hoverfly.GetAdminPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
			functional_tests.Run(hoverctlBinary, "targets", "delete", "-f", "--target", "default")
		})

		It("should log you in successfully with correct credentials", func() {
			output := functional_tests.Run(hoverctlBinary, "login", "--username", username, "--password", password)

			Expect(output).To(ContainSubstring("Login successful"))
		})

		It("should not log you with incorrect credentials", func() {
			output := functional_tests.Run(hoverctlBinary, "login", "--username", "incorrect", "--password", "incorrect")

			Expect(output).To(ContainSubstring("Failed to login to Hoverfly"))
		})

		It("should error nicely if username is missing", func() {
			output := functional_tests.Run(hoverctlBinary, "login", "-f", "--password", password)

			Expect(output).To(ContainSubstring("Missing username or password"))
		})

		It("should error nicely if password is missing", func() {
			output := functional_tests.Run(hoverctlBinary, "login", "-f", "--username", username)

			Expect(output).To(ContainSubstring("Missing username or password"))
		})
	})

	Context("logging into Hoverfly with no targets", func() {
		It("should error nicely if there are no targets", func() {
			functional_tests.Run(hoverctlBinary, "targets", "delete", "-f", "--target", "default")
			output := functional_tests.Run(hoverctlBinary, "login", "--username", username, "--password", password)

			Expect(output).To(ContainSubstring("Cannot login without a target"))
		})
	})
})
