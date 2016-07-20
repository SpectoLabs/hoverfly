package matching

import (
	"github.com/SpectoLabs/hoverfly/core/cache"
	"net/http"
	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/models"
	"io/ioutil"
	"bytes"
)

type RequestMatcher struct {
	RequestCache	cache.Cache
	TemplateStore	RequestTemplateStore
	Webserver	*bool

}

// getResponse returns stored response from cache
func (this *RequestMatcher) GetPayload(req *http.Request) (*models.Payload, *MatchingError) {

	if req.Body == nil {
		req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("")))
	}

	reqBody, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Got error when reading request body")
	}

	key := GetRequestFingerprint(req, reqBody, *this.Webserver)

	payloadBts, err := this.RequestCache.Get([]byte(key))

	if err != nil {
		log.WithFields(log.Fields{
			"key":         key,
			"error":       err.Error(),
			"query":       req.URL.RawQuery,
			"path":        req.URL.RawPath,
			"destination": req.Host,
			"method":      req.Method,
		}).Warn("Failed to retrieve response from cache")

		payload, err := this.TemplateStore.GetPayload()
		if err != nil {
			return nil, &MatchingError{
				StatusCode: 412,
				Description: "Could not find recorded request, please record it first!",
			}
		}
		return payload, nil
	}

	// getting cache response
	payload, err := models.NewPayloadFromBytes(payloadBts)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"value": string(payloadBts),
			"key":   key,
		}).Error("Failed to decode payload")
		return nil, &MatchingError{
			StatusCode: 500,
			Description: "Failed to decode payload",
		}
	}

	log.WithFields(log.Fields{
		"key":         key,
		"path":        req.URL.Path,
		"rawQuery":    req.URL.RawQuery,
		"method":      req.Method,
		"destination": req.Host,
		"status":      payload.Response.Status,
	}).Info("Payload found from cache")

	return payload, nil
}

type MatchingError struct {
	StatusCode int
	Description string
}

func (this MatchingError) Error() (string) {
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