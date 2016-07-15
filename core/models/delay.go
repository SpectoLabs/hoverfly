package models

import (
	"regexp"
	log "github.com/Sirupsen/logrus"
	"time"
	"errors"
	"fmt"
	"strings"
)

type ResponseDelay struct {
	UrlPattern string `json:"urlpattern"`
	HttpMethod string `json:"httpmethod"`
	Delay      int `json:"delay"`
}

type ResponseDelayJson struct {
	Data *ResponseDelayList `json:"data"`
}

type ResponseDelayList []ResponseDelay


func ValidateResponseDelayJson(j ResponseDelayJson) (err error) {
	if j.Data != nil {
		for _, delay := range *j.Data {
			if delay.UrlPattern != "" && delay.Delay != 0 {
				if _, err := regexp.Compile(delay.UrlPattern); err != nil {
					return errors.New(fmt.Sprintf("Response delay entry skipped due to invalid pattern : %s", delay.UrlPattern))
				}
			} else {
				return errors.New(fmt.Sprintf("Config error - Missing values found in: %v", delay))
			}
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

func (this *ResponseDelayList) GetDelay(url, httpMethod string) (*ResponseDelay) {
	for _, val := range *this {
		match := regexp.MustCompile(val.UrlPattern).MatchString(url)
		if match {
			if val.HttpMethod == "" || strings.EqualFold(val.HttpMethod, httpMethod) {
				log.Info("Found response delay setting for this request host: ", val)
				return &val
			}
		}
	}
	return nil
}