package templating

import (
	"fmt"
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
	Body       func(queryType, query string, options *raymond.Options) string
	body       string
	Method     string
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

func (*Templator) ParseTemplate(responseBody string) (*raymond.Template, error) {

	return raymond.Parse(responseBody)
}

func (*Templator) RenderTemplate(tpl *raymond.Template, requestDetails *models.RequestDetails, state map[string]string) (string, error) {
	if tpl == nil {
		return "", fmt.Errorf("template cannot be nil")
	}
	ctx := NewTemplatingDataFromRequest(requestDetails, state)
	return tpl.Exec(ctx)
}


func NewTemplatingDataFromRequest(requestDetails *models.RequestDetails, state map[string]string) *TemplatingData {
	return &TemplatingData{
		Request: Request{
			Path:       strings.Split(requestDetails.Path, "/")[1:],
			QueryParam: requestDetails.Query,
			Scheme:     requestDetails.Scheme,
			Body:       templateHelpers{}.requestBody,
			body:       requestDetails.Body,
			Method:     requestDetails.Method,
		},
		State: state,
		CurrentDateTime: func(a1, a2, a3 string) string {
			return a1 + " " + a2 + " " + a3
		},
	}

}
