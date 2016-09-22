package v2

import (
	"encoding/json"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"net/http"
	"github.com/SpectoLabs/hoverfly/core/handlers"
	"io/ioutil"
)

type HoverflyState interface {
	GetMode() string
	SetMode(string) error
}

type ModeView struct {
	Mode        string `json:"mode"`
}

type ErrorView struct {
	Error        string `json:"error"`
}

type HoverflyModeHandler struct {
	Hoverfly HoverflyState
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
		errorView := &ErrorView{Error: "Malformed JSON"}
		errorBytes, _ := json.Marshal(errorView)
		w.WriteHeader(400)
		w.Write(errorBytes)
		return
	}

	err = this.Hoverfly.SetMode(modeView.Mode)
	if err != nil {
		errorView := &ErrorView{Error: err.Error()}
		errorBytes, _ := json.Marshal(errorView)
		w.WriteHeader(422)
		w.Write(errorBytes)
		return
	}

	var responseModeView ModeView
	responseModeView.Mode = this.Hoverfly.GetMode()
	bytes, _ := json.Marshal(responseModeView)

	handlers.WriteResponse(w, bytes)
}