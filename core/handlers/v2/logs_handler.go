package v2

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"time"
)

type HoverflyLogs interface {
	GetLogs(limit int, from *time.Time) []*logrus.Entry
}

type LogsHandler struct {
	Hoverfly HoverflyLogs
}

var DefaultLimit = 500

func (this *LogsHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/v2/logs", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))
	mux.Options("/api/v2/logs", negroni.New(
		negroni.HandlerFunc(this.Options),
	))

	mux.Get("/api/v2/ws/logs", http.HandlerFunc(this.GetWS))
}

func (this *LogsHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	var logs []*logrus.Entry

	queryParams := req.URL.Query()
	limitQuery, _ := strconv.Atoi(queryParams.Get("limit"))
	if limitQuery == 0 {
		limitQuery = DefaultLimit
	}

	fromQuery, _ := strconv.Atoi(queryParams.Get("from"))

	var fromTime *time.Time
	if fromQuery != 0 {

		fromTimeValue := time.Unix(int64(fromQuery), 0)
		fromTime = &fromTimeValue
	}

	logs = this.Hoverfly.GetLogs(limitQuery, fromTime)

	if strings.Contains(req.Header.Get("Accept"), "text/plain") ||
		strings.Contains(req.Header.Get("Content-Type"), "text/plain") {
		handlers.WriteResponse(w, []byte(logsToPlainText(logs)))
	} else {
		bytes, _ := json.Marshal(logsToLogsView(logs))
		handlers.WriteResponse(w, bytes)
	}
}

func (this *LogsHandler) Options(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("Allow", "OPTIONS, GET")
	handlers.WriteResponse(w, []byte(""))
}

func logsToLogsView(logs []*logrus.Entry) LogsView {
	var logInterfaces []map[string]interface{}
	for _, entry := range logs {
		data := make(map[string]interface{}, len(entry.Data)+3)

		for k, v := range entry.Data {
			data[k] = v
		}

		data["time"] = entry.Time.Format(logrus.DefaultTimestampFormat)
		data["msg"] = entry.Message
		data["level"] = entry.Level.String()

		logInterfaces = append(logInterfaces, data)
	}

	return LogsView{
		Logs: logInterfaces,
	}
}

func logsToPlainText(logs []*logrus.Entry) string {

	var buffer bytes.Buffer
	for _, entry := range logs {
		entry.Logger = logrus.New()
		entry.Logger.Formatter = &logrus.TextFormatter{
			ForceColors:      true,
			DisableTimestamp: false,
			FullTimestamp:    true,
		}

		log, err := entry.String()
		if err == nil {
			buffer.WriteString(log)
		}
	}

	return buffer.String()
}

func (this *LogsHandler) GetWS(w http.ResponseWriter, r *http.Request) {

	var previousLogs LogsView

	handlers.NewWebsocket(func() ([]byte, error) {
		currentLogs := logsToLogsView(this.Hoverfly.GetLogs(500, nil))

		if !reflect.DeepEqual(currentLogs, previousLogs) {
			previousLogs = currentLogs
			return json.Marshal(currentLogs)
		}

		return nil, errors.New("No update needed")
	}, w, r)
}
