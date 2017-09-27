package templating

import (
	"strings"

	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/aymerick/raymond"
)

type TemplatingData struct {
	Request Request
	State      map[string]string
}

type Request struct {
	QueryParam map[string][]string
	Path       []string
	Scheme     string
}

func ApplyTemplate(requestDetails *models.RequestDetails, state map[string]string, responseBody string) (string, error) {

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
