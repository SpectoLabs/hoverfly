package models

import (
	json "encoding/json"
	"regexp"
	log "github.com/Sirupsen/logrus"
	"time"
)

type ResponseDelay struct {
	HostPattern string
	Delay int
}

type ResponseDelayJson struct {
	Data []ResponseDelay
}

func ParseResponseDelayJson(j []byte) []ResponseDelay {
	var responseDelayJson ResponseDelayJson
	result := make([]ResponseDelay,0)
	json.Unmarshal(j, &responseDelayJson)

	// filter any entries that don't meet the invariants
	for _, delay := range responseDelayJson.Data {
		if delay.HostPattern != "" && delay.Delay != 0 {
			if _, err := regexp.Compile(delay.HostPattern); err == nil {
				result = append(result, delay)
			} else {
				log.Warn("Response delay entry skipped due to invalid pattern : %s", delay.HostPattern)
			}
		} else {
			log.Warn("Response delay entry skipped due to missing values: %v", delay)
		}
	}
	return result
}

func (this *ResponseDelay) Execute() {
	// apply the delay - must be called from goroutine handling the request
	log.Info("Pausing before sending the response to simulate delays")
	time.Sleep(time.Duration(this.Delay) * time.Millisecond)
	log.Info("Response delay completed")
}

