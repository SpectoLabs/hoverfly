package main

import (
	"net/http"
	"testing"
)

func TestSynthesizeResponse(t *testing.T) {

	req, err := http.NewRequest("GET", "http://example.com", nil)
	expect(t, err, nil)

	sr := synthesizeResponse(req, "./examples/middleware/synthetic_service/synthetic.py")

	expect(t, sr.StatusCode, 200)
}
