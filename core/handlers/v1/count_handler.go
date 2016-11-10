package v1

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"net/http"
)

type HoverflyCount interface {
	GetRequestCacheCount() (int, error)
}

type CountHandler struct {
	Hoverfly HoverflyCount
}

func (this *CountHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/count", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))
}

// RecordsCount returns number of captured requests as a JSON payload
func (this *CountHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	count, err := this.Hoverfly.GetRequestCacheCount()

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
