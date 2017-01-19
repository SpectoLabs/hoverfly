package authentication

import (
	"encoding/json"
	"github.com/SpectoLabs/hoverfly/core/authentication/backends"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
)

type TokenAuthentication struct {
	Token string `json:"token" form:"token"`
}

func Login(requestUser *backends.User, ab backends.Authentication, secret []byte, exp int) (int, []byte) {
	authBackend := InitJWTAuthenticationBackend(ab, secret, exp)

	if authBackend.Authenticate(requestUser) {
		token, err := authBackend.GenerateToken(requestUser.UUID, requestUser.Username)
		if err != nil {
			return http.StatusInternalServerError, []byte("")
		} else {
			response, _ := json.Marshal(TokenAuthentication{token})
			return http.StatusOK, response
		}
	}

	return http.StatusUnauthorized, []byte("")
}

func RefreshToken(requestUser *backends.User, ab backends.Authentication, secret []byte, exp int) []byte {
	authBackend := InitJWTAuthenticationBackend(ab, secret, exp)
	token, err := authBackend.GenerateToken(requestUser.UUID, requestUser.Username)
	if err != nil {
		panic(err)
	}
	response, err := json.Marshal(TokenAuthentication{token})
	if err != nil {
		panic(err)
	}
	return response
}

func Logout(req *http.Request, ab backends.Authentication, secret []byte, exp int) error {
	authBackend := InitJWTAuthenticationBackend(ab, secret, exp)
	tokenRequest, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
		return authBackend.SecretKey, nil
	})
	if err != nil {
		return err
	}
	tokenString := req.Header.Get("Authorization")
	return authBackend.Logout(tokenString, tokenRequest)
}
