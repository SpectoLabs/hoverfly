package handlers

import (
	"github.com/go-zoo/bone"
	"net/http"
)

type AdminHandler interface {
	RegisterRoutes(*bone.Mux, *AuthHandler)
}

func WriteResponse(response http.ResponseWriter, bytes []byte) {
	response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response.Write(bytes)
}