package test

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"testing"
)

// ErrPlain is the default error that is returned for functions in this package.
var ErrPlain = errors.New("error")

////////////////////////////////////////////////////////////////

func fileline(i int) string {
	_, file, line, ok := runtime.Caller(i)
	if !ok {
		return ""
	}
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]
	return fmt.Sprintf("%s:%d", file, line)
}

func trace() string {
	trace2 := fileline(2)
	trace3 := fileline(3)
	return "\r\t" + strings.Repeat(" ", len(fmt.Sprintf("%s:", trace2))) + "\r\t" + trace3
}

func message(empty string, msgs ...interface{}) string {
	msg := fmt.Sprintln(msgs...)
	if len(msg) == 0 {
		msg = empty + "\n"
	}
	return msg
}

func printable(s string) string {
	s = strings.Replace(s, "\n", `\n`, -1)
	s = strings.Replace(s, "\r", `\r`, -1)
	s = strings.Replace(s, "\t", `\t`, -1)
	return s
}

////////////////////////////////////////////////////////////////

func That(t *testing.T, condition bool, msgs ...interface{}) {
	if !condition {
		t.Errorf("%s: %s", trace(), message("bad assertion", msgs...))
	}
}

func Error(t *testing.T, err, expected error, msgs ...interface{}) {
	if err != expected {
		t.Errorf("%s: %s   error: %v\nexpected: %v\n", trace(), message("", msgs...), err, expected)
	}
}

func String(t *testing.T, output, expected string, msgs ...interface{}) {
	if output != expected {
		t.Errorf("%s: %s  output: %s\nexpected: %s\n", trace(), message("", msgs...), printable(output), printable(expected))
	}
}

func Bytes(t *testing.T, output, expected []byte, msgs ...interface{}) {
	if !bytes.Equal(output, expected) {
		t.Errorf("%s: %s  output: %s\nexpected: %s\n", trace(), message("", msgs...), printable(string(output)), printable(string(expected)))
	}
}

func Minify(t *testing.T, input string, err error, output, expected string, msgs ...interface{}) {
	if err != nil {
		t.Errorf("%s: %s   given: %s\n   error: %v\n", trace(), message("", msgs...), printable(input), err)
	}
	if output != expected {
		t.Errorf("%s: %s   given: %s\nminified: %s\nexpected: %s\n", trace(), message("", msgs...), printable(input), printable(output), printable(expected))
	}
}
