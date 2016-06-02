package main

import (
	"fmt"
	"github.com/dghubble/sling"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"errors"
	"strings"
	"strconv"
	"path/filepath"
	"os"
	"os/exec"
)

type ApiStateResponse struct {
	Mode        string `json:"mode"`
	Destination string `json:"destination"`
}


type Hoverfly struct {
	Host       string
	AdminPort  string
	ProxyPort  string
	httpClient *http.Client
}

func (h *Hoverfly) Wipe() error {
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

func (h *Hoverfly) ImportSimulation(payload string) error {
	url := h.buildUrl("/api/records")
	request, err := sling.New().Post(url).Body(strings.NewReader(payload)).Request()

	if err != nil {
		return err
	}

	response, err := h.httpClient.Do(request)

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return errors.New("Import to Hoverfly failed")
	}

	return nil
}

func (h *Hoverfly) ExportSimulation() ([]byte, error) {
	url := h.buildUrl("/api/records")

	request, err := sling.New().Get(url).Request()
	if err != nil {
		return nil, err
	}

	response, err := h.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
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

/*
This isn't working as intended, its working, just not how I imagined it.
 */

func startHandler(hoverflyDirectory string, hoverfly Hoverfly) error {
	pidName := fmt.Sprintf("hoverfly.%v.%v.pid", hoverfly.AdminPort, hoverfly.ProxyPort)
	hoverflyPidFile := filepath.Join(hoverflyDirectory, pidName)

	if _, err := os.Stat(hoverflyPidFile); err != nil {
		if os.IsNotExist(err) {
			binaryUri := filepath.Join(hoverflyDirectory, "/hoverfly")
			cmd := exec.Command(binaryUri, "-db", "memory", "-ap", hoverfly.AdminPort, "-pp", hoverfly.ProxyPort)
			err = cmd.Start()

			if err != nil {
				return errors.New("Hoverfly did not start")
			}

			ioutil.WriteFile(hoverflyPidFile, []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
			fmt.Println("Hoverfly is now running")
		}
	} else {
		fmt.Println("Hoverfly is already running")
	}

	//WRITE A LOOP TO CHECK IF ITS RUNNING

	return nil
}

/*
This isn't working as intended, its working, just not how I imagined it.
 */

func stopHandler(hoverflyDirectory string, hoverfly Hoverfly) {
	pidName := fmt.Sprintf("hoverfly.%v.%v.pid", hoverfly.AdminPort, hoverfly.ProxyPort)
	hoverflyPidFile := filepath.Join(hoverflyDirectory, pidName)

	if _, err := os.Stat(hoverflyPidFile); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Hoverfly is not running")
		}
	} else {
		pidFileData, _ := ioutil.ReadFile(hoverflyPidFile)
		pid, _ := strconv.Atoi(string(pidFileData))
		hoverflyProcess := os.Process{Pid: pid}
		err := hoverflyProcess.Kill()
		if err == nil {
			fmt.Println("Hoverfly has been killed")
			os.Remove(hoverflyPidFile)
		} else {
			fmt.Println("Failed to kill Hoverfly")
			fmt.Println(err.Error())
			fmt.Printf("Pid: %#v", pid)
		}
	}
}