package models

import (
	"testing"
	. "github.com/onsi/gomega"
	"encoding/json"
)

func TestConvertJsonStringToResponseDelayConfig(t *testing.T) {
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
	Expect(err).To(BeNil())
}


func TestDelayIsIgnoredIfHostPatternNotSet(t *testing.T) {
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

func TestDelayIsIgnoredIfDelayNotSet(t *testing.T) {
	RegisterTestingT(t)

	jsonConf := `
	{
		"data": [{
				"hostPattern": "."
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
				"hostPattern": "*",
				"delay": 1
			}]
	}`
	var responseDelayJson ResponseDelayJson
	json.Unmarshal([]byte(jsonConf), &responseDelayJson)
	err := ValidateResponseDelayJson(responseDelayJson)
	Expect(err).To(Not(BeNil()))
}