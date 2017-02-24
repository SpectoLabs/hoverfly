package v2

import (
	"encoding/json"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type HoverflyCache interface {
	GetCache() ([]RequestResponsePairView, error)
	FlushCache() error
}

type CacheHandler struct {
	Hoverfly HoverflyCache
}

func (this *CacheHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/v2/cache", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))

	mux.Delete("/api/v2/cache", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Delete),
	))
}

func (this *CacheHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	pairViews, err := this.Hoverfly.GetCache()
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, _ := json.Marshal(CacheView{RequestResponsePairs: pairViews})

	handlers.WriteResponse(w, bytes)
}

func (this *CacheHandler) Delete(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	err := this.Hoverfly.FlushCache()
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	this.Get(w, req, next)
}
