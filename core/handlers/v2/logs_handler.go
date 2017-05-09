package v2

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type HoverflyLogs interface {
	GetLogsView() LogsView
	GetFilteredLogsView(int) LogsView
	GetLogs() string
	GetFilteredLogs(int) string
}

type LogsHandler struct {
	Hoverfly HoverflyLogs
}

func (this *LogsHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/v2/logs", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))

	mux.Get("/api/v2/ws/logs", http.HandlerFunc(this.GetWS))
}

func (this *LogsHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	var logs LogsView

	queryParams := req.URL.Query()
	limitQuery, err := strconv.Atoi(queryParams.Get("limit"))

	if strings.Contains(req.Header.Get("Content-Type"), "text/plain") {
		var logs string
		if err == nil {
			logs = this.Hoverfly.GetFilteredLogs(limitQuery)
		} else {
			logs = this.Hoverfly.GetLogs()
		}

		handlers.WriteResponse(w, []byte(logs))
	} else {
		if err == nil {
			logs = this.Hoverfly.GetFilteredLogsView(limitQuery)
		} else {
			logs = this.Hoverfly.GetLogsView()
		}

		bytes, _ := json.Marshal(logs)

		handlers.WriteResponse(w, bytes)
	}
}

func (this *LogsHandler) GetWS(w http.ResponseWriter, r *http.Request) {

	var previousLogs LogsView

	handlers.NewWebsocket(func() ([]byte, error) {
		currentLogs := this.Hoverfly.GetLogsView()

		if !reflect.DeepEqual(currentLogs, previousLogs) {
			previousLogs = currentLogs
			return json.Marshal(currentLogs)
		}

		return nil, errors.New("No update needed")
	}, w, r)
}
