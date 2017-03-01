package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/SpectoLabs/hoverfly/core/metrics"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

type HoverflyStats interface {
	GetStats() metrics.Stats
	GetSimulationPairsCount() int
}

type StatsHandler struct {
	Hoverfly HoverflyStats
}

func (this *StatsHandler) RegisterRoutes(mux *bone.Mux, am *handlers.AuthHandler) {
	mux.Get("/api/stats", negroni.New(
		negroni.HandlerFunc(am.RequireTokenAuthentication),
		negroni.HandlerFunc(this.Get),
	))
	// TODO: check auth for websocket connection
	mux.Get("/api/statsws", http.HandlerFunc(this.GetWS))
}

// StatsHandler - returns current stats about Hoverfly (request counts, record count)
func (this *StatsHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	stats := this.Hoverfly.GetStats()

	simulationCount := this.Hoverfly.GetSimulationPairsCount()

	var sr StatsResponse
	sr.Stats = stats
	sr.RecordsCount = simulationCount

	w.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(sr)

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Write(b)
		return
	}

}

// StatsWSHandler - returns current stats about Hoverfly (request counts, record count) through the websocket
func (this *StatsHandler) GetWS(w http.ResponseWriter, r *http.Request) {

	// defining counters for delta check
	var recordsCount int
	var statsCounters map[string]int64

	handlers.NewWebsocket(func() ([]byte, error) {
		count := this.Hoverfly.GetSimulationPairsCount()
		stats := this.Hoverfly.GetStats()

		if !reflect.DeepEqual(stats.Counters, statsCounters) || count != recordsCount {
			var sr StatsResponse
			sr.Stats = stats
			sr.RecordsCount = count

			recordsCount = count
			statsCounters = stats.Counters

			return json.Marshal(sr)
		}

		return nil, errors.New("No update needed")
	}, w, r)
}
