package models

import (
	json "encoding/json"
	"regexp"
	log "github.com/Sirupsen/logrus"
)

type ResponseDelay struct {
	HostPattern string
	Delay int
	DelayStdDev int
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
	// apply the delay
	log.Warn("execute delay not implemented")
}