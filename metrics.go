package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/rcrowley/go-metrics"

	"fmt"
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

	registry := metrics.DefaultRegistry

	c := &CounterByMode{
		counterVirtualize: metrics.NewCounter(),
		counterCapture:    metrics.NewCounter(),
		counterModify:     metrics.NewCounter(),
		counterSynthesize: metrics.NewCounter(),
		registry:          registry,
		flushInterval:     5 * time.Second,
	}

	orPanic(c.registry.Register(fmt.Sprintf(VirtualizeMode), c.counterVirtualize))
	orPanic(c.registry.Register(fmt.Sprintf(CaptureMode), c.counterCapture))
	orPanic(c.registry.Register(fmt.Sprintf(ModifyMode), c.counterModify))
	orPanic(c.registry.Register(fmt.Sprintf(SynthesizeMode), c.counterSynthesize))

	log.Info("new counter created, registration successful")

	return c
}

func (c *CounterByMode) count(mode string) {
	if mode == VirtualizeMode {
		c.counterVirtualize.Inc(1)
	} else if mode == CaptureMode {
		c.counterCapture.Inc(1)
	} else if mode == ModifyMode {
		c.counterModify.Inc(1)
	} else if mode == SynthesizeMode {
		c.counterModify.Inc(1)
	}
}
