package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"io/ioutil"
	"net/http"
)

type HoverflyState interface {
	GetMode() string
	SetMode(string) error
	GetDestination() string
	SetDestination(string) error
}

type StateHandler struct {
	Hoverfly HoverflyState
}

func (this *StateHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/state", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))
	mux.Post("/api/state", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Post),
	))
}

// CurrentStateHandler returns current state
func (this *StateHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	var resp StateRequest
	resp.Mode = this.Hoverfly.GetMode()
	resp.Destination = this.Hoverfly.GetDestination()

	b, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(b)
}

func (this *StateHandler) Post(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var sr StateRequest

	// this is mainly for testing, since when you create
	if r.Body == nil {
		r.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("")))
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		// failed to read response body
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Could not read response body!")
		http.Error(w, "Failed to read request body.", 400)
		return
	}

	err = json.Unmarshal(body, &sr)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400) // can't process this entity
		return
	}

	availableModes := map[string]bool{
		"simulate":   true,
		"capture":    true,
		"modify":     true,
		"synthesize": true,
	}

	if sr.Mode != "" {
		if !availableModes[sr.Mode] {
			log.WithFields(log.Fields{
				"suppliedMode": sr.Mode,
			}).Error("Wrong mode found, can't change state")
			http.Error(w, "Bad mode supplied, available modes: simulate, capture, modify, synthesize.", 400)
			return
		}
		log.WithFields(log.Fields{
			"newState":    sr.Mode,
			"body":        string(body),
			"destination": sr.Destination,
		}).Info("Handling state change request!")

		// setting new state
		err := this.Hoverfly.SetMode(sr.Mode)
		if err != nil {
			http.Error(w, "Hoverfly is currently configured to act as webserver, which can only operate in simulate mode", 403)
			return
		}

	}

	// checking whether we should update destination
	if sr.Destination != "" {
		err := this.Hoverfly.SetDestination(sr.Destination)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error while updating destination: %s", err.Error()), 500)
			return
		}
	}

	var resp StateRequest
	resp.Mode = this.Hoverfly.GetMode()
	resp.Destination = this.Hoverfly.GetDestination()
	b, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(b)

}
