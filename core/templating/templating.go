package templating

import (
	"strings"

	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/aymerick/raymond"
)

type TemplatingData struct {
	Request Request
}

type Request struct {
	QueryParam map[string][]string
	Path       []string
	Scheme     string
}

func ApplyTemplate(requestDetails *models.RequestDetails, responseBody string) (string, error) {

	t := NewTemplatingDataFromRequest(requestDetails)

	if rendered, err := raymond.Render(responseBody, t); err == nil {
		responseBody = rendered
		return responseBody, nil
	} else {
		return "", err
	}
}

func NewTemplatingDataFromRequest(requestDetails *models.RequestDetails) *TemplatingData {
	return &TemplatingData{
		Request: Request{
			Path:       strings.Split(requestDetails.Path, "/")[1:],
			QueryParam: requestDetails.Query,
			Scheme:     requestDetails.Scheme,
		},
	}

}
