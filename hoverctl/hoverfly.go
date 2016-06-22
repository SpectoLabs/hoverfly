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
	"time"
	"strconv"
)

type APIStateResponse struct {
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
	return Hoverfly{
		Host: config.HoverflyHost,
		AdminPort: config.HoverflyAdminPort,
		ProxyPort: config.HoverflyProxyPort,
		httpClient: http.DefaultClient,
	}
}

// Wipe will call the records endpoint in Hoverfly with a DELETE request, triggering Hoverfly to wipe the database
func (h *Hoverfly) Wipe() (error) {
	url := h.buildURL("/api/records")

	request, err := sling.New().Delete(url).Request()
	if err != nil {
		log.Debug(err.Error())
		return errors.New("Could not communicate with Hoverfly")
	}

	response, err := h.httpClient.Do(request)
	if err != nil {
		log.Debug(err.Error())
		return errors.New("Could not communicate with Hoverfly")
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Hoverfly did not wipe the database")
	}

	return nil
}

// GetMode will go the state endpoint in Hoverfly, parse the JSON response and return the mode of Hoverfly
func (h *Hoverfly) GetMode() (string, error) {
	url := h.buildURL("/api/state")

	request, err := sling.New().Get(url).Request()
	if err != nil {
		log.Debug(err.Error())
		return "", errors.New("Could not communicate with Hoverfly")
	}

	response, err := h.httpClient.Do(request)
	if err != nil {
		log.Debug(err.Error())
		return "", errors.New("Could not communicate with Hoverfly")
	}

	defer response.Body.Close()

	apiResponse := h.createAPIStateResponse(response)

	return apiResponse.Mode, nil
}

// Set will go the state endpoint in Hoverfly, sending JSON that will set the mode of Hoverfly
func (h *Hoverfly) SetMode(mode string) (string, error) {
	if mode != "simulate" && mode != "capture" && mode != "modify" && mode != "synthesize" {
		return "", errors.New(mode + " is not a valid mode")
	}

	url := h.buildURL("/api/state")
	request, err := sling.New().Post(url).Body(strings.NewReader(`{"mode":"` + mode + `"}`)).Request()
	if err != nil {
		log.Debug(err.Error())
		return "", errors.New("Could not communicate with Hoverfly")
	}

	response, err := h.httpClient.Do(request)
	if err != nil {
		log.Debug(err.Error())
		return "", errors.New("Could not communicate with Hoverfly")
	}

	apiResponse := h.createAPIStateResponse(response)

	return apiResponse.Mode, nil
}

func (h *Hoverfly) ImportSimulation(payload string) (error) {
	url := h.buildURL("/api/records")
	request, err := sling.New().Post(url).Body(strings.NewReader(payload)).Request()

	if err != nil {
		log.Debug(err.Error())
		return errors.New("Could not communicate with Hoverfly")
	}

	response, err := h.httpClient.Do(request)

	if err != nil {
		log.Debug(err.Error())
		return errors.New("Could not communicate with Hoverfly")
	}

	if response.StatusCode != 200 {
		return errors.New("Import to Hoverfly failed")
	}

	return nil
}

func (h *Hoverfly) ExportSimulation() ([]byte, error) {
	url := h.buildURL("/api/records")

	request, err := sling.New().Get(url).Request()
	if err != nil {
		log.Debug(err.Error())
		return nil, errors.New("Could not create a request to Hoverfly")
	}

	response, err := h.httpClient.Do(request)
	if err != nil {
		log.Debug(err.Error())
		return nil, errors.New("Could not communicate with Hoverfly")
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Debug(err.Error())
		return nil, errors.New("Could not export from Hoverfly")
	}

	return body, nil
}

func (h *Hoverfly) createAPIStateResponse(response *http.Response) (APIStateResponse) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Debug(err.Error())
	}

	var apiResponse APIStateResponse

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Debug(err.Error())
	}

	return apiResponse
}

func (h *Hoverfly) buildURL(endpoint string) (string) {
	return fmt.Sprintf("%v%v", h.buildBaseURL(), endpoint)
}

func (h *Hoverfly) buildBaseURL() string {
	return fmt.Sprintf("http://%v:%v", h.Host, h.AdminPort)
}

func (h *Hoverfly) isLocal() (bool) {
	return h.Host == "localhost" || h.Host == "127.0.0.1"
}
/*
This isn't working as intended, its working, just not how I imagined it.
 */

func (h *Hoverfly) start(hoverflyDirectory HoverflyDirectory) (error) {
	if !h.isLocal() {
		return errors.New("hoverctl can not start an instance of Hoverfly on a remote host")
	}

	pid, err := hoverflyDirectory.GetPid(h.AdminPort, h.ProxyPort)
	if err != nil {
		log.Debug(err.Error())
		return errors.New("Could not read Hoverfly pid file")
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
		log.Debug(err)
		return errors.New("Could not start Hoverfly")
	}

	timeout := time.After(10 * time.Second)
	tick := time.Tick(500 * time.Millisecond)
	statusCode := 0

	for {
		select {
			case <-timeout:
				if err != nil {
					log.Debug(err)
				}
				return errors.New(fmt.Sprintf("Timed out waiting for Hoverfly to become healthy, returns status: " + strconv.Itoa(statusCode)))
			case <-tick:
				resp, err := http.Get(fmt.Sprintf("http://localhost:%v/api/state", h.AdminPort))
				if err == nil {
					statusCode = resp.StatusCode
				} else {
					statusCode = 0
				}
			}

		if statusCode == 200 {
			break;
		}
	}

	err = hoverflyDirectory.WritePid(h.AdminPort, h.ProxyPort, cmd.Process.Pid)
	if err != nil {
		log.Debug(err.Error())
		return errors.New("Could not write a pid for Hoverfly")
	}

	return nil
}

func (h *Hoverfly) stop(hoverflyDirectory HoverflyDirectory) (error) {
	if !h.isLocal() {
		return errors.New("hoverctl can not stop an instance of Hoverfly on a remote host")
	}

	pid, err := hoverflyDirectory.GetPid(h.AdminPort, h.ProxyPort)

	if err != nil {
		log.Debug(err.Error())
		return errors.New("Could not read Hoverfly pid file")
	}

	if pid == 0 {
		return errors.New("Hoverfly is not running")
	}

	hoverflyProcess := os.Process{Pid: pid}
	err = hoverflyProcess.Kill()
	if err != nil {
		log.Info(err.Error())
		return errors.New("Could not kill Hoverfly")
	}

	err = hoverflyDirectory.DeletePid(h.AdminPort, h.ProxyPort)
	if err != nil {
		log.Debug(err.Error())
		return errors.New("Could not delete Hoverfly pid")
	}

	return nil
}