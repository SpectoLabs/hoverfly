package models

import (
	"testing"
	. "github.com/onsi/gomega"
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
	byteArr := []byte(jsonConf)
	responseDelayConf := ParseResponseDelayJson(byteArr)
	Expect(responseDelayConf[0].HostPattern).To(Equal("."))
	Expect(responseDelayConf[0].Delay).To((Equal(1)))
}

func TestDefaultResponseStdDevIsZero(t *testing.T) {
	RegisterTestingT(t)

	jsonConf := `
	{
		"data": [{
				"hostPattern": ".",
				"delay": 1
			}]
	}`
	byteArr := []byte(jsonConf)
	responseDelayConf := ParseResponseDelayJson(byteArr)
	Expect(responseDelayConf[0].HostPattern).To(Equal("."))
	Expect(responseDelayConf[0].Delay).To((Equal(1)))
}

func TestDelayIsIgnoredIfHostPatternNotSet(t *testing.T) {
	RegisterTestingT(t)

	jsonConf := `
	{
		"data": [{
				"delay": 2
			}]
	}`
	byteArr := []byte(jsonConf)
	responseDelayConf := ParseResponseDelayJson(byteArr)
	Expect(len(responseDelayConf)).To(Equal(0))
}

func TestDelayIsIgnoredIfDelayNotSet(t *testing.T) {
	RegisterTestingT(t)

	jsonConf := `
	{
		"data": [{
				"hostPattern": "."
			}]
	}`
	byteArr := []byte(jsonConf)
	responseDelayConf := ParseResponseDelayJson(byteArr)
	Expect(len(responseDelayConf)).To(Equal(0))
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
	byteArr := []byte(jsonConf)
	responseDelayConf := ParseResponseDelayJson(byteArr)
	Expect(len(responseDelayConf)).To(Equal(0))
}