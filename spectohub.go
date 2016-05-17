package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/dghubble/sling"
	"net/http"
	"strings"
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
	url := fmt.Sprintf("http://%v:%v/api/v1/users/%v/vendors/%v/apis/%v/versions/%v/%v", viper.Get("specto.hub.host"), viper.Get("specto.hub.port"), simulation.Vendor, simulation.Vendor, simulation.Api, simulation.Version, simulation.Name)
	authHeaderValue := fmt.Sprintf("Bearer %v", viper.Get("specto.hub.api.key"))

	request, _ := sling.New().Get(url).Add("Authorization", authHeaderValue).Request()
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

func (s *SpectoHub) buildUrl(endpoint string) string {
	return fmt.Sprintf("%v%v", s.buildBaseUrl(), endpoint)
}

func (s *SpectoHub) buildBaseUrl() string {
	return fmt.Sprintf("http://%v:%v", s.Host, s.Port)
}

func (s *SpectoHub) buildAuthorizationHeaderValue() string {
	return fmt.Sprintf("Bearer %v", s.ApiKey)
}







