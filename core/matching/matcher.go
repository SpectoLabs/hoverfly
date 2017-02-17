package matching

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/models"
)

type RequestMatcher struct {
	RequestCache cache.Cache
	Webserver    *bool
	Simulation   *models.Simulation
}

// getResponse returns stored response from cache
func (this *RequestMatcher) GetResponse(req *models.RequestDetails) (*models.ResponseDetails, *MatchingError) {

	var key string

	if *this.Webserver {
		key = req.HashWithoutHost()
	} else {
		key = req.Hash()
	}

	pairBytes, err := this.RequestCache.Get([]byte(key))

	if err != nil {
		log.WithFields(log.Fields{
			"key":         key,
			"error":       err.Error(),
			"query":       req.Query,
			"path":        req.Path,
			"destination": req.Destination,
			"method":      req.Method,
		}).Warn("Failed to retrieve response from cache")

		response, err := TemplateMatcher{}.Match(*req, *this.Webserver, this.Simulation)
		if err != nil {
			log.WithFields(log.Fields{
				"key":         key,
				"error":       err.Error(),
				"query":       req.Query,
				"path":        req.Path,
				"destination": req.Destination,
				"method":      req.Method,
			}).Warn("Failed to find matching request template from template store")

			return nil, &MatchingError{
				StatusCode:  412,
				Description: "Could not find recorded request, please record it first!",
			}
		}
		log.WithFields(log.Fields{
			"key":         key,
			"query":       req.Query,
			"path":        req.Path,
			"destination": req.Destination,
			"method":      req.Method,
		}).Info("Found template matching request from template store")
		return response, nil
	}

	// getting cache response
	pair, err := models.NewRequestResponsePairFromBytes(pairBytes)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"value": string(pairBytes),
			"key":   key,
		}).Error("Failed to decode payload")
		return nil, &MatchingError{
			StatusCode:  500,
			Description: "Failed to decode payload",
		}
	}

	log.WithFields(log.Fields{
		"key":         key,
		"path":        req.Path,
		"rawQuery":    req.Query,
		"method":      req.Method,
		"destination": req.Destination,
		"status":      pair.Response.Status,
	}).Info("Payload found from cache")

	return &pair.Response, nil
}

func (this *RequestMatcher) SaveRequestResponsePair(pair *models.RequestResponsePair) error {
	var key string

	if *this.Webserver {
		key = pair.IdWithoutHost()
	} else {
		key = pair.Id()
	}

	log.WithFields(log.Fields{
		"path":          pair.Request.Path,
		"rawQuery":      pair.Request.Query,
		"requestMethod": pair.Request.Method,
		"bodyLen":       len(pair.Request.Body),
		"destination":   pair.Request.Destination,
		"hashKey":       key,
	}).Debug("Capturing")

	pairBytes, err := pair.Encode()

	if err != nil {
		return err
	}

	return this.RequestCache.Set([]byte(key), pairBytes)
}

type MatchingError struct {
	StatusCode  int
	Description string
}

func (this MatchingError) Error() string {
	return this.Description
}

// getRequestFingerprint returns request hash
func GetRequestFingerprint(req *http.Request, requestBody []byte, webserver bool) string {
	var r models.RequestDetails

	r = models.RequestDetails{
		Path:        req.URL.Path,
		Method:      req.Method,
		Destination: req.Host,
		Query:       req.URL.RawQuery,
		Body:        string(requestBody),
		Headers:     req.Header,
	}

	if webserver {
		return r.HashWithoutHost()
	}

	return r.Hash()
}
