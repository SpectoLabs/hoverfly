package v1

import (
	"bytes"
	"encoding/json"
	"github.com/SpectoLabs/hoverfly/core/metrics"
)

// recordedRequests struct encapsulates payload data
type StoredMetadata struct {
	Data map[string]string `json:"data"`
}

type SetMetadata struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type RecordsCount struct {
	Count int `json:"count"`
}

type StatsResponse struct {
	Stats        metrics.Stats `json:"stats"`
	RecordsCount int           `json:"recordsCount"`
}

type StateRequest struct {
	Mode        string `json:"mode"`
	Destination string `json:"destination"`
}

type MiddlewareSchema struct {
	Middleware string `json:"middleware"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

func (m *MessageResponse) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(m)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type RequestTemplateResponsePairView struct {
	RequestTemplate RequestTemplateView `json:"requestTemplate"`
	Response        ResponseDetailsView `json:"response"`
}

type RequestTemplateResponsePairPayload struct {
	Data *[]RequestTemplateResponsePairView `json:"data"`
}

type RequestTemplateView struct {
	Path        *string             `json:"path"`
	Method      *string             `json:"method"`
	Destination *string             `json:"destination"`
	Scheme      *string             `json:"scheme"`
	Query       *string             `json:"query"`
	Body        *string             `json:"body"`
	Headers     map[string][]string `json:"headers"`
}

type ResponseDelayView struct {
	UrlPattern string `json:"urlPattern"`
	HttpMethod string `json:"httpMethod"`
	Delay      int    `json:"delay"`
}

type ResponseDelayPayloadView struct {
	Data []ResponseDelayView `json:"data"`
}

type RequestResponsePairPayload struct {
	Data []RequestResponsePairView `json:"data"`
}

// PayloadView is used when marshalling and unmarshalling payloads.
type RequestResponsePairView struct {
	Response ResponseDetailsView `json:"response"`
	Request  RequestDetailsView  `json:"request"`
}

// RequestDetailsView is used when marshalling and unmarshalling RequestDetails
type RequestDetailsView struct {
	RequestType *string             `json:"requestType"`
	Path        *string             `json:"path"`
	Method      *string             `json:"method"`
	Destination *string             `json:"destination"`
	Scheme      *string             `json:"scheme"`
	Query       *string             `json:"query"`
	Body        *string             `json:"body"`
	Headers     map[string][]string `json:"headers"`
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

//Gets Status - required for interfaces.Response
func (this ResponseDetailsView) GetStatus() int { return this.Status }

// Gets Body - required for interfaces.Response
func (this ResponseDetailsView) GetBody() string { return this.Body }

// Gets EncodedBody - required for interfaces.Response
func (this ResponseDetailsView) GetEncodedBody() bool { return this.EncodedBody }

// Gets Headers - required for interfaces.Response
func (this ResponseDetailsView) GetHeaders() map[string][]string { return this.Headers }
