package hoverfly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"time"

	// static assets
	_ "github.com/SpectoLabs/hoverfly/statik"
	"github.com/rakyll/statik/fs"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"github.com/gorilla/websocket"
	"github.com/meatballhat/negroni-logrus"

	// auth
	"github.com/SpectoLabs/hoverfly/authentication"
	"github.com/SpectoLabs/hoverfly/authentication/controllers"
	"github.com/SpectoLabs/hoverfly/metrics"
)

// recordedRequests struct encapsulates payload data
type recordedRequests struct {
	Data []Payload `json:"data"`
}

type storedMetadata struct {
	Data map[string]string `json:"data"`
}

type setMetadata struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type recordsCount struct {
	Count int `json:"count"`
}

type statsResponse struct {
	Stats        metrics.Stats `json:"stats"`
	RecordsCount int           `json:"recordsCount"`
}

type stateRequest struct {
	Mode        string `json:"mode"`
	Destination string `json:"destination"`
}

type messageResponse struct {
	Message string `json:"message"`
}

func (m *messageResponse) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(m)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// StartAdminInterface - starts admin interface web server
func (d *Hoverfly) StartAdminInterface() {

	// starting admin interface
	mux := getBoneRouter(*d)
	n := negroni.Classic()

	logLevel := log.ErrorLevel

	if d.Cfg.Verbose {
		logLevel = log.DebugLevel
	}

	n.Use(negronilogrus.NewCustomMiddleware(logLevel, &log.JSONFormatter{}, "admin"))
	n.UseHandler(mux)

	// admin interface starting message
	log.WithFields(log.Fields{
		"AdminPort": d.Cfg.AdminPort,
	}).Info("Admin interface is starting...")

	n.Run(fmt.Sprintf(":%s", d.Cfg.AdminPort))
}

// getBoneRouter returns mux for admin interface
func getBoneRouter(d Hoverfly) *bone.Mux {
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
		negroni.HandlerFunc(d.AllRecordsHandler),
	))

	mux.Delete("/api/records", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(d.DeleteAllRecordsHandler),
	))
	mux.Post("/api/records", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(d.ImportRecordsHandler),
	))

	mux.Get("/api/metadata", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(d.AllMetadataHandler),
	))

	mux.Put("/api/metadata", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(d.SetMetadataHandler),
	))

	mux.Delete("/api/metadata", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(d.DeleteMetadataHandler),
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
		negroni.HandlerFunc(d.CurrentStateHandler),
	))
	mux.Post("/api/state", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(d.StateHandler),
	))

	mux.Post("/api/add", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(d.ManualAddHandler),
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
		mux.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
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

// AllRecordsHandler returns JSON content type http response
func (d *Hoverfly) AllRecordsHandler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	records, err := d.RequestCache.GetAllValues()

	if err == nil {

		var payloads []Payload

		for _, v := range records {
			if payload, err := decodePayload(v); err == nil {
				payloads = append(payloads, *payload)
			} else {
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")

		var response recordedRequests
		response.Data = payloads
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

// RecordsCount returns number of captured requests as a JSON payload
func (d *Hoverfly) RecordsCount(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	count, err := d.RequestCache.RecordsCount()

	if err == nil {

		w.Header().Set("Content-Type", "application/json")

		var response recordsCount
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

	var sr statsResponse
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
				var sr statsResponse
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

// ImportRecordsHandler - accepts JSON payload and saves it to cache
func (d *Hoverfly) ImportRecordsHandler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	var requests recordedRequests

	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var response messageResponse

	if err != nil {
		// failed to read response body
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Could not read request body!")
		http.Error(w, "Failed to read request body.", 400)
		return
	}

	err = json.Unmarshal(body, &requests)

	if err != nil {
		w.WriteHeader(422) // can't process this entity
		return
	}

	err = d.ImportPayloads(requests.Data)

	if err != nil {
		response.Message = err.Error()
		w.WriteHeader(400)
	} else {
		response.Message = fmt.Sprintf("%d payloads import complete.", len(requests.Data))
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

	preq := RequestDetails{
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

	presp := ResponseDetails{
		Status:  sc,
		Headers: headers,
		Body:    respBody,
	}

	log.WithFields(log.Fields{
		"respBody":    respBody,
		"contentType": contentType,
	}).Info("manually adding request/response")

	p := Payload{Request: preq, Response: presp}

	var pls []Payload

	pls = append(pls, p)

	err = d.ImportPayloads(pls)

	w.Header().Set("Content-Type", "application/json")
	var response messageResponse

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

// DeleteAllRecordsHandler - deletes all captured requests
func (d *Hoverfly) DeleteAllRecordsHandler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	err := d.RequestCache.DeleteData()

	var en Entry
	en.ActionType = ActionTypeWipeDB
	en.Message = "wipe"
	en.Time = time.Now()

	if err := d.Hooks.Fire(ActionTypeWipeDB, &en); err != nil {
		log.WithFields(log.Fields{
			"error":      err.Error(),
			"message":    en.Message,
			"actionType": ActionTypeWipeDB,
		}).Error("failed to fire hook")
	}

	w.Header().Set("Content-Type", "application/json")

	var response messageResponse
	if err != nil {
		if err.Error() == "bucket not found" {
			response.Message = fmt.Sprintf("No records found")
			w.WriteHeader(200)
		} else {
			response.Message = fmt.Sprintf("Something went wrong: %s", err.Error())
			w.WriteHeader(500)
		}
	} else {
		response.Message = "Proxy cache deleted successfuly"
		w.WriteHeader(200)
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
	return
}

// CurrentStateHandler returns current state
func (d *Hoverfly) CurrentStateHandler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	var resp stateRequest
	resp.Mode = d.Cfg.GetMode()
	resp.Destination = d.Cfg.Destination

	b, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(b)
}

// StateHandler handles current proxy state
func (d *Hoverfly) StateHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var sr stateRequest

	// this is mainly for testing, since when you create
	if r.Body == nil {
		r.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("")))
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		// failed to read response body
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Could not read response body!")
		http.Error(w, "Failed to read request body.", 400)
		return
	}

	err = json.Unmarshal(body, &sr)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400) // can't process this entity
		return
	}

	availableModes := map[string]bool{
		"virtualize": true,
		"capture":    true,
		"modify":     true,
		"synthesize": true,
	}

	if sr.Mode != "" {
		if !availableModes[sr.Mode] {
			log.WithFields(log.Fields{
				"suppliedMode": sr.Mode,
			}).Error("Wrong mode found, can't change state")
			http.Error(w, "Bad mode supplied, available modes: virtualize, capture, modify, synthesize.", 400)
			return
		}
		log.WithFields(log.Fields{
			"newState":    sr.Mode,
			"body":        string(body),
			"destination": sr.Destination,
		}).Info("Handling state change request!")

		// setting new state
		d.Cfg.SetMode(sr.Mode)

	}

	// checking whether we should update destination
	if sr.Destination != "" {
		err := d.UpdateDestination(sr.Destination)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error while updating destination: %s", err.Error()), 500)
			return
		}
	}

	var en Entry
	en.ActionType = ActionTypeConfigurationChanged
	en.Message = "changed"
	en.Time = time.Now()
	en.Data = body

	if err := d.Hooks.Fire(ActionTypeConfigurationChanged, &en); err != nil {
		log.WithFields(log.Fields{
			"error":      err.Error(),
			"message":    en.Message,
			"actionType": ActionTypeConfigurationChanged,
		}).Error("failed to fire hook")
	}

	var resp stateRequest
	resp.Mode = d.Cfg.GetMode()
	resp.Destination = d.Cfg.Destination
	b, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(b)

}

// AllMetadataHandler returns JSON content type http response
func (d *Hoverfly) AllMetadataHandler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	entries, err := d.MetadataCache.GetAllEntries()

	metaData := make(map[string]string)

	for k, v := range entries {
		metaData[k] = string(v)
	}

	if err == nil {

		w.Header().Set("Content-Type", "application/json")

		var response storedMetadata
		response.Data = metaData
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
		}).Error("Failed to get metadata!")

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(500)
		return
	}
}

// SetMetadataHandler - sets new metadata
func (d *Hoverfly) SetMetadataHandler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	var sm setMetadata
	var mr messageResponse

	if req.Body == nil {
		req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("")))
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		// failed to read response body
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Could not read response body!")
		mr.Message = fmt.Sprintf("Failed to read request body. Error: %s", err.Error())
		w.WriteHeader(400)

		b, err := mr.Encode()
		if err != nil {
			// failed to read response body
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("Could not encode response body!")
			http.Error(w, "Failed to encode response", 500)
			return
		}
		w.Write(b)
		return
	}

	err = json.Unmarshal(body, &sm)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to unmarshal request body!")
		mr.Message = fmt.Sprintf("Failed to decode request body. Error: %s", err.Error())
		w.WriteHeader(400)

	} else if sm.Key == "" {
		mr.Message = "Key not provided."
		w.WriteHeader(400)

	} else {
		err = d.MetadataCache.Set([]byte(sm.Key), []byte(sm.Value))
		if err != nil {
			mr.Message = fmt.Sprintf("Failed to set metadata. Error: %s", err.Error())
			w.WriteHeader(500)
		} else {
			mr.Message = "Metadata set."
			w.WriteHeader(201)
		}
	}
	b, err := mr.Encode()
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

// DeleteMetadataHandler - deletes all metadata
func (d *Hoverfly) DeleteMetadataHandler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	err := d.MetadataCache.DeleteData()

	w.Header().Set("Content-Type", "application/json")

	var response messageResponse
	if err != nil {
		if err.Error() == "bucket not found" {
			response.Message = fmt.Sprintf("No metadata found.")
			w.WriteHeader(200)
		} else {
			response.Message = fmt.Sprintf("Something went wrong: %s", err.Error())
			w.WriteHeader(500)
		}
	} else {
		response.Message = "Metadata deleted successfuly"
		w.WriteHeader(200)
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
	return
}
