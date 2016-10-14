package v1

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/SpectoLabs/hoverfly/core/metrics"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"github.com/gorilla/websocket"
	"net/http"
	"reflect"
	"time"
)

type HoverflyStats interface {
	GetStats() metrics.Stats
	GetRequestCacheCount() (int, error)
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

	count, err := this.Hoverfly.GetRequestCacheCount()

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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// StatsWSHandler - returns current stats about Hoverfly (request counts, record count) through the websocket
func (this *StatsHandler) GetWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to upgrade websocket")
		return
	}

	// defining counters for delta check
	var recordsCount int
	var statsCounters map[string]int64

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		log.WithFields(log.Fields{
			"message": string(p),
		}).Info("Got message...")

		for _ = range time.Tick(1 * time.Second) {
			count, err := this.Hoverfly.GetRequestCacheCount()
			if err != nil {
				log.WithFields(log.Fields{
					"message": p,
					"error":   err.Error(),
				}).Error("got error while trying to get records count")
				continue
			}
			stats := this.Hoverfly.GetStats()

			// checking whether we should send an update
			if !reflect.DeepEqual(stats.Counters, statsCounters) || count != recordsCount {
				var sr StatsResponse
				sr.Stats = stats
				sr.RecordsCount = count

				b, err := json.Marshal(sr)

				if err = conn.WriteMessage(messageType, b); err != nil {
					log.WithFields(log.Fields{
						"message": p,
						"error":   err.Error(),
					}).Debug("Got error when writing message...")
					continue
				}
				recordsCount = count
				statsCounters = stats.Counters
			}
		}
	}
}
