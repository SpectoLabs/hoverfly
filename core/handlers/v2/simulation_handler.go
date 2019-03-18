package v2

import (
	"encoding/json"
	"net/http"

	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"

	log "github.com/sirupsen/logrus"
)

type HoverflySimulation interface {
	GetSimulation() (SimulationViewV5, error)
	GetFilteredSimulation(string) (SimulationViewV5, error)
	PutSimulation(SimulationViewV5) SimulationImportResult
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
	mux.Options("/api/v2/simulation", negroni.New(
		negroni.HandlerFunc(this.Options),
	))

	mux.Get("/api/v2/simulation/schema", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.GetSchema),
	))
	mux.Options("/api/v2/simulation/schema", negroni.New(
		negroni.HandlerFunc(this.Options),
	))
}

func (this *SimulationHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	urlPattern := req.URL.Query().Get("urlPattern")

	var err error
	var simulationView SimulationViewV5
	if urlPattern == "" {
		simulationView, err = this.Hoverfly.GetSimulation()
	} else {
		simulationView, err = this.Hoverfly.GetFilteredSimulation(urlPattern)
	}
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, _ := util.JSONMarshal(simulationView)

	handlers.WriteResponse(w, bytes)
}

func (this *SimulationHandler) Put(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	body, _ := ioutil.ReadAll(req.Body)

	simulationView, err := NewSimulationViewFromRequestBody(body)
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	this.Hoverfly.DeleteSimulation()

	result := this.Hoverfly.PutSimulation(simulationView)
	if result.err != nil {

		log.WithFields(log.Fields{
			"body": string(body),
		}).Debug(result.err.Error())

		handlers.WriteErrorResponse(w, "An error occurred: "+result.err.Error(), http.StatusInternalServerError)
		return
	}
	if len(result.WarningMessages) > 0 {
		bytes, _ := util.JSONMarshal(result)

		handlers.WriteResponse(w, bytes)
		return
	}

	this.Get(w, req, next)
}

func (this *SimulationHandler) Delete(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	this.Hoverfly.DeleteSimulation()

	this.Get(w, req, next)
}

func (this *SimulationHandler) Options(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("Allow", "OPTIONS, GET, PUT, DELETE")
	handlers.WriteResponse(w, []byte(""))
}

func (this *SimulationHandler) GetSchema(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	bytes, _ := json.Marshal(SimulationViewV5Schema)

	handlers.WriteResponse(w, bytes)
}

func (this *SimulationHandler) OptionsSchema(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("Allow", "OPTIONS, GET")
}
