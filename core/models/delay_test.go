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
				"delay": 100,
				"delayStdDev": 10
			}]
	}`
	byteArr := make([]byte, len(jsonConf))
	copy(byteArr[:], jsonConf)
	responseDelayConf := ParseResponseDelayJson(byteArr)
	Expect(responseDelayConf[0].HostPattern).To(Equal("."))
	Expect(responseDelayConf[0].Delay).To((Equal(100)))
	Expect(responseDelayConf[0].DelayStdDev).To(Equal(10))
}

func TestDefaultResponseStdDevIsZero(t *testing.T) {
	RegisterTestingT(t)

	jsonConf := `
	{
		"data": [{
				"hostPattern": ".",
				"delay": 100
			}]
	}`
	byteArr := make([]byte, len(jsonConf))
	copy(byteArr[:], jsonConf)
	responseDelayConf := ParseResponseDelayJson(byteArr)
	Expect(responseDelayConf[0].HostPattern).To(Equal("."))
	Expect(responseDelayConf[0].Delay).To((Equal(100)))
	Expect(responseDelayConf[0].DelayStdDev).To(Equal(0))
}

func TestDelayIsIgnoredIfHostPatternNotSet(t *testing.T) {
	RegisterTestingT(t)

	jsonConf := `
	{
		"data": [{
				"delay": 100
			}]
	}`
	byteArr := make([]byte, len(jsonConf))
	copy(byteArr[:], jsonConf)
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
	byteArr := make([]byte, len(jsonConf))
	copy(byteArr[:], jsonConf)
	responseDelayConf := ParseResponseDelayJson(byteArr)
	Expect(len(responseDelayConf)).To(Equal(0))
}

func TestHostPatternMustBeAValidRegexPattern(t *testing.T) {
	RegisterTestingT(t)

	jsonConf := `
	{
		"data": [{
				"hostPattern": "*",
				"delay": 100
			}]
	}`
	byteArr := make([]byte, len(jsonConf))
	copy(byteArr[:], jsonConf)
	responseDelayConf := ParseResponseDelayJson(byteArr)
	Expect(len(responseDelayConf)).To(Equal(0))
}