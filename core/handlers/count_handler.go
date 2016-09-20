package hoverfly

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"net/http"
)

type CountHandler struct {
	Hoverfly HoverflyRecords
}

// RecordsCount returns number of captured requests as a JSON payload
func (this *CountHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	count, err := this.Hoverfly.GetRequestCache().RecordsCount()

	if err == nil {

		w.Header().Set("Content-Type", "application/json")

		var response RecordsCount
		response.Count = count
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
