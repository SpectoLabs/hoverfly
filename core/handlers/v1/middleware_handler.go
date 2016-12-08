package v1

import (
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type MiddlewareHandler struct{}

func (this *MiddlewareHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/middleware", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Redirect),
	))

	mux.Post("/api/middleware", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Redirect),
	))
}

func (this *MiddlewareHandler) Redirect(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	http.Redirect(w, req, "/api/v2/hoverfly/middleware", http.StatusPermanentRedirect)
	return
}
