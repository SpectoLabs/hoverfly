package hoverfly

import (
	"testing"
)

func TestIsURLHTTP(t *testing.T) {
	url := "http://somehost.com"

	b := isURL(url)
	expect(t, b, true)
}

func TestIsURLHTTPS(t *testing.T) {
	url := "https://somehost.com"

	b := isURL(url)
	expect(t, b, true)
}

func TestIsURLWrong(t *testing.T) {
	url := "somehost.com"

	b := isURL(url)
	expect(t, b, false)
}

func TestIsURLWrongTLD(t *testing.T) {
	url := "http://somehost."

	b := isURL(url)
	expect(t, b, false)
}

func TestFileExists(t *testing.T) {
	fp := "examples/exports/readthedocs.json"

	ex, err := exists(fp)
	expect(t, ex, true)
	expect(t, err, nil)
}
