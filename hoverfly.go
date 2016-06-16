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

func (h *Hoverfly) start(hoverflyDirectory HoverflyDirectory) (error) {
	pid, err := hoverflyDirectory.GetPid(h.AdminPort, h.ProxyPort)
	if err != nil {
		return err
	}

	if pid != 0 {
		_, err := h.GetMode()
		if err == nil {
			return errors.New("Hoverfly is already running")
		}
		hoverflyDirectory.DeletePid(h.AdminPort, h.ProxyPort)
	}

	cmd := exec.Command("hoverfly", "-db", "memory", "-ap", h.AdminPort, "-pp", h.ProxyPort)

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = hoverflyDirectory.WritePid(h.AdminPort, h.ProxyPort, cmd.Process.Pid)
	if err != nil {
		return err
	}

	return nil
}

func (h *Hoverfly) stop(hoverflyDirectory HoverflyDirectory) (error) {
	pid, err := hoverflyDirectory.GetPid(h.AdminPort, h.ProxyPort)
	if err != nil {
		return err
	}

	if pid == 0 {
		return errors.New("Hoverfly is not running")

	} else {
		hoverflyProcess := os.Process{Pid: pid}
		err := hoverflyProcess.Kill()
		if err != nil {
			return err
		}

		err = hoverflyDirectory.DeletePid(h.AdminPort, h.ProxyPort)
		if err != nil {
			return err
		}
	}

	return nil
}