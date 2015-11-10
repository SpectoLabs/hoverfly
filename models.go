package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

type DBClient struct {
	cache Cache
	http *http.Client
}

// request holds structure for request
type request struct {
	URL    string
	Method string
}

// hash returns unique hash key for request
// TODO: match on destination, request body, etc..
func (r *request) hash() string {
	h := md5.New()
	io.WriteString(h, fmt.Sprintf("%s%s", r.URL, r.Method))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// res structure hold response body from external service, body is not decoded and is supposed
// to be bytes, however headers should provide all required information for later decoding
// by the client.
type res struct {
	Status  int               `json:"status"`
	Body    []byte            `json:"body"`
	Headers map[string]string `json:"headers"`
}

// recordRequest saves request for later playback
func (d *DBClient) recordRequest(req *http.Request) (*http.Response, error) {

	// forwarding request
	resp, err := d.doRequest(req)

	// record request here
	key := getRequestFingerprint(req)

	response := res{
		Status:  resp.StatusCode,
		Body:    resp.Body,
		Headers: getHeadersMap(resp.Header),
	}

	log.WithFields(log.Fields{
		"requestURL":    req.URL.Path,
		"requestMethod": req.Method,
		"hashKey":       key,
	}).Info("Recording")


	// return new response or error here
	return resp, err


}

// getRequestFingerprint returns request hash
func getRequestFingerprint(req *http.Request) string {
	r := request{URL: req.URL.Path, Method: req.Method}
	return r.hash()
}

// getHeadersMap converts map[string][]string to map[string]string structure
func getHeadersMap(hds map[string][]string) map[string]string {
	headers := make(map[string]string)
	for key, value := range hds {
		headers[key] = value[0]
	}
	return headers
}

func (d *DBClient) getResponse(r *http.Request) *http.Response {
	log.Info("Returning response")
	return nil
}

// doRequest performs original request and returns response that should be returned to client and error (if there is one)
func (d *DBClient) doRequest(request *http.Request) (*http.Response, error) {
	// We can't have this set. And it only contains "/pkg/net/http/" anyway
	request.RequestURI = ""

	resp, err := d.http.Do(request)

	if err != nil {
		log.WithFields(log.Fields{
			"error":  err.Error(),
			"host":   request.Host,
			"method": request.Method,
			"path":   request.URL.Path,
		}).Error("Could not forward request.")
		return nil, err
	}

	log.WithFields(log.Fields{}).Info("Request forwarded!")

	resp.Header.Set("Gen-proxy", "Was-Here")
	return resp, nil

}

/