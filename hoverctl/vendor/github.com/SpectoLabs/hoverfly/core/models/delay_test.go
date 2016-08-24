package models

import (
	"encoding/json"
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
	var responseDelayJson ResponseDelayJson
	json.Unmarshal([]byte(jsonConf), &responseDelayJson)
	err := ValidateResponseDelayJson(responseDelayJson)
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
	var responseDelayJson ResponseDelayJson
	json.Unmarshal([]byte(jsonConf), &responseDelayJson)
	err := ValidateResponseDelayJson(responseDelayJson)
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
	var responseDelayJson ResponseDelayJson
	json.Unmarshal([]byte(jsonConf), &responseDelayJson)
	err := ValidateResponseDelayJson(responseDelayJson)
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
	var responseDelayJson ResponseDelayJson
	json.Unmarshal([]byte(jsonConf), &responseDelayJson)
	err := ValidateResponseDelayJson(responseDelayJson)
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
	var responseDelayJson ResponseDelayJson
	json.Unmarshal([]byte(jsonConf), &responseDelayJson)
	err := ValidateResponseDelayJson(responseDelayJson)
	Expect(err).To(Not(BeNil()))
}

func TestGetDelayWithRegexMatch(t *testing.T) {
	RegisterTestingT(t)

	delay := ResponseDelay{
		UrlPattern: "example(.+)",
		Delay:      100,
	}
	delays := ResponseDelayList{delay}

	delayMatch := delays.GetDelay("delayexample.com", "method-dummy")
	Expect(*delayMatch).To(Equal(delay))

	delayMatch = delays.GetDelay("nodelay.com", "method-dummy")
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

	delayMatch := delays.GetDelay("delayexample.com", "method-dummy")
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

	delayMatch := delays.GetDelay("delayexample.com", "GET")
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

	delayMatch := delays.GetDelay("delayexample.com", "GET")
	Expect(*delayMatch).To(Equal(delay))
}

func TestIfDelayMethodBlankThenMatchesAnyMethod(t *testing.T) {
	RegisterTestingT(t)

	delay := ResponseDelay{
		UrlPattern: "example(.+)",
		Delay:      100,
	}
	delays := ResponseDelayList{delay}

	delayMatch := delays.GetDelay("delayexample.com", "method-dummy")
	Expect(*delayMatch).To(Equal(delay))
}
