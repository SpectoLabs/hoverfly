package models

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/interfaces"
	"github.com/SpectoLabs/hoverfly/core/util"
)

var (
	// mime types which will not be base 64 encoded when exporting as JSON
	supportedMimeTypes = [...]string{"text", "plain", "css", "html", "json", "xml", "js", "javascript"}
)

// Payload structure holds request and response structure
type RequestResponsePair struct {
	Response ResponseDetails
	Request  RequestDetails
}

func (this *RequestResponsePair) ConvertToRequestResponsePairView() v2.RequestResponsePairViewV1 {
	return v2.RequestResponsePairViewV1{Response: this.Response.ConvertToResponseDetailsView(), Request: this.Request.ConvertToRequestDetailsView()}
}

func NewRequestResponsePairFromRequestResponsePairView(pairView interfaces.RequestResponsePair) RequestResponsePair {
	return RequestResponsePair{
		Response: NewResponseDetailsFromResponse(pairView.GetResponse()),
		Request:  NewRequestDetailsFromRequest(pairView.GetRequest()),
	}
}

func NewRequestDetailsFromRequest(data interfaces.Request) RequestDetails {
	query, _ := url.ParseQuery(*data.GetQuery())
	return RequestDetails{
		Path:        util.PointerToString(data.GetPath()),
		Method:      util.PointerToString(data.GetMethod()),
		Destination: util.PointerToString(data.GetDestination()),
		Scheme:      util.PointerToString(data.GetScheme()),
		Query:       query,
		Body:        util.PointerToString(data.GetBody()),
		Headers:     data.GetHeaders(),
	}
}

// RequestDetails stores information about request, it's used for creating unique hash and also as a payload structure
type RequestDetails struct {
	Path        string
	Method      string
	Destination string
	Scheme      string
	Query       map[string][]string
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
		Query:       req.URL.Query(),
		Body:        string(reqBody),
		Headers:     req.Header,
	}
	return requestDetails, nil
}

func (this *RequestDetails) ConvertToRequestDetailsView() v2.RequestDetailsView {
	queryString := this.QueryString()

	return v2.RequestDetailsView{
		Path:        &this.Path,
		Method:      &this.Method,
		Destination: &this.Destination,
		Scheme:      &this.Scheme,
		Query:       &queryString,
		Body:        &this.Body,
		Headers:     this.Headers,
	}
}

// TODO: Remove this
// This only exists as there are parts of Hoverfly that still
// require the request query parameters to be a string and not
// a map
func (this *RequestDetails) QueryString() string {
	var buf bytes.Buffer
	keys := make([]string, 0, len(this.Query))
	for k := range this.Query {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := this.Query[k]
		prefix := k + "="
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(v)
		}
	}
	return util.SortQueryString(buf.String())
}

func (r *RequestDetails) concatenate(withHost bool) string {
	var buffer bytes.Buffer

	if withHost {
		buffer.WriteString(r.Destination)
	}

	buffer.WriteString(r.Path)
	buffer.WriteString(r.Method)
	buffer.WriteString(r.QueryString())
	if len(r.Body) > 0 {
		buffer.WriteString(r.Body)
	}

	return buffer.String()
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

// ResponseDetails structure hold response body from external service, body is not decoded and is supposed
// to be bytes, however headers should provide all required information for later decoding
// by the client.
type ResponseDetails struct {
	Status    int
	Body      string
	Headers   map[string][]string
	Templated bool
}

func NewResponseDetailsFromResponse(data interfaces.Response) ResponseDetails {
	body := data.GetBody()

	if data.GetEncodedBody() == true {
		decoded, _ := base64.StdEncoding.DecodeString(data.GetBody())
		body = string(decoded)
	}

	return ResponseDetails{Status: data.GetStatus(), Body: body, Headers: data.GetHeaders(), Templated: data.GetTemplated()}
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

	return v2.ResponseDetailsView{
		Status:      r.Status,
		Body:        body,
		Headers:     r.Headers,
		EncodedBody: needsEncoding,
	}
}

func (r *ResponseDetails) ConvertToResponseDetailsViewV3() v2.ResponseDetailsViewV3 {
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

	return v2.ResponseDetailsViewV3{
		Status:      r.Status,
		Body:        body,
		Headers:     r.Headers,
		EncodedBody: needsEncoding,
		Templated:   r.Templated,
	}
}
