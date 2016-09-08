package hoverfly

import (
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
	"net/http"
	"testing"
)

func TestSynthesizeResponse(t *testing.T) {
	RegisterTestingT(t)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	Expect(err).To(BeNil())

	requestDetails, err := models.NewRequestDetailsFromHttpRequest(req)
	Expect(err).To(BeNil())

	sr, err := SynthesizeResponse(req, requestDetails, "./examples/middleware/synthetic_service/synthetic.py")
	Expect(err).To(BeNil())

	Expect(sr.StatusCode).To(Equal(http.StatusOK))
}

func TestSynthesizeResponseWOMiddleware(t *testing.T) {
	RegisterTestingT(t)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	Expect(err).To(BeNil())

	requestDetails, err := models.NewRequestDetailsFromHttpRequest(req)
	Expect(err).To(BeNil())

	_, err = SynthesizeResponse(req, requestDetails, "")
	Expect(err).ToNot(BeNil())

	Expect(err).To(MatchError("Synthesize failed, middleware not provided"))
}

func TestSynthesizeMiddlewareFailure(t *testing.T) {
	RegisterTestingT(t)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	Expect(err).To(BeNil())

	requestDetails, err := models.NewRequestDetailsFromHttpRequest(req)
	Expect(err).To(BeNil())

	_, err = SynthesizeResponse(req, requestDetails, "./examples/middleware/this_is_not_there.py")
	Expect(err).ToNot(BeNil())
}
