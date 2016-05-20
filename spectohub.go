package main

import (
	"fmt"
	"github.com/dghubble/sling"
	"net/http"
	"strings"
	"io/ioutil"
	"errors"
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

func (s *SpectoHub) SimulationIsPresent(hoverfile Hoverfile) (bool, error) {
	url := s.buildUrl(fmt.Sprintf("/api/v1/users/%v/vendors/%v/apis/%v/versions/%v/%v", hoverfile.Vendor, hoverfile.Vendor, "build-pipeline", hoverfile.Version, hoverfile.Name))

	request, err := sling.New().Get(url).Add("Authorization", s.buildAuthorizationHeaderValue()).Request()
	if err != nil {
		return false, err
	}

	response, _ := http.DefaultClient.Do(request)
	defer response.Body.Close()

	return response.StatusCode == 200, nil
}

func (s *SpectoHub) CreateSimulation(hoverfile Hoverfile) int {
	simulation := SpectoHubSimulation{Vendor: hoverfile.Vendor, Api: "build-pipeline", Version: hoverfile.Version, Name: hoverfile.Name, Description: "test"}

	url := s.buildUrl("/api/v1/simulations")

	request, _ := sling.New().Post(url).Add("Authorization", s.buildAuthorizationHeaderValue()).BodyJSON(simulation).Request()
	response, _ := http.DefaultClient.Do(request)
	defer response.Body.Close()
	return response.StatusCode
}

func (s *SpectoHub) UploadSimulation(hoverfile Hoverfile, data []byte) (int, error) {
	simulationExists, _ := s.SimulationIsPresent(hoverfile)
	if !simulationExists {
		postStatusCode := s.CreateSimulation(hoverfile)
		if postStatusCode != 201 {
			return 0, errors.New("Failed to create a new simulation on the Specto Hub")
		}
	}

	url := s.buildUrl(fmt.Sprintf("/api/v1/users/%v/vendors/%v/apis/%v/versions/%v/%v/data", hoverfile.Vendor, hoverfile.Vendor, "build-pipeline", hoverfile.Version, hoverfile.Name))

	request, _ := sling.New().Put(url).Add("Authorization", s.buildAuthorizationHeaderValue()).Add("Content-Type", "application/json").Body(strings.NewReader(string(data))).Request()
	response, _ := http.DefaultClient.Do(request)
	defer response.Body.Close()

	return response.StatusCode, nil
}

func (s *SpectoHub) GetSimulation(hoverfile Hoverfile) []byte {
	url := s.buildUrl(fmt.Sprintf("/api/v1/users/%v/vendors/%v/apis/%v/versions/%v/%v/data", hoverfile.Vendor, hoverfile.Vendor, "build-pipeline", hoverfile.Version, hoverfile.Name))

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







