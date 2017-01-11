package modes

import (
	"errors"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/models"

	log "github.com/Sirupsen/logrus"
)

type HoverflySynthesize interface {
	ApplyMiddleware(models.RequestResponsePair) (models.RequestResponsePair, error)
	IsMiddlewareSet() bool
}

type SynthesizeMode struct {
	Hoverfly HoverflySynthesize
}

func (this SynthesizeMode) Process(request *http.Request, details models.RequestDetails) (*http.Response, error) {
	pair := models.RequestResponsePair{Request: details}

	log.WithFields(log.Fields{
		// "middleware":  this.hoverfly.Cfg.Middleware.toString(),
		"body":        details.Body,
		"destination": details.Destination,
	}).Debug("Synthesizing new response")

	if !this.Hoverfly.IsMiddlewareSet() {
		err := errors.New("Middleware not set")
		return errorResponse(request, err, "There was an error when creating a synthetic response"), err
	}

	pair, err := this.Hoverfly.ApplyMiddleware(pair)
	if err != nil {
		return errorResponse(request, err, "There was an error when creating a synthetic response"), err
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
