package authentication

import (
	"os"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/authentication/backends"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/dgrijalva/jwt-go"
)

// TestMain prepares database for testing and then performs a cleanup
func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	// delete test database
	teardown()
	// call with result of m.Run()
	os.Exit(retCode)
}

func TestGenerateToken(t *testing.T) {
	ab := backends.NewCacheBasedAuthBackend(cache.NewBoltDBCache(TestDB, []byte(backends.TokenBucketName)), cache.NewBoltDBCache(TestDB, []byte(backends.UserBucketName)))
	jwtBackend := InitJWTAuthenticationBackend(ab, []byte("verysecret"), 100)

	token, err := jwtBackend.GenerateToken("userUUIDhereVeryLong", "userx")
	expect(t, err, nil)
	expect(t, len(token) > 0, true)
}

func TestAuthenticate(t *testing.T) {
	ab := backends.NewCacheBasedAuthBackend(cache.NewBoltDBCache(TestDB, []byte(backends.TokenBucketName)), cache.NewBoltDBCache(TestDB, []byte(backends.UserBucketName)))
	username := "beloveduser"
	passw := "12345"
	ab.AddUser(username, passw, true)

	jwtBackend := InitJWTAuthenticationBackend(ab, []byte("verysecret"), 100)
	user := &backends.User{
		Username: string(username),
		Password: string(passw),
		UUID:     "uuid_here",
		IsAdmin:  true}

	success := jwtBackend.Authenticate(user)
	expect(t, success, true)
}

func TestAuthenticateFail(t *testing.T) {
	ab := backends.NewCacheBasedAuthBackend(cache.NewBoltDBCache(TestDB, []byte(backends.TokenBucketName)), cache.NewBoltDBCache(TestDB, []byte(backends.UserBucketName)))

	jwtBackend := InitJWTAuthenticationBackend(ab, []byte("verysecret"), 100)
	user := &backends.User{
		Username: "shouldntbehere",
		Password: "secret",
		UUID:     "uuid_here",
		IsAdmin:  true}

	success := jwtBackend.Authenticate(user)
	expect(t, success, false)
}

func TestLogout(t *testing.T) {
	ab := backends.NewCacheBasedAuthBackend(cache.NewBoltDBCache(TestDB, []byte(backends.TokenBucketName)), cache.NewBoltDBCache(TestDB, []byte(backends.UserBucketName)))

	jwtBackend := InitJWTAuthenticationBackend(ab, []byte("verysecret"), 100)

	tokenString := "exampletokenstring"
	token := jwt.New(jwt.SigningMethodHS512)

	err := jwtBackend.Logout(tokenString, token)
	expect(t, err, nil)

	// checking whether token is in blacklist

	blacklisted := jwtBackend.IsInBlacklist(tokenString)
	expect(t, blacklisted, true)
}

func TestNotBlacklisted(t *testing.T) {
	ab := backends.NewCacheBasedAuthBackend(cache.NewBoltDBCache(TestDB, []byte(backends.TokenBucketName)), cache.NewBoltDBCache(TestDB, []byte(backends.UserBucketName)))

	jwtBackend := InitJWTAuthenticationBackend(ab, []byte("verysecret"), 100)

	tokenString := "exampleTokenStringThatIsNotBlacklisted"

	blacklisted := jwtBackend.IsInBlacklist(tokenString)
	expect(t, blacklisted, false)
}
