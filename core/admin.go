package hoverfly

import (
	"fmt"
	"io"
	"net/http"

	// static assets
	_ "github.com/SpectoLabs/hoverfly/core/statik"
	"github.com/rakyll/statik/fs"

	log "github.com/sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
)

type AdminApi struct{}

// Starts the Admin API on a new HTTP port. Port is chosen by
// hoverfly.Cfg.AdminPort.
func (this *AdminApi) StartAdminInterface(hoverfly *Hoverfly) {
	router := bone.New()

	mux := this.addAdminApiRoutes(router, hoverfly)
	mux = this.addDashboardRoutes(router)
	n := negroni.New(negroni.NewRecovery(), negroni.NewStatic(http.Dir("public")))

	n.UseHandler(mux)

	// admin interface starting message
	log.WithFields(log.Fields{
		"AdminPort": hoverfly.Cfg.AdminPort,
	}).Info("Admin interface is starting...")

	http.ListenAndServe(fmt.Sprintf("%s:%s", hoverfly.Cfg.ListenOnHost, hoverfly.Cfg.AdminPort), n)
}

// Will add the handlers to the router.
func (this *AdminApi) addAdminApiRoutes(router *bone.Mux, d *Hoverfly) *bone.Mux {
	authHandler := &handlers.AuthHandler{
		AB:                 d.Authentication,
		SecretKey:          d.Cfg.SecretKey,
		JWTExpirationDelta: d.Cfg.JWTExpirationDelta,
		Enabled:            d.Cfg.AuthEnabled,
	}

	authHandler.RegisterRoutes(router)

	handlers := getAllHandlers(d)
	for _, handler := range handlers {
		handler.RegisterRoutes(router, authHandler)
	}

	return router
}

// Will add the dashboard front-end to the router.
// To update the front-end, please run `make build-ui`.
func (this *AdminApi) addDashboardRoutes(router *bone.Mux) *bone.Mux {
	// preparing static assets for embedded admin
	statikFS, err := fs.New()

	if err != nil {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Error("Failed to load statikFS")
	}

	indexHandler := func(w http.ResponseWriter, r *http.Request) {
		file, err := statikFS.Open("/index.html")
		if err != nil {
			w.WriteHeader(500)
			log.WithFields(log.Fields{
				"error": err,
			}).Error("got error while opening index file")
			return
		}
		io.Copy(w, file)
		w.WriteHeader(200)
	}

	router.HandleFunc("/dashboard", indexHandler)
	router.HandleFunc("/login", indexHandler)
	router.Handle("/", http.FileServer(statikFS))
	router.Handle("/*.js", http.FileServer(statikFS))
	router.Handle("/*.css", http.FileServer(statikFS))
	router.Handle("/*.ico", http.FileServer(statikFS))

	return router
}

func getAllHandlers(hoverfly *Hoverfly) []handlers.AdminHandler {
	list := []handlers.AdminHandler{
		&handlers.HealthHandler{},

		&v2.HoverflyHandler{Hoverfly: hoverfly},
		&v2.HoverflyDestinationHandler{Hoverfly: hoverfly},
		&v2.HoverflyModeHandler{Hoverfly: hoverfly},
		&v2.HoverflyMiddlewareHandler{Hoverfly: hoverfly},
		&v2.HoverflyUsageHandler{Hoverfly: hoverfly},
		&v2.HoverflyVersionHandler{Hoverfly: hoverfly},
		&v2.HoverflyUpstreamProxyHandler{Hoverfly: hoverfly},
		&v2.HoverflyPACHandler{Hoverfly: hoverfly},
		&v2.SimulationHandler{Hoverfly: hoverfly},
		&v2.CacheHandler{Hoverfly: hoverfly},
		&v2.LogsHandler{Hoverfly: hoverfly.StoreLogsHook},
		&v2.JournalHandler{Hoverfly: hoverfly.Journal},
		&v2.ShutdownHandler{},
		&v2.StateHandler{Hoverfly: hoverfly},
		&v2.DiffHandler{Hoverfly: hoverfly},
	}

	return list
}
