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

	_ = json.Unmarshal(body, &modeView)


	availableModes := map[string]bool{
		"simulate":   true,
		"capture":    true,
		"modify":     true,
		"synthesize": true,
	}

	if modeView.Mode != "" {
		if !availableModes[modeView.Mode] {
			http.Error(w, "Bad mode supplied, available modes: simulate, capture, modify, synthesize.", 400)
			return
		}
		// setting new state
		err := this.Hoverfly.SetMode(modeView.Mode)
		if err != nil {
			http.Error(w, "Hoverfly is currently configured to act as webserver, which can only operate in simulate mode", 403)
			return
		}

	}

	var responseModeView ModeView
	responseModeView.Mode = this.Hoverfly.GetMode()
	bytes, _ := json.Marshal(responseModeView)

	handlers.WriteResponse(w, bytes)
}