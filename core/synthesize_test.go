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

	err = middleware.SetScript("#!/usr/bin/env python\n" +
		"import sys\n" +
		"import json\n" +
		"\n" +
		"def main():\n" +
		"	data = sys.stdin.readlines()\n" +
		"	payload = data[0]\n" +
		"\n" +
		"	payload_dict = json.loads(payload)\n" +
		"\n" +
		"	payload_dict['response']['status'] = 200" +
		"\n" +
		"	print(json.dumps(payload_dict))\n" +
		"\n" +
		"if __name__ == \"__main__\":\n" +
		"	main()")

	Expect(err).To(BeNil())

	sr, err := SynthesizeResponse(req, requestDetails, middleware)
	Expect(err).To(BeNil())

	Expect(sr.StatusCode).To(Equal(http.StatusOK))
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
