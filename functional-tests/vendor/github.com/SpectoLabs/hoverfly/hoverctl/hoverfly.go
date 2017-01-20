package main

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/dghubble/sling"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type APIStateSchema struct {
	Mode        string `json:"mode"`
	Destination string `json:"destination"`
}

type APIDelaySchema struct {
	Data []ResponseDelaySchema `json:"data"`
}

type ResponseDelaySchema struct {
	UrlPattern string `json:"urlpattern"`
	Delay      int    `json:"delay"`
	HttpMethod string `json:"httpmethod"`
}

type HoverflyAuthSchema struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type HoverflyAuthTokenSchema struct {
	Token string `json:"token"`
}

type MiddlewareSchema struct {
	Middleware string `json:"middleware"`
}

type Hoverfly struct {
	Host       string
	AdminPort  string
	ProxyPort  string
	Username   string
	Password   string
	authToken  string
	httpClient *http.Client
}

func NewHoverfly(config Config) Hoverfly {
	return Hoverfly{
		Host:       config.HoverflyHost,
		AdminPort:  config.HoverflyAdminPort,
		ProxyPort:  config.HoverflyProxyPort,
		Username:   config.HoverflyUsername,
		Password:   config.HoverflyPassword,
		httpClient: http.DefaultClient,
	}
}

// Wipe will call the records endpoint in Hoverfly with a DELETE request, triggering Hoverfly to wipe the database
func (h *Hoverfly) DeleteSimulations() error {
	url := h.buildURL("/api/records")

	slingRequest := sling.New().Delete(url)
	slingRequest, err := h.addAuthIfNeeded(slingRequest)
	if err != nil {
		log.Debug(err.Error())
		return errors.New("Could not authenticate  with Hoverfly")
	}

	request, err := slingRequest.Request()

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

	if response.StatusCode == 401 {
		return errors.New("Hoverfly requires authentication")
	}

	if response.StatusCode != 200 {
		return errors.New("Simulations were not deleted from Hoverfly")
	}

	return nil
}

func (h *Hoverfly) DeleteDelays() error {
	url := h.buildURL("/api/delays")

	slingRequest := sling.New().Delete(url)
	slingRequest, err := h.addAuthIfNeeded(slingRequest)
	if err != nil {
		log.Debug(err.Error())
		return errors.New("Could not authenticate  with Hoverfly")
	}

	request, err := slingRequest.Request()

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

	if response.StatusCode == 401 {
		return errors.New("Hoverfly requires authentication")
	}

	if response.StatusCode != 200 {
		return errors.New("Delays were not deleted from Hoverfly")
	}

	return nil
}

// GetMode will go the state endpoint in Hoverfly, parse the JSON response and return the mode of Hoverfly
func (h *Hoverfly) GetMode() (string, error) {
	url := h.buildURL("/api/state")

	slingRequest := sling.New().Get(url)

	slingRequest, err := h.addAuthIfNeeded(slingRequest)
	if err != nil {
		log.Debug(err.Error())
		return "", errors.New("Could not authenticate with Hoverfly")
	}

	request, err := slingRequest.Request()

	if err != nil {
		log.Debug(err.Error())
		return "", errors.New("Could not communicate with Hoverfly")
	}

	response, err := h.httpClient.Do(request)
	if err != nil {
		log.Debug(err.Error())
		return "", errors.New("Could not communicate with Hoverfly")
	}

	if response.StatusCode == 401 {
		return "", errors.New("Hoverfly requires authentication")
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

	slingRequest := sling.New().Post(url).Body(strings.NewReader(`{"mode":"` + mode + `"}`))

	slingRequest, err := h.addAuthIfNeeded(slingRequest)
	if err != nil {
		log.Debug(err.Error())
		return "", errors.New("Could not authenticate  with Hoverfly")
	}

	request, err := slingRequest.Request()
	if err != nil {
		log.Debug(err.Error())
		return "", errors.New("Could not communicate with Hoverfly")
	}

	response, err := h.httpClient.Do(request)
	if err != nil {
		log.Debug(err.Error())
		return "", errors.New("Could not communicate with Hoverfly")
	}

	if response.StatusCode == 401 {
		return "", errors.New("Hoverfly requires authentication")
	}

	if response.StatusCode == 403 {
		return "", errors.New("Cannot change the mode of Hoverfly when running as a webserver")
	}

	apiResponse := h.createAPIStateResponse(response)

	return apiResponse.Mode, nil
}

// GetMiddle will go the middleware endpoint in Hoverfly, parse the JSON response and return the middleware of Hoverfly
func (h *Hoverfly) GetMiddleware() (string, error) {
	url := h.buildURL("/api/middleware")

	slingRequest := sling.New().Get(url)

	slingRequest, err := h.addAuthIfNeeded(slingRequest)
	if err != nil {
		log.Debug(err.Error())
		return "", errors.New("Could not authenticate with Hoverfly")
	}

	request, err := slingRequest.Request()

	if err != nil {
		log.Debug(err.Error())
		return "", errors.New("Could not communicate with Hoverfly")
	}

	response, err := h.httpClient.Do(request)
	if err != nil {
		log.Debug(err.Error())
		return "", errors.New("Could not communicate with Hoverfly")
	}

	if response.StatusCode == 401 {
		return "", errors.New("Hoverfly requires authentication")
	}

	defer response.Body.Close()

	middlewareResponse := h.createMiddlewareSchema(response)

	return middlewareResponse.Middleware, nil
}

func (h *Hoverfly) SetMiddleware(middleware string) (string, error) {
	url := h.buildURL("/api/middleware")

	slingRequest := sling.New().Post(url).Body(strings.NewReader(`{"middleware":"` + middleware + `"}`))

	slingRequest, err := h.addAuthIfNeeded(slingRequest)
	if err != nil {
		log.Debug(err.Error())
		return "", errors.New("Could not authenticate  with Hoverfly")
	}

	request, err := slingRequest.Request()
	if err != nil {
		log.Debug(err.Error())
		return "", errors.New("Could not communicate with Hoverfly")
	}

	response, err := h.httpClient.Do(request)
	if err != nil {
		log.Debug(err.Error())
		return "", errors.New("Could not communicate with Hoverfly")
	}

	if response.StatusCode == 401 {
		return "", errors.New("Hoverfly requires authentication")
	}

	if response.StatusCode == 403 {
		return "", errors.New("Cannot change the mode of Hoverfly when running as a webserver")
	}

	if response.StatusCode != 200 {
		defer response.Body.Close()
		errorMessage, _ := ioutil.ReadAll(response.Body)
		trimmedError := strings.TrimSpace(string(errorMessage))
		log.Debug(trimmedError)
		return "", errors.New("Hoverfly could not execute this middleware")
	}

	apiResponse := h.createMiddlewareSchema(response)

	return apiResponse.Middleware, nil
}

func (h *Hoverfly) ImportSimulation(simulationData string) error {
	url := h.buildURL("/api/records")

	slingRequest := sling.New().Post(url).Body(strings.NewReader(simulationData))
	slingRequest, err := h.addAuthIfNeeded(slingRequest)
	if err != nil {
		log.Debug(err.Error())
		return errors.New("Could not authenticate  with Hoverfly")
	}

	request, err := slingRequest.Request()
	if err != nil {
		log.Debug(err.Error())
		return errors.New("Could not communicate with Hoverfly")
	}

	response, err := h.httpClient.Do(request)

	if err != nil {
		log.Debug(err.Error())
		return errors.New("Could not communicate with Hoverfly")
	}

	if response.StatusCode == 401 {
		return errors.New("Hoverfly requires authentication")
	}

	if response.StatusCode != 200 {
		return errors.New("Import to Hoverfly failed")
	}

	return nil
}

func (h *Hoverfly) ExportSimulation() ([]byte, error) {
	url := h.buildURL("/api/records")

	slingRequest := sling.New().Get(url)
	slingRequest, err := h.addAuthIfNeeded(slingRequest)
	if err != nil {
		log.Debug(err.Error())
		return nil, errors.New("Could not authenticate  with Hoverfly")
	}

	request, err := slingRequest.Request()
	if err != nil {
		log.Debug(err.Error())
		return nil, errors.New("Could not create a request to Hoverfly")
	}

	response, err := h.httpClient.Do(request)
	if err != nil {
		log.Debug(err.Error())
		return nil, errors.New("Could not communicate with Hoverfly")
	}

	if response.StatusCode == 401 {
		return nil, errors.New("Hoverfly requires authentication")
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Debug(err.Error())
		return nil, errors.New("Could not export from Hoverfly")
	}

	return body, nil
}

func (h *Hoverfly) createAPIStateResponse(response *http.Response) APIStateSchema {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Debug(err.Error())
	}

	var apiResponse APIStateSchema

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Debug(err.Error())
	}

	return apiResponse
}

func (h *Hoverfly) createMiddlewareSchema(response *http.Response) MiddlewareSchema {
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	if err != nil {
		log.Debug(err.Error())
	}

	var middleware MiddlewareSchema

	err = json.Unmarshal(body, &middleware)
	if err != nil {
		log.Debug(err.Error())
	}

	return middleware
}

func (h *Hoverfly) addAuthIfNeeded(sling *sling.Sling) (*sling.Sling, error) {
	if len(h.Username) > 0 || len(h.Password) > 0 && len(h.authToken) == 0 {
		var err error

		h.authToken, err = h.generateAuthToken()
		if err != nil {
			return nil, err
		}
	}

	if len(h.authToken) > 0 {
		sling.Add("Authorization", h.buildAuthorizationHeaderValue())
	}

	return sling, nil
}

func (h *Hoverfly) generateAuthToken() (string, error) {
	credentials := HoverflyAuthSchema{
		Username: h.Username,
		Password: h.Password,
	}

	jsonCredentials, err := json.Marshal(credentials)
	if err != nil {
		return "", err
	}

	request, err := sling.New().Post(h.buildURL("/api/token-auth")).Body(strings.NewReader(string(jsonCredentials))).Request()
	if err != nil {
		return "", err
	}

	response, err := h.httpClient.Do(request)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var authToken HoverflyAuthTokenSchema
	err = json.Unmarshal(body, &authToken)
	if err != nil {
		return "", err
	}

	return authToken.Token, nil
}

func (h *Hoverfly) buildURL(endpoint string) string {
	return fmt.Sprintf("%v%v", h.buildBaseURL(), endpoint)
}

func (h *Hoverfly) buildBaseURL() string {
	return fmt.Sprintf("http://%v:%v", h.Host, h.AdminPort)
}

func (h *Hoverfly) isLocal() bool {
	return h.Host == "localhost" || h.Host == "127.0.0.1"
}

func (h *Hoverfly) buildAuthorizationHeaderValue() string {
	return fmt.Sprintf("Bearer %v", h.authToken)
}

/*
This isn't working as intended, its working, just not how I imagined it.
*/

func (h *Hoverfly) start(hoverflyDirectory HoverflyDirectory) error {
	return h.startWithFlags(hoverflyDirectory, "")
}

func (h *Hoverfly) startWithFlags(hoverflyDirectory HoverflyDirectory, flags string) error {

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

	cmd := exec.Command("hoverfly", "-db", "memory", "-ap", h.AdminPort, "-pp", h.ProxyPort, flags)

	file, err := os.Create(hoverflyDirectory.Path + "/hoverfly." + h.AdminPort + "." + h.ProxyPort + ".log")
	if err != nil {
		log.Debug(err)
		return errors.New("Could not create log file")
	}

	cmd.Stdout = file
	cmd.Stderr = file
	defer file.Close()

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
			break
		}
	}

	err = hoverflyDirectory.WritePid(h.AdminPort, h.ProxyPort, cmd.Process.Pid)
	if err != nil {
		log.Debug(err.Error())
		return errors.New("Could not write a pid for Hoverfly")
	}

	return nil
}

func (h *Hoverfly) stop(hoverflyDirectory HoverflyDirectory) error {
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

// GetMode will go the state endpoint in Hoverfly, parse the JSON response and return the mode of Hoverfly
func (h *Hoverfly) GetDelays() (rd []ResponseDelaySchema, err error) {
	url := h.buildURL("/api/delays")

	slingRequest := sling.New().Get(url)

	slingRequest, err = h.addAuthIfNeeded(slingRequest)
	if err != nil {
		log.Debug(err.Error())
		return rd, errors.New("Could not authenticate  with Hoverfly")
	}

	request, err := slingRequest.Request()

	if err != nil {
		log.Debug(err.Error())
		return rd, errors.New("Could not communicate with Hoverfly")
	}

	response, err := h.httpClient.Do(request)
	if err != nil {
		log.Debug(err.Error())
		return rd, errors.New("Could not communicate with Hoverfly")
	}

	defer response.Body.Close()

	apiResponse := createAPIDelaysResponse(response)

	return apiResponse.Data, nil
}

// Set will go the state endpoint in Hoverfly, sending JSON that will set the mode of Hoverfly
func (h *Hoverfly) SetDelays(path string) (rd []ResponseDelaySchema, err error) {

	conf, err := ioutil.ReadFile(path)
	if err != nil {
		return rd, err
	}

	url := h.buildURL("/api/delays")

	slingRequest := sling.New().Put(url).Body(strings.NewReader(string(conf)))

	slingRequest, err = h.addAuthIfNeeded(slingRequest)
	if err != nil {
		log.Debug(err.Error())
		return rd, errors.New("Could not authenticate  with Hoverfly")
	}

	request, err := slingRequest.Request()
	if err != nil {
		log.Debug(err.Error())
		return rd, errors.New("Could not communicate with Hoverfly")
	}

	response, err := h.httpClient.Do(request)
	if err != nil {
		log.Debug(err.Error())
		return rd, errors.New("Could not communicate with Hoverfly")
	}

	slingRequest = sling.New().Get(url).Body(strings.NewReader(string(conf)))

	slingRequest, err = h.addAuthIfNeeded(slingRequest)
	if err != nil {
		log.Debug(err.Error())
		return rd, errors.New("Could not authenticate  with Hoverfly")
	}

	request, err = slingRequest.Request()
	if err != nil {
		log.Debug(err.Error())
		return rd, errors.New("Could not communicate with Hoverfly")
	}

	response, err = h.httpClient.Do(request)
	if err != nil {
		log.Debug(err.Error())
		return rd, errors.New("Could not communicate with Hoverfly")
	}
	apiResponse := createAPIDelaysResponse(response)

	return apiResponse.Data, nil
}

func createAPIDelaysResponse(response *http.Response) APIDelaySchema {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Debug(err.Error())
	}

	var apiResponse APIDelaySchema

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Debug(err.Error())
	}

	return apiResponse
}
