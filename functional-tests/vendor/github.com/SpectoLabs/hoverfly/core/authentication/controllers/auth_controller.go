package controllers

import (
	log "github.com/Sirupsen/logrus"

	"encoding/json"
	"github.com/SpectoLabs/hoverfly/core/authentication"
	"github.com/SpectoLabs/hoverfly/core/authentication/backends"
	"net/http"
)

type AllUsersResponse struct {
	Users []backends.User `json:"users"`
}

type AuthController struct {
	AB                 backends.Authentication
	SecretKey          []byte
	JWTExpirationDelta int
	Enabled            bool
}

// GetNewAuthenticationController - returns a pointer to initialised AuthController
func GetNewAuthenticationController(authBackend backends.Authentication, secretKey []byte, exp int, enabled bool) *AuthController {
	return &AuthController{AB: authBackend, SecretKey: secretKey, JWTExpirationDelta: exp, Enabled: enabled}
}

func (a *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if !a.Enabled {
		w.WriteHeader(http.StatusOK)
		// returning dummy token
		token := `{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkR1bW15IHRva2VuIiwiYWRtaW4iOnRydWV9.sKfJparPo3LUmkYoGboBjVfOV3K1qWKUzqx9XFDEsAs"}`
		w.Write([]byte(token))
		return
	}
	requestUser := new(backends.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	responseStatus, token := authentication.Login(requestUser, a.AB, a.SecretKey, a.JWTExpirationDelta)

	w.WriteHeader(responseStatus)
	w.Write(token)
}

func (a *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Content-Type", "application/json")

	requestUser := new(backends.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	w.Write(authentication.RefreshToken(requestUser, a.AB, a.SecretKey, a.JWTExpirationDelta))
}

func (a *AuthController) Logout(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Content-Type", "application/json")
	if !a.Enabled {
		w.WriteHeader(http.StatusOK)
		return
	}

	err := authentication.Logout(r, a.AB, a.SecretKey, a.JWTExpirationDelta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// GetAllUsersHandler - returns a list of all users
func (a *AuthController) GetAllUsersHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	users, err := a.AB.GetAllUsers()

	if err == nil {

		w.Header().Set("Content-Type", "application/json")

		var response AllUsersResponse
		response.Users = users
		b, err := json.Marshal(response)

		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Write(b)
			return
		}
	} else {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Error("Failed to get data from authentication backend!")

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(500)
		return
	}
}
