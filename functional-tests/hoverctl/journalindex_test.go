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

	Describe("set journal index", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start()
			functional_tests.Run(hoverctlBinary, "targets", "update", "local", "--admin-port", hoverfly.GetAdminPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("should return success on setting journal index", func() {
			output := functional_tests.Run(hoverctlBinary, "journal-index", "set", "--name", "Request.QueryParam.id")

			Expect(output).To(ContainSubstring("Success"))
		})

	})

	Describe("delete journal index", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start()
			functional_tests.Run(hoverctlBinary, "targets", "update", "local", "--admin-port", hoverfly.GetAdminPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("should return success on deleting journal", func() {
			output := functional_tests.Run(hoverctlBinary, "journal-index", "delete", "--name", "Request.QueryParam.id")

			Expect(output).To(ContainSubstring("Success"))
		})

	})

	Describe("get journal index", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start()
			functional_tests.Run(hoverctlBinary, "targets", "update", "local", "--admin-port", hoverfly.GetAdminPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("should return template data source", func() {
			output := functional_tests.Run(hoverctlBinary, "journal-index", "set", "--name", "Request.QueryParam.id")

			Expect(output).To(ContainSubstring("Success"))

			output = functional_tests.Run(hoverctlBinary, "journal-index", "set", "--name", "Request.Body 'jsonpath' '$.id'")

			output = functional_tests.Run(hoverctlBinary, "journal-index", "get-all")
			Expect(output).To(ContainSubstring("Request.QueryParam.id"))
			Expect(output).To(ContainSubstring("Request.Body jsonpath $.id"))
		})

	})
})
