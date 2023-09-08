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

	Describe("set template-data-source", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start()
			functional_tests.Run(hoverctlBinary, "targets", "update", "local", "--admin-port", hoverfly.GetAdminPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("should return success on setting template-data-source", func() {
			output := functional_tests.Run(hoverctlBinary, "templating-data-source", "set", "--name", "test-csv", "--filePath", "testdata/test-student-data.csv")

			Expect(output).To(ContainSubstring("Success"))
		})

		It("should override existing template-data-source", func() {
			output1 := functional_tests.Run(hoverctlBinary, "templating-data-source", "set", "--name", "test-csv", "--filePath", "testdata/test-student-data.csv")

			Expect(output1).To(ContainSubstring("Success"))

			output2 := functional_tests.Run(hoverctlBinary, "templating-data-source", "set", "--name", "test-csv", "--filePath", "testdata/test-student-data1.csv")
			Expect(output2).To(ContainSubstring("Success"))

			output3 := functional_tests.Run(hoverctlBinary, "templating-data-source", "get-all")
			Expect(output3).To(ContainSubstring("test-csv"))
			Expect(output3).To(ContainSubstring("1,Test1,20"))

		})

	})

	Describe("delete template-data-source", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start()
			functional_tests.Run(hoverctlBinary, "targets", "update", "local", "--admin-port", hoverfly.GetAdminPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("should return success on deleting template-data-source", func() {
			output := functional_tests.Run(hoverctlBinary, "templating-data-source", "set", "--name", "test-csv", "--filePath", "testdata/test-student-data.csv")
			Expect(output).To(ContainSubstring("Success"))
			output = functional_tests.Run(hoverctlBinary, "templating-data-source", "delete", "--name", "test-csv")
			Expect(output).To(ContainSubstring("Success"))
		})

	})

	Describe("get templating-data-source", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start()
			functional_tests.Run(hoverctlBinary, "targets", "update", "local", "--admin-port", hoverfly.GetAdminPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("should return template data source", func() {
			output := functional_tests.Run(hoverctlBinary, "templating-data-source", "set", "--name", "test-csv", "--filePath", "testdata/test-student-data.csv")

			Expect(output).To(ContainSubstring("Success"))

			output = functional_tests.Run(hoverctlBinary, "templating-data-source", "get-all")
			Expect(output).To(ContainSubstring("test-csv"))
			Expect(output).To(ContainSubstring("test-csv"))
			Expect(output).To(ContainSubstring("1,Test1,45"))
		})

	})
})
