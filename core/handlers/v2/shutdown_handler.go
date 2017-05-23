package v2

import (
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type ShutdownHandler struct {
}

func (this *ShutdownHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/v2/shutdown", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))
}

func (this *ShutdownHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	handlers.WriteResponse(w, []byte(""))
	go func() {
		log.Warning("Hoverfly will shut down in 10 seconds")
		time.Sleep(time.Second * 10)
		log.Warning("Shutting down")
		os.Exit(0)
	}()
}
