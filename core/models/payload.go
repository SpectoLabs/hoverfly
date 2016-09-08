package models

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	log "github.com/Sirupsen/logrus"
	. "github.com/SpectoLabs/hoverfly/core/util"
	"github.com/SpectoLabs/hoverfly/core/views"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify/xml"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const (
	contentTypeJSON = "application/json"
	contentTypeXML  = "application/xml"
	otherType       = "otherType"
)

var (
	rxJSON = regexp.MustCompile("[/+]json$")
	rxXML  = regexp.MustCompile("[/+]xml$")
	// mime types which will not be base 64 encoded when exporting as JSON
	supportedMimeTypes = [...]string{"text", "plain", "css", "html", "json", "xml", "js", "javascript"}
	minifiers          *minify.M
)

func init() {
	// GetNewMinifiers - sets minify.M with prepared xml/json minifiers
	minifiers = minify.New()
	minifiers.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
	minifiers.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
}

// Payload structure holds request and response structure
type RequestResponsePair struct {
	Response ResponseDetails `json:"response"`
	Request  RequestDetails  `json:"request"`
}

func (this RequestResponsePair) Id() string {
	return this.Request.Hash()
}

func (this RequestResponsePair) IdWithoutHost() string {
	return this.Request.HashWithoutHost()
}

// Encode method encodes all exported Payload fields to bytes
func (this *RequestResponsePair) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(this)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (this *RequestResponsePair) ConvertToRequestResponsePairView() *views.RequestResponsePairView {
	return &views.RequestResponsePairView{Response: this.Response.ConvertToResponseDetailsView(), Request: this.Request.ConvertToRequestDetailsView()}
}

// NewPayloadFromBytes decodes supplied bytes into Payload structure
func NewRequestResponsePairFromBytes(data []byte) (*RequestResponsePair, error) {
	var pair *RequestResponsePair
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&pair)
	if err != nil {
		return nil, err
	}
	return pair, nil
}

func NewRequestResponsePairFromRequestResponsePairView(pairView views.RequestResponsePairView) RequestResponsePair {
	return RequestResponsePair{
		Response: NewResponseDetailsFromResponseDetailsView(pairView.Response),
		Request:  NewRequestDetailsFromRequestDetailsView(pairView.Request),
	}
}

// RequestDetails stores information about request, it's used for creating unique hash and also as a payload structure
type RequestDetails struct {
	Path        string              `json:"path"`
	Method      string              `json:"method"`
	Destination string              `json:"destination"`
	Scheme      string              `json:"scheme"`
	Query       string              `json:"query"`
	Body        string              `json:"body"`
	Headers     map[string][]string `json:"headers"`
}

func NewRequestDetailsFromHttpRequest(req *http.Request) (RequestDetails, error) {
	if req.Body == nil {
		req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("")))
	}

	reqBody, err := extractRequestBody(req)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"mode":  "capture",
		}).Error("Got error while reading request body")
		return RequestDetails{}, err
	}

	requestDetails := RequestDetails{
		Path:        req.URL.Path,
		Method:      req.Method,
		Destination: req.Host,
		Scheme:      req.URL.Scheme,
		Query:       req.URL.RawQuery,
		Body:        string(reqBody),
		Headers:     req.Header,
	}
	return requestDetails, nil
}

func extractRequestBody(req *http.Request) (extract []byte, err error) {
	save := req.Body
	savecl := req.ContentLength

	save, req.Body, err = CopyBody(req.Body)
	if err != nil {
		return
	}

	defer req.Body.Close()
	extract, err = ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	req.Body = save
	req.ContentLength = savecl
	return extract, nil
}

func CopyBody(body io.ReadCloser) (resp1, resp2 io.ReadCloser, err error) {
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(body); err != nil {
		return nil, nil, err
	}
	if err = body.Close(); err != nil {
		return nil, nil, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

func NewRequestDetailsFromRequestDetailsView(data views.RequestDetailsView) RequestDetails {
	return RequestDetails{
		Path:        PointerToString(data.Path),
		Method:      PointerToString(data.Method),
		Destination: PointerToString(data.Destination),
		Scheme:      PointerToString(data.Scheme),
		Query:       PointerToString(data.Query),
		Body:        PointerToString(data.Body),
		Headers:     data.Headers,
	}
}

func (this *RequestDetails) ConvertToRequestDetailsView() views.RequestDetailsView {
	s := "recording"
	return views.RequestDetailsView{
		RequestType: &s,
		Path:        &this.Path,
		Method:      &this.Method,
		Destination: &this.Destination,
		Scheme:      &this.Scheme,
		Query:       &this.Query,
		Body:        &this.Body,
		Headers:     this.Headers,
	}
}

func (r *RequestDetails) concatenate(withHost bool) string {
	var buffer bytes.Buffer

	if withHost {
		buffer.WriteString(r.Destination)
	}

	buffer.WriteString(r.Path)
	buffer.WriteString(r.Method)
	buffer.WriteString(r.Query)
	if len(r.Body) > 0 {
		ct := r.getContentType()

		if ct == contentTypeJSON || ct == contentTypeXML {
			buffer.WriteString(r.minifyBody(ct))
		} else {
			log.WithFields(log.Fields{
				"content-type": r.Headers["Content-Type"],
			}).Debug("unknown content type")

			buffer.WriteString(r.Body)
		}
	}

	return buffer.String()
}

func (r *RequestDetails) minifyBody(mediaType string) (minified string) {
	var err error
	minified, err = minifiers.String(mediaType, r.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"error":       err.Error(),
			"destination": r.Destination,
			"path":        r.Path,
			"method":      r.Method,
		}).Errorf("failed to minify request body, media type given: %s. Request matching might fail", mediaType)
		return r.Body
	}
	log.Debugf("body minified, mediatype: %s", mediaType)
	return minified
}

func (r *RequestDetails) Hash() string {
	h := md5.New()
	io.WriteString(h, r.concatenate(true))
	return fmt.Sprintf("%x", h.Sum(nil))
}
func (r *RequestDetails) HashWithoutHost() string {
	h := md5.New()
	io.WriteString(h, r.concatenate(false))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (r *RequestDetails) getContentType() string {
	for _, v := range r.Headers["Content-Type"] {
		if rxJSON.MatchString(v) {
			return contentTypeJSON
		}
		if rxXML.MatchString(v) {
			return contentTypeXML
		}
	}
	return otherType
}

// ResponseDetails structure hold response body from external service, body is not decoded and is supposed
// to be bytes, however headers should provide all required information for later decoding
// by the client.
type ResponseDetails struct {
	Status  int                 `json:"status"`
	Body    string              `json:"body"`
	Headers map[string][]string `json:"headers"`
}

func NewResponseDetailsFromResponseDetailsView(data views.ResponseDetailsView) ResponseDetails {
	body := data.Body

	if data.EncodedBody == true {
		decoded, _ := base64.StdEncoding.DecodeString(data.Body)
		body = string(decoded)
	}

	return ResponseDetails{Status: data.Status, Body: body, Headers: data.Headers}
}

func (r *ResponseDetails) ConvertToResponseDetailsView() views.ResponseDetailsView {
	needsEncoding := false

	// Check headers for gzip
	contentEncodingValues := r.Headers["Content-Encoding"]
	if len(contentEncodingValues) > 0 {
		needsEncoding = true
	} else {
		mimeType := http.DetectContentType([]byte(r.Body))
		needsEncoding = true
		for _, v := range supportedMimeTypes {
			if strings.Contains(mimeType, v) {
				needsEncoding = false
				break
			}
		}
	}

	// If contains gzip, base64 encode
	body := r.Body
	if needsEncoding {
		body = base64.StdEncoding.EncodeToString([]byte(r.Body))
	}

	return views.ResponseDetailsView{Status: r.Status, Body: body, Headers: r.Headers, EncodedBody: needsEncoding}
}
