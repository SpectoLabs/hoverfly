package v2

import (
	"encoding/json"
	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"net/http"
)

type HoverflyTemplateDataSource interface {
	SetCsvDataSource(string, string) error
	DeleteDataSource(string)
	GetAllDataSources() TemplateDataSourceView
}

type HoverflyTemplateDataSourceHandler struct {
	Hoverfly HoverflyTemplateDataSource
}

func (templateDataSourceHandler HoverflyTemplateDataSourceHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {

	mux.Get("/api/v2/hoverfly/templating-data-source/csv", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(templateDataSourceHandler.Get),
	))

	mux.Put("/api/v2/hoverfly/templating-data-source/csv", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(templateDataSourceHandler.Put),
	))

	mux.Delete("/api/v2/hoverfly/templating-data-source/csv/:dataSourceName", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(templateDataSourceHandler.Delete),
	))
}

func (templateDataSourceHandler HoverflyTemplateDataSourceHandler) Put(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	var templateDataSourceRequest CSVDataSourceView
	err := handlers.ReadFromRequest(req, &templateDataSourceRequest)
	if err != nil {
		handlers.WriteErrorResponse(rw, err.Error(), 400)
		return
	}
	if err := templateDataSourceHandler.Hoverfly.SetCsvDataSource(templateDataSourceRequest.Name, templateDataSourceRequest.Data); err != nil {
		handlers.WriteErrorResponse(rw, err.Error(), 400)
		return
	}
	templateDataSourceHandler.Get(rw, req, next)
}

func (templateDataSourceHandler HoverflyTemplateDataSourceHandler) Delete(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	dataSourceName := bone.GetValue(req, "dataSourceName")
	templateDataSourceHandler.Hoverfly.DeleteDataSource(dataSourceName)
	templateDataSourceHandler.Get(rw, req, next)
}

func (templateDataSourceHandler HoverflyTemplateDataSourceHandler) Get(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	templateDataSourceView := templateDataSourceHandler.Hoverfly.GetAllDataSources()
	bytes, _ := json.Marshal(templateDataSourceView)

	handlers.WriteResponse(rw, bytes)
}
