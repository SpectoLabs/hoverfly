package v2

import (
	"encoding/json"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"strconv"
	"time"
)

const DefaultJournalLimit = 25

type HoverflyJournal interface {
	GetEntries(offset int, limit int, from *time.Time, to *time.Time, sort string) (JournalView, error)
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

	queryParams := request.URL.Query()
	offset, _ := strconv.Atoi(queryParams.Get("offset"))
	limit, _ := strconv.Atoi(queryParams.Get("limit"))
	fromTime := util.GetUnixTimeQueryParam(request, "from")
	toTime := util.GetUnixTimeQueryParam(request, "to")
	sort := queryParams.Get("sort")

	if limit == 0 {
		limit = DefaultJournalLimit
	}

	journalView, err := this.Hoverfly.GetEntries(offset, limit, fromTime, toTime, sort)
	if err != nil {
		handlers.WriteErrorResponse(response, err.Error(), http.StatusInternalServerError)
		return
	}

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
