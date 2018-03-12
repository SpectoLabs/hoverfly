package hoverctl_suite

import (
	"strconv"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phayes/freeport"
)

var _ = Describe("When I use hoverctl", func() {

	var (
		adminPort = strconv.Itoa(freeport.GetPort())
		proxyPort = strconv.Itoa(freeport.GetPort())
	)

	BeforeEach(func() {
		functional_tests.Run(hoverctlBinary, "targets", "create", "local", "--admin-port", adminPort, "--proxy-port", proxyPort)
	})

	AfterEach(func() {
		functional_tests.Run(hoverctlBinary, "stop")
	})

	Context("I can get the logs using the log command in JSON format", func() {

		It("should return the logs", func() {
			functional_tests.Run(hoverctlBinary, "start", "--admin-port="+adminPort, "--proxy-port="+proxyPort)

			output := functional_tests.Run(hoverctlBinary, "logs", "--json")

			Expect(output).To(ContainSubstring(`"Destination":".","Mode":"simulate","ProxyPort":"` + proxyPort + `","level":"info","msg":"Proxy prepared..."`))
		})

		It("should return an error if the logs don't exist", func() {
			functional_tests.Run(hoverctlBinary, "start", "--admin-port="+adminPort, "--proxy-port="+proxyPort)
			functional_tests.Run(hoverctlBinary, "targets", "create", "incorrect", "--admin-port", "12345", "--proxy-port", "65432")

			output := functional_tests.Run(hoverctlBinary, "logs", "--json", "-t", "incorrect")

			Expect(output).To(ContainSubstring("Could not connect to Hoverfly at localhost:12345"))
		})
	})

	Context("I can get the logs using the log command in plaintext format", func() {

		It("should return the logs", func() {
			functional_tests.Run(hoverctlBinary, "start", "--admin-port="+adminPort, "--proxy-port="+proxyPort)

			output := functional_tests.Run(hoverctlBinary, "logs")

			Expect(output).To(ContainSubstring("Proxy prepared..."))
			Expect(output).To(ContainSubstring("=."))
			Expect(output).To(ContainSubstring("=simulate"))
			Expect(output).To(ContainSubstring("=" + proxyPort))
		})

		It("should return an error if the logs don't exist", func() {
			functional_tests.Run(hoverctlBinary, "start", "--admin-port="+adminPort, "--proxy-port="+proxyPort)
			functional_tests.Run(hoverctlBinary, "targets", "create", "incorrect", "--admin-port", "12345", "--proxy-port", "65432")

			output := functional_tests.Run(hoverctlBinary, "logs", "-t", "incorrect")

			Expect(output).To(ContainSubstring("Could not connect to Hoverfly at localhost:12345"))
		})
	})

	Context("with a target that doesn't exist", func() {
		It("should error", func() {
			output := functional_tests.Run(hoverctlBinary, "logs", "--target", "test-target")

			Expect(output).To(ContainSubstring("test-target is not a target"))
			Expect(output).To(ContainSubstring("Run `hoverctl targets create test-target`"))
		})
	})
})
