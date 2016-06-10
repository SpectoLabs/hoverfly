package main

import (
	"fmt"
	"github.com/dghubble/sling"
	"net/http"
	"strings"
	"io/ioutil"
	"errors"
)

type SpectoLabSimulation struct {
	Vendor      string `json:"vendor"`
	Api         string `json:"api"`
	Version     string `json:"version"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SpectoLab struct {
	Host   string
	Port   string
	ApiKey string
}

func (s *SpectoLab) SimulationIsPresent(simulation Simulation) (bool, error) {
	url := s.buildUrl(fmt.Sprintf("/api/v1/users/%v/vendors/%v/apis/%v/versions/%v/%v", simulation.Vendor, simulation.Vendor, "build-pipeline", simulation.Version, simulation.Name))

	request, err := sling.New().Get(url).Add("Authorization", s.buildAuthorizationHeaderValue()).Request()
	if err != nil {
		return false, err
	}

	response, _ := http.DefaultClient.Do(request)
	defer response.Body.Close()

	return response.StatusCode == 200, nil
}

func (s *SpectoLab) CreateSimulation(simulationName Simulation) int {
	simulation := SpectoLabSimulation{Vendor: simulationName.Vendor, Api: "build-pipeline", Version: simulationName.Version, Name: simulationName.Name, Description: "test"}

	url := s.buildUrl("/api/v1/simulations")

	request, _ := sling.New().Post(url).Add("Authorization", s.buildAuthorizationHeaderValue()).BodyJSON(simulation).Request()
	response, _ := http.DefaultClient.Do(request)
	defer response.Body.Close()
	return response.StatusCode
}

func (s *SpectoLab) UploadSimulation(simulation Simulation, data []byte) (int, error) {
	simulationExists, _ := s.SimulationIsPresent(simulation)
	if !simulationExists {
		postStatusCode := s.CreateSimulation(simulation)
		if postStatusCode != 201 {
			return 0, errors.New("Failed to create a new simulation on the Specto Lab")
		}
	}

	url := s.buildUrl(fmt.Sprintf("/api/v1/users/%v/vendors/%v/apis/%v/versions/%v/%v/data", simulation.Vendor, simulation.Vendor, "build-pipeline", simulation.Version, simulation.Name))

	request, _ := sling.New().Put(url).Add("Authorization", s.buildAuthorizationHeaderValue()).Add("Content-Type", "application/json").Body(strings.NewReader(string(data))).Request()
	response, _ := http.DefaultClient.Do(request)
	defer response.Body.Close()

	return response.StatusCode, nil
}

func (s *SpectoLab) GetSimulation(simulation Simulation, overrideHost string) []byte {
	var url string

	if len(overrideHost) > 0 {
		url = s.buildUrl(fmt.Sprintf("/api/v1/users/%v/vendors/%v/apis/%v/versions/%v/%v/data?override-host=%v", simulation.Vendor, simulation.Vendor, "build-pipeline", simulation.Version, simulation.Name, overrideHost))
	} else {
		url = s.buildUrl(fmt.Sprintf("/api/v1/users/%v/vendors/%v/apis/%v/versions/%v/%v/data", simulation.Vendor, simulation.Vendor, "build-pipeline", simulation.Version, simulation.Name))
	}

	request, _ := sling.New().Get(url).Add("Authorization", s.buildAuthorizationHeaderValue()).Add("Content-Type", "application/json").Request()
	response, _ := http.DefaultClient.Do(request)
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	return body
}

func (s *SpectoLab) buildUrl(endpoint string) string {
	return fmt.Sprintf("%v%v", s.buildBaseUrl(), endpoint)
}

func (s *SpectoLab) buildBaseUrl() string {
	return fmt.Sprintf("http://%v:%v", s.Host, s.Port)
}

func (s *SpectoLab) buildAuthorizationHeaderValue() string {
	return fmt.Sprintf("Bearer %v", s.ApiKey)
}







