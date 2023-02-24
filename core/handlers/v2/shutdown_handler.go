package v2

import (
	"net/http"
	"os"

	"github.com/SpectoLabs/hoverfly/v2/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	log "github.com/sirupsen/logrus"
)

type ShutdownHandler struct {
}

func (this *ShutdownHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Delete("/api/v2/shutdown", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Delete),
	))
	mux.Options("/api/v2/shutdown", negroni.New(
		negroni.HandlerFunc(this.Options),
	))
}

func (this *ShutdownHandler) Delete(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	handlers.WriteResponse(w, []byte(""))
	go func() {
		log.Warning("Shutting down")
		os.Exit(0)
	}()
}

func (this *ShutdownHandler) Options(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("Allow", "OPTIONS, GET")
	handlers.WriteResponse(w, []byte(""))
}
