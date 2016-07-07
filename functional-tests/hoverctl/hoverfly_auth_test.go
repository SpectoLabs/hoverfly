package hoverfly_end_to_end_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os/exec"
	"strings"
	"strconv"
	"github.com/phayes/freeport"
	"io/ioutil"
	"bytes"
	"fmt"
)

var _ = Describe("When I use hoverctl with a running an authenticated hoverfly", func() {
	var (
		hoverflyCmd *exec.Cmd

		adminPort = freeport.GetPort()
		adminPortAsString = strconv.Itoa(adminPort)

		proxyPort = freeport.GetPort()
		proxyPortAsString = strconv.Itoa(proxyPort)


		username = "ft_user"
		password = "ft_password"
	)

	Describe("and the credentials are in the hoverctl config", func() {

		BeforeEach(func() {
			hoverflyCmd = startHoverflyWithAuth(adminPort, proxyPort, workingDirectory, username, password)
			WriteConfigurationWithAuth("localhost", adminPortAsString, proxyPortAsString, username, password)
		})

		AfterEach(func() {
			hoverflyCmd.Process.Kill()
		})

		Context("you can get the mode", func() {

			It("and it returns the correct mode", func() {
				out, _ := exec.Command(hoverctlBinary, "mode").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Hoverfly is set to simulate mode"))
			})
		})

		Context("you can set the mode", func() {

			It("and it correctly sets it", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "mode", "capture").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly has been set to capture mode"))

				getOutput, _ := exec.Command(hoverctlBinary, "mode").Output()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(ContainSubstring("Hoverfly is set to capture mode"))
			})
		})

		Context("you can manage simulations", func() {

			It("by importing and exporting data", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "import", "mogronalol/twitter:latest").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("mogronalol/twitter:latest imported successfully"))

				setOutput, _ = exec.Command(hoverctlBinary, "export", "benjih/twitter:latest").Output()

				output = strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("benjih/twitter:latest exported successfully"))

				importFile, err1 := ioutil.ReadFile(workingDirectory + "/.hoverfly/cache/mogronalol.twitter.latest.json")
				if err1 != nil {
					Fail("Failed reading test data")
				}

				exportFile, err2 := ioutil.ReadFile(workingDirectory + "/.hoverfly/cache/benjih.twitter.latest.json")
				if err2 != nil {
					Fail("Failed reading test data")
				}
				fmt.Println(string(importFile))
				fmt.Println(string(exportFile))
				Expect(bytes.Equal(importFile, exportFile)).To(BeTrue())
			})

			It("and then wiping hoverfly", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "wipe").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly has been wiped"))
			})
		})
	})
})