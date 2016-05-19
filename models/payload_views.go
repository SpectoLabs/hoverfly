package models

import (
	"bytes"
	"encoding/gob"
	"encoding/base64"
)

type PayloadViewData struct {
	Data []PayloadView `json:"data"`
}

// PayloadView is used when marshalling and unmarshalling payloads.
type PayloadView struct {
	Response ResponseDetailsView `json:"response"`
	Request  RequestDetailsView  `json:"request"`
}

func (r *PayloadView) ConvertToPayload() (Payload) {
	return Payload{Response: r.Response.ConvertToResponseDetails(), Request: r.Request.ConvertToRequestDetails()}
}

// Encode method encodes all exported Payload fields to bytes
func (p *PayloadView) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(p)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// RequestDetailsView is used when marshalling and unmarshalling RequestDetails
type RequestDetailsView struct {
	Path        string              `json:"path"`
	Method      string              `json:"method"`
	Destination string              `json:"destination"`
	Scheme      string              `json:"scheme"`
	Query       string              `json:"query"`
	Body        string              `json:"body"`
	Headers     map[string][]string `json:"headers"`
}

func (r *RequestDetailsView) ConvertToRequestDetails() (RequestDetails) {
	return RequestDetails{
		Path: r.Path,
		Method: r.Method,
		Destination: r.Destination,
		Scheme: r.Scheme,
		Query: r.Query,
		Body: r.Body,
		Headers: r.Headers,
	}
}

// ResponseDetailsView is used when marshalling and
// unmarshalling requests. This struct's Body may be Base64
// encoded based on the EncodedBody field.
type ResponseDetailsView struct {
	Status      int                 `json:"status"`
	Body        string              `json:"body"`
	EncodedBody bool                `json:"encodedBody"`
	Headers     map[string][]string `json:"headers"`
}

func (r *ResponseDetailsView) ConvertToResponseDetails() (ResponseDetails) {
	body := r.Body

	if r.EncodedBody == true {
		decoded, _ := base64.StdEncoding.DecodeString(r.Body)
		body = string(decoded)
	}

	return ResponseDetails{Status: r.Status, Body: body, Headers: r.Headers}
}
