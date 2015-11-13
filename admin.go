package main

import (
	"io/ioutil"
	"net/http"

	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/go-zoo/bone"
)

// jsonResponse struct encapsulates payload data
type jsonResponse struct {
	Data []Payload `json:"data"`
}

type StateRequest struct {
	Record      bool   `json:"record"`
	Destination string `json:"destination"`
}

// getBoneRouter returns mux for admin interface
func getBoneRouter(d DBClient) *bone.Mux {
	mux := bone.New()
	mux.Get("/records", http.HandlerFunc(d.AllRecordsHandler))
	mux.Get("/state", http.HandlerFunc(d.CurrentStateHandler))
	mux.Post("/state", http.HandlerFunc(d.stateHandler))

	return mux
}

// AllRecordsHandler returns JSON content type http response
func (d *DBClient) AllRecordsHandler(w http.ResponseWriter, req *http.Request) {
	records, err := d.getAllRecords()

	if err == nil {

		w.Header().Set("Content-Type", "application/json")

		var response jsonResponse
		response.Data = records
		b, err := json.Marshal(response)

		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Write(b)
			return
		}
	} else {
		log.WithFields(log.Fields{
			"Error":        err.Error(),
			"PasswordUsed": AppConfig.redisPassword,
		}).Error("Failed to authenticate to Redis!")

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(500) // can't process this entity
		return
	}
}

// CurrentStateHandler returns current state
func (d *DBClient) CurrentStateHandler(w http.ResponseWriter, req *http.Request) {
	var resp StateRequest
	resp.Record = AppConfig.recordState
	resp.Destination = AppConfig.destination

	b, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(b)
}

// stateHandler handles current proxy state
func (d *DBClient) stateHandler(w http.ResponseWriter, r *http.Request) {
	var stateRequest StateRequest

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

	err = json.Unmarshal(body, &stateRequest)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // can't process this entity
		return
	}
	log.WithFields(log.Fields{
		"newState": stateRequest.Record,
		"body":     string(body),
	}).Info("Handling state change request!")

	// setting new state
	AppConfig.recordState = stateRequest.Record

	var resp StateRequest
	resp.Record = stateRequest.Record
	resp.Destination = AppConfig.destination
	b, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(b)

}
