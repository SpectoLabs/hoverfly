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
)

var _ = Describe("When I use hoverfly-cli", func() {
	var (
		hoverflyCmd *exec.Cmd

		workingDir, _ = os.Getwd()
		adminPort = "8888"
	)

	Describe("with a running hoverfly", func() {

		BeforeEach(func() {
			hoverflyBinaryUri := filepath.Join(workingDir, "bin/hoverfly")
			hoverflyCmd = exec.Command(hoverflyBinaryUri, "-db", "memory", "-ap", adminPort, "-pp", "8500")

			err := hoverflyCmd.Start()

			if err != nil {
				fmt.Println("Unable to start Hoverfly")
				fmt.Println(hoverflyBinaryUri)
				fmt.Println("Is the binary there?")
				os.Exit(1)
			}

			Eventually(func() int {
				resp, err := http.Get("http://localhost:8888/api/state")
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

				out, _ := exec.Command(cliBinaryUri, "mode").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(Equal("Hoverfly is set to simulate mode"))
			})

			It("when hoverfly is in capture mode", func() {
				SetHoverflyMode("capture", adminPort)

				out, _ := exec.Command(cliBinaryUri, "mode").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(Equal("Hoverfly is set to capture mode"))
			})

			It("when hoverfly is in synthesize mode", func() {
				SetHoverflyMode("synthesize", adminPort)

				out, _ := exec.Command(cliBinaryUri, "mode").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(Equal("Hoverfly is set to synthesize mode"))
			})

			It("when hoverfly is in modify mode", func() {
				SetHoverflyMode("modify", adminPort)

				out, _ := exec.Command(cliBinaryUri, "mode").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(Equal("Hoverfly is set to modify mode"))
			})
		})

		Context("I can set hoverfly's mode", func() {
			cliBinaryUri := filepath.Join(workingDir, "bin/hoverctl")

			It("to simulate mode", func() {
				setOutput, _ := exec.Command(cliBinaryUri, "mode", "simulate").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(Equal("Hoverfly has been set to simulate mode"))

				getOutput, _ := exec.Command(cliBinaryUri, "mode").Output()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(Equal("Hoverfly is set to simulate mode"))
			})

			It("to capture mode", func() {
				setOutput, _ := exec.Command(cliBinaryUri, "mode", "capture").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(Equal("Hoverfly has been set to capture mode"))

				getOutput, _ := exec.Command(cliBinaryUri, "mode").Output()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(Equal("Hoverfly is set to capture mode"))
			})

			It("to synthesize mode", func() {
				setOutput, _ := exec.Command(cliBinaryUri, "mode", "synthesize").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(Equal("Hoverfly has been set to synthesize mode"))

				getOutput, _ := exec.Command(cliBinaryUri, "mode").Output()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(Equal("Hoverfly is set to synthesize mode"))
			})

			It("to modify mode", func() {
				setOutput, _ := exec.Command(cliBinaryUri, "mode", "modify").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(Equal("Hoverfly has been set to modify mode"))

				getOutput, _ := exec.Command(cliBinaryUri, "mode").Output()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(Equal("Hoverfly is set to modify mode"))
			})
		})


	})
})