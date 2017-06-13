package v2

import (
	"encoding/json"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type HoverflyJournal interface {
	GetEntries() []JournalEntryView
}

type JournalHandler struct {
	Hoverfly HoverflyJournal
}

func (this *JournalHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/v2/journal", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))
	mux.Options("/api/v2/journal", negroni.New(
		negroni.HandlerFunc(this.Options),
	))
}

func (this *JournalHandler) Get(response http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
	bytes, _ := json.Marshal(this.Hoverfly.GetEntries())
	handlers.WriteResponse(response, bytes)
}

func (this *JournalHandler) Options(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("Allow", "OPTIONS, GET")
	handlers.WriteResponse(w, []byte(""))
}
