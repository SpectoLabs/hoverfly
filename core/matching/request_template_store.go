package matching

import (
	"errors"

	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/ryanuber/go-glob"
)

func GetResponse(req models.RequestDetails, webserver bool, simulation *models.Simulation) (*models.ResponseDetails, error) {
	// iterate through the request templates, looking for template to match request
	for _, entry := range simulation.Templates {
		// TODO: not matching by default on URL and body - need to enable this
		// TODO: need to enable regex matches
		// TODO: enable matching on scheme

		template := entry.RequestTemplate

		if template.Body != nil && !glob.Glob(*template.Body, req.Body) {
			continue
		}

		if !webserver {
			if template.Destination != nil && !glob.Glob(*template.Destination, req.Destination) {
				continue
			}
		}
		if template.Path != nil && !glob.Glob(*template.Path, req.Path) {
			continue
		}
		if template.Query != nil && !glob.Glob(*template.Query, req.Query) {
			continue
		}
		if !headerMatch(template.Headers, req.Headers) {
			continue
		}
		if template.Method != nil && !glob.Glob(*template.Method, req.Method) {
			continue
		}

		// return the first template to match
		return &entry.Response, nil
	}
	return nil, errors.New("No match found")
}

func ConvertToRequestTemplateResponsePair(pairView v1.RequestTemplateResponsePairView) models.RequestTemplateResponsePair {
	return models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Path:        pairView.RequestTemplate.Path,
			Method:      pairView.RequestTemplate.Method,
			Destination: pairView.RequestTemplate.Destination,
			Scheme:      pairView.RequestTemplate.Scheme,
			Query:       pairView.RequestTemplate.Query,
			Body:        pairView.RequestTemplate.Body,
			Headers:     pairView.RequestTemplate.Headers,
		},
		Response: models.NewResponseDetailsFromResponse(pairView.Response),
	}
}
