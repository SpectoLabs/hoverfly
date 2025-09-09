package v2

import (
	"encoding/json"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type HoverflyMiddleware interface {
	GetMiddleware() (string, string, string)
	SetMiddleware(string, string, string) error
}

type HoverflyMiddlewareHandler struct {
	Hoverfly HoverflyMiddleware
	Enabled  bool
}

func (this *HoverflyMiddlewareHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/v2/hoverfly/middleware", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))

	if this.Enabled {
		mux.Put("/api/v2/hoverfly/middleware", negroni.New(
			negroni.HandlerFunc(am.RequireTokenAuthentication),
			negroni.HandlerFunc(this.Put),
		))
	}
	mux.Options("/api/v2/hoverfly/middleware", negroni.New(
		negroni.HandlerFunc(this.Options),
	))
}

func (this *HoverflyMiddlewareHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	var middlewareView MiddlewareView
	middlewareView.Binary, middlewareView.Script, middlewareView.Remote = this.Hoverfly.GetMiddleware()

	middlewareBytes, _ := json.Marshal(middlewareView)

	handlers.WriteResponse(w, middlewareBytes)
}

func (this *HoverflyMiddlewareHandler) Put(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	var middlewareReq MiddlewareView
	err := handlers.ReadFromRequest(req, &middlewareReq)
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), 400)
		return
	}

	err = this.Hoverfly.SetMiddleware(middlewareReq.Binary, middlewareReq.Script, middlewareReq.Remote)
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), 422)
		return
	}

	this.Get(w, req, next)
}

func (this *HoverflyMiddlewareHandler) Options(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	allow := "OPTIONS, GET"
	if this.Enabled {
		allow += ", PUT"
	}
	w.Header().Add("Allow", allow)
	handlers.WriteResponse(w, []byte(""))
}
