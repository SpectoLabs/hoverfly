package middleware

import (
	"testing"

	. "github.com/onsi/gomega"
)

func Test_MiddlewareError_ToString(t *testing.T) {
	RegisterTestingT(t)

	err := MiddlewareError{
		Message: "Just a message",
	}
	Expect(err.Error()).To(Equal("Just a message"))

	err.Command = "just -a 'command'"
	Expect(err.Error()).To(Equal("Just a message\nCommand: just -a 'command'"))

	err.Url = "http://justa.url"
	Expect(err.Error()).To(Equal("Just a message\nCommand: just -a 'command'\nURL: http://justa.url"))

	err.Stdin = "{stdin}"
	Expect(err.Error()).To(Equal("Just a message\nCommand: just -a 'command'\nURL: http://justa.url\n\nSTDIN:\n{stdin}"))

	err.Stdout = "{stdout}"
	Expect(err.Error()).To(Equal("Just a message\nCommand: just -a 'command'\nURL: http://justa.url\n\nSTDIN:\n{stdin}\n\nSTDOUT:\n{stdout}"))

	err.Stderr = "{stderr}"
	Expect(err.Error()).To(Equal("Just a message\nCommand: just -a 'command'\nURL: http://justa.url\n\nSTDIN:\n{stdin}\n\nSTDOUT:\n{stdout}\n\nSTDERR:\n{stderr}"))
}
