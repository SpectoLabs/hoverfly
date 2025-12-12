package matchers_test

import (
    "testing"

    "github.com/SpectoLabs/hoverfly/core/matching/matchers"
    . "github.com/onsi/gomega"
)

// Token with payload: {"sub":"1234567890","user_name":"stuart.kelly","aud":["svc-a","svc-b"]}
// Note: signature is not validated by matcher; only header/payload are decoded.
const sampleJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwidXNlcl9uYW1lIjoic3R1YXJ0LmtlbGx5IiwiYXVkIjpbInN2Yy1hIiwic3ZjLWIiXX0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

func TestJwtJsonPathMatch_InvalidToken(t *testing.T) {
    RegisterTestingT(t)
    Expect(matchers.JwtJsonPathMatch("$.user_name", "not-a-jwt")).To(BeFalse())
}

func TestJwtJsonPathMatch_FindsPayloadField_WithShorthand(t *testing.T) {
    RegisterTestingT(t)
    Expect(matchers.JwtJsonPathMatch("$.user_name", sampleJWT)).To(BeTrue())
}

func TestJwtJsonPathMatch_FindsPayloadField_WithExplicitPayload(t *testing.T) {
    RegisterTestingT(t)
    Expect(matchers.JwtJsonPathMatch("$.payload.user_name", sampleJWT)).To(BeTrue())
}

func TestJwtJsonPathMatchValueGenerator_ExtractsValue_ForChaining(t *testing.T) {
    RegisterTestingT(t)
    gen := matchers.Matchers[matchers.JWTJsonPath].MatchValueGenerator
    value := gen("$.user_name", sampleJWT)
    Expect(value).To(Equal("stuart.kelly"))
}
