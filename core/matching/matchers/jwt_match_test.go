package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func TestJwtMatcher_ReturnsFalseForInvalidJwtToken(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JwtMatcher("value", "apple")).To(BeFalse())
}

func TestJwtMatcher_ReturnsTrue_ValidTokenWithMatchingValue(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JwtMatcher(`{"header":{"alg":"HS256"},"payload":{"sub":"1234567890","name":"John Doe"}}`,
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")).To(BeTrue())
}

func TestJwtMatcher_ReturnsFalse_ValidTokenWithDifferentMatcherValue(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JwtMatcher(`{"header":{"alg":"HS256"},"payload":{"sub":"123416767890","name":"JohnDoe"}}`,
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")).To(BeFalse())
}
