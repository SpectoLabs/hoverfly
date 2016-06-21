package hoverfly_end_to_end_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os/exec"
	"strings"
	"io/ioutil"
	"strconv"
	"fmt"
)

var _ = Describe("When I use hoverctl", func() {

		Describe("without a running hoverfly", func() {

			BeforeEach(func() {
				exec.Command(hoverctlBinary, "stop", "-v").Run()
			})

			AfterEach(func() {
				exec.Command(hoverctlBinary, "stop", "-v").Run()
			})

			Context("I can control a process of hoverfly", func() {

				It("by starting hoverfly", func() {
					setOutput, _ := exec.Command(hoverctlBinary, "start", "-v").CombinedOutput()

					output := strings.TrimSpace(string(setOutput))
					Expect(output).To(ContainSubstring("Hoverfly is now running"))

					data, err := ioutil.ReadFile("./.hoverfly/hoverfly.8888.8500.pid")

					if err != nil {
						Fail("Could not find pid file")
					}

					Expect(data).ToNot(BeEmpty())

					if _, err := strconv.Atoi(string(data)); err != nil {
						Fail("Pid file not have an integer in it")
					}

				})

				It("by stopping  hoverfly", func() {
					exec.Command(hoverctlBinary, "start").Run()

					setOutput, _ := exec.Command(hoverctlBinary, "stop").CombinedOutput()

					output := strings.TrimSpace(string(setOutput))
					Expect(output).To(ContainSubstring("Hoverfly has been stopped"))

					_, err := ioutil.ReadFile("./.hoverfly/hoverfly.8888.8500.pid")

					if err == nil {
						Fail("Could not find pid file")
					}
 				})

				It("but you cannot start hoverfly if already running", func() {
					setOutput, _ := exec.Command(hoverctlBinary, "start").CombinedOutput()

					output := strings.TrimSpace(string(setOutput))
					Expect(output).To(ContainSubstring("Hoverfly is now running"))

					setOutput, _ = exec.Command(hoverctlBinary, "start").CombinedOutput()

					output = strings.TrimSpace(string(setOutput))
					Expect(output).To(ContainSubstring("Hoverfly is already running"))
				})
				fmt.Println("end")
				setOutput, _ := exec.Command(hoverctlBinary, "stop", "-v").CombinedOutput()
				fmt.Println(string(setOutput))

			})
		})
	})