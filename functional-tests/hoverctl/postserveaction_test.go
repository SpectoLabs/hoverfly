package hoverctl_suite

import (
	functional_tests "github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I use hoverctl", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	Describe("set post-serve-action", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start()
			functional_tests.Run(hoverctlBinary, "targets", "update", "local", "--admin-port", hoverfly.GetAdminPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("should return success on setting post-serve-action", func() {
			output := functional_tests.Run(hoverctlBinary, "post-serve-action", "set", "--binary", "python", "--script", "testdata/add_random_delay.py", "--delay", "1500", "--name", "test-callback")

			Expect(output).To(ContainSubstring("Success"))
		})

	})

	Describe("delete post-serve-action", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start()
			functional_tests.Run(hoverctlBinary, "targets", "update", "local", "--admin-port", hoverfly.GetAdminPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("should return error on deleting invalid post-serve-action", func() {
			output := functional_tests.Run(hoverctlBinary, "post-serve-action", "delete", "--name", "test-callback")

			Expect(output).To(ContainSubstring("invalid action name passed"))
		})

		It("should return success on deleting post-serve-action after setting it", func() {
			output := functional_tests.Run(hoverctlBinary, "post-serve-action", "set", "--binary", "python", "--script", "testdata/add_random_delay.py", "--delay", "1500", "--name", "test-callback")
			Expect(output).To(ContainSubstring("Success"))
			output = functional_tests.Run(hoverctlBinary, "post-serve-action", "delete", "--name", "test-callback")
			Expect(output).To(ContainSubstring("Success"))
		})

	})

	Describe("get post-serve-action", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start()
			functional_tests.Run(hoverctlBinary, "targets", "update", "local", "--admin-port", hoverfly.GetAdminPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("should return post-serve-action", func() {
			output := functional_tests.Run(hoverctlBinary, "post-serve-action", "set", "--binary", "python", "--script", "testdata/add_random_delay.py", "--delay", "1300", "--name", "test-callback")

			Expect(output).To(ContainSubstring("Success"))

			output = functional_tests.Run(hoverctlBinary, "post-serve-action", "get-all")
			Expect(output).To(ContainSubstring("test-callback"))
			Expect(output).To(ContainSubstring("python"))
			Expect(output).To(ContainSubstring("1300"))
		})

	})
})
