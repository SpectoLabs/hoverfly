package v2

import (
	"encoding/json"
	"net/http"

	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"

	log "github.com/Sirupsen/logrus"
)

type HoverflySimulation interface {
	GetSimulation() (SimulationView, error)
	PutSimulation(SimulationView) error
	DeleteSimulation()
}

type SimulationHandler struct {
	Hoverfly HoverflySimulation
}

func (this *SimulationHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/v2/simulation", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))

	mux.Put("/api/v2/simulation", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Put),
	))

	mux.Delete("/api/v2/simulation", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Delete),
	))
}

func (this *SimulationHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	simulationView, err := this.Hoverfly.GetSimulation()
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, _ := json.Marshal(simulationView)

	handlers.WriteResponse(w, bytes)
}

func (this *SimulationHandler) Put(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	body, _ := ioutil.ReadAll(req.Body)

	var jsonMap map[string]interface{}
	err := json.Unmarshal(body, &jsonMap)
	if err != nil {
		log.WithFields(log.Fields{
			"body": string(body),
		}).Debug(err.Error())

		handlers.WriteErrorResponse(w, "Invalid json", http.StatusBadRequest)
		return
	}

	var simulationView SimulationView

	if path, err := simulationView.GetValidationSchema().Validate(jsonMap); err != nil {
		log.WithFields(log.Fields{
			"body": string(body),
		}).Debug(err.Error())

		handlers.WriteErrorResponse(w, "Json did not match schema: "+path, http.StatusUnprocessableEntity)
		return
	}

	json.Unmarshal(body, &simulationView)

	this.Hoverfly.DeleteSimulation()

	err = this.Hoverfly.PutSimulation(simulationView)
	if err != nil {

		log.WithFields(log.Fields{
			"body": string(body),
		}).Debug(err.Error())

		handlers.WriteErrorResponse(w, "An error occured: "+err.Error(), http.StatusInternalServerError)
		return
	}

	this.Get(w, req, next)
}

func (this *SimulationHandler) Delete(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	this.Hoverfly.DeleteSimulation()

	this.Get(w, req, next)
}
