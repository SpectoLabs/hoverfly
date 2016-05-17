package main

import (
	"fmt"
	"github.com/dghubble/sling"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"errors"
	"strings"
)

type ApiStateResponse struct {
	Mode        string `json:"mode"`
	Destination string `json"destination"`
}


type Hoverfly struct {
	Host       string
	AdminPort  string
	ProxyPort  string
	httpClient *http.Client
}

func (h *Hoverfly) WipeDatabase() error {
	url := h.buildUrl("/api/records")
	request, _ := sling.New().Delete(url).Request()
	response, _ := h.httpClient.Do(request)
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Hoverfly did not wipe the database")
	}

	return nil
}

func (h *Hoverfly) GetMode() (string, error) {
	url := h.buildUrl("/api/state")
	request, err := sling.New().Get(url).Request()

	if err != nil {
		return "", err
	}

	response, err := h.httpClient.Do(request)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	apiResponse := h.createApiStateResponse(response)

	return apiResponse.Mode, nil
}

func (h *Hoverfly) SetMode(mode string) (string, error) {
	if mode != "simulate" && mode != "capture" && mode != "modify" && mode != "synthesize" {
		return "", errors.New(mode + " is not a valid mode")
	}

	url := h.buildUrl("/api/state")
	request, err := sling.New().Post(url).Body(strings.NewReader(`{"mode":"` + mode + `"}`)).Request()

	if err != nil {
		return "", err
	}

	response, err := h.httpClient.Do(request)

	if err != nil {
		return "", err
	}

	apiResponse := h.createApiStateResponse(response)

	return apiResponse.Mode, nil

}

func (h *Hoverfly) createApiStateResponse(response *http.Response) ApiStateResponse {
	body, _ := ioutil.ReadAll(response.Body)
	var apiResponse ApiStateResponse
	json.Unmarshal(body, &apiResponse)
	return apiResponse
}

func (h * Hoverfly) buildUrl(endpoint string) string {
	return fmt.Sprintf("%v%v", h.buildBaseUrl(), endpoint)
}

func (h * Hoverfly) buildBaseUrl() string {
	return fmt.Sprintf("http://%v:%v", h.Host, h.AdminPort)
}