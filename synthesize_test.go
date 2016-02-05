package hoverfly

import (
	"net/http"
	"testing"
)

func TestSynthesizeResponse(t *testing.T) {

	req, err := http.NewRequest("GET", "http://example.com", nil)
	expect(t, err, nil)

	sr, err := synthesizeResponse(req, "./examples/middleware/synthetic_service/synthetic.py")
	expect(t, err, nil)

	expect(t, sr.StatusCode, 200)
}

func TestSynthesizeResponseWOMiddleware(t *testing.T) {

	req, err := http.NewRequest("GET", "http://example.com", nil)
	expect(t, err, nil)

	_, err = synthesizeResponse(req, "")
	refute(t, err, nil)

	expect(t, err.Error(), "Synthesize failed, middleware not provided")
}

func TestSynthesizeMiddlewareFailure(t *testing.T) {

	req, err := http.NewRequest("GET", "http://example.com", nil)
	expect(t, err, nil)

	_, err = synthesizeResponse(req, "./examples/middleware/this_is_not_there.py")
	refute(t, err, nil)
}
