package modes

import (
	"bytes"
	"fmt"
	"github.com/SpectoLabs/hoverfly/core/util"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/SpectoLabs/goproxy"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/sirupsen/logrus"
)

// SimulateMode - default mode when Hoverfly looks for captured requests to respond
const Simulate = "simulate"

// SynthesizeMode - all requests are sent to middleware to create response
const Synthesize = "synthesize"

// ModifyMode - middleware is applied to outgoing and incoming traffic
const Modify = "modify"

// CaptureMode - requests are captured and stored in cache
const Capture = "capture"

// SpyMode - simulateMode but will call real service when cache miss
const Spy = "spy"

// DiffMode - calls real service and compares response with simulation
const Diff = "diff"

type Mode interface {
	Process(*http.Request, models.RequestDetails) (ProcessResult, error)
	SetArguments(arguments ModeArguments)
	View() v2.ModeView
}

type ModeArguments struct {
	Headers            []string
	MatchingStrategy   *string
	Stateful           bool
	OverwriteDuplicate bool
}

type ProcessResult struct {
	Response       *http.Response
	FixedDelay     int
	LogNormalDelay *models.ResponseDetailsLogNormal
}

func (p ProcessResult) IsResponseDelayable() bool {
	return p.FixedDelay > 0 || p.LogNormalDelay != nil
}

func newProcessResult(response *http.Response, fixedDelay int, logNormalDelay *models.ResponseDetailsLogNormal) ProcessResult {
	return ProcessResult{Response: response, FixedDelay: fixedDelay, LogNormalDelay: logNormalDelay}
}

// ReconstructRequest replaces original request with details provided in Constructor Payload.RequestMatcher
func ReconstructRequest(pair models.RequestResponsePair) (*http.Request, error) {
	if pair.Request.Destination == "" {
		return nil, fmt.Errorf("failed to reconstruct request, destination not specified")
	}

	bodyBytes := []byte(pair.Request.Body)
	// recompress the request body if the original was encoded
	if values, found := pair.Request.Headers["Content-Encoding"]; found {
		for _, value := range values {
			// Only gzip is supported at the moment
			if value == "gzip" {
				compressedBody, err := util.CompressGzip(bodyBytes)
				if err == nil {
					bodyBytes = compressedBody
				} else {
					// Fail to compress, we should remove the encoding header
					delete(pair.Request.Headers, "Content-Encoding")
				}
				break
			}
		}
	}

	newRequest, err := http.NewRequest(
		pair.Request.Method,
		fmt.Sprintf("%s://%s%s", pair.Request.Scheme, pair.Request.Destination, pair.Request.Path),
		bytes.NewBuffer(bodyBytes))

	if err != nil {
		return nil, err
	}


	newRequest.Method = pair.Request.Method
	newRequest.Header = pair.Request.Headers

	if pair.Request.GetRawQuery() == "" {
		// rawQuery is empty if middleware is applied, as unexported fields are not marshal, hence re-encoding of the query params is needed here
		t := &url.URL{Path: pair.Request.QueryString()}
		newRequest.URL.RawQuery = t.String()
	} else {
		// otherwise we use the original raw query for pass-through
		newRequest.URL.RawQuery = pair.Request.GetRawQuery()
	}

	return newRequest, nil
}

// ReconstructResponse changes original response with details provided in Constructor Payload.Response
func ReconstructResponse(request *http.Request, pair models.RequestResponsePair) *http.Response {
	response := &http.Response{}
	response.Request = request

	response.ContentLength = int64(len(pair.Response.Body))
	response.Body = ioutil.NopCloser(strings.NewReader(pair.Response.Body))
	response.StatusCode = pair.Response.Status
	response.Status = http.StatusText(pair.Response.Status)

	headers := make(http.Header)

	// Make copy to prevent modifying the simulation
	for k, v := range pair.Response.Headers {
		headers[k] = v
	}

	if keys, present := headers["Trailer"]; present {
		response.Trailer = make(http.Header)
		for _, key := range keys {
			response.Trailer[key] = headers[key]
			delete(headers, key)
		}
		delete(headers, "Trailer")
	}

	response.Header = headers

	if response.ContentLength > 0 && response.Header.Get("Content-Length") == "" && response.Header.Get("Transfer-Encoding") == "" {
		response.Header.Set("Content-Length", fmt.Sprintf("%v", response.ContentLength))
	}

	return response
}

func GetRequestLogFields(request *models.RequestDetails) *logrus.Fields {
	if request == nil {
		return &log.Fields{
			"error": "nil request",
		}
	}

	return &log.Fields{
		"method":      request.Method,
		"scheme":      request.Scheme,
		"destination": request.Destination,
		"path":        request.Path,
		"query":       request.Query,
		"headers":     request.Headers,
		"body":        request.Body,
	}
}

func GetResponseLogFields(response *models.ResponseDetails) *logrus.Fields {
	if response == nil || response.Status == 0 {
		return &log.Fields{
			"error": "nil response",
		}
	}

	return &log.Fields{
		"body":    response.Body,
		"headers": response.Headers,
		"status":  response.Status,
	}
}

func ReturnErrorAndLog(request *http.Request, err error, pair *models.RequestResponsePair, msg, mode string) (ProcessResult, error) {
	log.WithFields(log.Fields{
		"error":    err.Error(),
		"mode":     mode,
		"request":  GetRequestLogFields(&pair.Request),
		"response": GetResponseLogFields(&pair.Response),
	}).Error(msg)

	return ErrorResponse(request, err, msg), err
}

func ErrorResponse(req *http.Request, err error, msg string) ProcessResult {
	return newProcessResult(goproxy.NewResponse(req,
		goproxy.ContentTypeText, http.StatusBadGateway,
		fmt.Sprintf("Hoverfly Error!\n\n%s\n\nGot error: %s", msg, err.Error())), 0, nil)
}
