package hoverfly

import (
	"net/http"
	"testing"

	. "github.com/onsi/gomega"
)

func Test_authFromHeader_ShouldRemoveProxyAuthorizationHeader(t *testing.T) {
	RegisterTestingT(t)
	req, _ := http.NewRequest(http.MethodGet, "localhost:8888", nil)
	req.Header.Add("Proxy-Authorization", "something")

	authFromHeader(req, nil, nil)
	Expect(req.Header).ToNot(HaveKey("Proxy-Authorization"))
}

func Test_authFromHeader_ShouldReturnFalseIfBasicOrBearer(t *testing.T) {
	RegisterTestingT(t)
	req, _ := http.NewRequest(http.MethodGet, "localhost:8888", nil)
	req.Header.Add("Proxy-Authorization", "Something YmVuamloOlBhc3N3b3JkMTIz")

	Expect(authFromHeader(req, nil, nil)).To(BeFalse())
}

func Test_authFromHeader_Basic_ShouldBase64DecodeUsernameAndPassword(t *testing.T) {
	RegisterTestingT(t)
	req, _ := http.NewRequest(http.MethodGet, "localhost:8888", nil)
	req.Header.Add("Proxy-Authorization", "Basic YmVuamloOlBhc3N3b3JkMTIz")

	var basicUsername, basicPassword string

	Expect(authFromHeader(req, func(username, password string) bool {
		basicUsername = username
		basicPassword = password
		return true
	}, nil)).To(BeTrue())

	Expect(basicUsername).To(Equal("benjih"))
	Expect(basicPassword).To(Equal("Password123"))
}

func Test_authFromHeader_Basic_ShouldReturnFalseIfNotBase64Encoded(t *testing.T) {
	RegisterTestingT(t)
	req, _ := http.NewRequest(http.MethodGet, "localhost:8888", nil)
	req.Header.Add("Proxy-Authorization", "Basic benjih:Password123")

	Expect(authFromHeader(req, nil, nil)).To(BeFalse())
}

func Test_authFromHeader_Basic_ShouldReturnFalseIfDecodedBasicCredentialsArentFormattedCorrectly(t *testing.T) {
	RegisterTestingT(t)
	req, _ := http.NewRequest(http.MethodGet, "localhost:8888", nil)
	req.Header.Add("Proxy-Authorization", "Basic YmVuamlo")

	Expect(authFromHeader(req, nil, nil)).To(BeFalse())
}

func Test_authFromHeader_Bearer_ShouldPassJwtTokenOntoFunction(t *testing.T) {
	RegisterTestingT(t)
	req, _ := http.NewRequest(http.MethodGet, "localhost:8888", nil)
	req.Header.Add("Proxy-Authorization", "Bearer gregg.EEewGREQ.GDSG")

	var bearerToken string

	Expect(authFromHeader(req, nil, func(token string) bool {
		bearerToken = token
		return true
	})).To(BeTrue())

	Expect(bearerToken).To(Equal("gregg.EEewGREQ.GDSG"))
}
