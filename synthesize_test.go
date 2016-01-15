package main

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestSynthesizeResponse(t *testing.T) {

	req, err := http.NewRequest("GET", "http://example.com", nil)
	expect(t, err, nil)

	sr := synthesizeResponse(req, "./examples/middleware/synthetic_service/synthetic.py")

	expect(t, sr.StatusCode, 200)
}

func TestSynthesizeResponseWOMiddleware(t *testing.T) {

	req, err := http.NewRequest("GET", "http://example.com", nil)
	expect(t, err, nil)

	sr := synthesizeResponse(req, "")

	body, err := ioutil.ReadAll(sr.Body)
	bodys := string(body)
	expect(t, err, nil)

	expect(t, sr.StatusCode, 428)
	expect(t, len(bodys) > 0, true)
}
