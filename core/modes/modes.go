package modes

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/rusenask/goproxy"
)

type Mode interface {
	Process(*http.Request, models.RequestDetails) (*http.Response, error)
}

type Hoverfly interface {
	GetResponse(models.RequestDetails) (*models.ResponseDetails, *matching.MatchingError)
	ApplyMiddleware(models.RequestResponsePair) (models.RequestResponsePair, error)
	DoRequest(*http.Request) (*http.Response, error)
	IsMiddlewareSet() bool
	Save(*models.RequestDetails, *models.ResponseDetails)
}

// ReconstructRequest replaces original request with details provided in Constructor Payload.Request
func ReconstructRequest(pair models.RequestResponsePair) (*http.Request, error) {
	if pair.Request.Destination == "" {
		return nil, fmt.Errorf("failed to reconstruct request, destination not specified")
	}

	newRequest, err := http.NewRequest(
		pair.Request.Method,
		fmt.Sprintf("%s://%s", pair.Request.Scheme, pair.Request.Destination),
		ioutil.NopCloser(bytes.NewBuffer([]byte(pair.Request.Body))))

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Request reconstruction failed...")
		return nil, err
	}

	newRequest.Method = pair.Request.Method
	newRequest.URL.Path = pair.Request.Path
	newRequest.URL.RawQuery = pair.Request.Query
	newRequest.Header = pair.Request.Headers

	return newRequest, nil
}

// ReconstructResponse changes original response with details provided in Constructor Payload.Response
func ReconstructResponse(request *http.Request, pair models.RequestResponsePair) *http.Response {
	response := &http.Response{}
	response.Request = request

	// adding headers
	response.Header = make(http.Header)

	// applying payload
	if len(pair.Response.Headers) > 0 {
		for k, values := range pair.Response.Headers {
			// headers is a map, appending each value
			for _, v := range values {
				response.Header.Add(k, v)
			}

		}
	}

	// adding body, length, status code
	buf := bytes.NewBufferString(pair.Response.Body)
	response.ContentLength = int64(buf.Len())
	response.Body = ioutil.NopCloser(buf)
	response.StatusCode = pair.Response.Status

	return response
}

func ErrorResponse(req *http.Request, err error, msg string) *http.Response {
	return goproxy.NewResponse(req,
		goproxy.ContentTypeText, http.StatusBadGateway,
		fmt.Sprintf("Hoverfly Error! \n\n%s\n\nGot error: %s", msg, err.Error()))
}
