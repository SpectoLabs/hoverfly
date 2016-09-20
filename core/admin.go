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

	// auth
	"github.com/SpectoLabs/hoverfly/core/authentication"
	"github.com/SpectoLabs/hoverfly/core/authentication/controllers"
	handlers "github.com/SpectoLabs/hoverfly/core/handlers"
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
	// getting auth controllers and middleware
	ac := controllers.GetNewAuthenticationController(
		d.Authentication,
		d.Cfg.SecretKey,
		d.Cfg.JWTExpirationDelta,
		d.Cfg.AuthEnabled)

	am := authentication.GetNewAuthenticationMiddleware(
		d.Authentication,
		d.Cfg.SecretKey,
		d.Cfg.JWTExpirationDelta,
		d.Cfg.AuthEnabled)

	healthHandler := handlers.HealthHandler{}
	healthHandler.RegisterRoutes(mux, am)

	middlewareHandler := handlers.MiddlewareHandler{Hoverfly: d}
	middlewareHandler.RegisterRoutes(mux, am)

	recordsHandler := handlers.RecordsHandler{Hoverfly: d}
	recordsHandler.RegisterRoutes(mux, am)

	templatesHandler := handlers.TemplatesHandler{Hoverfly: d}
	templatesHandler.RegisterRoutes(mux, am)

	metadataHandler := handlers.MetadataHandler{Hoverfly: d}
	metadataHandler.RegisterRoutes(mux, am)

	stateHandler := handlers.StateHandler{Hoverfly: d}
	stateHandler.RegisterRoutes(mux, am)

	delaysHandler := handlers.DelaysHandler{Hoverfly: d}
	delaysHandler.RegisterRoutes(mux, am)

	addHandler := handlers.AddHandler{Hoverfly: d}
	addHandler.RegisterRoutes(mux, am)

	countHandler := handlers.CountHandler{Hoverfly: d}
	countHandler.RegisterRoutes(mux, am)

	statsHandler := handlers.StatsHandler{Hoverfly: d}
	statsHandler.RegisterRoutes(mux, am)

	mux.Post("/api/token-auth", http.HandlerFunc(ac.Login))

	mux.Get("/api/refresh-token-auth", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(ac.RefreshToken),
	))
	mux.Get("/api/logout", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(ac.Logout),
	))

	mux.Get("/api/users", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(ac.GetAllUsersHandler),
	))

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
