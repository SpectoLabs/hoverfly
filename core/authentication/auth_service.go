package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"

	"time"

	"github.com/SpectoLabs/hoverfly/core/authentication/backends"
	"github.com/dgrijalva/jwt-go"
)

type TokenAuthentication struct {
	Token string `json:"token" form:"token"`
}

type FailedAttempts struct {
	Count      int
	LastFailed time.Time
}

var Attempts FailedAttempts

func HasReachedFailedAttemptsLimit(limit int, timeout string) bool {
	if Attempts.Count >= limit {
		failureTimeout, _ := time.ParseDuration(timeout)

		if time.Now().Sub(Attempts.LastFailed) > failureTimeout {
			Attempts.Count = 0
		} else {
			updateFailedAttempts()
			return true
		}
	}

	return false
}

func updateFailedAttempts() {
	Attempts.Count++
	Attempts.LastFailed = time.Now()
}

func Login(requestUser *backends.User, ab backends.Authentication, secret []byte, exp int) (int, []byte) {
	authBackend := InitJWTAuthenticationBackend(ab, secret, exp)

	if HasReachedFailedAttemptsLimit(3, "10m") {
		return http.StatusTooManyRequests, []byte("")
	}

	if authBackend.Authenticate(requestUser) {
		token, err := authBackend.GenerateToken(requestUser.UUID, requestUser.Username)
		if err != nil {
			return http.StatusInternalServerError, []byte("")
		} else {
			response, _ := json.Marshal(TokenAuthentication{token})
			return http.StatusOK, response
		}
	}

	updateFailedAttempts()
	return http.StatusUnauthorized, []byte("")
}

func IsJwtTokenValid(token string, ab backends.Authentication, secret []byte, exp int) bool {
	authBackend := InitJWTAuthenticationBackend(ab, secret, exp)

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		} else {
			return authBackend.SecretKey, nil
		}
	})

	if err == nil && jwtToken.Valid && !authBackend.IsInBlacklist(token) {
		return true
	} else {
		return false
	}
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
