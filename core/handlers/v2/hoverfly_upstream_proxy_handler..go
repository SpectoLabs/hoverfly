package v2

import (
	"encoding/json"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type HoverflyUpstreamProxy interface {
	GetUpstreamProxy() string
}

type HoverflyUpstreamProxyHandler struct {
	Hoverfly HoverflyUpstreamProxy
}

func (this *HoverflyUpstreamProxyHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/v2/hoverfly/upstream-proxy", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))
	mux.Options("/api/v2/hoverfly/upstream-proxy", negroni.New(
		negroni.HandlerFunc(this.Options),
	))
}

func (this *HoverflyUpstreamProxyHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	upstreamProxyView := UpstreamProxyView{
		UpstreamProxy: this.Hoverfly.GetUpstreamProxy(),
	}

	bytes, _ := json.Marshal(upstreamProxyView)

	handlers.WriteResponse(w, bytes)
}

func (this *HoverflyUpstreamProxyHandler) Options(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("Allow", "OPTIONS, GET")
	handlers.WriteResponse(w, []byte(""))
}
