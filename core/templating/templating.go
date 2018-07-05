package templating

import (
	"strings"
	"time"

	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/aymerick/raymond"
)

type TemplatingData struct {
	Request         Request
	State           map[string]string
	CurrentDateTime func(string, string, string) string
}

type Request struct {
	QueryParam map[string][]string
	Path       []string
	Scheme     string
}

type Templator struct {
}

var helpersRegistered = false

func NewTemplator() *Templator {
	t := templateHelpers{
		now: time.Now,
	}

	if !helpersRegistered {
		raymond.RegisterHelper("iso8601DateTime", t.iso8601DateTime)
		raymond.RegisterHelper("iso8601DateTimePlusDays", t.iso8601DateTimePlusDays)
		raymond.RegisterHelper("currentDateTime", t.currentDateTime)
		raymond.RegisterHelper("currentDateTimeAdd", t.currentDateTimeAdd)
		raymond.RegisterHelper("currentDateTimeSubtract", t.currentDateTimeSubtract)
		raymond.RegisterHelper("randomString", t.randomString)
		raymond.RegisterHelper("randomStringLength", t.randomStringLength)
		raymond.RegisterHelper("randomBoolean", t.randomBoolean)
		raymond.RegisterHelper("randomInteger", t.randomInteger)
		raymond.RegisterHelper("randomIntegerRange", t.randomIntegerRange)
		raymond.RegisterHelper("randomFloat", t.randomFloat)
		raymond.RegisterHelper("randomFloatRange", t.randomFloatRange)
		raymond.RegisterHelper("randomEmail", t.randomEmail)
		raymond.RegisterHelper("randomIPv4", t.randomIPv4)
		raymond.RegisterHelper("randomIPv6", t.randomIPv6)
		raymond.RegisterHelper("randomUuid", t.randomUuid)

		helpersRegistered = true
	}

	return &Templator{}
}

func (*Templator) ApplyTemplate(requestDetails *models.RequestDetails, state map[string]string, responseBody string) (string, error) {

	t := NewTemplatingDataFromRequest(requestDetails, state)

	if rendered, err := raymond.Render(responseBody, t); err == nil {
		responseBody = rendered
		return responseBody, nil
	} else {
		return "", err
	}
}

func NewTemplatingDataFromRequest(requestDetails *models.RequestDetails, state map[string]string) *TemplatingData {
	return &TemplatingData{
		Request: Request{
			Path:       strings.Split(requestDetails.Path, "/")[1:],
			QueryParam: requestDetails.Query,
			Scheme:     requestDetails.Scheme,
		},
		State: state,
		CurrentDateTime: func(a1, a2, a3 string) string {
			return a1 + " " + a2 + " " + a3
		},
	}

}
