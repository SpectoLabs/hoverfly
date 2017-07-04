package templating

import (
	"github.com/aymerick/raymond"
	"net/http"
	"strings"
)

type TemplatingData struct {
	Request Request
}

type Request struct {
	QueryParam map[string][]string
	Path  []string
	Scheme     string
}

func ApplyTemplate(request *http.Request, responseBody string) (string, error) {

	t := NewTemplatingDataFromRequest(request)

	if rendered, err := raymond.Render(responseBody, t); err == nil {
		responseBody = rendered
		return responseBody, nil
	} else {
		return "", err
	}
}

func NewTemplatingDataFromRequest(request *http.Request) * TemplatingData {

	requestPath := request.URL.Path

	Path := strings.Split(requestPath, "/")[1:]

	return &TemplatingData{
		Request: Request{
			Path: Path,
			QueryParam: request.URL.Query(),
			Scheme: request.URL.Scheme,
		},
	}

}