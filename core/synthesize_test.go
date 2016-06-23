package hoverfly

import (
	"github.com/SpectoLabs/hoverfly/core/testutil"
	"net/http"
	"testing"
)

func TestSynthesizeResponse(t *testing.T) {

	req, err := http.NewRequest("GET", "http://example.com", nil)
	testutil.Expect(t, err, nil)

	sr, err := SynthesizeResponse(req, "./examples/middleware/synthetic_service/synthetic.py")
	testutil.Expect(t, err, nil)

	testutil.Expect(t, sr.StatusCode, 200)
}

func TestSynthesizeResponseWOMiddleware(t *testing.T) {

	req, err := http.NewRequest("GET", "http://example.com", nil)
	testutil.Expect(t, err, nil)

	_, err = SynthesizeResponse(req, "")
	testutil.Refute(t, err, nil)

	testutil.Expect(t, err.Error(), "Synthesize failed, middleware not provided")
}

func TestSynthesizeMiddlewareFailure(t *testing.T) {

	req, err := http.NewRequest("GET", "http://example.com", nil)
	testutil.Expect(t, err, nil)

	_, err = SynthesizeResponse(req, "./examples/middleware/this_is_not_there.py")
	testutil.Refute(t, err, nil)
}
