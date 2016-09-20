package hoverfly

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/authentication"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"io/ioutil"
	"net/http"
)

type HoverflyMiddleware interface {
	GetMiddleware() string
	SetMiddleware(string) error
}

type MiddlewareHandler struct {
	Hoverfly HoverflyMiddleware
}

func (this *MiddlewareHandler) RegisterRoutes(mux *bone.Mux, am *authentication.AuthMiddleware) {
	mux.Get("/api/middleware", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))

	mux.Post("/api/middleware", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Post),
	))
}

func (this *MiddlewareHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	var resp MiddlewareSchema

	resp.Middleware = this.Hoverfly.GetMiddleware()

	jsonResp, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
}

func (this *MiddlewareHandler) Post(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// failed to read response body
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Could not read response body!")
		http.Error(w, "Failed to read request body.", 400)
		return
	}

	var middlewareReq MiddlewareSchema

	err = json.Unmarshal(body, &middlewareReq)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Could not deserialize middleware")
		http.Error(w, "Unable to deserialize request body.", 400)
		return
	}

	err = this.Hoverfly.SetMiddleware(middlewareReq.Middleware)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Could not execute middleware")
		http.Error(w, "Invalid middleware: "+err.Error(), 400)
		return
	}

	this.Get(w, req, next)
}
