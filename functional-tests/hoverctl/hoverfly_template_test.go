package hoverfly_end_to_end_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phayes/freeport"
	"os/exec"
	"strconv"
	"strings"
)

var _ = Describe("When I use hoverfly-cli", func() {

	var (
		hoverflyCmd *exec.Cmd

		adminPort         = freeport.GetPort()
		adminPortAsString = strconv.Itoa(adminPort)

		proxyPort         = freeport.GetPort()
		proxyPortAsString = strconv.Itoa(proxyPort)
	)

	Context("with a running hoverfly", func() {

		BeforeEach(func() {
			hoverflyCmd = startHoverfly(adminPort, proxyPort, workingDirectory)
			WriteConfiguration("localhost", adminPortAsString, proxyPortAsString)
		})

		AfterEach(func() {
			hoverflyCmd.Process.Kill()
		})

		Context("With no templates imported", func() {

			It("getting templates will print out null data", func() {

				out, _ := exec.Command(hoverctlBinary, "templates").Output()
				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("\"data\": null"))
			})

			It("is possible to set request templates by import", func() {

				out, _ := exec.Command(hoverctlBinary, "templates", "testdata/request-template.json").Output()
				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Request template data set in Hoverfly"))
				Expect(output).To(ContainSubstring("\"path\": \"/path1\""))
				Expect(output).To(ContainSubstring("\"path\": \"/path2\""))
			})

		})

		Context("With some templates already imported", func() {

			BeforeEach(func() {
				_, err := exec.Command(hoverctlBinary, "templates", "testdata/request-template.json").Output()
				if err != nil {
					Fail("Template import failed: " + err.Error())
				}
			})

			It("will print out the existing request template data when getting templates", func() {

				out, _ := exec.Command(hoverctlBinary, "templates").Output()
				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("\"path\": \"/path1\""))
				Expect(output).To(ContainSubstring("\"path\": \"/path2\""))
			})

			It("adds the extra request templates when calling import", func() {

				out, _ := exec.Command(hoverctlBinary, "templates", "testdata/request-template.json").Output()
				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Request template data set in Hoverfly"))
				Expect(output).To(ContainSubstring("\"path\": \"/path1\""))
				Expect(output).To(ContainSubstring("\"path\": \"/path2\""))
				Expect(strings.Count(output, "\"path\": \"/path1\"")).To(Equal(2))
				Expect(strings.Count(output, "\"path\": \"/path2\"")).To(Equal(2))
			})
		})
	})
})
