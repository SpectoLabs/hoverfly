package main

import (
	"net/http"

	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/go-zoo/bone"
)

type RecordedRequests struct {
	Data []Payload `json:"data"`
}

func getBoneRouter(d DBClient) *bone.Mux {
	mux := bone.New()
	mux.Get("/records", http.HandlerFunc(d.AllRecordsHandler))
	// handling static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return mux
}

// AllRecordsHandler returns JSON content type http response
func (d *DBClient) AllRecordsHandler(w http.ResponseWriter, req *http.Request) {
	records, err := d.getAllRecords()

	if err == nil {

		w.Header().Set("Content-Type", "application/json")

		var response RecordedRequests
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
