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
	"github.com/meatballhat/negroni-logrus"

	handlers "github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
)

type AdminApi struct{}

// StartAdminInterface - starts admin interface web server
func (this *AdminApi) StartAdminInterface(hoverfly *Hoverfly) {

	// starting admin interface
	mux := this.getBoneRouter(hoverfly)
	n := negroni.Classic()

	logLevel := log.ErrorLevel

	if hoverfly.Cfg.Verbose {
		logLevel = log.DebugLevel
	}

	n.Use(negronilogrus.NewCustomMiddleware(logLevel, &log.JSONFormatter{}, "admin"))
	n.UseHandler(mux)

	// admin interface starting message
	log.WithFields(log.Fields{
		"AdminPort": hoverfly.Cfg.AdminPort,
	}).Info("Admin interface is starting...")

	n.Run(fmt.Sprintf(":%s", hoverfly.Cfg.AdminPort))
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

	if d.Cfg.Development {
		// since hoverfly is not started from cmd/hoverfly/hoverfly
		// we have to target to that directory
		log.Warn("Hoverfly is serving files from /static/admin/dist instead of statik binary!")
		mux.Handle("/js/*", http.StripPrefix("/js/", http.FileServer(http.Dir("../../static/admin/dist/js"))))

		mux.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../../static/admin/dist/index.html")
		})

	} else {
		// preparing static assets for embedded admin
		statikFS, err := fs.New()

		if err != nil {
			log.WithFields(log.Fields{
				"Error": err.Error(),
			}).Error("Failed to load statikFS, admin UI might not work :(")
		}
		mux.Handle("/js/*", http.FileServer(statikFS))
		mux.Handle("/app.32dc9945fd902da8ed2cccdc8703129f.css", http.FileServer(statikFS))
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
		})
	}
	return mux
}

func GetAllHandlers(hoverfly *Hoverfly) []handlers.AdminHandler {
	var list []handlers.AdminHandler

	list = append(list, &v1.CountHandler{Hoverfly: hoverfly})
	list = append(list, &v1.DelaysHandler{Hoverfly: hoverfly})
	list = append(list, &v1.HealthHandler{})
	list = append(list, &v1.MetadataHandler{Hoverfly: hoverfly})
	list = append(list, &v1.RecordsHandler{Hoverfly: hoverfly})
	list = append(list, &v1.StateHandler{Hoverfly: hoverfly})
	list = append(list, &v1.StatsHandler{Hoverfly: hoverfly})
	list = append(list, &v2.HoverflyHandler{Hoverfly: hoverfly})
	list = append(list, &v2.HoverflyDestinationHandler{Hoverfly: hoverfly})
	list = append(list, &v2.HoverflyModeHandler{Hoverfly: hoverfly})
	list = append(list, &v2.HoverflyMiddlewareHandler{Hoverfly: hoverfly})
	list = append(list, &v2.HoverflyUsageHandler{Hoverfly: hoverfly})
	list = append(list, &v2.HoverflyVersionHandler{Hoverfly: hoverfly})
	list = append(list, &v2.HoverflyUpstreamProxyHandler{Hoverfly: hoverfly})
	list = append(list, &v2.SimulationHandler{Hoverfly: hoverfly})

	return list
}
