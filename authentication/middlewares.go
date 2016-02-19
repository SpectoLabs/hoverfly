package authentication

import (
	"fmt"
	"github.com/SpectoLabs/hoverfly/authentication/backends"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
)

type AuthMiddleware struct {
	AB backends.AuthBackend
}

func GetNewAuthenticationMiddleware(authBackend backends.AuthBackend) *AuthMiddleware {
	return &AuthMiddleware{AB: authBackend}
}

func (a *AuthMiddleware) RequireTokenAuthentication(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	authBackend := InitJWTAuthenticationBackend(a.AB)

	token, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		} else {
			return authBackend.PublicKey, nil
		}
	})

	if err == nil && token.Valid && !authBackend.IsInBlacklist(req.Header.Get("Authorization")) {
		next(rw, req)
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
	}
}
