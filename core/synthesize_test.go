package hoverfly

import (
	"net/http"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
)

func Test_SynthesizeResponse_WorksWhenGivenMiddleware(t *testing.T) {
	RegisterTestingT(t)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	Expect(err).To(BeNil())

	requestDetails, err := models.NewRequestDetailsFromHttpRequest(req)
	Expect(err).To(BeNil())

	middleware := &Middleware{}

	err = middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = middleware.SetScript(pythonReflectBody)

	Expect(err).To(BeNil())

	sr, err := SynthesizeResponse(req, requestDetails, middleware)
	Expect(err).To(BeNil())

	Expect(sr.StatusCode).To(Equal(http.StatusCreated))
}

func Test_SynthesizeResponse_WithoutProperlyConfiguredMiddleware(t *testing.T) {
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

func Test_SynthesizeResponse_MiddlewareFailure(t *testing.T) {
	RegisterTestingT(t)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	Expect(err).To(BeNil())

	requestDetails, err := models.NewRequestDetailsFromHttpRequest(req)
	Expect(err).To(BeNil())

	middleware := &Middleware{
		Script: nil,
	}

	_, err = SynthesizeResponse(req, requestDetails, middleware)
	Expect(err).ToNot(BeNil())
}
