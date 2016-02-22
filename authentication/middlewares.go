package authentication

import (
	log "github.com/Sirupsen/logrus"

	"fmt"
	"github.com/SpectoLabs/hoverfly/authentication/backends"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
)

type AuthMiddleware struct {
	AB                 backends.AuthBackend
	SecretKey          []byte
	JWTExpirationDelta int
}

func GetNewAuthenticationMiddleware(authBackend backends.AuthBackend, secretKey []byte, exp int) *AuthMiddleware {
	return &AuthMiddleware{AB: authBackend, SecretKey: secretKey, JWTExpirationDelta: exp}
}

func (a *AuthMiddleware) RequireTokenAuthentication(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	authBackend := InitJWTAuthenticationBackend(a.AB, a.SecretKey, a.JWTExpirationDelta)

	token, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		} else {
			return authBackend.SecretKey, nil
		}
	})
	log.WithFields(log.Fields{
		"valid":  token.Valid,
		"string": token.Raw,
	}).Warn("token information")

	if err == nil && token.Valid && !authBackend.IsInBlacklist(req.Header.Get("Authorization")) {
		next(rw, req)
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
	}
}
