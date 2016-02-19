package authentication

type JWTAuthenticationBackend struct {
	privateKey  *rsa.PrivateKey
	PublicKey   *rsa.PublicKey
	AuthBackend backends.AuthBackend
}

const (
	tokenDuration = 72
	expireOffset  = 3600
)

var authBackendInstance *JWTAuthenticationBackend = nil

