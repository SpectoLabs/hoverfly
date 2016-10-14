package v2

import (
	"encoding/json"
	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"io/ioutil"
	"net/http"
)

type HoverflyMode interface {
	GetMode() string
	SetMode(string) error
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
}

func (this *HoverflyModeHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	var modeView ModeView
	modeView.Mode = this.Hoverfly.GetMode()

	bytes, _ := json.Marshal(modeView)

	handlers.WriteResponse(w, bytes)
}

func (this *HoverflyModeHandler) Put(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer r.Body.Close()

	var modeView ModeView

	body, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(body, &modeView)
	if err != nil {
		handlers.WriteErrorResponse(w, "Malformed JSON", 400)
		return
	}

	err = this.Hoverfly.SetMode(modeView.Mode)
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), 422)
		return
	}

	this.Get(w, r, next)
}
