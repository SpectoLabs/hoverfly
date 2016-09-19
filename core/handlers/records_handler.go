package hoverfly

import (
	"net/http"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/views"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"io/ioutil"
	"fmt"
)

type HoverflyRecords interface {
	GetRequestCache() cache.Cache
	GetTemplateCache() matching.RequestTemplateStore
	ImportRequestResponsePairViews(pairViews []views.RequestResponsePairView) error
}


type RecordsHandler struct {
	Hoverfly HoverflyRecords
}

// AllRecordsHandler returns JSON content type http response
func (this *RecordsHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	records, err := this.Hoverfly.GetRequestCache().GetAllEntries()

	if err != nil {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Error("Failed to get data from cache!")

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(500) // can't process this entity
		return
	}

	var pairViews []views.RequestResponsePairView

	for _, v := range records {
		if pair, err := models.NewRequestResponsePairFromBytes(v); err == nil {
			pairView := pair.ConvertToRequestResponsePairView()
			pairViews = append(pairViews, *pairView)
		} else {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	for _, v := range this.Hoverfly.GetTemplateCache() {
		pairView := v.ConvertToRequestResponsePairView()
		pairViews = append(pairViews, pairView)
	}

	w.Header().Set("Content-Type", "application/json")

	var response views.RequestResponsePairPayload
	response.Data = pairViews
	b, err := json.Marshal(response)

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(b)
	return
}

// ImportRecordsHandler - accepts JSON payload and saves it to cache
func (this *RecordsHandler) Post(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	var requests views.RequestResponsePairPayload

	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var response MessageResponse

	if err != nil {
		// failed to read response body
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Could not read request body!")
		http.Error(w, "Failed to read request body.", 400)
		return
	}

	err = json.Unmarshal(body, &requests)

	if err != nil {
		w.WriteHeader(422) // can't process this entity
		return
	}

	err = this.Hoverfly.ImportRequestResponsePairViews(requests.Data)

	if err != nil {
		response.Message = err.Error()
		w.WriteHeader(400)
	} else {
		response.Message = fmt.Sprintf("%d payloads import complete.", len(requests.Data))
	}

	b, err := response.Encode()
	if err != nil {
		// failed to read response body
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Could not encode response body!")
		http.Error(w, "Failed to encode response", 500)
		return
	}

	w.Write(b)

}

// DeleteAllRecordsHandler - deletes all captured requests
func (this *RecordsHandler) Delete(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	err := this.Hoverfly.GetRequestCache().DeleteData()

	w.Header().Set("Content-Type", "application/json")

	var response MessageResponse
	if err != nil {
		if err.Error() == "bucket not found" {
			response.Message = fmt.Sprintf("No records found")
			w.WriteHeader(200)
		} else {
			response.Message = fmt.Sprintf("Something went wrong: %s", err.Error())
			w.WriteHeader(500)
		}
	} else {
		response.Message = "Proxy cache deleted successfuly"
		w.WriteHeader(200)
	}

	b, err := response.Encode()
	if err != nil {
		// failed to read response body
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Could not encode response body!")
		http.Error(w, "Failed to encode response", 500)
		return
	}
	w.Write(b)
	return
}
