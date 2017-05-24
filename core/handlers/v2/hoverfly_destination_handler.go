package v2

import (
	"encoding/json"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type HoverflyDestination interface {
	GetDestination() string
	SetDestination(string) error
}

type HoverflyDestinationHandler struct {
	Hoverfly HoverflyDestination
}

func (this *HoverflyDestinationHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/v2/hoverfly/destination", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))
	mux.Put("/api/v2/hoverfly/destination", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Put),
	))
	mux.Options("/api/v2/hoverfly/destination", negroni.New(
		negroni.HandlerFunc(this.Options),
	))
}

func (this *HoverflyDestinationHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	var destinationView DestinationView
	destinationView.Destination = this.Hoverfly.GetDestination()

	bytes, _ := json.Marshal(destinationView)

	handlers.WriteResponse(w, bytes)
}

func (this *HoverflyDestinationHandler) Put(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var destinationView DestinationView
	err := handlers.ReadFromRequest(r, &destinationView)
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), 400)
		return
	}

	err = this.Hoverfly.SetDestination(destinationView.Destination)
	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), 422)
		return
	}

	this.Get(w, r, next)
}

func (this *HoverflyDestinationHandler) Options(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("Allow", "OPTIONS, GET, PUT")
	handlers.WriteResponse(w, []byte(""))
}
