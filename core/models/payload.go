package models

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/interfaces"
	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify/xml"
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
	Response ResponseDetails
	Request  RequestDetails
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

func (this *RequestResponsePair) ConvertToV1RequestResponsePairView() *v1.RequestResponsePairView {
	return &v1.RequestResponsePairView{Response: this.Response.ConvertToV1ResponseDetailsView(), Request: this.Request.ConvertToV1RequestDetailsView()}
}

func (this *RequestResponsePair) ConvertToRequestResponsePairView() v2.RequestResponsePairView {
	return v2.RequestResponsePairView{Response: this.Response.ConvertToResponseDetailsView(), Request: this.Request.ConvertToRequestDetailsView()}
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

func NewRequestResponsePairFromRequestResponsePairView(pairView interfaces.RequestResponsePair) RequestResponsePair {
	return RequestResponsePair{
		Response: NewResponseDetailsFromResponse(pairView.GetResponse()),
		Request:  NewRequestDetailsFromRequest(pairView.GetRequest()),
	}
}

// RequestDetails stores information about request, it's used for creating unique hash and also as a payload structure
type RequestDetails struct {
	Path        string
	Method      string
	Destination string
	Scheme      string
	Query       string
	Body        string
	Headers     map[string][]string
}

func NewRequestDetailsFromHttpRequest(req *http.Request) (RequestDetails, error) {
	if req.Body == nil {
		req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("")))
	}

	reqBody, err := util.GetRequestBody(req)

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
		Destination: strings.ToLower(req.Host),
		Scheme:      req.URL.Scheme,
		Query:       util.SortQueryString(req.URL.RawQuery),
		Body:        string(reqBody),
		Headers:     req.Header,
	}
	return requestDetails, nil
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

func NewRequestDetailsFromRequest(data interfaces.Request) RequestDetails {
	return RequestDetails{
		Path:        util.PointerToString(data.GetPath()),
		Method:      util.PointerToString(data.GetMethod()),
		Destination: util.PointerToString(data.GetDestination()),
		Scheme:      util.PointerToString(data.GetScheme()),
		Query:       util.PointerToString(data.GetQuery()),
		Body:        util.PointerToString(data.GetBody()),
		Headers:     data.GetHeaders(),
	}
}

func (this *RequestDetails) ConvertToV1RequestDetailsView() v1.RequestDetailsView {
	s := "recording"
	return v1.RequestDetailsView{
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

func (this *RequestDetails) ConvertToRequestDetailsView() v2.RequestDetailsView {
	s := "recording"
	return v2.RequestDetailsView{
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
	Status  int
	Body    string
	Headers map[string][]string
}

func NewResponseDetailsFromResponse(data interfaces.Response) ResponseDetails {
	body := data.GetBody()

	if data.GetEncodedBody() == true {
		decoded, _ := base64.StdEncoding.DecodeString(data.GetBody())
		body = string(decoded)
	}

	return ResponseDetails{Status: data.GetStatus(), Body: body, Headers: data.GetHeaders()}
}

// This function will create a JSON appriopriate version of ResponseDetails for the v1 API
// If the response headers indicate that the content is encoded, or it has a non-matching
// supported mimetype, we base64 encode it.
func (r *ResponseDetails) ConvertToV1ResponseDetailsView() v1.ResponseDetailsView {
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

	return v1.ResponseDetailsView{Status: r.Status, Body: body, Headers: r.Headers, EncodedBody: needsEncoding}
}

// This function will create a JSON appriopriate version of ResponseDetails for the v2 API
// If the response headers indicate that the content is encoded, or it has a non-matching
// supported mimetype, we base64 encode it.
func (r *ResponseDetails) ConvertToResponseDetailsView() v2.ResponseDetailsView {
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

	return v2.ResponseDetailsView{Status: r.Status, Body: body, Headers: r.Headers, EncodedBody: needsEncoding}
}
