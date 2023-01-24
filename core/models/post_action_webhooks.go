package models

import (
	"bytes"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type PostActionHook struct {
	Name       string                `json:"name"`
	Parameters *ActionHookParameters `json:"parameters"`
}

type ActionHookParameters struct {
	Method      string                 `json:"method"`
	Destination string                 `json:"destination"`
	Scheme      string                 `json:"scheme"`
	Path        string                 `json:"path"`
	Query       map[string]interface{} `json:"query"`
	Headers     map[string][]string    `json:"headers"`
	Body        string                 `json:"body"`
	EncodedBody bool                   `json:"encodedBody"`
}

type PostActionHooks []PostActionHook

func (postActionHook *PostActionHook) Execute() (*http.Response, error) {
	bodyBytes := []byte(postActionHook.Parameters.Body)
	url := fmt.Sprintf("%s://%s%s", postActionHook.Parameters.Scheme, postActionHook.Parameters.Destination, postActionHook.Parameters.Path)
	req, err := http.NewRequest(postActionHook.Parameters.Method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.WithFields(log.Fields{
			"destination":  url,
			"method":       postActionHook.Parameters.Method,
			"postHookName": postActionHook.Name,
			"scheme":       postActionHook.Parameters.Scheme,
		}).Error("Failed to create request for post action hook")
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"destination":  url,
			"method":       postActionHook.Parameters.Method,
			"postHookName": postActionHook.Name,
			"scheme":       postActionHook.Parameters.Scheme,
		}).Error("Failed to send request for post action hook")
		return nil, err
	}
	log.Debug("Successfully attempted request for post action hook")
	return res, nil
}
