package hoverfly

import (
	"bytes"
	"encoding/json"
)

type messageResponse struct {
	Message string `json:"message"`
}

func (m *messageResponse) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(m)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}