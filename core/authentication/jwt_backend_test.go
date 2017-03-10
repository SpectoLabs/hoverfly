package authentication_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/authentication"
	"github.com/SpectoLabs/hoverfly/core/authentication/backends"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/dgrijalva/jwt-go"
	. "github.com/onsi/gomega"
)

func TestGenerateToken(t *testing.T) {
	RegisterTestingT(t)

	ab := backends.NewCacheBasedAuthBackend(cache.NewInMemoryCache(), cache.NewInMemoryCache())
	jwtBackend := authentication.InitJWTAuthenticationBackend(ab, []byte("verysecret"), 100)

	Expect(jwtBackend.GenerateToken("userUUIDhereVeryLong", "userx")).ToNot(BeEmpty())
}

func TestAuthenticate(t *testing.T) {
	RegisterTestingT(t)

	ab := backends.NewCacheBasedAuthBackend(cache.NewInMemoryCache(), cache.NewInMemoryCache())
	username := "beloveduser"
	passw := "12345"
	ab.AddUser(username, passw, true)
	jwtBackend := authentication.InitJWTAuthenticationBackend(ab, []byte("verysecret"), 100)
	user := &backends.User{
		Username: string(username),
		Password: string(passw),
		UUID:     "uuid_here",
		IsAdmin:  true}

	Expect(jwtBackend.Authenticate(user)).To(BeTrue())
}

func TestAuthenticateFail(t *testing.T) {
	RegisterTestingT(t)

	ab := backends.NewCacheBasedAuthBackend(cache.NewInMemoryCache(), cache.NewInMemoryCache())

	jwtBackend := authentication.InitJWTAuthenticationBackend(ab, []byte("verysecret"), 100)
	user := &backends.User{
		Username: "shouldntbehere",
		Password: "secret",
		UUID:     "uuid_here",
		IsAdmin:  true}

	Expect(jwtBackend.Authenticate(user)).To(BeFalse())
}

func TestLogout(t *testing.T) {
	RegisterTestingT(t)

	ab := backends.NewCacheBasedAuthBackend(cache.NewInMemoryCache(), cache.NewInMemoryCache())

	jwtBackend := authentication.InitJWTAuthenticationBackend(ab, []byte("verysecret"), 100)

	tokenString := "exampletokenstring"
	token := jwt.New(jwt.SigningMethodHS512)

	Expect(jwtBackend.Logout(tokenString, token)).To(BeNil())

	// checking whether token is in blacklist
	Expect(jwtBackend.IsInBlacklist(tokenString)).To(BeTrue())
}

func TestNotBlacklisted(t *testing.T) {
	RegisterTestingT(t)

	ab := backends.NewCacheBasedAuthBackend(cache.NewInMemoryCache(), cache.NewInMemoryCache())
	jwtBackend := authentication.InitJWTAuthenticationBackend(ab, []byte("verysecret"), 100)

	tokenString := "exampleTokenStringThatIsNotBlacklisted"

	Expect(jwtBackend.IsInBlacklist(tokenString)).To(BeFalse())
}
