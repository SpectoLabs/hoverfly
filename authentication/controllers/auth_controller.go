package controllers

import (
	log "github.com/Sirupsen/logrus"

	"encoding/json"
	"github.com/SpectoLabs/hoverfly/authentication"
	"github.com/SpectoLabs/hoverfly/authentication/backends"
	"net/http"
)

type AllUsersResponse struct {
	Users []backends.User `json:"users"`
}

type AuthController struct {
	AB backends.AuthBackend
}

func GetNewAuthenticationController(authBackend backends.AuthBackend) *AuthController {
	return &AuthController{AB: authBackend}
}

func (a *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	requestUser := new(backends.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	responseStatus, token := authentication.Login(requestUser, a.AB)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(token)
}

func (a *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestUser := new(backends.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	w.Header().Set("Content-Type", "application/json")
	w.Write(authentication.RefreshToken(requestUser, a.AB))
}

func (a *AuthController) Logout(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	err := authentication.Logout(r, a.AB)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (a *AuthController) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
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
