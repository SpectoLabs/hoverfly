package v1

import (
	"bytes"
	"encoding/json"
	"github.com/SpectoLabs/hoverfly/core/metrics"
	"github.com/SpectoLabs/hoverfly/core/views"
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
	RequestTemplate RequestTemplateView       `json:"requestTemplate"`
	Response        views.ResponseDetailsView `json:"response"`
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
	Data *[]ResponseDelayView `json:"data"`
}
