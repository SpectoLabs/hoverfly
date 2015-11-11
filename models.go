package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

type DBClient struct {
	cache Cache
	http  *http.Client
}

// request holds structure for request
type request struct {
	details requestDetails
}

type requestDetails struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	Destination string `json:"destination"`
	Query       string `json:"query"`
}

// hash returns unique hash key for request
// TODO: match on destination, request body, etc..
func (r *request) hash() string {
	h := md5.New()
	io.WriteString(h, fmt.Sprintf("%s%s%s%s", r.details.Destination, r.details.Path, r.details.Method, r.details.Query))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// res structure hold response body from external service, body is not decoded and is supposed
// to be bytes, however headers should provide all required information for later decoding
// by the client.
type response struct {
	Status  int               `json:"status"`
	Body    []byte            `json:"body"`
	Headers map[string]string `json:"headers"`
}

// Payload structure holds request and response structure
type Payload struct {
	Response response       // `json:"response"`
	Request  requestDetails // `json:"request"`
}

// recordRequest saves request for later playback
func (d *DBClient) recordRequest(req *http.Request) (*http.Response, error) {

	// forwarding request
	resp, err := d.doRequest(req)

	go d.save(req, resp)

	// return new response or error here
	return resp, err
}

// save gets request fingerprint, extracts request body, status code and headers, then saves it to cache
func (d *DBClient) save(req *http.Request, resp *http.Response) {
	// record request here
	key := getRequestFingerprint(req)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	responseObj := response{
		Status:  resp.StatusCode,
		Body:    body,
		Headers: getHeadersMap(resp.Header),
	}

	log.WithFields(log.Fields{
		"path":          req.URL.Path,
		"rawQuery":      req.URL.RawQuery,
		"requestMethod": req.Method,
		"destination":   req.Host,
		"hashKey":       key,
	}).Info("Recording")

	requestObj := requestDetails{
		Path:        req.URL.Path,
		Method:      req.Method,
		Destination: req.Host,
		Query:       req.URL.RawQuery}

	payload := Payload{
		Response: responseObj,
		Request:  requestObj,
	}
	// converting it to json bytes
	bts, err := json.Marshal(payload)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to marshal json")
	} else {
		d.cache.set(key, bts)
	}

}

// getRequestFingerprint returns request hash
func getRequestFingerprint(req *http.Request) string {
	details := requestDetails{Path: req.URL.Path, Method: req.Method, Destination: req.Host, Query: req.URL.RawQuery}
	r := request{details: details}
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
