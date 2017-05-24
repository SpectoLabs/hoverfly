package v2

import (
	"encoding/json"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/SpectoLabs/hoverfly/core/metrics"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type HoverflyUsage interface {
	GetStats() metrics.Stats
}

type HoverflyUsageHandler struct {
	Hoverfly HoverflyUsage
}

func (this *HoverflyUsageHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/v2/hoverfly/usage", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))
	mux.Options("/api/v2/hoverfly/usage", negroni.New(
		negroni.HandlerFunc(this.Options),
	))
}

func (this *HoverflyUsageHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	var metricsView UsageView
	metricsView.Usage = this.Hoverfly.GetStats()

	bytes, _ := json.Marshal(metricsView)

	handlers.WriteResponse(w, bytes)
}

func (this *HoverflyUsageHandler) Options(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("Allow", "OPTIONS, GET")
	handlers.WriteResponse(w, []byte(""))
}
