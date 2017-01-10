package modes

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

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
	DoRequest(*http.Request) (*http.Request, *http.Response, error)
	IsMiddlewareSet() bool
	Save(*models.RequestDetails, *models.ResponseDetails)
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

func errorResponse(req *http.Request, err error, msg string) *http.Response {
	return goproxy.NewResponse(req,
		goproxy.ContentTypeText, http.StatusBadGateway,
		fmt.Sprintf("Hoverfly Error! \n\n%s\n\nGot error: %s", msg, err.Error()))
}
