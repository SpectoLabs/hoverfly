package v2

import (
	"io/ioutil"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type HoverflyPAC interface {
	GetPACFile() []byte
	SetPACFile([]byte)
	DeletePACFile()
}

type HoverflyPACHandler struct {
	Hoverfly HoverflyPAC
}

func (this *HoverflyPACHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/v2/hoverfly/pac", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))

	mux.Put("/api/v2/hoverfly/pac", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Put),
	))
	mux.Delete("/api/v2/hoverfly/pac", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Delete),
	))
	mux.Options("/api/v2/hoverfly/pac", negroni.New(
		negroni.HandlerFunc(this.Options),
	))
}

func (this *HoverflyPACHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	pacFile := this.Hoverfly.GetPACFile()
	if pacFile == nil {
		handlers.WriteErrorResponse(w, "Not found", 404)
	}
	handlers.WriteResponse(w, pacFile)
	w.Header().Set("Content-Type", "application/x-ns-proxy-autoconfig")
}

func (this *HoverflyPACHandler) Put(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), 400)
		return
	}
	this.Hoverfly.SetPACFile(bodyBytes)

	this.Get(w, req, next)
}

func (this *HoverflyPACHandler) Delete(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	this.Hoverfly.DeletePACFile()
}

func (this *HoverflyPACHandler) Options(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("Allow", "OPTIONS, GET, PUT, DELETE")
	handlers.WriteResponse(w, []byte(""))
}
