package hoverfly

import (
	"net/http"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
)

func TestSynthesizeResponse(t *testing.T) {
	RegisterTestingT(t)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	Expect(err).To(BeNil())

	requestDetails, err := models.NewRequestDetailsFromHttpRequest(req)
	Expect(err).To(BeNil())

	middleware := &Middleware{
		Script: "./examples/middleware/synthetic_service/synthetic.py",
	}

	sr, err := SynthesizeResponse(req, requestDetails, middleware)
	Expect(err).To(BeNil())

	Expect(sr.StatusCode).To(Equal(http.StatusOK))
}

func TestSynthesizeResponseWOMiddleware(t *testing.T) {
	RegisterTestingT(t)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	Expect(err).To(BeNil())

	requestDetails, err := models.NewRequestDetailsFromHttpRequest(req)
	Expect(err).To(BeNil())

	middleware := &Middleware{}

	_, err = SynthesizeResponse(req, requestDetails, middleware)
	Expect(err).ToNot(BeNil())

	Expect(err).To(MatchError("Synthesize failed, middleware not provided"))
}

func TestSynthesizeMiddlewareFailure(t *testing.T) {
	RegisterTestingT(t)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	Expect(err).To(BeNil())

	requestDetails, err := models.NewRequestDetailsFromHttpRequest(req)
	Expect(err).To(BeNil())

	middleware := &Middleware{
		Script: "./examples/middleware/this_is_not_there.py",
	}

	_, err = SynthesizeResponse(req, requestDetails, middleware)
	Expect(err).ToNot(BeNil())
}
