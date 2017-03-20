package models_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/models"
	// . "github.com/SpectoLabs/hoverflsy/core/util"
	. "github.com/onsi/gomega"
)

func Test_CachedResponse_Encode_AndDecodeIntoAndOutOfBytes(t *testing.T) {
	RegisterTestingT(t)

	originalUnit := &models.CachedResponse{
		Request: models.RequestDetails{
			Body:        "test",
			Destination: "test.com",
			Headers: map[string][]string{
				"test": []string{"header"},
			},
			Method: "GET",
			Path:   "/test",
			Query:  "?test=query",
			Scheme: "http",
		},
		MatchingPair: &models.RequestTemplateResponsePair{
			RequestTemplate: models.RequestTemplate{},
			Response:        models.ResponseDetails{},
		},
	}

	encodedBytes, err := originalUnit.Encode()
	Expect(err).To(BeNil())

	Expect(encodedBytes).ToNot(BeEmpty())

	unit, err := models.NewCachedResponseFromBytes(encodedBytes)
	Expect(err).To(BeNil())

	Expect(unit).To(Equal(originalUnit))
}
