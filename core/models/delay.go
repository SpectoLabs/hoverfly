package models

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

type ResponseDelay struct {
	UrlPattern string `json:"urlPattern"`
	HttpMethod string `json:"httpMethod"`
	Delay      int    `json:"delay"`
}

type ResponseDelayPayload struct {
	Data *ResponseDelayList `json:"data"`
}

type ResponseDelayList []ResponseDelay

type ResponseDelays interface {
	Json() []byte
	GetDelay(request RequestDetails) *ResponseDelay
	Len() int
}

func ValidateResponseDelayJson(j ResponseDelayPayload) (err error) {
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

func (this *ResponseDelayList) GetDelay(request RequestDetails) *ResponseDelay {
	for _, val := range *this {
		match := regexp.MustCompile(val.UrlPattern).MatchString(request.Destination + request.Path)
		if match {
			if val.HttpMethod == "" || strings.EqualFold(val.HttpMethod, request.Method) {
				log.Info("Found response delay setting for this request host: ", val)
				return &val
			}
		}
	}
	return nil
}

func (this *ResponseDelayList) Json() []byte {
	resp := ResponseDelayPayload{
		Data: this,
	}
	b, _ := json.Marshal(resp)
	return b
}

func (this *ResponseDelayList) Len() int {
	list := []ResponseDelay{}
	if this != nil {
		list = append(list, *this...)
	}
	return len(list)
}
