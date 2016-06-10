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
	url := s.buildUrl(fmt.Sprintf("/api/v1/users/%v/simulations/%v/versions/%v", simulation.Vendor, simulation.Name, simulation.Version))

	request, err := sling.New().Get(url).Add("Authorization", s.buildAuthorizationHeaderValue()).Request()

	if err != nil {
		return false, err
	}

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return false, err
	}

	defer response.Body.Close()
	return response.StatusCode == 200, nil
}

func (s *SpectoLab) CreateSimulation(simulationName Simulation) (int, error) {
	simulation := SpectoLabSimulation{Version: simulationName.Version, Name: simulationName.Name, Description: "A description could go here"}

	url := s.buildUrl("/api/v1/simulations")
	request, err := sling.New().Post(url).BodyJSON(simulation).Add("Authorization", s.buildAuthorizationHeaderValue()).Request()

	if err != nil {
		return 0, err
	}

	response, err := http.DefaultClient.Do(request)
	defer response.Body.Close()

	if err != nil {
		return 0, err
	}

	return response.StatusCode, nil
}

func (s *SpectoLab) UploadSimulation(simulation Simulation, data []byte) (int, error) {
	simulationExists, err := s.SimulationIsPresent(simulation)

	if err != nil {
		return 0, err
	}

	if !simulationExists {
		postStatusCode, err := s.CreateSimulation(simulation)

		if err != nil {
			return 0, err
		}
		if postStatusCode != 201 {
			return 0, errors.New("Failed to create a new simulation on the Specto Lab")
		}
	}

	url := s.buildUrl(fmt.Sprintf("/api/v1/users/%v/simulations/%v/versions/%v/data", simulation.Vendor,  simulation.Name, simulation.Version))

	request, _ := sling.New().Put(url).Add("Authorization", s.buildAuthorizationHeaderValue()).Add("Content-Type", "application/json").Body(strings.NewReader(string(data))).Request()
	response, _ := http.DefaultClient.Do(request)
	defer response.Body.Close()

	return response.StatusCode, nil
}

func (s *SpectoLab) GetSimulation(simulation Simulation, overrideHost string) []byte {
	var url string
	if len(overrideHost) > 0 {
		url = s.buildUrl(fmt.Sprintf("/api/v1/users/%v/simulations/%v/versions/%v/data?override-host=%v", simulation.Vendor, simulation.Name, simulation.Version, overrideHost))
	} else {
		url = s.buildUrl(fmt.Sprintf("/api/v1/users/%v/simulations/%v/versions/%v/data", simulation.Vendor, simulation.Name, simulation.Version))
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
	if len(s.Port) > 0 {
		return fmt.Sprintf("http://%v:%v", s.Host, s.Port)
	} else {
		return fmt.Sprintf("http://%v", s.Host)
	}
}

func (s *SpectoLab) buildAuthorizationHeaderValue() string {
	return fmt.Sprintf("Bearer %v", s.ApiKey)
}







