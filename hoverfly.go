package main

import (
	log "github.com/Sirupsen/logrus"
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

func NewHoverfly(config Config) (Hoverfly) {
	return Hoverfly {
		Host: config.HoverflyHost,
		AdminPort: config.HoverflyAdminPort,
		ProxyPort: config.HoverflyProxyPort,
		httpClient: http.DefaultClient,
	}
}

func (h *Hoverfly) Wipe() (error) {
	url := h.buildUrl("/api/records")

	request, err := sling.New().Delete(url).Request()
	if err != nil {
		log.Debug(err.Error())
		return err
	}

	response, err := h.httpClient.Do(request)
	if err != nil {
		log.Debug(err.Error())
		return err
	}
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
		log.Debug(err.Error())
		return "", err
	}

	response, err := h.httpClient.Do(request)

	if err != nil {
		log.Debug(err.Error())
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
		log.Debug(err.Error())
		return "", err
	}

	response, err := h.httpClient.Do(request)

	if err != nil {
		log.Debug(err.Error())
		return "", err
	}

	apiResponse := h.createApiStateResponse(response)

	return apiResponse.Mode, nil
}

func (h *Hoverfly) ImportSimulation(payload string) (error) {
	url := h.buildUrl("/api/records")
	request, err := sling.New().Post(url).Body(strings.NewReader(payload)).Request()

	if err != nil {
		log.Debug(err.Error())
		return err
	}

	response, err := h.httpClient.Do(request)

	if err != nil {
		log.Debug(err.Error())
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
		log.Debug(err.Error())
		return nil, err
	}

	response, err := h.httpClient.Do(request)
	if err != nil {
		log.Debug(err.Error())
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Debug(err.Error())
		return nil, err
	}

	return body, nil
}

func (h *Hoverfly) createApiStateResponse(response *http.Response) (ApiStateResponse) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Debug(err.Error())
	}

	var apiResponse ApiStateResponse

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Debug(err.Error())
	}

	return apiResponse
}

func (h *Hoverfly) buildUrl(endpoint string) (string) {
	return fmt.Sprintf("%v%v", h.buildBaseUrl(), endpoint)
}

func (h *Hoverfly) buildBaseUrl() string {
	return fmt.Sprintf("http://%v:%v", h.Host, h.AdminPort)
}

/*
This isn't working as intended, its working, just not how I imagined it.
 */

func (h *Hoverfly) start(hoverflyDir string) (error) {
	hoverflyPidFile := h.buildPidFilePath(hoverflyDir)

	if _, err := os.Stat(hoverflyPidFile); err != nil {
		if os.IsNotExist(err) {
			cmd := exec.Command("hoverfly", "-db", "memory", "-ap", h.AdminPort, "-pp", h.ProxyPort)
			err = cmd.Start()

			if err != nil {
				log.Debug(err.Error())
				return err
			}

			ioutil.WriteFile(hoverflyPidFile, []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
			log.Info("Hoverfly is now running")
		}
	} else {
		log.Info("Hoverfly is already running")
	}

	//WRITE A LOOP TO CHECK IF ITS RUNNING

	return nil
}

/*
This isn't working as intended, its working, just not how I imagined it.
 */

func (h *Hoverfly) stop(hoverflyDir string) {
	hoverflyPidFile := h.buildPidFilePath(hoverflyDir)

	if _, err := os.Stat(hoverflyPidFile); err != nil {
		if os.IsNotExist(err) {
			log.Fatal("Hoverfly is not running")
		}
	} else {
		pidFileData, _ := ioutil.ReadFile(hoverflyPidFile)
		pid, _ := strconv.Atoi(string(pidFileData))
		hoverflyProcess := os.Process{Pid: pid}
		err := hoverflyProcess.Kill()
		if err == nil {
			log.Info("Hoverfly has been killed")
			os.Remove(hoverflyPidFile)
		} else {
			log.Debug(err.Error())
			log.Debug("Pid: %#v", pid)
			log.Fatal("Failed to kill Hoverfly")
		}
	}
}

func (h *Hoverfly) buildPidFilePath(hoverflyDir string) (string) {
	pidName := fmt.Sprintf("hoverfly.%v.%v.pid", h.AdminPort, h.ProxyPort)
	return filepath.Join(hoverflyDir, pidName)
}