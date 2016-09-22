package handlers

import (
	"encoding/json"
	"github.com/go-zoo/bone"
	"net/http"
)

type ErrorView struct {
	Error string `json:"error"`
}

type AdminHandler interface {
	RegisterRoutes(*bone.Mux, *AuthHandler)
}

func WriteResponse(response http.ResponseWriter, bytes []byte) {
	response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response.Write(bytes)
}

func WriteErrorResponse(response http.ResponseWriter, message string, code int) {
	errorView := &ErrorView{Error: message}
	errorBytes, err := json.Marshal(errorView)
	if err != nil {
		response.WriteHeader(500)
		return
	}
	response.WriteHeader(code)
	WriteResponse(response, errorBytes)
}
