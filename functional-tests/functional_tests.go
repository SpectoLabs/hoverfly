package functional_tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"io"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/dghubble/sling"
	. "github.com/onsi/gomega"
	"github.com/phayes/freeport"
)

func DoRequest(r *sling.Sling) *http.Response {
	req, err := r.Request()
	Expect(err).To(BeNil())
	response, err := http.DefaultClient.Do(req)

	Expect(err).To(BeNil())
	return response
}

type Hoverfly struct {
	adminPort int
	adminUrl  string
	proxyPort int
	proxyUrl  string
	process   *exec.Cmd
}

func NewHoverfly() *Hoverfly {
	return &Hoverfly{
		adminPort: freeport.GetPort(),
		proxyPort: freeport.GetPort(),
	}
}

func (this *Hoverfly) Start(commands ...string) {
	this.process = startHoverflyInternal(this.adminPort, this.proxyPort, commands...)
	this.adminUrl = fmt.Sprintf("http://localhost:%v", this.adminPort)
	this.proxyUrl = fmt.Sprintf("http://localhost:%v", this.proxyPort)
}

func (this Hoverfly) Stop() error {
	return this.process.Process.Kill()
}

func (this Hoverfly) GetMode() string {
	currentState := &v2.ModeView{}
	resp := DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/v2/hoverfly/mode", this.adminPort)))

	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).To(BeNil())

	json.Unmarshal(body, currentState)

	return currentState.Mode
}
func (this Hoverfly) SetMode(mode string) {
	newMode := &v2.ModeView{
		Mode: mode,
	}

	DoRequest(sling.New().Put(this.adminUrl + "/api/v2/hoverfly/mode").BodyJSON(newMode))
}

func (this Hoverfly) SetMiddleware(binary, script string) {
	newMiddleware := v2.MiddlewareView{
		Binary: binary,
		Script: script,
	}

	DoRequest(sling.New().Put(fmt.Sprintf("http://localhost:%v/api/v2/hoverfly/middleware", this.adminPort)).BodyJSON(newMiddleware))
}

func (this Hoverfly) GetSimulation() io.Reader {
	res := sling.New().Get(this.adminUrl + "/api/records")
	req := DoRequest(res)
	Expect(req.StatusCode).To(Equal(200))
	return req.Body
}

func (this Hoverfly) Proxy(r *sling.Sling) *http.Response {
	req, err := r.Request()
	Expect(err).To(BeNil())

	proxy, _ := url.Parse(this.proxyUrl)
	proxyHttpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}, CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }}
	response, err := proxyHttpClient.Do(req)

	Expect(err).To(BeNil())

	return response
}

func (this Hoverfly) GetAdminPort() string {
	return strconv.Itoa(this.adminPort)
}

func (this Hoverfly) GetProxyPort() string {
	return strconv.Itoa(this.proxyPort)
}

func startHoverflyInternal(adminPort, proxyPort int, additionalCommands ...string) *exec.Cmd {
	hoverflyBinaryUri := BuildBinaryPath()

	commands := []string{
		"-ap",
		strconv.Itoa(adminPort),
		"-pp",
		strconv.Itoa(proxyPort),
	}

	commands = append(commands, additionalCommands...)

	hoverflyCmd := exec.Command(hoverflyBinaryUri, commands...)

	err := hoverflyCmd.Start()

	BinaryErrorCheck(err, hoverflyBinaryUri)

	healthCheckNeeded := true
	for _, command := range commands {
		if command == "-add" {
			healthCheckNeeded = false
		}
	}

	if healthCheckNeeded {
		Healthcheck(adminPort)
	}

	return hoverflyCmd
}

func BuildBinaryPath() string {
	workingDirectory, _ := os.Getwd()
	return filepath.Join(workingDirectory, "bin/hoverfly")
}

func BinaryErrorCheck(err error, binaryPath string) {
	if err != nil {
		fmt.Println("Unable to start Hoverfly")
		fmt.Println(binaryPath)
		fmt.Println("Is the binary there?")
		os.Exit(1)
	}
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
	}, time.Second*3).Should(BeNumerically("==", http.StatusOK))

}
