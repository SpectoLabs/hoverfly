package util

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	log "github.com/sirupsen/logrus"
)

// ParseJWTComposite builds a JSON string: {"header":{...},"payload":{...}}
// Does NOT verify signature. Skips sections that fail to decode.
func ParseJWTComposite(raw string) (string, error) {
	token := strings.TrimSpace(raw)
	if token == "" {
		return "", ErrInvalidJWT("empty token")
	}
	lower := strings.ToLower(token)
	if strings.HasPrefix(lower, "bearer ") {
		token = strings.TrimSpace(token[7:])
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", ErrInvalidJWT("token must have 3 sections")
	}

	composite := make(map[string]interface{})
	if h, err := decodeSegment(parts[0]); err == nil {
		composite["header"] = h
	} else {
		log.Error("failed to decode jwt header: ", err)
	}
	if p, err := decodeSegment(parts[1]); err == nil {
		composite["payload"] = p
	} else {
		log.Error("failed to decode jwt payload: ", err)
	}

	bytes, err := json.Marshal(composite)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func decodeSegment(seg string) (interface{}, error) {
	bytes, err := base64.RawURLEncoding.DecodeString(seg)
	if err != nil {
		return nil, err
	}
	var out interface{}
	if err := json.Unmarshal(bytes, &out); err != nil {
		return nil, err
	}
	return out, nil
}

type ErrInvalidJWT string

func (e ErrInvalidJWT) Error() string { return string(e) }
