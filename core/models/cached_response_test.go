package models_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/models"
	// . "github.com/SpectoLabs/hoverflsy/core/util"
	. "github.com/onsi/gomega"
)

func Test_CachedResponse_EncodeAndDecodeIntoAndOutOfBytes(t *testing.T) {
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
			Query: map[string][]string{
				"test": []string{"query"},
			},
			Scheme: "http",
		},
		MatchingPair: &models.RequestMatcherResponsePair{
			RequestMatcher: models.RequestMatcher{},
			Response:       models.ResponseDetails{},
		},
	}

	encodedBytes, err := originalUnit.Encode()
	Expect(err).To(BeNil())

	Expect(encodedBytes).ToNot(BeEmpty())

	unit, err := models.NewCachedResponseFromBytes(encodedBytes)
	Expect(err).To(BeNil())

	Expect(unit).To(Equal(originalUnit))
}

func Test_CachedResponse_EncodeAndDecode_NilMatchingPair(t *testing.T) {
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
			Query: map[string][]string{
				"test": []string{"query"},
			},
			Scheme: "http",
		},
		MatchingPair: nil,
	}

	encodedBytes, err := originalUnit.Encode()
	Expect(err).To(BeNil())

	Expect(encodedBytes).ToNot(BeEmpty())

	unit, err := models.NewCachedResponseFromBytes(encodedBytes)
	Expect(err).To(BeNil())

	Expect(unit).To(Equal(originalUnit))
}
