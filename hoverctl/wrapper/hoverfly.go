package wrapper

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
	"github.com/kardianos/osext"
)

const (
	v2ApiSimulation  = "/api/v2/simulation"
	v2ApiMode        = "/api/v2/hoverfly/mode"
	v2ApiDestination = "/api/v2/hoverfly/destination"
	v2ApiState       = "/api/v2/state"
	v2ApiMiddleware  = "/api/v2/hoverfly/middleware"
	v2ApiPac         = "/api/v2/hoverfly/pac"
	v2ApiCache       = "/api/v2/cache"
	v2ApiLogs        = "/api/v2/logs"
	v2ApiHoverfly    = "/api/v2/hoverfly"
	v2ApiDiff        = "/api/v2/diff"

	v2ApiShutdown = "/api/v2/shutdown"
	v2ApiHealth   = "/api/health"
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

func UnmarshalToInterface(response *http.Response, v interface{}) error {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}

func Login(target configuration.Target, username, password string) (string, error) {
	credentials := HoverflyAuthSchema{
		Username: username,
		Password: password,
	}

	jsonCredentials, err := json.Marshal(credentials)
	if err != nil {
		return "", fmt.Errorf("There was an error when preparing to login")
	}

	request, err := http.NewRequest("POST", BuildURL(target, "/api/token-auth"), strings.NewReader(string(jsonCredentials)))
	if err != nil {
		return "", fmt.Errorf("There was an error when preparing to login")
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("There was an error when logging in")
	}

	if response.StatusCode == http.StatusTooManyRequests {
		return "", fmt.Errorf("Too many failed login attempts, please wait 10 minutes")
	}

	if response.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("Incorrect username or password")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("There was an error when logging in")
	}

	var authToken HoverflyAuthTokenSchema
	err = json.Unmarshal(body, &authToken)
	if err != nil {
		return "", fmt.Errorf("There was an error when logging in")
	}

	return authToken.Token, nil
}

func BuildURL(target configuration.Target, endpoint string) string {
	if !strings.HasPrefix(target.Host, "http://") && !strings.HasPrefix(target.Host, "https://") {
		return fmt.Sprintf("http://%v:%v%v", target.Host, target.AdminPort, endpoint)
	}
	return fmt.Sprintf("%v:%v%v", target.Host, target.AdminPort, endpoint)
}

func IsLocal(url string) bool {
	return strings.Contains(url, "localhost") || strings.Contains(url, "127.0.0.1")
}

/*
This isn't working as intended, its working, just not how I imagined it.
*/

func runBinary(target *configuration.Target, path string) (*exec.Cmd, error) {
	flags := target.BuildFlags()

	cmd := exec.Command(path, flags...)
	log.Debug(cmd.Args)

	err := cmd.Start()
	if err != nil {
		log.Debug(err)
		return nil, errors.New("Could not start Hoverfly")
	}

	return cmd, nil
}

func Start(target *configuration.Target) error {
	// TODO only check port if is it localhost
	err := checkPorts(target.AdminPort, target.ProxyPort)
	if err != nil {
		return err
	}

	binaryLocation, err := osext.ExecutableFolder()
	if err != nil {
		log.Debug(err)
		return errors.New("Could not start Hoverfly")
	}

	_, err = runBinary(target, binaryLocation+"/hoverfly")
	if err != nil {
		_, err = runBinary(target, "hoverfly")
		if err != nil {
			return errors.New("Could not start Hoverfly")
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
			resp, err := http.Get(fmt.Sprintf("http://localhost:%v/api/health", target.AdminPort))
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

	if target.PACFile != "" {
		SetPACFile(*target)
	}

	return nil
}

func Stop(target configuration.Target) error {
	response, err := doRequest(target, "DELETE", v2ApiShutdown, "", nil)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	err = handleResponseError(response, "Could not stop Hoverfly")
	if err != nil {
		return err
	}

	return nil
}

func CheckIfRunning(target configuration.Target) error {
	_, err := doRequest(target, http.MethodGet, v2ApiHealth, "", nil)
	if err != nil {
		return fmt.Errorf("Target Hoverfly is not running\n\nRun `hoverctl start -t %s` to start it", target.Name)
	}

	return nil
}

// GetHoverfly will get the Hoverfly API which contains current configurations
func GetHoverfly(target configuration.Target) (*v2.HoverflyView, error) {
	response, err := doRequest(target, http.MethodGet, v2ApiHoverfly, "", nil)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	err = handleResponseError(response, "Could not retrieve hoverfly information")
	if err != nil {
		return nil, err
	}

	var hoverflyView v2.HoverflyView

	err = UnmarshalToInterface(response, &hoverflyView)
	if err != nil {
		return nil, err
	}

	return &hoverflyView, nil
}

func doRequest(target configuration.Target, method, url, body string, headers map[string]string) (*http.Response, error) {
	url = BuildURL(target, url)

	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("Could not connect to Hoverfly at %v:%v", target.Host, target.AdminPort)
	}

	if headers != nil {
		for key, value := range headers {
			request.Header.Add(key, value)
		}
	}

	if target.AuthToken != "" {
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", target.AuthToken))
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Could not connect to Hoverfly at %v:%v", target.Host, target.AdminPort)
	}

	if response.StatusCode == 401 {
		return nil, errors.New("Hoverfly requires authentication\n\nRun `hoverctl login -t " + target.Name + "`")
	}

	return response, nil
}

func checkPorts(ports ...int) error {
	for _, port := range ports {
		server, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
		if err != nil {
			return fmt.Errorf("Could not start Hoverfly\n\nPort %v was not free", port)
		}
		server.Close()
	}

	return nil
}

func handleResponseError(response *http.Response, errorMessage string) error {
	if response.StatusCode != 200 {
		defer response.Body.Close()
		responseError, _ := ioutil.ReadAll(response.Body)

		errSchema := &ErrorSchema{}

		err := json.Unmarshal(responseError, errSchema)
		if err != nil {
			return errors.New(errorMessage + "\n\n" + string(responseError))
		}
		return errors.New(errorMessage + "\n\n" + errSchema.ErrorMessage)
	}

	return nil
}
