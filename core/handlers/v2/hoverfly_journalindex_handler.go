package v2

import (
	"encoding/json"
	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"net/http"
)

type HoverflyJournalIndex interface {
	GetAllIndexes() []JournalIndexView
	AddJournalIndex(string) error
	DeleteJournalIndex(string)
}

type HoverflyJournalIndexHandler struct {
	Hoverfly HoverflyJournalIndex
}

func (hoverflyJournalIndexHandler HoverflyJournalIndexHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {

	mux.Get("/api/v2/journal/index", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(hoverflyJournalIndexHandler.GetAll),
	))

	mux.Post("/api/v2/journal/index", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(hoverflyJournalIndexHandler.Post),
	))

	mux.Delete("/api/v2/journal/index/:indexName", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(hoverflyJournalIndexHandler.Delete),
	))
}

func (hoverflyJournalIndexHandler HoverflyJournalIndexHandler) Post(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var journalIndexRequestView JournalIndexRequestView
	err := handlers.ReadFromRequest(r, &journalIndexRequestView)
	if err != nil {
		handlers.WriteErrorResponse(rw, err.Error(), 400)
		return
	}
	err = hoverflyJournalIndexHandler.Hoverfly.AddJournalIndex(journalIndexRequestView.Name)
	if err != nil {
		handlers.WriteErrorResponse(rw, err.Error(), 400)
		return
	}
	hoverflyJournalIndexHandler.GetAll(rw, r, next)
}

func (hoverflyJournalIndexHandler HoverflyJournalIndexHandler) Delete(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	indexName := bone.GetValue(r, "indexName")
	hoverflyJournalIndexHandler.Hoverfly.DeleteJournalIndex(indexName)
	hoverflyJournalIndexHandler.GetAll(rw, r, next)
}

func (hoverflyJournalIndexHandler HoverflyJournalIndexHandler) GetAll(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	journalIndexViews := hoverflyJournalIndexHandler.Hoverfly.GetAllIndexes()
	bytes, _ := json.Marshal(journalIndexViews)

	handlers.WriteResponse(rw, bytes)
}
