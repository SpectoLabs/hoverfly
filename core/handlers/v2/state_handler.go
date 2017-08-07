package v2

import (
	"encoding/json"
	"net/http"

	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type StateHandler struct {
	Hoverfly Hoverfly
}

func (this *StateHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/v2/state", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))
	mux.Delete("/api/v2/state", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Delete),
	))
	mux.Put("/api/v2/state", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Put),
	))
	mux.Patch("/api/v2/state", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Patch),
	))
	mux.Options("/api/v2/state", negroni.New(
		negroni.HandlerFunc(this.Options),
	))
}

func (this *StateHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	marshal, err := json.Marshal(StateView{
		State: this.Hoverfly.GetState(),
	})

	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handlers.WriteResponse(w, marshal)
	w.WriteHeader(http.StatusOK)
}

func (this *StateHandler) Delete(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	this.Hoverfly.ClearState()
	w.WriteHeader(http.StatusOK)
}

func (this *StateHandler) Put(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	toPut := &StateView{}

	responseBody, err := ioutil.ReadAll(req.Body)

	err = json.Unmarshal(responseBody, &toPut)

	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	this.Hoverfly.SetState(toPut.State)

	marshal, _ := json.Marshal(StateView{
		State: this.Hoverfly.GetState(),
	})

	handlers.WriteResponse(w, marshal)
	w.WriteHeader(http.StatusOK)
}

func (this *StateHandler) Patch(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	toPatch := &StateView{}

	body, err := ioutil.ReadAll(req.Body)

	err = json.Unmarshal(body, &toPatch)

	if err != nil {
		handlers.WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	this.Hoverfly.PatchState(toPatch.State)

	marshal, _ := json.Marshal(StateView{
		State: this.Hoverfly.GetState(),
	})

	handlers.WriteResponse(w, marshal)

	w.WriteHeader(http.StatusOK)
}

func (this *StateHandler) Options(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("Allow", "OPTIONS, GET, DELETE, PUT, PATCH")
	handlers.WriteResponse(w, []byte(""))
}
