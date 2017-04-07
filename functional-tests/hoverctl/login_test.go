package hoverctl_suite

import (
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("hoverctl login", func() {

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

	Context("logging into Hoverfly with hoverfly running", func() {
		It("should error nicely if it cannot connect", func() {
			output := functional_tests.Run(hoverctlBinary, "login", "--username", username, "--password", password)

			Expect(output).To(ContainSubstring("Failed to login to Hoverfly"))
		})
	})

	Context("needing to log in", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start("-auth", "-username", username, "-password", password)

			functional_tests.Run(hoverctlBinary, "targets", "create", "--target", "no-auth", "--admin-port", hoverfly.GetAdminPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("should error when getting the mode", func() {
			output := functional_tests.Run(hoverctlBinary, "mode", "-t", "no-auth")
			Expect(output).To(ContainSubstring("Hoverfly requires authentication"))
			Expect(output).To(ContainSubstring("Run `hoverctl login -t no-auth`"))

			functional_tests.Run(hoverctlBinary, "login", "-t", "no-auth", "--username", username, "--password", password)

			output = functional_tests.Run(hoverctlBinary, "mode", "-t", "no-auth")
			Expect(output).To(ContainSubstring("Hoverfly is currently set to simulate mode"))
		})

		It("should error when setting the mode", func() {
			output := functional_tests.Run(hoverctlBinary, "mode", "-t", "no-auth", "capture")
			Expect(output).To(ContainSubstring("Hoverfly requires authentication"))
			Expect(output).To(ContainSubstring("Run `hoverctl login -t no-auth`"))

			functional_tests.Run(hoverctlBinary, "login", "-t", "no-auth", "--username", username, "--password", password)

			output = functional_tests.Run(hoverctlBinary, "mode", "-t", "no-auth", "capture")
			Expect(output).To(ContainSubstring("Hoverfly has been set to capture mode"))
		})

		It("should error when flushing", func() {
			output := functional_tests.Run(hoverctlBinary, "flush", "-f", "-t", "no-auth")
			Expect(output).To(ContainSubstring("Hoverfly requires authentication"))
			Expect(output).To(ContainSubstring("Run `hoverctl login -t no-auth`"))

			functional_tests.Run(hoverctlBinary, "login", "-t", "no-auth", "--username", username, "--password", password)

			output = functional_tests.Run(hoverctlBinary, "flush", "-f", "-t", "no-auth")
			Expect(output).ToNot(ContainSubstring("Hoverfly requires authentication"))
			Expect(output).ToNot(ContainSubstring("Run `hoverctl login -t no-auth`"))
		})

		It("should error when importing", func() {
			filePath := functional_tests.GenerateFileName()
			ioutil.WriteFile(filePath, []byte(functional_tests.JsonPayload), 0644)

			output := functional_tests.Run(hoverctlBinary, "import", "-t", "no-auth", filePath)
			Expect(output).To(ContainSubstring("Hoverfly requires authentication"))
			Expect(output).To(ContainSubstring("Run `hoverctl login -t no-auth`"))

			functional_tests.Run(hoverctlBinary, "login", "-t", "no-auth", "--username", username, "--password", password)

			output = functional_tests.Run(hoverctlBinary, "import", "-t", "no-auth", filePath)
			Expect(output).ToNot(ContainSubstring("Hoverfly requires authentication"))
			Expect(output).ToNot(ContainSubstring("Run `hoverctl login -t no-auth`"))
		})

		It("should error when exporting", func() {
			filePath := functional_tests.GenerateFileName()

			output := functional_tests.Run(hoverctlBinary, "export", "-t", "no-auth", filePath)
			Expect(output).To(ContainSubstring("Hoverfly requires authentication"))
			Expect(output).To(ContainSubstring("Run `hoverctl login -t no-auth`"))

			functional_tests.Run(hoverctlBinary, "login", "-t", "no-auth", "--username", username, "--password", password)

			output = functional_tests.Run(hoverctlBinary, "export", "-t", "no-auth", filePath)
			Expect(output).To(ContainSubstring("Successfully exported simulation to " + filePath))
		})

		It("should error when deleting", func() {
			output := functional_tests.Run(hoverctlBinary, "delete", "-t", "no-auth", "--force")
			Expect(output).To(ContainSubstring("Hoverfly requires authentication"))
			Expect(output).To(ContainSubstring("Run `hoverctl login -t no-auth`"))

			functional_tests.Run(hoverctlBinary, "login", "-t", "no-auth", "--username", username, "--password", password)

			output = functional_tests.Run(hoverctlBinary, "delete", "-t", "no-auth", "--force")
			Expect(output).To(ContainSubstring("Simulation data has been deleted from Hoverfly"))
		})

		It("should error when changing destination", func() {
			output := functional_tests.Run(hoverctlBinary, "destination", "-t", "no-auth", "example.org")
			Expect(output).To(ContainSubstring("Hoverfly requires authentication"))
			Expect(output).To(ContainSubstring("Run `hoverctl login -t no-auth`"))

			functional_tests.Run(hoverctlBinary, "login", "-t", "no-auth", "--username", username, "--password", password)

			output = functional_tests.Run(hoverctlBinary, "destination", "-t", "no-auth", "example.org")
			Expect(output).To(ContainSubstring("Hoverfly destination has been set to example.org"))
		})

		It("should error when getting middleware", func() {
			output := functional_tests.Run(hoverctlBinary, "middleware", "-t", "no-auth")
			Expect(output).To(ContainSubstring("Hoverfly requires authentication"))
			Expect(output).To(ContainSubstring("Run `hoverctl login -t no-auth`"))

			functional_tests.Run(hoverctlBinary, "login", "-t", "no-auth", "--username", username, "--password", password)

			output = functional_tests.Run(hoverctlBinary, "middleware", "-t", "no-auth")
			Expect(output).To(ContainSubstring("Hoverfly middleware configuration is currently set to"))
		})
	})

	Context("with a target that doesn't exist", func() {
		It("should error", func() {
			output := functional_tests.Run(hoverctlBinary, "login", "--target", "test-target")

			Expect(output).To(ContainSubstring("test-target is not a target"))
			Expect(output).To(ContainSubstring("Run `hoverctl targets new test-target`"))
		})
	})
})
