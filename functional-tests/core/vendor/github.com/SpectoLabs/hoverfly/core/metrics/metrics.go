package metrics

import (
	log "github.com/Sirupsen/logrus"
	"github.com/rcrowley/go-metrics"

	"time"
)

// CounterByMode - container for mode counters, registry and flush interval
type CounterByMode struct {
	Counters      map[string]metrics.Counter
	registry      metrics.Registry
	flushInterval time.Duration
}

// NewModeCounter - returns new counter instance
func NewModeCounter(modes []string) *CounterByMode {

	registry := metrics.NewRegistry()
	counters := make(map[string]metrics.Counter)

	for _, v := range modes {
		counter := metrics.NewCounter()
		counters[v] = counter
		registry.GetOrRegister(v, counter)
	}

	c := &CounterByMode{
		Counters:      counters,
		registry:      registry,
		flushInterval: 5 * time.Second,
	}

	log.Debug("new counter created, registration successful")

	return c
}

// Count - counts requests based on mode
func (c *CounterByMode) Count(mode string) {
	c.Counters[mode].Inc(1)
}

// Init initializes logging
func (c *CounterByMode) Init() {
	go func() {
		for _ = range time.Tick(c.flushInterval) {
			m := c.Flush()
			log.WithFields(log.Fields{"counters": m.Counters}).Info("hoverfly metrics")
		}
	}()
}

// Stats - holds information about various system metrics like requests counts
type Stats struct {
	Counters    map[string]int64   `json:"counters"`
	Gauges      map[string]int64   `json:"gauges,omitempty"`
	GaugesFloat map[string]float64 `json:"gaugesFloat,omitempty"`
}

// Flush gets current metrics from stats registry
func (c *CounterByMode) Flush() (h Stats) {

	counters := make(map[string]int64)
	gauges := make(map[string]int64)
	gaugesFloat := make(map[string]float64)

	c.registry.Each(func(name string, i interface{}) {
		switch metric := i.(type) {
		case metrics.Counter:
			counters[name] = metric.Count()
		case metrics.Gauge:
			gauges[name] = metric.Value()
		case metrics.GaugeFloat64:
			gaugesFloat[name] = metric.Value()
		}
	})

	h.Counters = counters
	h.Gauges = gauges
	h.GaugesFloat = gaugesFloat
	return
}
