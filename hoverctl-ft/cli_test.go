package hoverfly_end_to_end_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"fmt"
	"net/http"
	"time"
	"path/filepath"
	"os/exec"
	"os"
	"strings"
	"strconv"
	"github.com/phayes/freeport"
)

const (
	simulate = "simulate"
	capture = "capture"
	synthesize = "synthesize"
	modify = "modify"
)

var _ = Describe("When I use hoverfly-cli", func() {
	var (
		hoverflyCmd *exec.Cmd

		workingDir, _ = os.Getwd()
		adminPort = freeport.GetPort()
		adminPortAsString = strconv.Itoa(adminPort)

		proxyPort = freeport.GetPort()
		proxyPortAsString = strconv.Itoa(proxyPort)
	)

	Describe("with a running hoverfly", func() {

		BeforeEach(func() {
			hoverflyBinaryUri := filepath.Join(workingDir, "bin/hoverfly")
			hoverflyCmd = exec.Command(hoverflyBinaryUri, "-db", "memory", "-ap", adminPortAsString, "-pp", proxyPortAsString)


			err := hoverflyCmd.Start()

			if err != nil {
				fmt.Println("Unable to start Hoverfly")
				fmt.Println(hoverflyBinaryUri)
				fmt.Println("Is the binary there?")
				os.Exit(1)
			}

			Eventually(func() int {
				resp, err := http.Get(fmt.Sprintf("http://localhost:%v/api/state", adminPort))
				if err == nil {
					return resp.StatusCode
				} else {
					fmt.Println(err.Error())
					return 0
				}
			}, time.Second * 3).Should(BeNumerically("==", http.StatusOK))
		})

		AfterEach(func() {
			hoverflyCmd.Process.Kill()
		})


		Context("I can get the hoverfly's mode", func() {
			cliBinaryUri := filepath.Join(workingDir, "bin/hoverctl")

			It("when hoverfly is in simulate mode", func() {
				SetHoverflyMode("simulate", adminPort)

				out, _ := exec.Command(cliBinaryUri, "mode", "--admin-port=" + adminPortAsString).CombinedOutput()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Hoverfly is set to simulate mode"))
			})

			It("when hoverfly is in capture mode", func() {
				SetHoverflyMode("capture", adminPort)

				out, _ := exec.Command(cliBinaryUri, "mode", "--admin-port=" + adminPortAsString).CombinedOutput()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Hoverfly is set to capture mode"))
			})

			It("when hoverfly is in synthesize mode", func() {
				SetHoverflyMode("synthesize", adminPort)

				out, _ := exec.Command(cliBinaryUri, "mode", "--admin-port=" + adminPortAsString).CombinedOutput()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Hoverfly is set to synthesize mode"))
			})

			It("when hoverfly is in modify mode", func() {
				SetHoverflyMode("modify", adminPort)

				out, _ := exec.Command(cliBinaryUri, "mode", "--admin-port=" + adminPortAsString).CombinedOutput()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Hoverfly is set to modify mode"))
			})
		})

		Context("I can set hoverfly's mode", func() {
			cliBinaryUri := filepath.Join(workingDir, "bin/hoverctl")

			It("to simulate mode", func() {
				setOutput, _ := exec.Command(cliBinaryUri, "mode", "simulate", "--admin-port=" + adminPortAsString).CombinedOutput()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly has been set to simulate mode"))

				getOutput, _ := exec.Command(cliBinaryUri, "mode", "--admin-port=" + adminPortAsString).CombinedOutput()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(ContainSubstring("Hoverfly is set to simulate mode"))
				Expect(GetHoverflyMode(adminPort)).To(Equal(simulate))
			})

			It("to capture mode", func() {
				setOutput, _ := exec.Command(cliBinaryUri, "mode", "capture", "--admin-port=" + adminPortAsString).CombinedOutput()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly has been set to capture mode"))

				getOutput, _ := exec.Command(cliBinaryUri, "mode", "--admin-port=" + adminPortAsString).CombinedOutput()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(ContainSubstring("Hoverfly is set to capture mode"))
				Expect(GetHoverflyMode(adminPort)).To(Equal(capture))
			})

			It("to synthesize mode", func() {
				setOutput, _ := exec.Command(cliBinaryUri, "mode", "synthesize", "--admin-port=" + adminPortAsString).CombinedOutput()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly has been set to synthesize mode"))

				getOutput, _ := exec.Command(cliBinaryUri, "mode", "--admin-port=" + adminPortAsString).CombinedOutput()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(ContainSubstring("Hoverfly is set to synthesize mode"))
				Expect(GetHoverflyMode(adminPort)).To(Equal(synthesize))
			})

			It("to modify mode", func() {
				setOutput, _ := exec.Command(cliBinaryUri, "mode", "modify", "--admin-port=" + adminPortAsString).CombinedOutput()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly has been set to modify mode"))

				getOutput, _ := exec.Command(cliBinaryUri, "mode", "--admin-port=" + adminPortAsString).CombinedOutput()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(ContainSubstring("Hoverfly is set to modify mode"))
				Expect(GetHoverflyMode(adminPort)).To(Equal(modify))
			})
		})
	})
})