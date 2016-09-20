package hoverfly

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"time"

	// static assets
	_ "github.com/SpectoLabs/hoverfly/core/statik"
	"github.com/rakyll/statik/fs"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"github.com/gorilla/websocket"
	"github.com/meatballhat/negroni-logrus"

	// auth
	"github.com/SpectoLabs/hoverfly/core/authentication"
	"github.com/SpectoLabs/hoverfly/core/authentication/controllers"
	handlers "github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/views"
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
	middlewareHandler := handlers.MiddlewareHandler{Hoverfly: d}
	recordsHandler := handlers.RecordsHandler{Hoverfly: d}
	templatesHandler := handlers.TemplatesHandler{Hoverfly: d}
	metadataHandler := handlers.MetadataHandler{Hoverfly: d}
	stateHandler := handlers.StateHandler{Hoverfly: d}
	delaysHandler := handlers.DelaysHandler{Hoverfly: d}

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

	mux.Get("/api/records", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(recordsHandler.Get),
	))

	mux.Delete("/api/records", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(recordsHandler.Delete),
	))

	mux.Post("/api/records", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(recordsHandler.Post),
	))

	mux.Get("/api/templates", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(templatesHandler.Get),
	))

	mux.Delete("/api/templates", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(templatesHandler.Delete),
	))

	mux.Post("/api/templates", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(templatesHandler.Post),
	))

	mux.Get("/api/metadata", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(metadataHandler.Get),
	))

	mux.Put("/api/metadata", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(metadataHandler.Put),
	))

	mux.Delete("/api/metadata", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(metadataHandler.Delete),
	))

	mux.Get("/api/count", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(d.RecordsCount),
	))
	mux.Get("/api/stats", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(d.StatsHandler),
	))
	// TODO: check auth for websocket connection
	mux.Get("/api/statsws", http.HandlerFunc(d.StatsWSHandler))

	mux.Get("/api/state", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(stateHandler.Get),
	))
	mux.Post("/api/state", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(stateHandler.Post),
	))

	mux.Get("/api/middleware", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(middlewareHandler.Get),
	))

	mux.Post("/api/middleware", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(middlewareHandler.Post),
	))

	mux.Post("/api/add", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(d.ManualAddHandler),
	))

	mux.Get("/api/health", negroni.New(
		negroni.HandlerFunc(healthHandler.Get),
	))

	mux.Get("/api/delays", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(delaysHandler.Get),
	))

	mux.Put("/api/delays", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(delaysHandler.Put),
	))

	mux.Delete("/api/delays", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(delaysHandler.Delete),
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

// RecordsCount returns number of captured requests as a JSON payload
func (d *Hoverfly) RecordsCount(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	count, err := d.RequestCache.RecordsCount()

	if err == nil {

		w.Header().Set("Content-Type", "application/json")

		var response handlers.RecordsCount
		response.Count = count
		b, err := json.Marshal(response)

		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Write(b)
			return
		}
	} else {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Error("Failed to get data from cache!")

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(500) // can't process this entity
		return
	}
}

// StatsHandler - returns current stats about Hoverfly (request counts, record count)
func (d *Hoverfly) StatsHandler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	stats := d.Counter.Flush()

	count, err := d.RequestCache.RecordsCount()

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var sr handlers.StatsResponse
	sr.Stats = stats
	sr.RecordsCount = count

	w.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(sr)

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Write(b)
		return
	}

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// StatsWSHandler - returns current stats about Hoverfly (request counts, record count) through the websocket
func (d *Hoverfly) StatsWSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to upgrade websocket")
		return
	}

	// defining counters for delta check
	var recordsCount int
	var statsCounters map[string]int64

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		log.WithFields(log.Fields{
			"message": string(p),
		}).Info("Got message...")

		for _ = range time.Tick(1 * time.Second) {
			count, err := d.RequestCache.RecordsCount()
			if err != nil {
				log.WithFields(log.Fields{
					"message": p,
					"error":   err.Error(),
				}).Error("got error while trying to get records count")
				continue
			}
			stats := d.Counter.Flush()

			// checking whether we should send an update
			if !reflect.DeepEqual(stats.Counters, statsCounters) || count != recordsCount {
				var sr handlers.StatsResponse
				sr.Stats = stats
				sr.RecordsCount = count

				b, err := json.Marshal(sr)

				if err = conn.WriteMessage(messageType, b); err != nil {
					log.WithFields(log.Fields{
						"message": p,
						"error":   err.Error(),
					}).Debug("Got error when writing message...")
					continue
				}
				recordsCount = count
				statsCounters = stats.Counters
			}
		}
	}
}

// ManualAddHandler - manually add new request/responses, using a form
func (d *Hoverfly) ManualAddHandler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	err := req.ParseForm()

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Got error while parsing form")
	}

	// request details
	destination := req.PostFormValue("inputDestination")
	method := req.PostFormValue("inputMethod")
	path := req.PostFormValue("inputPath")
	query := req.PostFormValue("inputQuery")
	reqBody := req.PostFormValue("inputRequestBody")

	preq := models.RequestDetails{
		Destination: destination,
		Method:      method,
		Path:        path,
		Query:       query,
		Body:        reqBody}

	// response
	respStatusCode := req.PostFormValue("inputResponseStatusCode")
	respBody := req.PostFormValue("inputResponseBody")
	contentType := req.PostFormValue("inputContentType")

	headers := make(map[string][]string)

	// getting content type
	if contentType == "xml" {
		headers["Content-Type"] = []string{"application/xml"}
	} else if contentType == "json" {
		headers["Content-Type"] = []string{"application/json"}
	} else {
		headers["Content-Type"] = []string{"text/html"}
	}

	sc, _ := strconv.Atoi(respStatusCode)

	presp := models.ResponseDetails{
		Status:  sc,
		Headers: headers,
		Body:    respBody,
	}

	log.WithFields(log.Fields{
		"respBody":    respBody,
		"contentType": contentType,
	}).Info("manually adding request/response")

	p := models.RequestResponsePair{Request: preq, Response: presp}

	var pairViews []views.RequestResponsePairView

	pairViews = append(pairViews, *p.ConvertToRequestResponsePairView())

	err = d.ImportRequestResponsePairViews(pairViews)

	w.Header().Set("Content-Type", "application/json")
	var response handlers.MessageResponse

	if err != nil {
		response.Message = fmt.Sprintf("Got error: %s", err.Error())
		w.WriteHeader(400)

	} else {
		// redirecting to home
		response.Message = "Record added successfuly"
		w.WriteHeader(201)
	}
	b, err := response.Encode()
	if err != nil {
		// failed to read response body
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Could not encode response body!")
		http.Error(w, "Failed to encode response", 500)
		return
	}
	w.Write(b)

}