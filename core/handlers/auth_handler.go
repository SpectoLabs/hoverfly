package handlers

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"encoding/json"

	"github.com/SpectoLabs/hoverfly/core/authentication"
	"github.com/SpectoLabs/hoverfly/core/authentication/backends"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type AuthHandler struct {
	AB                 backends.Authentication
	SecretKey          []byte
	JWTExpirationDelta int
	Enabled            bool
}

func (this *AuthHandler) RegisterRoutes(mux *bone.Mux) {

	mux.Post("/api/token-auth", http.HandlerFunc(this.Login))

	mux.Options("/api/token-auth", negroni.New(
		negroni.HandlerFunc(this.OptionsLogin),
	))

	mux.Get("/api/refresh-token-auth", negroni.New(
		negroni.HandlerFunc(this.RequireTokenAuthentication),
		negroni.HandlerFunc(this.RefreshToken),
	))
	mux.Get("/api/logout", negroni.New(
		negroni.HandlerFunc(this.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Logout),
	))

	mux.Get("/api/users", negroni.New(
		negroni.HandlerFunc(this.RequireTokenAuthentication),
		negroni.HandlerFunc(this.GetAllUsersHandler),
	))
}

func (a *AuthHandler) RequireTokenAuthentication(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// if auth is disabled - do not check token
	if !a.Enabled {
		next(w, req)
		return
	}

	if authorizationValue := req.Header.Get("Authorization"); authorizationValue != "" {
		// Should be a bearer token
		if len(authorizationValue) > 6 && strings.ToUpper(authorizationValue[0:7]) == "BEARER " {
			if authentication.IsJwtTokenValid(authorizationValue[7:], a.AB, a.SecretKey, a.JWTExpirationDelta) {
				next(w, req)
			}
		}
	}

	WriteErrorResponse(w, "", http.StatusUnauthorized)
}

type AllUsersResponse struct {
	Users []backends.User `json:"users"`
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if !a.Enabled {
		// returning dummy token
		token := `{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkR1bW15IHRva2VuIiwiYWRtaW4iOnRydWV9.sKfJparPo3LUmkYoGboBjVfOV3K1qWKUzqx9XFDEsAs"}`
		WriteResponse(w, []byte(token))
		return
	}
	requestUser := new(backends.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	responseStatus, token := authentication.Login(requestUser, a.AB, a.SecretKey, a.JWTExpirationDelta)

	if responseStatus == http.StatusOK {
		WriteResponse(w, token)
	} else {
		WriteErrorResponse(w, "", responseStatus)
	}
}

func (a *AuthHandler) OptionsLogin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("Allow", "OPTIONS, POST")
	WriteResponse(w, []byte(""))
}

func (a *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestUser := new(backends.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	WriteResponse(w, authentication.RefreshToken(requestUser, a.AB, a.SecretKey, a.JWTExpirationDelta))
}

func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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

	WriteResponse(w, []byte(""))
}

// GetAllUsersHandler - returns a list of all users
func (a *AuthHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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
			WriteResponse(w, b)
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
