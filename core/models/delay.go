package models

import (
	"regexp"
	log "github.com/Sirupsen/logrus"
	"time"
	"errors"
	"fmt"
)

type ResponseDelay struct {
	HostPattern string
	Delay int
}

type ResponseDelayJson struct {
	Data []ResponseDelay
}

func ValidateResponseDelayJson(j ResponseDelayJson) (err error) {
	// filter any entries that don't meet the invariants
	for _, delay := range j.Data {
		if delay.HostPattern != "" && delay.Delay != 0 {
			if _, err := regexp.Compile(delay.HostPattern); err != nil {
				return errors.New(fmt.Sprintf("Response delay entry skipped due to invalid pattern : %s", delay.HostPattern))
			}
		} else {
			return errors.New(fmt.Sprintf("Config error - Missing values found in: %v", delay))
		}
	}
	return nil
}

func (this *ResponseDelay) Execute() {
	// apply the delay - must be called from goroutine handling the request
	log.Info("Pausing before sending the response to simulate delays")
	time.Sleep(time.Duration(this.Delay) * time.Millisecond)
	log.Info("Response delay completed")
}

