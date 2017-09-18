package handlers

import (
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type HealthView struct {
	Message string `json:"message"`
}

type HealthHandler struct{}

func (this *HealthHandler) RegisterRoutes(mux *bone.Mux, am *AuthHandler) {
	mux.Get("/api/health", negroni.New(
		negroni.HandlerFunc(this.Get),
	))
}

func (this *HealthHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	w.Header().Set("Content-Type", "application/json")

	var response HealthView
	response.Message = "Hoverfly is healthy"

	bytes, err := util.JSONMarshal(response)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteResponse(w, bytes)
}
