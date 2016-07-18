package models

import (
	"testing"
	. "github.com/onsi/gomega"
	"encoding/json"
	"github.com/SpectoLabs/hoverfly/core/testutil"
	"fmt"
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
	delay := ResponseDelay{
		UrlPattern: "example(.+)",
		Delay:       100,
	}
	delays := ResponseDelayList{delay}

	delayMatch := delays.GetDelay("delayexample.com", "method-dummy")
	testutil.Expect(t, *delayMatch, delay)

	var nilDelay *ResponseDelay
	delayMatch = delays.GetDelay("nodelay.com", "method-dummy")
	testutil.Expect(t, delayMatch, nilDelay)
}

func TestMultipleMatchingDelaysReturnsTheFirst(t *testing.T) {
	delayOne := ResponseDelay{
		UrlPattern: "example.com",
		Delay:       100,
	}
	delayTwo := ResponseDelay{
		UrlPattern: "example",
		Delay:       100,
	}
	delays := ResponseDelayList{delayOne, delayTwo}

	delayMatch := delays.GetDelay("delayexample.com", "method-dummy")
	testutil.Expect(t, *delayMatch, delayOne)
}

func TestNoMatchIfMethodsDontMatch(t *testing.T) {
	delay := ResponseDelay{
		UrlPattern: "example.com",
		Delay:       100,
		HttpMethod: "PURPLE",
	}
	delays := ResponseDelayList{delay}

	var nilDelay *ResponseDelay
	delayMatch := delays.GetDelay("delayexample.com", "GET")
	testutil.Expect(t, delayMatch, nilDelay)
	if (delayMatch!=nil) {
		t.Fail()
	}
}

func TestReturnMatchIfMethodsMatch(t *testing.T) {
	delay := ResponseDelay{
		UrlPattern: "example.com",
		Delay:       100,
		HttpMethod: "GET",
	}
	delays := ResponseDelayList{delay}

	delayMatch := delays.GetDelay("delayexample.com", "GET")
	testutil.Expect(t, *delayMatch, delay)
}

func TestIfDelayMethodBlankThenMatchesAnyMethod(t *testing.T) {
	delay := ResponseDelay{
		UrlPattern: "example(.+)",
		Delay:       100,
	}
	delays := ResponseDelayList{delay}

	delayMatch := delays.GetDelay("delayexample.com", "method-dummy")
	testutil.Expect(t, *delayMatch, delay)
}

func TestMarshalToJSONWorks(t *testing.T) {
	delay := ResponseDelay{
		UrlPattern: "example(.+)",
		Delay:       100,
	}
	delays := ResponseDelayList{delay}

	resp := ResponseDelayJson{
		Data: &delays,
	}
	b, _ := json.Marshal(resp)
	fmt.Print(string(b))
}