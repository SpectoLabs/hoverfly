package v2

import (
	"encoding/json"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type HoverflyJournal interface {
	GetEntries() ([]JournalEntryView, error)
	GetFilteredEntries(journalEntryFilterView JournalEntryFilterView) ([]JournalEntryView, error)
	DeleteEntries() error
}

type JournalHandler struct {
	Hoverfly HoverflyJournal
}

func (this *JournalHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/v2/journal", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))
	mux.Post("/api/v2/journal", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Post),
	))
	mux.Delete("/api/v2/journal", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Delete),
	))
	mux.Options("/api/v2/journal", negroni.New(
		negroni.HandlerFunc(this.Options),
	))
}

func (this *JournalHandler) Get(response http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
	var journalView JournalView

	entries, err := this.Hoverfly.GetEntries()
	if err != nil {
		handlers.WriteErrorResponse(response, err.Error(), http.StatusInternalServerError)
		return
	}

	journalView.Journal = entries

	bytes, _ := json.Marshal(journalView)
	handlers.WriteResponse(response, bytes)
}

func (this *JournalHandler) Post(response http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
	var journalView JournalView

	var journalEntryFilterView JournalEntryFilterView

	err := handlers.ReadFromRequest(request, &journalEntryFilterView)
	if err != nil {
		handlers.WriteErrorResponse(response, err.Error(), http.StatusBadRequest)
		return
	} else if journalEntryFilterView.Request == nil {
		handlers.WriteErrorResponse(response, "No \"request\" object in search parameters", http.StatusBadRequest)
		return
	}

	entries, err := this.Hoverfly.GetFilteredEntries(journalEntryFilterView)
	if err != nil {
		handlers.WriteErrorResponse(response, err.Error(), http.StatusInternalServerError)
		return
	}

	journalView.Journal = entries

	bytes, _ := json.Marshal(journalView)
	handlers.WriteResponse(response, bytes)
}

func (this *JournalHandler) Delete(response http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
	err := this.Hoverfly.DeleteEntries()
	if err != nil {
		handlers.WriteErrorResponse(response, err.Error(), http.StatusInternalServerError)
		return
	}

	this.Get(response, request, next)
}

func (this *JournalHandler) Options(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("Allow", "OPTIONS, GET, DELETE, POST")
	handlers.WriteResponse(w, []byte(""))
}
