package templating

import (
	"strings"

	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/aymerick/raymond"
	"time"
	"strconv"
)

type TemplatingData struct {
	Request Request
	State   map[string]string
}

type Request struct {
	QueryParam map[string][]string
	Path       []string
	Scheme     string
}

type Templator struct{

}

var helpersRegistered = false

func NewTemplator() *Templator {

	if !helpersRegistered {
		raymond.RegisterHelper("iso8601DateTime", func() string {
			return time.Now().UTC().Format("2006-01-02T15:04:05Z07:00")
		})

		raymond.RegisterHelper("iso8601DateTimePlusDays", func(days string) string {
			atoi, _ := strconv.Atoi(days)
			return time.Now().AddDate(0, 0, atoi).UTC().Format("2006-01-02T15:04:05Z07:00")
		})

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
	}

}
