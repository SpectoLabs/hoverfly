package v2

import (
	"encoding/json"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type HoverflyCORS interface {
	GetCORS() CORSView
}

type HoverflyCORSHandler struct {
	Hoverfly HoverflyCORS
}

func (h *HoverflyCORSHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/v2/hoverfly/cors", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(h.Get),
	))

	mux.Options("/api/v2/hoverfly/cors", negroni.New(
		negroni.HandlerFunc(h.Options),
	))
}

func (h *HoverflyCORSHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	corsView := h.Hoverfly.GetCORS()

	corsBytes, _ := json.Marshal(corsView)

	handlers.WriteResponse(w, corsBytes)
}

func (h *HoverflyCORSHandler) Options(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("Allow", "OPTIONS, GET, PUT")
	handlers.WriteResponse(w, []byte(""))
}
