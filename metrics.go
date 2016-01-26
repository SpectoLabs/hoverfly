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
