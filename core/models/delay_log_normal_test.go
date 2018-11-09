package models_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
)

func TestDelayLogNormalConvertJsonStringToResponseDelayLogNormalConfig(t *testing.T) {
	RegisterTestingT(t)

	jsonConf := `
	{
		"data": [{
				"urlPattern": ".",
				"min": 100,
				"max": 10000,
				"mean": 5000,
				"median": 500
			}]
	}`
	var responseDelayLogNormalJson v1.ResponseDelayLogNormalPayloadView
	json.Unmarshal([]byte(jsonConf), &responseDelayLogNormalJson)
	err := models.ValidateResponseDelayLogNormalPayload(responseDelayLogNormalJson)
	Expect(err).To(BeNil())
}

func TestDelayLogNormalErrorIfHostPatternNotSet(t *testing.T) {
	RegisterTestingT(t)

	jsonConf := `
	{
		"data": [{
				"min": 100,
				"max": 10000,
				"mean": 5000,
				"median": 500
			}]
	}`
	var responseDelayLogNormalJson v1.ResponseDelayLogNormalPayloadView
	json.Unmarshal([]byte(jsonConf), &responseDelayLogNormalJson)
	err := models.ValidateResponseDelayLogNormalPayload(responseDelayLogNormalJson)
	Expect(err).To(Not(BeNil()))
}

func TestDelayLogNormalHostPatternMustBeAValidRegexPattern(t *testing.T) {
	RegisterTestingT(t)

	jsonConf := `
	{
		"data": [{
				"urlPattern": "*",
				"min": 100,
				"max": 10000,
				"mean": 5000,
				"median": 500
			}]
	}`
	var responseDelayLogNormalJson v1.ResponseDelayLogNormalPayloadView
	json.Unmarshal([]byte(jsonConf), &responseDelayLogNormalJson)
	err := models.ValidateResponseDelayLogNormalPayload(responseDelayLogNormalJson)
	Expect(err).To(Not(BeNil()))
}

func TestDelayLogNormalErrorIfHostPatternUsed(t *testing.T) {
	RegisterTestingT(t)

	jsonConf := `
	{
		"data": [{
				"hostPattern": ".",
				"min": 100,
				"max": 10000,
				"mean": 5000,
				"median": 500
			}]
	}`
	var responseDelayLogNormalJson v1.ResponseDelayLogNormalPayloadView
	json.Unmarshal([]byte(jsonConf), &responseDelayLogNormalJson)
	err := models.ValidateResponseDelayLogNormalPayload(responseDelayLogNormalJson)
	Expect(err).To(Not(BeNil()))
}

func TestDelayLogNormalErrorWrongTimeParams(t *testing.T) {
	RegisterTestingT(t)
	cases := []struct {
		jsonConf     string
		errorMessage string
	}{
		{
			`{
				"data": [{
					"urlPattern": ".",
					"min": -1
				}]
			}`,
			"Config error - delay min and max can't be less than 0",
		},
		{
			`{
				"data": [{
					"urlPattern": ".",
					"max": -1
				}]
			}`,
			"Config error - delay min and max can't be less than 0",
		},
		{
			`{
				"data": [{
					"urlPattern": ".",
					"mean": 0,
					"median": 1
				}]
			}`,
			"Config error - delay mean and median params can't be less or equals 0",
		},
		{
			`{
				"data": [{
					"urlPattern": ".",
					"mean": 1,
					"median": 0
				}]
			}`,
			"Config error - delay mean and median params can't be less or equals 0",
		},
		{
			`{
				"data": [{
					"urlPattern": ".",
					"min": 2,
					"max": 1,
					"mean": 2,
					"median": 1
				}]
			}`,
			"Config error - min delay must be less than max one",
		},
		{
			`{
				"data": [{
					"urlPattern": ".",
					"min": 1,
					"max": 20,
					"mean": 30,
					"median": 1
				}]
			}`,
			"Config error - mean delay can't be greather than max one",
		},
		{
			`{
				"data": [{
					"urlPattern": ".",
					"min": 1,
					"max": 30,
					"mean": 20,
					"median": 40
				}]
			}`,
			"Config error - median delay can't be and greather than max one",
		},
		{
			`{
				"data": [{
					"urlPattern": ".",
					"min": 10,
					"max": 20,
					"mean": 2,
					"median": 15
				}]
			}`,
			"Config error - mean delay can't be less than min one",
		},
		{
			`{
				"data": [{
					"urlPattern": ".",
					"min": 10,
					"max": 20,
					"mean": 15,
					"median": 2
				}]
			}`,
			"Config error - median delay can't be less than min one",
		},
		{
			`{
				"data": [{
					"urlPattern": ".",
					"min": 10,
					"max": 40,
					"mean": 15,
					"median": 20
				}]
			}`,
			"Config error - mean delay can't be less than median one",
		},
	}
	for i, tc := range cases {
		var responseDelayLogNormalJson v1.ResponseDelayLogNormalPayloadView
		json.Unmarshal([]byte(tc.jsonConf), &responseDelayLogNormalJson)
		err := models.ValidateResponseDelayLogNormalPayload(responseDelayLogNormalJson)
		Expect(err).To(Not(BeNil()))
		Expect(err.Error()).To(Equal(tc.errorMessage), fmt.Sprintf("Case #%d", i))
	}

}

type DelayGeneratorMock struct {
	counter int
}

func (g *DelayGeneratorMock) GenerateDelay() int {
	g.counter++
	return 0
}

func TestResponseDelayLogNormal_Execute(t *testing.T) {
	RegisterTestingT(t)
	delayGeneratorMock := &DelayGeneratorMock{}
	delay := models.ResponseDelayLogNormal{
		DelayGenerator: delayGeneratorMock,
	}
	delay.Execute()
	Expect(delayGeneratorMock.counter).To(Equal(1))
}

func TestDelayLogNormalGetDelayLogNormalWithRegexMatch(t *testing.T) {
	RegisterTestingT(t)

	delayLogNormal := models.ResponseDelayLogNormal{
		UrlPattern: "example(.+)",
		Mean:       5000,
		Median:     500,
	}
	delayLogNormals := models.ResponseDelayLogNormalList{delayLogNormal}

	request1 := models.RequestDetails{
		Destination: "delayLogNormalexample.com",
		Method:      "method-dummy",
	}

	delayLogNormalMatch := delayLogNormals.GetDelay(request1)
	Expect(*delayLogNormalMatch).To(Equal(delayLogNormal))

	request2 := models.RequestDetails{
		Destination: "nodelayLogNormal.com",
		Method:      "method-dummy",
	}

	delayLogNormalMatch = delayLogNormals.GetDelay(request2)
	Expect(delayLogNormalMatch).To(BeNil())
}

func TestDelayLogNormalMultipleMatchingDelayLogNormalsReturnsTheFirst(t *testing.T) {
	RegisterTestingT(t)

	delayLogNormalOne := models.ResponseDelayLogNormal{
		UrlPattern: "example.com",
		Mean:       5000,
		Median:     500,
	}
	delayLogNormalTwo := models.ResponseDelayLogNormal{
		UrlPattern: "example",
		Mean:       5000,
		Median:     500,
	}
	delayLogNormals := models.ResponseDelayLogNormalList{delayLogNormalOne, delayLogNormalTwo}

	request1 := models.RequestDetails{
		Destination: "delayLogNormalexample.com",
		Method:      "method-dummy",
	}

	delayLogNormalMatch := delayLogNormals.GetDelay(request1)
	Expect(*delayLogNormalMatch).To(Equal(delayLogNormalOne))
}

func TestDelayLogNormalNoMatchIfMethodsDontMatch(t *testing.T) {
	RegisterTestingT(t)

	delayLogNormal := models.ResponseDelayLogNormal{
		UrlPattern: "example.com",
		Mean:       5000,
		Median:     500,
		HttpMethod: "PURPLE",
	}
	delayLogNormals := models.ResponseDelayLogNormalList{delayLogNormal}

	request := models.RequestDetails{
		Destination: "delayLogNormalexample.com",
		Method:      "GET",
	}

	delayLogNormalMatch := delayLogNormals.GetDelay(request)
	Expect(delayLogNormalMatch).To(BeNil())
}

func TestDelayLogNormalReturnMatchIfMethodsMatch(t *testing.T) {
	RegisterTestingT(t)

	delayLogNormal := models.ResponseDelayLogNormal{
		UrlPattern: "example.com",
		Mean:       5000,
		Median:     500,
		HttpMethod: "GET",
	}
	delayLogNormals := models.ResponseDelayLogNormalList{delayLogNormal}

	request := models.RequestDetails{
		Destination: "delayLogNormalexample.com",
		Method:      "GET",
	}

	delayLogNormalMatch := delayLogNormals.GetDelay(request)
	Expect(*delayLogNormalMatch).To(Equal(delayLogNormal))
}

func TestDelayLogNormalIfDelayLogNormalMethodBlankThenMatchesAnyMethod(t *testing.T) {
	RegisterTestingT(t)

	delayLogNormal := models.ResponseDelayLogNormal{
		UrlPattern: "example(.+)",
		Mean:       5000,
		Median:     500,
	}
	delayLogNormals := models.ResponseDelayLogNormalList{delayLogNormal}

	request := models.RequestDetails{
		Destination: "delayLogNormalexample.com",
		Method:      "method-dummy",
	}

	delayLogNormalMatch := delayLogNormals.GetDelay(request)
	Expect(*delayLogNormalMatch).To(Equal(delayLogNormal))
}

func TestDelayLogNormalResponseDelayLogNormalList_ConvertToPayloadView(t *testing.T) {
	RegisterTestingT(t)

	delayLogNormal := models.ResponseDelayLogNormal{
		UrlPattern: "example(.+)",
		Mean:       5000,
		Median:     500,
	}
	delayLogNormals := models.ResponseDelayLogNormalList{delayLogNormal}

	payloadView := delayLogNormals.ConvertToResponseDelayLogNormalPayloadView()

	Expect(payloadView.Data[0].UrlPattern).To(Equal("example(.+)"))
	Expect(payloadView.Data[0].Mean).To(Equal(5000))
	Expect(payloadView.Data[0].Median).To(Equal(500))

}
