package handlers

import (
	"github.com/go-zoo/bone"
	"net/http"
	"encoding/json"
)

type ErrorView struct {
	Error        string `json:"error"`
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
	errorBytes, _ := json.Marshal(errorView)
	response.WriteHeader(code)
	WriteResponse(response, errorBytes)
}