package models

import (
	"encoding/json"
	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	. "github.com/onsi/gomega"
	"testing"
)

func TestConvertJsonStringToResponseDelayConfig(t *testing.T) {
	RegisterTestingT(t)

	jsonConf := `
	{
		"data": [{
				"urlPattern": ".",
				"delay": 1
			}]
	}`
	var responseDelayJson v1.ResponseDelayPayloadView
	json.Unmarshal([]byte(jsonConf), &responseDelayJson)
	err := ValidateResponseDelayPayload(responseDelayJson)
	Expect(err).To(BeNil())
}

func TestErrorIfHostPatternNotSet(t *testing.T) {
	RegisterTestingT(t)

	jsonConf := `
	{
		"data": [{
				"delay": 2
			}]
	}`
	var responseDelayJson v1.ResponseDelayPayloadView
	json.Unmarshal([]byte(jsonConf), &responseDelayJson)
	err := ValidateResponseDelayPayload(responseDelayJson)
	Expect(err).To(Not(BeNil()))
}

func TestErrprIfDelayNotSet(t *testing.T) {
	RegisterTestingT(t)

	jsonConf := `
	{
		"data": [{
				"urlPattern": "."
			}]
	}`
	var responseDelayJson v1.ResponseDelayPayloadView
	json.Unmarshal([]byte(jsonConf), &responseDelayJson)
	err := ValidateResponseDelayPayload(responseDelayJson)
	Expect(err).To(Not(BeNil()))
}

func TestHostPatternMustBeAValidRegexPattern(t *testing.T) {
	RegisterTestingT(t)

	jsonConf := `
	{
		"data": [{
				"urlPattern": "*",
				"delay": 1
			}]
	}`
	var responseDelayJson v1.ResponseDelayPayloadView
	json.Unmarshal([]byte(jsonConf), &responseDelayJson)
	err := ValidateResponseDelayPayload(responseDelayJson)
	Expect(err).To(Not(BeNil()))
}

func TestErrorIfHostPatternUsed(t *testing.T) {
	RegisterTestingT(t)

	jsonConf := `
	{
		"data": [{
				"hostPattern": ".",
				"delay": 1
			}]
	}`
	var responseDelayJson v1.ResponseDelayPayloadView
	json.Unmarshal([]byte(jsonConf), &responseDelayJson)
	err := ValidateResponseDelayPayload(responseDelayJson)
	Expect(err).To(Not(BeNil()))
}

func TestGetDelayWithRegexMatch(t *testing.T) {
	RegisterTestingT(t)

	delay := ResponseDelay{
		UrlPattern: "example(.+)",
		Delay:      100,
	}
	delays := ResponseDelayList{delay}

	request1 := RequestDetails{
		Destination: "delayexample.com",
		Method:      "method-dummy",
	}

	delayMatch := delays.GetDelay(request1)
	Expect(*delayMatch).To(Equal(delay))

	request2 := RequestDetails{
		Destination: "nodelay.com",
		Method:      "method-dummy",
	}

	delayMatch = delays.GetDelay(request2)
	Expect(delayMatch).To(BeNil())
}

func TestMultipleMatchingDelaysReturnsTheFirst(t *testing.T) {
	RegisterTestingT(t)

	delayOne := ResponseDelay{
		UrlPattern: "example.com",
		Delay:      100,
	}
	delayTwo := ResponseDelay{
		UrlPattern: "example",
		Delay:      100,
	}
	delays := ResponseDelayList{delayOne, delayTwo}

	request1 := RequestDetails{
		Destination: "delayexample.com",
		Method:      "method-dummy",
	}

	delayMatch := delays.GetDelay(request1)
	Expect(*delayMatch).To(Equal(delayOne))
}

func TestNoMatchIfMethodsDontMatch(t *testing.T) {
	RegisterTestingT(t)

	delay := ResponseDelay{
		UrlPattern: "example.com",
		Delay:      100,
		HttpMethod: "PURPLE",
	}
	delays := ResponseDelayList{delay}

	request := RequestDetails{
		Destination: "delayexample.com",
		Method:      "GET",
	}

	delayMatch := delays.GetDelay(request)
	Expect(delayMatch).To(BeNil())
}

func TestReturnMatchIfMethodsMatch(t *testing.T) {
	RegisterTestingT(t)

	delay := ResponseDelay{
		UrlPattern: "example.com",
		Delay:      100,
		HttpMethod: "GET",
	}
	delays := ResponseDelayList{delay}

	request := RequestDetails{
		Destination: "delayexample.com",
		Method:      "GET",
	}

	delayMatch := delays.GetDelay(request)
	Expect(*delayMatch).To(Equal(delay))
}

func TestIfDelayMethodBlankThenMatchesAnyMethod(t *testing.T) {
	RegisterTestingT(t)

	delay := ResponseDelay{
		UrlPattern: "example(.+)",
		Delay:      100,
	}
	delays := ResponseDelayList{delay}

	request := RequestDetails{
		Destination: "delayexample.com",
		Method:      "method-dummy",
	}

	delayMatch := delays.GetDelay(request)
	Expect(*delayMatch).To(Equal(delay))
}

func TestResponseDelayList_ConvertToPayloadView(t *testing.T) {
	RegisterTestingT(t)

	delay := ResponseDelay{
		UrlPattern: "example(.+)",
		Delay:      100,
	}
	delays := ResponseDelayList{delay}

	payloadView := delays.ConvertToResponseDelayPayloadView()

	Expect(payloadView.Data[0].UrlPattern).To(Equal("example(.+)"))
	Expect(payloadView.Data[0].Delay).To(Equal(100))

}
