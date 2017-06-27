package hoverfly

import (
	"fmt"
	"io"
	"net/http"

	// static assets
	_ "github.com/SpectoLabs/hoverfly/core/statik"
	"github.com/rakyll/statik/fs"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"

	handlers "github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
)

type AdminApi struct{}

// StartAdminInterface - starts admin interface web server
func (this *AdminApi) StartAdminInterface(hoverfly *Hoverfly) {

	// starting admin interface
	mux := this.getBoneRouter(hoverfly)
	n := negroni.New(negroni.NewRecovery(), negroni.NewStatic(http.Dir("public")))

	n.UseHandler(mux)

	// admin interface starting message
	log.WithFields(log.Fields{
		"AdminPort": hoverfly.Cfg.AdminPort,
	}).Info("Admin interface is starting...")

	http.ListenAndServe(fmt.Sprintf(":%s", hoverfly.Cfg.AdminPort), n)
}

// getBoneRouter returns mux for admin interface
func (this *AdminApi) getBoneRouter(d *Hoverfly) *bone.Mux {
	mux := bone.New()

	authHandler := &handlers.AuthHandler{
		d.Authentication,
		d.Cfg.SecretKey,
		d.Cfg.JWTExpirationDelta,
		d.Cfg.AuthEnabled,
	}

	authHandler.RegisterRoutes(mux)

	handlers := GetAllHandlers(d)
	for _, handler := range handlers {
		handler.RegisterRoutes(mux, authHandler)
	}

	// preparing static assets for embedded admin
	statikFS, err := fs.New()

	if err != nil {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Error("Failed to load statikFS, admin UI might not work :(")
	}

	indexHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("index.html")
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

	mux.HandleFunc("/dashboard", indexHandler)
	mux.HandleFunc("/login", indexHandler)
	mux.Handle("/", http.FileServer(statikFS))
	mux.Handle("/*.js", http.FileServer(statikFS))
	mux.Handle("/*.css", http.FileServer(statikFS))
	mux.Handle("/*.ico", http.FileServer(statikFS))

	return mux
}

func GetAllHandlers(hoverfly *Hoverfly) []handlers.AdminHandler {
	var list []handlers.AdminHandler

	list = append(list, &v1.HealthHandler{})
	list = append(list, &v1.MetadataHandler{Hoverfly: hoverfly})
	list = append(list, &v2.HoverflyHandler{Hoverfly: hoverfly})
	list = append(list, &v2.HoverflyDestinationHandler{Hoverfly: hoverfly})
	list = append(list, &v2.HoverflyModeHandler{Hoverfly: hoverfly})
	list = append(list, &v2.HoverflyMiddlewareHandler{Hoverfly: hoverfly})
	list = append(list, &v2.HoverflyUsageHandler{Hoverfly: hoverfly})
	list = append(list, &v2.HoverflyVersionHandler{Hoverfly: hoverfly})
	list = append(list, &v2.HoverflyUpstreamProxyHandler{Hoverfly: hoverfly})
	list = append(list, &v2.SimulationHandler{Hoverfly: hoverfly})
	list = append(list, &v2.CacheHandler{Hoverfly: hoverfly})
	list = append(list, &v2.LogsHandler{Hoverfly: hoverfly.StoreLogsHook})
	list = append(list, &v2.JournalHandler{Hoverfly: hoverfly.Journal})
	list = append(list, &v2.ShutdownHandler{})

	return list
}
