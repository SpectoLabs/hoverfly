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
	"github.com/SpectoLabs/hoverfly/core/handlers"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/kardianos/osext"
)

const (
	v1ApiDelays     = "/api/delays"
	v1ApiSimulation = "/api/records"

	v2ApiSimulation  = "/api/v2/simulation"
	v2ApiMode        = "/api/v2/hoverfly/mode"
	v2ApiDestination = "/api/v2/hoverfly/destination"
	v2ApiMiddleware  = "/api/v2/hoverfly/middleware"
	v2ApiCache       = "/api/v2/cache"
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
func DeleteSimulations(target Target) error {
	response, err := doRequest(target, "DELETE", v2ApiSimulation, "")
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
func GetMode(target Target) (string, error) {
	response, err := doRequest(target, "GET", v2ApiMode, "")
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	apiResponse := createAPIStateResponse(response)

	return apiResponse.Mode, nil
}

// Set will go the state endpoint in Hoverfly, sending JSON that will set the mode of Hoverfly
func SetModeWithArguments(target Target, modeView v2.ModeView) (string, error) {
	if modeView.Mode != "simulate" && modeView.Mode != "capture" &&
		modeView.Mode != "modify" && modeView.Mode != "synthesize" {
		return "", errors.New(modeView.Mode + " is not a valid mode")
	}
	bytes, err := json.Marshal(modeView)
	if err != nil {
		return "", err
	}

	response, err := doRequest(target, "PUT", v2ApiMode, string(bytes))
	if err != nil {
		return "", err
	}

	if response.StatusCode == http.StatusBadRequest {
		return "", handlerError(response)
	}

	apiResponse := createAPIStateResponse(response)

	return apiResponse.Mode, nil
}

// GetDestination will go the destination endpoint in Hoverfly, parse the JSON response and return the destination of Hoverfly
func GetDestination(target Target) (string, error) {
	response, err := doRequest(target, "GET", v2ApiDestination, "")
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	apiResponse := createAPIStateResponse(response)

	return apiResponse.Destination, nil
}

// SetDestination will go the destination endpoint in Hoverfly, sending JSON that will set the destination of Hoverfly
func SetDestination(target Target, destination string) (string, error) {

	response, err := doRequest(target, "PUT", v2ApiDestination, `{"destination":"`+destination+`"}`)
	if err != nil {
		return "", err
	}

	apiResponse := createAPIStateResponse(response)

	return apiResponse.Destination, nil
}

// GetMiddle will go the middleware endpoint in Hoverfly, parse the JSON response and return the middleware of Hoverfly
func GetMiddleware(target Target) (v2.MiddlewareView, error) {
	response, err := doRequest(target, "GET", v2ApiMiddleware, "")
	if err != nil {
		return v2.MiddlewareView{}, err
	}

	defer response.Body.Close()

	middlewareResponse := createMiddlewareSchema(response)

	return middlewareResponse, nil
}

func SetMiddleware(target Target, binary, script, remote string) (v2.MiddlewareView, error) {
	middlewareRequest := &v2.MiddlewareView{
		Binary: binary,
		Script: script,
		Remote: remote,
	}

	marshalledMiddleware, err := json.Marshal(middlewareRequest)
	if err != nil {
		return v2.MiddlewareView{}, err
	}

	response, err := doRequest(target, "PUT", v2ApiMiddleware, string(marshalledMiddleware))
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
		return v2.MiddlewareView{}, errors.New("Hoverfly could not execute this middleware\n\n" + error.ErrorMessage)
	}

	apiResponse := createMiddlewareSchema(response)

	return apiResponse, nil
}

func ImportSimulation(target Target, simulationData string) error {
	response, err := doRequest(target, "PUT", v2ApiSimulation, simulationData)
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

func FlushCache(target Target) error {
	response, err := doRequest(target, "DELETE", v2ApiCache, "")
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return errors.New("Cache was not set on Hoverfly")
	}

	return nil
}

func ExportSimulation(target Target) ([]byte, error) {
	response, err := doRequest(target, "GET", v2ApiSimulation, "")
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

func createAPIStateResponse(response *http.Response) APIStateSchema {
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

func createMiddlewareSchema(response *http.Response) v2.MiddlewareView {
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

func Login(target Target, username, password string) (string, error) {
	credentials := HoverflyAuthSchema{
		Username: username,
		Password: password,
	}

	jsonCredentials, err := json.Marshal(credentials)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", buildURL(target, "/api/token-auth"), strings.NewReader(string(jsonCredentials)))
	if err != nil {
		return "", err
	}

	response, err := http.DefaultClient.Do(request)
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

func buildURL(target Target, endpoint string) string {
	return fmt.Sprintf("http://%v:%v%v", target.Host, target.AdminPort, endpoint)
}

func (h *Hoverfly) buildURL(endpoint string) string {
	return fmt.Sprintf("http://%v:%v%v", h.Host, h.AdminPort, endpoint)
}

func isLocal(url string) bool {
	return url == "localhost" || url == "127.0.0.1"
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

func (h *Hoverfly) Start(target Target, hoverflyDirectory HoverflyDirectory) error {

	if !isLocal(h.Host) {
		return errors.New("hoverctl can not start an instance of Hoverfly on a remote host")
	}

	pid, err := hoverflyDirectory.GetPid(h.AdminPort, h.ProxyPort)
	if err != nil {
		log.Debug(err.Error())
		return errors.New("Could not read Hoverfly pid file")
	}

	if pid != 0 {
		_, err := GetMode(target)
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
	if !isLocal(h.Host) {
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

func doRequest(target Target, method, url, body string) (*http.Response, error) {
	url = fmt.Sprintf("http://%v:%v%v", target.Host, target.AdminPort, url)

	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	if target.AuthToken != "" {
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", target.AuthToken))
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == 401 {
		return nil, errors.New("Hoverfly requires authentication\n\nRun `hoverctl login -t " + target.Name + "`")
	}

	return response, nil
}

func handlerError(response *http.Response) error {
	responseBody, err := util.GetResponseBody(response)
	if err != nil {
		return errors.New("Error when communicating with Hoverfly")
	}

	var errorView handlers.ErrorView
	err = json.Unmarshal([]byte(responseBody), &errorView)
	if err != nil {
		return errors.New("Error when communicating with Hoverfly")
	}

	return errors.New(errorView.Error)
}
