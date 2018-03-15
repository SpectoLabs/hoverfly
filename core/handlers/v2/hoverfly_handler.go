package v2

import (
	"encoding/json"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/SpectoLabs/hoverfly/core/metrics"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type Hoverfly interface {
	GetDestination() string
	GetMiddleware() (string, string, string)
	GetMode() ModeView
	GetStats() metrics.Stats
	GetVersion() string
	GetState() map[string]string
	SetState(map[string]string)
	PatchState(map[string]string)
	ClearState()
	GetUpstreamProxy() string
	IsWebServer() bool
	GetDiff() map[SimpleRequestDefinitionView][]DiffReport
	ClearDiff()
}

type HoverflyHandler struct {
	Hoverfly Hoverfly
}

func (this *HoverflyHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/v2/hoverfly", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))
	mux.Options("/api/v2/hoverfly", negroni.New(
		negroni.HandlerFunc(this.Options),
	))
}

func (this *HoverflyHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	var hoverflyView HoverflyView

	hoverflyView.Destination = this.Hoverfly.GetDestination()
	hoverflyView.Mode = this.Hoverfly.GetMode().Mode
	hoverflyView.Arguments = this.Hoverfly.GetMode().Arguments
	hoverflyView.Binary, hoverflyView.Script, hoverflyView.Remote = this.Hoverfly.GetMiddleware()
	hoverflyView.Usage = this.Hoverfly.GetStats()
	hoverflyView.Version = this.Hoverfly.GetVersion()
	hoverflyView.UpstreamProxy = this.Hoverfly.GetUpstreamProxy()
	hoverflyView.IsWebServer = this.Hoverfly.IsWebServer()

	bytes, _ := json.Marshal(hoverflyView)

	handlers.WriteResponse(w, bytes)
}

func (this *HoverflyHandler) Options(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("Allow", "OPTIONS, GET")
	handlers.WriteResponse(w, []byte(""))
}
