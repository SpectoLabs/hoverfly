package authentication

import (
	"fmt"
	"net/http"

	"github.com/SpectoLabs/hoverfly/authentication/backends"
	jwt "github.com/dgrijalva/jwt-go"
)

type AuthMiddleware struct {
	AB                 backends.AuthBackend
	SecretKey          []byte
	JWTExpirationDelta int
	Enabled            bool
}

func GetNewAuthenticationMiddleware(authBackend backends.AuthBackend, secretKey []byte, exp int, enabled bool) *AuthMiddleware {
	return &AuthMiddleware{AB: authBackend, SecretKey: secretKey, JWTExpirationDelta: exp, Enabled: enabled}
}

func (a *AuthMiddleware) RequireTokenAuthentication(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// if auth is disabled - do not check token
	if !a.Enabled {
		next(w, req)
	}

	authBackend := InitJWTAuthenticationBackend(a.AB, a.SecretKey, a.JWTExpirationDelta)

	token, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		} else {
			return authBackend.SecretKey, nil
		}
	})

	if err == nil && token.Valid && !authBackend.IsInBlacklist(req.Header.Get("Authorization")) {
		next(w, req)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
