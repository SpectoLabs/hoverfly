package v2

import (
	"encoding/json"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"io/ioutil"
)

type CurrentStateHandler struct {
	Hoverfly Hoverfly
}

func (this *CurrentStateHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/v2/hoverfly/current-state", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))
	mux.Delete("/api/v2/hoverfly/current-state", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Delete),
	))
	mux.Put("/api/v2/hoverfly/current-state", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Put),
	))
	mux.Patch("/api/v2/hoverfly/current-state", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Patch),
	))
	mux.Options("/api/v2/hoverfly/current-state", negroni.New(
		negroni.HandlerFunc(this.Options),
	))
}

func (this *CurrentStateHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	marshal, err := json.Marshal(this.Hoverfly.GetState())

	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handlers.WriteResponse(w, marshal)
	w.WriteHeader(http.StatusOK)
}

func (this *CurrentStateHandler) Delete(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	this.Hoverfly.ClearState()
	w.WriteHeader(http.StatusOK)
}

func (this *CurrentStateHandler) Put(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	var toPut map[string]string

	body, err := ioutil.ReadAll(req.Body)

	err = json.Unmarshal(body, &toPut)

	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	this.Hoverfly.SetState(toPut)

	marshal, _ := json.Marshal(this.Hoverfly.GetState())

	handlers.WriteResponse(w, marshal)
	w.WriteHeader(http.StatusOK)
}

func (this *CurrentStateHandler) Patch(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	var toPatch map[string]string

	body, err := ioutil.ReadAll(req.Body)

	err = json.Unmarshal(body, &toPatch)

	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	this.Hoverfly.PatchState(toPatch)

	marshal, _ := json.Marshal(this.Hoverfly.GetState())

	handlers.WriteResponse(w, marshal)

	w.WriteHeader(http.StatusOK)
}

func (this *CurrentStateHandler) Options(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("Allow", "OPTIONS, GET, DELETE, PUT, PATCH")
	handlers.WriteResponse(w, []byte(""))
}
