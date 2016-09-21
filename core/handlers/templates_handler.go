package hoverfly

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"io/ioutil"
	"net/http"
)

type HoverflyTemplates interface {
	GetTemplateCache() matching.RequestTemplateStore
	ImportTemplates(pairPayload matching.RequestTemplateResponsePairPayload) error
	DeleteTemplateCache()
}

type TemplatesHandler struct {
	Hoverfly HoverflyTemplates
}

func (this *TemplatesHandler) RegisterRoutes(mux *bone.Mux, am *AuthHandler) {
	mux.Get("/api/templates", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))

	mux.Delete("/api/templates", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Delete),
	))

	mux.Post("/api/templates", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Post),
	))
}

// AllRecordsHandler returns JSON content type http response
func (this *TemplatesHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	requestTemplatePayload := this.Hoverfly.GetTemplateCache().GetPayload()

	w.Header().Set("Content-Type", "application/json")

	requestTemplateJson, err := json.Marshal(requestTemplatePayload)

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Write(requestTemplateJson)
		return
	}
}

func (this *TemplatesHandler) Post(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	var requestTemplatePayload matching.RequestTemplateResponsePairPayload

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

	err = json.Unmarshal(body, &requestTemplatePayload)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Could not read request body as request template JSON!")
		w.WriteHeader(422) // can't process this entity
		return
	}

	err = this.Hoverfly.ImportTemplates(requestTemplatePayload)

	if err != nil {
		response.Message = err.Error()
		w.WriteHeader(400)
	} else {
		response.Message = fmt.Sprintf("%d payloads import complete.", len(*requestTemplatePayload.Data))
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
func (this *TemplatesHandler) Delete(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	this.Hoverfly.DeleteTemplateCache()

	w.Header().Set("Content-Type", "application/json")

	var response MessageResponse
	response.Message = "Template store wiped successfuly"
	w.WriteHeader(200)

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
