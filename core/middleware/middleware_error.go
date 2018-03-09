package middleware

import (
	"fmt"
)

type MiddlewareError struct {
	OriginalError error
	Message       string
	Command       string
	Stdin         string
	Stdout        string
	Stderr        string
}

func (m MiddlewareError) Error() string {
	errorString := fmt.Sprintf("%s", m.Message)
	if m.Command != "" {
		errorString = fmt.Sprintf(errorString+"\nCommand: %s", m.Command)
	}
	if m.Stdin != "" {
		errorString = fmt.Sprintf(errorString+"\n\nSTDIN:\n%s", m.Stdin)
	}
	if m.Stdout != "" {
		errorString = fmt.Sprintf(errorString+"\n\nSTDOUT:\n%s", m.Stdout)
	}
	if m.Stderr != "" {
		errorString = fmt.Sprintf(errorString+"\n\nSTDERR:\n%s", m.Stderr)
	}
	return errorString
}
