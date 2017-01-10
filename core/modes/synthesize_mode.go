package modes

import (
	"fmt"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/models"

	log "github.com/Sirupsen/logrus"
)

type SynthesizeMode struct {
	Hoverfly Hoverfly
}

func (this SynthesizeMode) Process(request *http.Request, details models.RequestDetails) (*http.Response, error) {
	pair := models.RequestResponsePair{Request: details}

	log.WithFields(log.Fields{
		// "middleware":  this.hoverfly.Cfg.Middleware.toString(),
		"body":        details.Body,
		"destination": details.Destination,
	}).Debug("Synthesizing new response")

	if !this.Hoverfly.IsMiddlewareSet() {
		err := fmt.Errorf("Middleware not set")
		return errorResponse(request, err, "Synthesize failed, middleware not provided"), err
	}

	pair, err := this.Hoverfly.ApplyMiddleware(pair)

	if err != nil {
		return errorResponse(request, err, "Could not create synthetic response!"), err
	}

	log.WithFields(log.Fields{
		"mode": "synthesize",
		// "middleware":  this.hoverfly.Cfg.Middleware,
		"path":        request.URL.Path,
		"rawQuery":    request.URL.RawQuery,
		"method":      request.Method,
		"destination": request.Host,
	}).Info("synthetic response created successfuly")

	return ReconstructResponse(request, pair), nil
}
