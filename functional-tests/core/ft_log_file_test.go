package hoverfly_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	functional_tests "github.com/SpectoLabs/hoverfly/functional-tests"
)

var _ = Describe("When running Hoverfly as a webserver", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()

	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("with flag -logs-output=file", func() {

		It("and it should write logs into file hoverfly.log in plaintext", func() {

			hoverfly.Start("-webserver", "-logs-output=file")

			out, err := hoverfly.GetStdOut()

			Expect(out).To(Equal(""))
			Expect(err).To(BeNil())

			text, err := hoverfly.GetLogFile("hoverfly.log")
			Expect(text).NotTo(ContainSubstring("[36mINFO"))
			Expect(text).To(ContainSubstring("level=info msg=\"Webserver prepared...\""))
			Expect(err).To(BeNil())
		})

		It("and it should write logs into specified log file in plaintext", func() {

			hoverfly.Start("-webserver", "-logs-output=file", "-logs-file=log-test.log")
			out, err := hoverfly.GetStdOut()

			Expect(err).To(BeNil())
			Expect(out).To(Equal(""))

			text, err := hoverfly.GetLogFile("log-test.log")
			Expect(text).NotTo(ContainSubstring("[36mINFO"))
			Expect(text).To(ContainSubstring("level=info msg=\"Webserver prepared...\""))
			Expect(err).To(BeNil())

		})

	})

	Context("with flag -logs-output=console", func() {

		It("and it should show logs on console with colors", func() {

			hoverfly.Start("-webserver", "-logs-output=console")
			out, err := hoverfly.GetStdOut()

			Expect(err).To(BeNil())
			// assert ANSI escape code
			Expect(out).To(ContainSubstring("[36mINFO"))
			Expect(out).To(ContainSubstring("Webserver prepared..."))
		})

	})

	Context("with flag -logs-output=console and -logs-out=file", func() {

		It("and it should show logs on console with colors and write logs into file in plaintext", func() {

			hoverfly.Start("-webserver", "-logs-output=console", "-logs-output=file")
			out, err := hoverfly.GetStdOut()

			Expect(err).To(BeNil())
			Expect(out).To(ContainSubstring("[36mINFO"))
			Expect(out).To(ContainSubstring("Webserver prepared..."))

			text, err := hoverfly.GetLogFile("hoverfly.log")
			Expect(text).NotTo(ContainSubstring("[36mINFO"))
			Expect(text).To(ContainSubstring("level=info msg=\"Webserver prepared...\""))
			Expect(err).To(BeNil())
		})

		It("and it should show logs on console with colors and write logs into specified file in plaintext", func() {

			hoverfly.Start("-webserver", "-logs-output=console", "-logs-output=file", "-logs-file=test.log")
			out, err := hoverfly.GetStdOut()

			Expect(err).To(BeNil())
			Expect(out).To(ContainSubstring("[36mINFO"))
			Expect(out).To(ContainSubstring("Webserver prepared..."))

			text, err := hoverfly.GetLogFile("test.log")
			Expect(text).NotTo(ContainSubstring("[36mINFO"))
			Expect(text).To(ContainSubstring("level=info msg=\"Webserver prepared...\""))
			Expect(err).To(BeNil())
		})

	})
})
