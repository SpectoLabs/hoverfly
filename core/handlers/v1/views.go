package v1

import (
	"bytes"
	"encoding/json"
)

// recordedRequests struct encapsulates payload data
type StoredMetadata struct {
	Data map[string]string `json:"data"`
}

type SetMetadata struct {
	Key   string `json:"key"`
	Value string `json:"value"`
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

type ResponseDelayView struct {
	UrlPattern string `json:"urlPattern"`
	HttpMethod string `json:"httpMethod"`
	Delay      int    `json:"delay"`
}

type ResponseDelayPayloadView struct {
	Data []ResponseDelayView `json:"data"`
}
