package main

import (
	"fmt"
	"github.com/dghubble/sling"
	"net/http"
	"strings"
	"io/ioutil"
)

type SpectoHubSimulation struct {
	Vendor      string `json:"vendor"`
	Api         string `json:"api"`
	Version     string `json:"version"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SpectoHub struct {
	Host   string
	Port   string
	ApiKey string
}

func (s *SpectoHub) CheckSimulation(simulation SpectoHubSimulation) int {
	url := s.buildUrl(fmt.Sprintf("/api/v1/users/%v/vendors/%v/apis/%v/versions/%v/%v", simulation.Vendor, simulation.Vendor, simulation.Api, simulation.Version, simulation.Name))

	request, _ := sling.New().Get(url).Add("Authorization", s.buildAuthorizationHeaderValue()).Request()
	response, _ := http.DefaultClient.Do(request)
	defer response.Body.Close()

	return response.StatusCode
}

func (s *SpectoHub) CreateSimulation(simulation SpectoHubSimulation) int {
	url := s.buildUrl("/api/v1/simulations")

	request, _ := sling.New().Post(url).Add("Authorization", s.buildAuthorizationHeaderValue()).BodyJSON(simulation).Request()
	response, _ := http.DefaultClient.Do(request)
	defer response.Body.Close()
	return response.StatusCode
}

func (s *SpectoHub) UploadSimulation(simulation SpectoHubSimulation, body string) int {
	url := s.buildUrl(fmt.Sprintf("/api/v1/users/%v/vendors/%v/apis/%v/versions/%v/%v/data", simulation.Vendor, simulation.Vendor, simulation.Api, simulation.Version, simulation.Name))

	request, _ := sling.New().Put(url).Add("Authorization", s.buildAuthorizationHeaderValue()).Add("Content-Type", "application/json").Body(strings.NewReader(body)).Request()
	response, _ := http.DefaultClient.Do(request)
	defer response.Body.Close()

	return response.StatusCode
}

func (s *SpectoHub) GetSimulation(key string) []byte {
	vendor, name := splitHoverfileName(key)

	url := s.buildUrl(fmt.Sprintf("/api/v1/users/%v/vendors/%v/apis/%v/versions/%v/%v/data", vendor, vendor, "build-pipeline", "none", name))

	request, _ := sling.New().Get(url).Add("Authorization", s.buildAuthorizationHeaderValue()).Add("Content-Type", "application/json").Request()
	response, _ := http.DefaultClient.Do(request)
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	return body
}

func (s *SpectoHub) buildUrl(endpoint string) string {
	return fmt.Sprintf("%v%v", s.buildBaseUrl(), endpoint)
}

func (s *SpectoHub) buildBaseUrl() string {
	return fmt.Sprintf("http://%v:%v", s.Host, s.Port)
}

func (s *SpectoHub) buildAuthorizationHeaderValue() string {
	return fmt.Sprintf("Bearer %v", s.ApiKey)
}







