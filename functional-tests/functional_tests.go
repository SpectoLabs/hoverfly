package functional_tests

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"io"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/dghubble/sling"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phayes/freeport"
)

var HoverflyUsername = "benjih"
var HoverflyPassword = "password"
var BinaryPrefix = ""

func DoRequest(r *sling.Sling) *http.Response {
	response, err := doRequest(r)
	Expect(err).To(BeNil())
	return response
}

func doRequest(r *sling.Sling) (*http.Response, error) {
	req, err := r.Request()
	if err != nil {
		return nil, err
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func Unmarshal(data []byte, to interface{}) {
	Expect(json.Unmarshal(data, to)).To(BeNil())
}

func UnmarshalFromResponse(res *http.Response, to interface{}) {
	responseJson, err := ioutil.ReadAll(res.Body)
	Expect(err).To(BeNil())
	Unmarshal(responseJson, to)
}

type Hoverfly struct {
	adminPort int
	adminUrl  string
	proxyPort int
	proxyUrl  string
	process   *exec.Cmd
	commands  []string
}

func NewHoverfly() *Hoverfly {
	return &Hoverfly{
		adminPort: freeport.GetPort(),
		proxyPort: freeport.GetPort(),
	}
}

func NewHoverflyWithAdminPort(adminPort int) *Hoverfly {
	return &Hoverfly{
		adminPort: adminPort,
		adminUrl:  fmt.Sprintf("http://localhost:%v", adminPort),
	}
}

func (this *Hoverfly) Start(commands ...string) {
	this.process = this.startHoverflyInternal(this.adminPort, this.proxyPort, commands...)
	this.adminUrl = fmt.Sprintf("http://localhost:%v", this.adminPort)
	this.proxyUrl = fmt.Sprintf("http://localhost:%v", this.proxyPort)
}

func (this Hoverfly) Stop() {
	this.process.Process.Kill()
}

func (this Hoverfly) StopAPIAuthenticated(username, password string) {
	token, err := this.getAPIToken(username, password)
	if err != nil {
		return
	}
	_, err = doRequest(sling.New().Delete(this.adminUrl+"/api/v2/shutdown").Add("Authorization", "Bearer "+token))
	if err != nil {
		panic(err)
	}
}

func (this Hoverfly) DeleteBoltDb() {
	workingDirectory, _ := os.Getwd()
	Expect(os.Remove(workingDirectory + "requests.db")).To(BeNil())
}

func (this Hoverfly) GetMode() *v2.ModeView {
	currentState := &v2.ModeView{}
	resp := DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/v2/hoverfly/mode", this.adminPort)))

	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).To(BeNil())

	err = json.Unmarshal(body, currentState)
	Expect(err).To(BeNil())

	return currentState
}

func (this Hoverfly) SetMode(mode string) {
	this.SetModeWithArgs(mode, v2.ModeArgumentsView{})
}

func (this Hoverfly) SetModeWithArgs(mode string, arguments v2.ModeArgumentsView) {
	newMode := &v2.ModeView{
		Mode:      mode,
		Arguments: arguments,
	}

	DoRequest(sling.New().Put(this.adminUrl + "/api/v2/hoverfly/mode").BodyJSON(newMode))
}

func (this Hoverfly) SetDestination(destination string) {
	newDestination := &v2.DestinationView{
		Destination: destination,
	}
	DoRequest(sling.New().Put(this.adminUrl + "/api/v2/hoverfly/destination").BodyJSON(newDestination))
}

func (this Hoverfly) SetMiddleware(binary, script string) {
	newMiddleware := v2.MiddlewareView{
		Binary: binary,
		Script: script,
	}

	DoRequest(sling.New().Put(fmt.Sprintf("http://localhost:%v/api/v2/hoverfly/middleware", this.adminPort)).BodyJSON(newMiddleware))
}

func (this Hoverfly) GetSimulation() io.Reader {
	res := sling.New().Get(this.adminUrl + "/api/v2/simulation")
	req := DoRequest(res)
	Expect(req.StatusCode).To(Equal(200))
	return req.Body
}

func (this Hoverfly) ImportSimulation(simulation string) {
	req := sling.New().Put(this.adminUrl + "/api/v2/simulation").Body(bytes.NewBufferString(simulation))
	response := DoRequest(req)
	Expect(response.StatusCode).To(Equal(http.StatusOK), "Failed to import simulation")
	importedSimulationBytes, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())
	ginkgo.GinkgoWriter.Write(importedSimulationBytes)
}

// Used for debugging when trying to find out why a functional test is failing
func (this Hoverfly) WriteLogsIfError() {
	req := sling.New().Get(this.adminUrl+"/api/v2/logs").Add("Accept", "text/plain")
	res := DoRequest(req)
	Expect(res.StatusCode).To(Equal(200))

	logs, err := ioutil.ReadAll(res.Body)
	Expect(err).To(BeNil())
	ginkgo.GinkgoWriter.Write(logs) // Only writes when test fails
}

func (this Hoverfly) ExportSimulation() v2.SimulationViewV5 {
	reader := this.GetSimulation()
	simulationBytes, err := ioutil.ReadAll(reader)
	Expect(err).To(BeNil())

	var simulation v2.SimulationViewV5

	err = json.Unmarshal(simulationBytes, &simulation)
	Expect(err).To(BeNil())

	return simulation
}

func (this Hoverfly) GetCache() v2.CacheView {
	req := sling.New().Get(this.adminUrl + "/api/v2/cache")
	response := DoRequest(req)
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	cacheBytes, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())

	var cache v2.CacheView

	err = json.Unmarshal(cacheBytes, &cache)
	Expect(err).To(BeNil())

	return cache
}

func (this Hoverfly) FlushCache() v2.CacheView {
	req := sling.New().Delete(this.adminUrl + "/api/v2/cache")
	res := DoRequest(req)
	Expect(res.StatusCode).To(Equal(200))

	cacheBytes, err := ioutil.ReadAll(res.Body)
	Expect(err).To(BeNil())

	var cache v2.CacheView

	err = json.Unmarshal(cacheBytes, &cache)
	Expect(err).To(BeNil())

	return cache
}

func (this Hoverfly) SetPACFile(pacFile string) {
	req := sling.New().Put(this.adminUrl + "/api/v2/hoverfly/pac").Body(bytes.NewBufferString(pacFile))
	response := DoRequest(req)
	Expect(response.StatusCode).To(Equal(http.StatusOK), "Failed to set PAC file")
	_, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())
}

func (this Hoverfly) Proxy(r *sling.Sling) *http.Response {
	req, err := r.Request()
	Expect(err).To(BeNil())

	return this.ProxyRequest(req)
}

func (this Hoverfly) ProxyRequest(req *http.Request) *http.Response {

	proxy, _ := url.Parse(this.proxyUrl)
	proxyHttpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	response, err := proxyHttpClient.Do(req)

	Expect(err).To(BeNil())

	return response
}

func (this Hoverfly) ProxyWithAuth(r *sling.Sling, user, password string) *http.Response {
	req, err := r.Request()
	Expect(err).To(BeNil())

	proxy, _ := url.Parse(fmt.Sprintf("http://%s:%s@localhost:%v", user, password, this.proxyPort))
	proxyHttpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy), TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }}
	response, err := proxyHttpClient.Do(req)

	Expect(err).To(BeNil())

	return response
}

func (this Hoverfly) GetAPIToken(username, password string) string {
	token, err := this.getAPIToken(username, password)
	Expect(err).To(BeNil())

	return token
}

func (this Hoverfly) getAPIToken(username, password string) (string, error) {
	response, err := doRequest(
		sling.New().Post(this.adminUrl + "/api/token-auth").BodyJSON(map[string]interface{}{
			"username": username,
			"password": password,
		}),
	)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	bodyMap := make(map[string]interface{})
	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		return "", err
	}

	return bodyMap["token"].(string), nil
}

func (this Hoverfly) GetAdminPort() string {
	return strconv.Itoa(this.adminPort)
}

func (this Hoverfly) GetProxyPort() string {
	return strconv.Itoa(this.proxyPort)
}

func (this Hoverfly) GetPid() int {
	return this.process.Process.Pid
}

func (this Hoverfly) startHoverflyInternal(adminPort, proxyPort int, additionalCommands ...string) *exec.Cmd {
	hoverflyBinaryUri := BuildBinaryPath()

	commands := []string{
		"-ap",
		strconv.Itoa(adminPort),
		"-pp",
		strconv.Itoa(proxyPort),
	}

	commands = append(commands, additionalCommands...)
	this.commands = commands
	hoverflyCmd := exec.Command(hoverflyBinaryUri, commands...)
	err := hoverflyCmd.Start()

	BinaryErrorCheck(err, hoverflyBinaryUri)

	for _, command := range commands {
		if command == "-add" {
			time.Sleep(time.Second * 3)
			return hoverflyCmd
		}
	}

	this.healthcheck()

	return hoverflyCmd
}

func BuildBinaryPath() string {
	workingDirectory, _ := os.Getwd()
	return filepath.Join(workingDirectory, BinaryPrefix, "bin/hoverfly")
}

func BinaryErrorCheck(err error, binaryPath string) {
	if err != nil {
		fmt.Println("Unable to start Hoverfly")
		fmt.Println(os.Getwd())
		fmt.Println(binaryPath)
		fmt.Println("Is the binary there?")
		os.Exit(1)
	}
}

func (this Hoverfly) healthcheck() {
	Eventually(func() int {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%v/api/health", this.adminPort))
		if err == nil {
			return resp.StatusCode
		} else {
			return 0
		}
	}, time.Second*5).Should(BeNumerically("==", http.StatusOK), "Hoverfly not running on %d", this.adminPort, this.commands)
}

func Healthcheck(adminPort int) {
	var err error
	var resp *http.Response

	Eventually(func() int {
		resp, err = http.Get(fmt.Sprintf("http://localhost:%v/api/health", adminPort))
		if err == nil {
			return resp.StatusCode
		} else {
			return 0
		}
	}, time.Second*5).Should(BeNumerically("==", http.StatusOK), "Hoverfly not running on %d but have no extra information", adminPort)
}

func Run(binary string, commands ...string) string {
	cmd := exec.Command(binary, commands...)
	out, err := cmd.Output()
	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok {
			return string(exitError.Stderr)
		}
	}

	return strings.TrimSpace(string(out))
}

func GenerateFileName() string {

	rb := make([]byte, 6)
	rand.Read(rb)

	rs := base64.URLEncoding.EncodeToString(rb)

	return "testdata-gen/" + rs + ".json"
}

func TableToSliceMapStringString(table string) map[string]map[string]string {
	results := map[string]map[string]string{}

	tableRows := strings.Split(table, "\n")
	headings := []string{}

	for _, heading := range strings.Split(tableRows[1], "|") {
		headings = append(headings, strings.TrimSpace(heading))
	}

	for _, row := range tableRows[2:] {
		if !strings.Contains(row, "-+-") {
			rowValues := strings.Split(row, "|")

			result := map[string]string{}
			for i, value := range rowValues {
				if value != "" {
					result[headings[i]] = strings.TrimSpace(value)
				}
			}

			results[result["TARGET NAME"]] = result
		}
	}

	return results
}
