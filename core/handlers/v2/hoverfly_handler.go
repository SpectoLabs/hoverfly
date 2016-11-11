package v2

import (
	"encoding/json"
	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/SpectoLabs/hoverfly/core/metrics"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"net/http"
)

type Hoverfly interface {
	GetDestination() string
	GetMiddleware() string
	GetMode() string
	GetStats() metrics.Stats
}

type HoverflyHandler struct {
	Hoverfly Hoverfly
}

func (this *HoverflyHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/v2/hoverfly", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))
}

func (this *HoverflyHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	var hoverflyView HoverflyView

	hoverflyView.Destination = this.Hoverfly.GetDestination()
	hoverflyView.Mode = this.Hoverfly.GetMode()
	hoverflyView.Middleware = this.Hoverfly.GetMiddleware()
	hoverflyView.Usage = this.Hoverfly.GetStats()

	bytes, _ := json.Marshal(hoverflyView)

	handlers.WriteResponse(w, bytes)
}
