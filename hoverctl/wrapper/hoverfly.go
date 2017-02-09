package wrapper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/dghubble/sling"
	"github.com/kardianos/osext"
)

const (
	v1ApiDelays     = "/api/delays"
	v1ApiSimulation = "/api/records"

	v2ApiSimulation  = "/api/v2/simulation"
	v2ApiMode        = "/api/v2/hoverfly/mode"
	v2ApiDestination = "/api/v2/hoverfly/destination"
	v2ApiMiddleware  = "/api/v2/hoverfly/middleware"
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

type ErrorSchema struct {
	ErrorMessage string `json:"error"`
}

type Hoverfly struct {
	Host       string
	AdminPort  string
	ProxyPort  string
	Username   string
	Password   string
	authToken  string
	config     Config
	httpClient *http.Client
}

func NewHoverfly(config Config) Hoverfly {
	return Hoverfly{
		Host:       config.HoverflyHost,
		AdminPort:  config.HoverflyAdminPort,
		ProxyPort:  config.HoverflyProxyPort,
		Username:   config.HoverflyUsername,
		Password:   config.HoverflyPassword,
		config:     config,
		httpClient: http.DefaultClient,
	}
}

// Wipe will call the records endpoint in Hoverfly with a DELETE request, triggering Hoverfly to wipe the database
func (h *Hoverfly) DeleteSimulations() error {
	slingRequest, err := h.buildDeleteRequest(v2ApiSimulation)
	if err != nil {
		return err
	}

	response, err := h.doRequest(slingRequest)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Simulations were not deleted from Hoverfly")
	}

	return nil
}

// GetMode will go the state endpoint in Hoverfly, parse the JSON response and return the mode of Hoverfly
func (h *Hoverfly) GetMode() (string, error) {
	slingRequest, err := h.buildGetRequest(v2ApiMode)
	if err != nil {
		return "", err
	}

	response, err := h.doRequest(slingRequest)
	if err != nil {
		return "", err
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

	slingRequest, err := h.buildPutRequest(v2ApiMode, `{"mode":"`+mode+`"}`)
	if err != nil {
		return "", err
	}

	response, err := h.doRequest(slingRequest)
	if err != nil {
		return "", err
	}

	if response.StatusCode == 422 {
		return "", errors.New("Cannot change the mode of Hoverfly to capture when running as a webserver")
	}

	apiResponse := h.createAPIStateResponse(response)

	return apiResponse.Mode, nil
}

// GetDestination will go the destination endpoint in Hoverfly, parse the JSON response and return the destination of Hoverfly
func (h *Hoverfly) GetDestination() (string, error) {
	slingRequest, err := h.buildGetRequest(v2ApiDestination)
	if err != nil {
		return "", err
	}

	response, err := h.doRequest(slingRequest)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	apiResponse := h.createAPIStateResponse(response)

	return apiResponse.Destination, nil
}

// SetDestination will go the destination endpoint in Hoverfly, sending JSON that will set the destination of Hoverfly
func (h *Hoverfly) SetDestination(destination string) (string, error) {

	slingRequest, err := h.buildPutRequest(v2ApiDestination, `{"destination":"`+destination+`"}`)
	if err != nil {
		return "", err
	}

	response, err := h.doRequest(slingRequest)
	if err != nil {
		return "", err
	}

	apiResponse := h.createAPIStateResponse(response)

	return apiResponse.Destination, nil
}

// GetMiddle will go the middleware endpoint in Hoverfly, parse the JSON response and return the middleware of Hoverfly
func (h *Hoverfly) GetMiddleware() (v2.MiddlewareView, error) {
	slingRequest, err := h.buildGetRequest(v2ApiMiddleware)
	if err != nil {
		return v2.MiddlewareView{}, err
	}

	response, err := h.doRequest(slingRequest)
	if err != nil {
		return v2.MiddlewareView{}, err
	}

	defer response.Body.Close()

	middlewareResponse := h.createMiddlewareSchema(response)

	return middlewareResponse, nil
}

func (h *Hoverfly) SetMiddleware(binary, script, remote string) (v2.MiddlewareView, error) {
	middlewareRequest := &v2.MiddlewareView{
		Binary: binary,
		Script: script,
		Remote: remote,
	}

	marshalledMiddleware, err := json.Marshal(middlewareRequest)
	if err != nil {
		return v2.MiddlewareView{}, err
	}

	slingRequest, err := h.buildPutRequest(v2ApiMiddleware, string(marshalledMiddleware))
	if err != nil {
		return v2.MiddlewareView{}, err
	}

	response, err := h.doRequest(slingRequest)
	if err != nil {
		return v2.MiddlewareView{}, err
	}

	if response.StatusCode == 403 {
		return v2.MiddlewareView{}, errors.New("Cannot change the mode of Hoverfly when running as a webserver")
	}

	if response.StatusCode != 200 {
		defer response.Body.Close()
		errorMessage, _ := ioutil.ReadAll(response.Body)

		error := &ErrorSchema{}

		json.Unmarshal(errorMessage, error)
		log.Debug(error.ErrorMessage)
		return v2.MiddlewareView{}, errors.New("Hoverfly could not execute this middleware")
	}

	apiResponse := h.createMiddlewareSchema(response)

	return apiResponse, nil
}

func (h *Hoverfly) ImportSimulation(simulationData string, v1 bool) error {
	slingRequest := sling.New().Body(strings.NewReader(simulationData))

	if v1 {
		slingRequest = slingRequest.Post(h.buildURL(v1ApiSimulation))
	} else {
		slingRequest = slingRequest.Put(h.buildURL(v2ApiSimulation))
	}
	slingRequest, err := h.addAuthIfNeeded(slingRequest)
	if err != nil {
		log.Debug(err.Error())
		return errors.New("Could not authenticate  with Hoverfly")
	}

	response, err := h.doRequest(slingRequest)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		body, _ := ioutil.ReadAll(response.Body)
		var errorView ErrorSchema
		json.Unmarshal(body, &errorView)
		return errors.New("Import to Hoverfly failed: " + errorView.ErrorMessage)
	}

	return nil
}

func (h *Hoverfly) ExportSimulation() ([]byte, error) {
	slingRequest, err := h.buildGetRequest(v2ApiSimulation)
	if err != nil {
		return nil, err
	}

	response, err := h.doRequest(slingRequest)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Debug(err.Error())
		return nil, errors.New("Could not export from Hoverfly")
	}

	var jsonBytes bytes.Buffer
	err = json.Indent(&jsonBytes, body, "", "\t")
	if err != nil {
		log.Debug(err.Error())
		return nil, errors.New("Could not export from Hoverfly")
	}

	return jsonBytes.Bytes(), nil
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

func (h *Hoverfly) createMiddlewareSchema(response *http.Response) v2.MiddlewareView {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Debug(err.Error())
	}

	var middleware v2.MiddlewareView

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

func (h *Hoverfly) buildGetRequest(endpoint string) (*sling.Sling, error) {
	url := h.buildURL(endpoint)
	request := sling.New().Get(url)
	return h.buildRequest(request)
}

func (h *Hoverfly) buildPutRequest(endpoint, body string) (*sling.Sling, error) {
	url := h.buildURL(endpoint)
	request := sling.New().Put(url).Body(strings.NewReader(body))
	return h.buildRequest(request)
}

func (h *Hoverfly) buildDeleteRequest(endpoint string) (*sling.Sling, error) {
	url := h.buildURL(endpoint)
	request := sling.New().Delete(url)
	return h.buildRequest(request)
}

func (h *Hoverfly) buildRequest(request *sling.Sling) (*sling.Sling, error) {
	request, err := h.addAuthIfNeeded(request)

	if err != nil {
		log.Debug(err.Error())
		return nil, errors.New("Could not authenticate with Hoverfly")
	}

	return request, nil
}

func (h *Hoverfly) doRequest(request *sling.Sling) (*http.Response, error) {
	httpRequest, err := request.Request()
	if err != nil {
		log.Debug(err.Error())
		return nil, errors.New("Could not communicate with Hoverfly")
	}

	response, err := h.httpClient.Do(httpRequest)
	if err != nil {
		log.Debug(err.Error())
		return nil, errors.New("Could not communicate with Hoverfly")
	}

	if response.StatusCode == 401 {
		return nil, errors.New("Hoverfly requires authentication")
	}

	return response, nil
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

func (h *Hoverfly) runBinary(path string, hoverflyDirectory HoverflyDirectory) (*exec.Cmd, error) {
	flags := h.config.BuildFlags()

	cmd := exec.Command(path, flags...)
	log.Debug(cmd.Args)
	file, err := os.Create(hoverflyDirectory.Path + "/hoverfly." + h.AdminPort + "." + h.ProxyPort + ".log")
	if err != nil {
		log.Debug(err)
		return nil, errors.New("Could not create log file")
	}

	cmd.Stdout = file
	cmd.Stderr = file
	defer file.Close()

	err = cmd.Start()
	if err != nil {
		log.Debug(err)
		return nil, errors.New("Could not start Hoverfly")
	}

	return cmd, nil
}

func (h *Hoverfly) Start(hoverflyDirectory HoverflyDirectory) error {

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

	binaryLocation, err := osext.ExecutableFolder()
	if err != nil {
		log.Debug(err)
		return errors.New("Could not start Hoverfly")
	}

	cmd, err := h.runBinary(binaryLocation+"/hoverfly", hoverflyDirectory)
	if err != nil {
		cmd, err = h.runBinary("hoverfly", hoverflyDirectory)
		if err != nil {
			return errors.New("Could not read Hoverfly pid file")
		}
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
			return errors.New(fmt.Sprintf("Timed out waiting for Hoverfly to become healthy, returns status: %v", statusCode))
		case <-tick:
			resp, err := http.Get(fmt.Sprintf("http://localhost:%v/api/v2/hoverfly/mode", h.AdminPort))

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

func (h *Hoverfly) Stop(hoverflyDirectory HoverflyDirectory) error {
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
