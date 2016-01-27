package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/rcrowley/go-metrics"

	"time"
)

// CounterByMode - container for mode counters, registry and flush interval
type CounterByMode struct {
	counterVirtualize, counterCapture, counterModify, counterSynthesize metrics.Counter
	registry                                                            metrics.Registry
	flushInterval                                                       time.Duration
}

// NewModeCounter - returns new counter instance
func NewModeCounter() *CounterByMode {

	registry := metrics.NewRegistry()

	c := &CounterByMode{
		counterVirtualize: metrics.NewCounter(),
		counterCapture:    metrics.NewCounter(),
		counterModify:     metrics.NewCounter(),
		counterSynthesize: metrics.NewCounter(),
		registry:          registry,
		flushInterval:     5 * time.Second,
	}

	c.registry.GetOrRegister(VirtualizeMode, c.counterVirtualize)
	c.registry.GetOrRegister(CaptureMode, c.counterCapture)
	c.registry.GetOrRegister(ModifyMode, c.counterModify)
	c.registry.GetOrRegister(SynthesizeMode, c.counterSynthesize)

	log.Debug("new counter created, registration successful")

	return c
}

// Count - counts requests based on mode
func (c *CounterByMode) Count(mode string) {
	if mode == VirtualizeMode {
		c.counterVirtualize.Inc(1)
	} else if mode == CaptureMode {
		c.counterCapture.Inc(1)
	} else if mode == ModifyMode {
		c.counterModify.Inc(1)
	} else if mode == SynthesizeMode {
		c.counterSynthesize.Inc(1)
	}
}

// Init initializes logging
func (c *CounterByMode) Init() {
	for _ = range time.Tick(c.flushInterval) {
		m := c.Flush()
		log.WithFields(log.Fields{"counters": m.Counters}).Info("hoverfly metrics")
	}
}

// HoverflyStats - holds information about various system metrics like requests counts
type HoverflyStats struct {
	Counters    map[string]int64   `json:"counters"`
	Gauges      map[string]int64   `json:"gauges,omitempty"`
	GaugesFloat map[string]float64 `json:"gaugesFloat,omitempty"`
}

// Flush gets current metrics from stats registry
func (c *CounterByMode) Flush() (h HoverflyStats) {

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
