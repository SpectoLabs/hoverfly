package hoverfly

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/metrics"
	"net/http"
)

type HoverflyStats interface {
	GetStats() metrics.Stats
	GetCounter() metrics.CounterByMode
	GetRequestCache() cache.Cache
}

type StatsHandler struct {
	Hoverfly HoverflyStats
}

// StatsHandler - returns current stats about Hoverfly (request counts, record count)
func (this *StatsHandler) Get(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	stats := this.Hoverfly.GetStats()

	count, err := this.Hoverfly.GetRequestCache().RecordsCount()

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var sr StatsResponse
	sr.Stats = stats
	sr.RecordsCount = count

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
