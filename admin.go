package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"github.com/meatballhat/negroni-logrus"
)

// recordedRequests struct encapsulates payload data
type recordedRequests struct {
	Data []Payload `json:"data"`
}

type recordsCount struct {
	Count int `json:"count"`
}

type StateRequest struct {
	Mode        string `json:"mode"`
	Destination string `json:"destination"`
}

type messageResponse struct {
	Message string `json:"message"`
}

func (d *DBClient) startAdminInterface() {
	// starting admin interface
	mux := getBoneRouter(*d)
	n := negroni.Classic()
	n.Use(negronilogrus.NewMiddleware())
	n.UseHandler(mux)

	// admin interface starting message
	log.WithFields(log.Fields{
		"AdminPort": d.cfg.adminInterface,
	}).Info("Admin interface is starting...")

	n.Run(d.cfg.adminInterface)
}

// getBoneRouter returns mux for admin interface
func getBoneRouter(d DBClient) *bone.Mux {
	mux := bone.New()

	mux.Get("/records", http.HandlerFunc(d.AllRecordsHandler))
	mux.Delete("/records", http.HandlerFunc(d.DeleteAllRecordsHandler))
	mux.Post("/records", http.HandlerFunc(d.ImportRecordsHandler))

	mux.Get("/count", http.HandlerFunc(d.RecordsCount))

	mux.Get("/state", http.HandlerFunc(d.CurrentStateHandler))
	mux.Post("/state", http.HandlerFunc(d.stateHandler))

	mux.Handle("/*", http.FileServer(http.Dir("static")))

	return mux
}

// AllRecordsHandler returns JSON content type http response
func (d *DBClient) AllRecordsHandler(w http.ResponseWriter, req *http.Request) {
	records, err := d.cache.GetAllRequests()

	if err == nil {

		w.Header().Set("Content-Type", "application/json")

		var response recordedRequests
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
			"Error": err.Error(),
		}).Error("Failed to get data from cache!")

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(500) // can't process this entity
		return
	}
}

func (d *DBClient) RecordsCount(w http.ResponseWriter, req *http.Request) {
	records, err := d.cache.GetAllRequests()

	if err == nil {

		w.Header().Set("Content-Type", "application/json")

		var response recordsCount
		response.Count = len(records)
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
			"Error": err.Error(),
		}).Error("Failed to get data from cache!")

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(500) // can't process this entity
		return
	}
}

func (d *DBClient) ImportRecordsHandler(w http.ResponseWriter, req *http.Request) {

	var requests recordedRequests

	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var response messageResponse

	if err != nil {
		// failed to read response body
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Could not read response body!")
		response.Message = "Bad request. Nothing to import!"
		http.Error(w, "Failed to read request body.", 400)
		return
	}

	err = json.Unmarshal(body, &requests)

	if err != nil {
		w.WriteHeader(422) // can't process this entity
		return
	}

	payloads := requests.Data
	if len(payloads) > 0 {
		for _, pl := range payloads {
			bts, err := json.Marshal(pl)

			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Error("Failed to marshal json")
			} else {
				// recalculating request hash and storing it in database
				r := request{details: pl.Request}
				d.cache.Set([]byte(r.hash()), bts)
			}
		}
		response.Message = fmt.Sprintf("%d requests imported successfully", len(payloads))
	} else {
		response.Message = "Bad request. Nothing to import!"
		w.WriteHeader(400)
	}

	b, err := json.Marshal(response)
	w.Write(b)

}

func (d *DBClient) DeleteAllRecordsHandler(w http.ResponseWriter, req *http.Request) {
	err := d.cache.DeleteBucket(d.cache.requestsBucket)

	w.Header().Set("Content-Type", "application/json")

	var response messageResponse
	if err != nil {
		response.Message = fmt.Sprintf("Something went wrong: %s", err.Error())
		w.WriteHeader(500)
	} else {
		response.Message = "Proxy cache deleted successfuly"
		w.WriteHeader(200)
	}
	b, err := json.Marshal(response)

	w.Write(b)
	return
}

// CurrentStateHandler returns current state
func (d *DBClient) CurrentStateHandler(w http.ResponseWriter, req *http.Request) {
	var resp StateRequest
	resp.Mode = d.cfg.GetMode()
	resp.Destination = d.cfg.destination

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
		"newState": stateRequest.Mode,
		"body":     string(body),
	}).Info("Handling state change request!")

	// setting new state
	d.cfg.SetMode(stateRequest.Mode)

	var resp StateRequest
	resp.Mode = d.cfg.GetMode()
	resp.Destination = d.cfg.destination
	b, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(b)

}
