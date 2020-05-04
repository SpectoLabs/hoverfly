package v1

import (
	"bytes"
	"encoding/json"
	"github.com/SpectoLabs/hoverfly/core/delay"
)

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

type ResponseDelayLogNormalPayloadView struct {
	Data []ResponseDelayLogNormalView `json:"data"`
}

type ResponseDelayLogNormalView struct {
	UrlPattern string `json:"urlPattern"`
	HttpMethod string `json:"httpMethod"`
	*delay.LogNormalDelayOptions
}
