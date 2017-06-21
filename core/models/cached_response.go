package models

import (
	"bytes"
	"encoding/gob"
)

type CachedResponse struct {
	Request      RequestDetails
	MatchingPair *RequestMatcherResponsePair
	ClosestMiss  *ClosestMiss
}

func (this CachedResponse) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(this)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func NewCachedResponseFromBytes(data []byte) (*CachedResponse, error) {
	var cachedResponse *CachedResponse
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&cachedResponse)
	if err != nil {
		return nil, err
	}
	return cachedResponse, nil
}
