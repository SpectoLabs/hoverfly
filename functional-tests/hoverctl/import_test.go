package hoverctl_end_to_end

import (
	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I import with hoverctl", func() {

	Context("without providing a path to write to", func() {

		It("it should fail nicely", func() {
			output := functional_tests.Run(hoverctlBinary, "import")

			Expect(output).To(ContainSubstring("You have not provided a path to simulation"))
			Expect(output).To(ContainSubstring("Try hoverctl import --help for more information"))
		})
	})
})
