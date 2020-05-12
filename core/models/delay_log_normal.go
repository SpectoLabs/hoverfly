package models

import (
	"errors"
	"fmt"
	coreDelay "github.com/SpectoLabs/hoverfly/core/delay"
	"regexp"
	"strings"
	"time"

	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	log "github.com/sirupsen/logrus"
)

type ResponseDelayLogNormal struct {
	UrlPattern     string         `json:"urlPattern"`
	HttpMethod     string         `json:"httpMethod"`
	DelayGenerator DelayGenerator `json:"-"`
	Min            int            `json:"min"`
	Max            int            `json:"max"`
	Mean           int            `json:"mean"`
	Median         int            `json:"median"`
}

type ResponseDelayLogNormalList []ResponseDelayLogNormal

type ResponseDelaysLogNormal interface {
	GetDelay(request RequestDetails) *ResponseDelayLogNormal
	ConvertToResponseDelayLogNormalPayloadView() v1.ResponseDelayLogNormalPayloadView
}

type DelayGenerator interface {
	GenerateDelay() int
}

func ValidateResponseDelayLogNormalPayload(j v1.ResponseDelayLogNormalPayloadView) (err error) {
	if j.Data != nil {
		for _, delay := range j.Data {
			if delay.UrlPattern == "" {
				return errors.New("Config error - Missing urlPattern")
			}

			if _, err := regexp.Compile(delay.UrlPattern); err != nil {
				return errors.New(fmt.Sprintf("Response delay entry skipped due to invalid pattern : %s", delay.UrlPattern))
			}

			if err := coreDelay.ValidateLogNormalDelayOptions(delay.Min, delay.Max, delay.Mean, delay.Median); err != nil {
				return err
			}
		}
	}
	return nil
}

func (this *ResponseDelayLogNormal) Execute() {
	// apply the delay - must be called from goroutine handling the request
	log.Info("Pausing before sending the response to simulate delays")
	time.Sleep(time.Duration(this.DelayGenerator.GenerateDelay()) * time.Millisecond)
	log.Info("Response delay completed")
}

func (this *ResponseDelayLogNormalList) GetDelay(request RequestDetails) *ResponseDelayLogNormal {
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

func (this ResponseDelayLogNormalList) ConvertToResponseDelayLogNormalPayloadView() v1.ResponseDelayLogNormalPayloadView {
	payloadView := v1.ResponseDelayLogNormalPayloadView{
		Data: []v1.ResponseDelayLogNormalView{},
	}

	for _, responseDelayLogNormal := range this {
		responseDelayLogNormalView := v1.ResponseDelayLogNormalView{
			UrlPattern: responseDelayLogNormal.UrlPattern,
			HttpMethod: responseDelayLogNormal.HttpMethod,
			Min:        responseDelayLogNormal.Min,
			Max:        responseDelayLogNormal.Max,
			Mean:       responseDelayLogNormal.Mean,
			Median:     responseDelayLogNormal.Median,
		}

		payloadView.Data = append(payloadView.Data, responseDelayLogNormalView)
	}

	return payloadView
}
