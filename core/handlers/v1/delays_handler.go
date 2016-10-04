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

type HoverflyDelays interface {
	GetResponseDelays() ResponseDelayPayloadView
	SetResponseDelays(ResponseDelayPayloadView) error
	DeleteResponseDelays()
}

type DelaysHandler struct {
	Hoverfly HoverflyDelays
}

func (this *DelaysHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/delays", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))

	mux.Put("/api/delays", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Put),
	))

	mux.Delete("/api/delays", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Delete),
	))
}

func (this *DelaysHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	payloadView := this.Hoverfly.GetResponseDelays()
	bytes, err := json.Marshal(payloadView)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bytes)
}

func (this *DelaysHandler) Put(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	var responseDelaysView ResponseDelayPayloadView
	var mr MessageResponse

	if req.Body == nil {
		req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("")))
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		// failed to read response body
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Could not read response body!")
		mr.Message = fmt.Sprintf("Failed to read request body. Error: %s", err.Error())
		w.WriteHeader(400)

		b, err := mr.Encode()
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

	err = json.Unmarshal(body, &responseDelaysView)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to unmarshal request body!")
		mr.Message = fmt.Sprintf("Failed to decode request body. Error: %s", err.Error())
		w.WriteHeader(400)
	} else if responseDelaysView.Data == nil {
		log.Error("No delay data in the request body!")
		mr.Message = fmt.Sprintf("Failed to get data from the request body.")
		w.WriteHeader(422)
	} else {
		err = this.Hoverfly.SetResponseDelays(responseDelaysView)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("Error validating response delays config supplied")
			mr.Message = fmt.Sprintf("Failed to validate response delays config. Error: %s", err.Error())
			w.WriteHeader(422)
		} else {
			mr.Message = "Response delays updated."
			w.WriteHeader(201)
		}
	}

	b, err := mr.Encode()
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

func (this *DelaysHandler) Delete(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	this.Hoverfly.DeleteResponseDelays()

	var response MessageResponse
	response.Message = "Delays deleted successfuly"
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
