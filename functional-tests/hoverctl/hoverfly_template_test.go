package hoverctl_end_to_end

import (
	"strings"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I use hoverfly-cli", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	Context("with a running hoverfly", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start()

			WriteConfiguration("localhost", hoverfly.GetAdminPort(), hoverfly.GetProxyPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		Context("With no templates imported", func() {

			It("getting templates will print out null data", func() {

				output := functional_tests.Run(hoverctlBinary, "templates")

				Expect(output).To(ContainSubstring("\"data\": null"))
			})

			It("is possible to set request templates by import", func() {

				output := functional_tests.Run(hoverctlBinary, "templates", "testdata/request-template.json")

				Expect(output).To(ContainSubstring("Request template data set in Hoverfly"))
				Expect(output).To(ContainSubstring("\"path\": \"/path1\""))
				Expect(output).To(ContainSubstring("\"path\": \"/path2\""))
			})

		})

		Context("With some templates already imported", func() {

			BeforeEach(func() {
				functional_tests.Run(hoverctlBinary, "templates", "testdata/request-template.json")
			})

			It("will print out the existing request template data when getting templates", func() {

				output := functional_tests.Run(hoverctlBinary, "templates")

				Expect(output).To(ContainSubstring("\"path\": \"/path1\""))
				Expect(output).To(ContainSubstring("\"path\": \"/path2\""))
			})

			It("adds the extra request templates when calling import", func() {

				output := functional_tests.Run(hoverctlBinary, "templates", "testdata/request-template.json")

				Expect(output).To(ContainSubstring("Request template data set in Hoverfly"))
				Expect(output).To(ContainSubstring("\"path\": \"/path1\""))
				Expect(output).To(ContainSubstring("\"path\": \"/path2\""))
				Expect(strings.Count(output, "\"path\": \"/path1\"")).To(Equal(2))
				Expect(strings.Count(output, "\"path\": \"/path2\"")).To(Equal(2))
			})
		})
	})
})
