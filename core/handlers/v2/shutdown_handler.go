package v2

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type ShutdownHandler struct {
}

func (this *ShutdownHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Delete("/api/v2/shutdown", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Delete),
	))
}

func (this *ShutdownHandler) Delete(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	handlers.WriteResponse(w, []byte(""))
	go func() {
		log.Warning("Shutting down")
		os.Exit(0)
	}()
}
