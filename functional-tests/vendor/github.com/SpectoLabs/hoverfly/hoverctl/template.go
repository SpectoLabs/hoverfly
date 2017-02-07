package main

import (
	"encoding/json"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/dghubble/sling"
	"io/ioutil"
	"net/http"
	"strings"
)

func (h *Hoverfly) performAPIRequest(slingRequest *sling.Sling) (*http.Response, error) {
	slingRequest, err := h.addAuthIfNeeded(slingRequest)
	if err != nil {
		log.Warn(err.Error())
		return nil, errors.New("Could not authenticate  with Hoverfly")
	}

	request, err := slingRequest.Request()

	if err != nil {
		log.Warn(err.Error())
		return nil, errors.New("Could not communicate with Hoverfly")
	}

	response, err := h.httpClient.Do(request)
	if err != nil {
		log.Warn(err.Error())
		return nil, errors.New("Could not communicate with Hoverfly")
	}
	return response, nil

}

func (h *Hoverfly) GetRequestTemplates() (*matching.RequestTemplateResponsePairPayload, error) {
	url := h.buildURL("/api/templates")
	slingRequest := sling.New().Get(url)
	response, err := h.performAPIRequest(slingRequest)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	requestTemplates, err := unmarshalRequestTemplates(response)
	if err != nil {
		return nil, err
	}

	return requestTemplates, nil
}

func (h *Hoverfly) SetRequestTemplates(path string) (responseTemplates *matching.RequestTemplateResponsePairPayload, err error) {

	conf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	url := h.buildURL("/api/templates")

	slingRequest := sling.New().Post(url).Body(strings.NewReader(string(conf)))
	_, err = h.performAPIRequest(slingRequest)
	if err != nil {
		return nil, err
	}

	slingRequest = sling.New().Get(url).Body(strings.NewReader(string(conf)))
	getResponse, err := h.performAPIRequest(slingRequest)
	if err != nil {
		return nil, err
	}

	requestTemplates, err := unmarshalRequestTemplates(getResponse)
	if err != nil {
		return nil, err
	}

	return requestTemplates, nil
}

func unmarshalRequestTemplates(response *http.Response) (*matching.RequestTemplateResponsePairPayload, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error("Error reading request templates response body: " + err.Error())
		return nil, err
	}

	var requestTemplates matching.RequestTemplateResponsePairPayload

	err = json.Unmarshal(body, &requestTemplates)
	if err != nil {
		log.Error("Error unmarshalling JSON for request templates: " + err.Error())
		return nil, err
	}

	return &requestTemplates, nil
}

func (h *Hoverfly) DeleteRequestTemplates() error {
	url := h.buildURL("/api/templates")

	slingRequest := sling.New().Delete(url)
	response, err := h.performAPIRequest(slingRequest)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode == 401 {
		return errors.New("Hoverfly requires authentication")
	}

	if response.StatusCode != 200 {
		return errors.New("Templates were not deleted from Hoverfly")
	}

	return nil
}
