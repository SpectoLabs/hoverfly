package v2

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type HoverflyLogs interface {
	GetLogsView() LogsView
	GetFilteredLogsView(int) LogsView
}

type LogsHandler struct {
	Hoverfly HoverflyLogs
}

func (this *LogsHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/v2/logs", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))
}

func (this *LogsHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	var logs LogsView

	queryParams := req.URL.Query()
	limitQuery, err := strconv.Atoi(queryParams.Get("limit"))
	if err == nil {
		logs = this.Hoverfly.GetFilteredLogsView(limitQuery)
	} else {
		logs = this.Hoverfly.GetLogsView()
	}

	bytes, _ := json.Marshal(logs)

	handlers.WriteResponse(w, bytes)
}
