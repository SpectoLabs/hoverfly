package v2

import (
	"encoding/json"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type HoverflyMode interface {
	GetMode() ModeView
	SetModeWithArguments(ModeView) error
}

type HoverflyModeHandler struct {
	Hoverfly HoverflyMode
}

func (this *HoverflyModeHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/v2/hoverfly/mode", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))
	mux.Put("/api/v2/hoverfly/mode", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Put),
	))
	mux.Options("/api/v2/hoverfly/mode", negroni.New(
		negroni.HandlerFunc(this.Options),
	))
}

func (this *HoverflyModeHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	bytes, _ := json.Marshal(this.Hoverfly.GetMode())

	handlers.WriteResponse(w, bytes)
}

func (this *HoverflyModeHandler) Put(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var modeView ModeView
	err := handlers.ReadFromRequest(r, &modeView)
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), 400)
		return
	}

	err = this.Hoverfly.SetModeWithArguments(modeView)

	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	this.Get(w, r, next)
}

func (this *HoverflyModeHandler) Options(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("Allow", "OPTIONS, GET, PUT")
	handlers.WriteResponse(w, []byte(""))
}
