package v2

import (
	"encoding/json"
	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"net/http"
)

type HoverflySimulation interface {
	GetSimulation() (SimulationView, error)
}

type SimulationHandler struct {
	Hoverfly HoverflySimulation
}

func (this *SimulationHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/v2/simulation", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))
}

func (this *SimulationHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	simulationView, err := this.Hoverfly.GetSimulation()
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}

	bytes, _ := json.Marshal(simulationView)

	handlers.WriteResponse(w, bytes)
}
