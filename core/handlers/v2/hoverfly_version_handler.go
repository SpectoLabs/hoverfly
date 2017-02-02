package v2

import (
	"encoding/json"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type HoverflyVersion interface {
	GetVersion() string
}

type HoverflyVersionHandler struct {
	Hoverfly HoverflyVersion
}

func (this *HoverflyVersionHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/v2/hoverfly/version", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))
}

func (this *HoverflyVersionHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	var versionView VersionView
	versionView.Version = this.Hoverfly.GetVersion()

	bytes, _ := json.Marshal(versionView)

	handlers.WriteResponse(w, bytes)
}
