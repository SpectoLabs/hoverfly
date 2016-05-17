package main

import (
	"fmt"
	"github.com/dghubble/sling"
	"net/http"
)

type Hoverfly struct {
	Host       string
	AdminPort  string
	ProxyPort  string
	httpClient *http.Client
}

func (h *Hoverfly) WipeDatabase() int {
	url := fmt.Sprintf("http://%v:%v/api/records", h.Host, h.AdminPort)
	request, _ := sling.New().Delete(url).Request()
	response, _ := h.httpClient.Do(request)
	defer response.Body.Close()
	return response.StatusCode
}