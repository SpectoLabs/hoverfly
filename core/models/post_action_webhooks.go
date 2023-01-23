package models

import (
	"bytes"
	"fmt"
	"net/http"
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
}

type PostActionHooks []PostActionHook

func (postActionHook *PostActionHook) Execute() (*http.Request, error) {
	bodyBytes := []byte(postActionHook.Parameters.Body)
	url := fmt.Sprintf("%s://%s%s", postActionHook.Parameters.Scheme, postActionHook.Parameters.Destination, postActionHook.Parameters.Path)
	req, err := http.NewRequest(postActionHook.Parameters.Method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	return req, nil
}
