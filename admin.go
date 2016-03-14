package hoverfly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	Stats        Stats `json:"stats"`
	RecordsCount int   `json:"recordsCount"`
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
func (d *DBClient) StartAdminInterface() {

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
func getBoneRouter(d DBClient) *bone.Mux {
	mux := bone.New()

	// getting auth controllers and middleware
	ac := controllers.GetNewAuthenticationController(d.AB, d.Cfg.SecretKey, d.Cfg.JWTExpirationDelta)
	am := authentication.GetNewAuthenticationMiddleware(d.AB,
		d.Cfg.SecretKey,
		d.Cfg.JWTExpirationDelta,
		d.Cfg.AuthEnabled)

	mux.Post("/token-auth", http.HandlerFunc(ac.Login))
	mux.Get("/refresh-token-auth", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(ac.RefreshToken),
	))
	mux.Get("/logout", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(ac.Logout),
	))

	mux.Get("/users", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(ac.GetAllUsersHandler),
	))

	mux.Get("/records", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(d.AllRecordsHandler),
	))

	mux.Delete("/records", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(d.DeleteAllRecordsHandler),
	))
	mux.Post("/records", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(d.ImportRecordsHandler),
	))

	mux.Get("/metadata", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(d.AllMetadataHandler),
	))

	mux.Put("/metadata", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(d.SetMetadataHandler),
	))

	mux.Delete("/metadata", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(d.DeleteMetadataHandler),
	))

	mux.Get("/count", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(d.RecordsCount),
	))
	mux.Get("/stats", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(d.StatsHandler),
	))
	// TODO: check auth for websocket connection
	mux.Get("/statsws", http.HandlerFunc(d.StatsWSHandler))

	mux.Get("/state", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(d.CurrentStateHandler),
	))
	mux.Post("/state", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(d.StateHandler),
	))

	mux.Post("/add", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(d.ManualAddHandler),
	))

	if d.Cfg.Development {
		// since hoverfly is not started from cmd/hoverfly/hoverfly
		// we have to target to that directory
		log.Warn("Hoverfly is serving files from /static/dist instead of statik binary!")
		mux.Handle("/*", http.FileServer(http.Dir("../../static/dist")))
	} else {
		// preparing static assets for embedded admin
		statikFS, err := fs.New()

		if err != nil {
			log.WithFields(log.Fields{
				"Error": err.Error(),
			}).Error("Failed to load statikFS, admin UI might not work :(")
		}

		mux.Handle("/*", http.FileServer(statikFS))
	}
	return mux
}

// AllRecordsHandler returns JSON content type http response
func (d *DBClient) AllRecordsHandler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	records, err := d.Cache.GetAllRequests()

	if err == nil {

		w.Header().Set("Content-Type", "application/json")

		var response recordedRequests
		response.Data = records
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
func (d *DBClient) RecordsCount(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	count, err := d.Cache.RecordsCount()

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
func (d *DBClient) StatsHandler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	stats := d.Counter.Flush()

	count, err := d.Cache.RecordsCount()

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
func (d *DBClient) StatsWSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		log.WithFields(log.Fields{
			"message": string(p),
		}).Info("Got message...")

		for _ = range time.Tick(1 * time.Second) {

			count, err := d.Cache.RecordsCount()

			if err != nil {
				log.WithFields(log.Fields{
					"message": p,
					"error":   err.Error(),
				}).Error("got error while trying to get records count")
				return
			}

			stats := d.Counter.Flush()

			var sr statsResponse
			sr.Stats = stats
			sr.RecordsCount = count

			b, err := json.Marshal(sr)

			if err = conn.WriteMessage(messageType, b); err != nil {
				log.WithFields(log.Fields{
					"message": p,
					"error":   err.Error(),
				}).Debug("Got error when writing message...")
				return
			}
		}

	}

}

// ImportRecordsHandler - accepts JSON payload and saves it to cache
func (d *DBClient) ImportRecordsHandler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

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
func (d *DBClient) ManualAddHandler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
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
func (d *DBClient) DeleteAllRecordsHandler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	err := d.Cache.DeleteData()

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
func (d *DBClient) CurrentStateHandler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	var resp stateRequest
	resp.Mode = d.Cfg.GetMode()
	resp.Destination = d.Cfg.Destination

	b, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(b)
}

// StateHandler handles current proxy state
func (d *DBClient) StateHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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
		} else {
			log.WithFields(log.Fields{
				"newState":    sr.Mode,
				"body":        string(body),
				"destination": sr.Destination,
			}).Info("Handling state change request!")

			// setting new state
			d.Cfg.SetMode(sr.Mode)
		}
	}

	// checking whether we should update destination
	if sr.Destination != "" {
		d.UpdateDestination(sr.Destination)
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
func (d *DBClient) AllMetadataHandler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	metadata, err := d.MD.GetAll()

	if err == nil {

		w.Header().Set("Content-Type", "application/json")

		var response storedMetadata
		respMap := make(map[string]string)
		for _, v := range metadata {
			respMap[v.Key] = v.Value
		}
		response.Data = respMap
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
func (d *DBClient) SetMetadataHandler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
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
		err = d.MD.Set([]byte(sm.Key), []byte(sm.Value))
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
func (d *DBClient) DeleteMetadataHandler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	err := d.MD.DeleteData()

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
