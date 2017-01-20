package functional_tests

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"io"

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
	ProxyPort int
	process   *exec.Cmd
}

func NewHoverfly() *Hoverfly {
	return &Hoverfly{
		adminPort: freeport.GetPort(),
		ProxyPort: freeport.GetPort(),
	}
}

func (this *Hoverfly) Start(commands ...string) {
	this.process = startHoverflyInternal(this.adminPort, this.ProxyPort, commands...)
}

func (this Hoverfly) Stop() error {
	return this.process.Process.Kill()
}
func (this Hoverfly) SetMode(mode string) {
	req := sling.New().Put("http://localhost:" + strconv.Itoa(this.adminPort) + "/api/v2/hoverfly/mode").Body(strings.NewReader(`{"mode":"capture"}`))
	DoRequest(req)
}

func (this Hoverfly) GetSimulation() io.Reader {
	res := sling.New().Get("http://localhost:" + strconv.Itoa(this.adminPort) + "/api/records")
	req := DoRequest(res)
	Expect(req.StatusCode).To(Equal(200))
	return req.Body
}

func (this Hoverfly) Proxy(r *sling.Sling) *http.Response {
	req, err := r.Request()
	Expect(err).To(BeNil())

	proxy, _ := url.Parse("http://localhost:" + strconv.Itoa(this.ProxyPort))
	proxyHttpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}, CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }}
	response, err := proxyHttpClient.Do(req)

	Expect(err).To(BeNil())

	return response
}

func startHoverflyInternal(adminPort, proxyPort int, commands ...string) *exec.Cmd {
	hoverflyBinaryUri := buildBinaryPath()

	commands = append(commands, "-ap")
	commands = append(commands, strconv.Itoa(adminPort))
	commands = append(commands, "-pp")
	commands = append(commands, strconv.Itoa(proxyPort))

	hoverflyCmd := exec.Command(hoverflyBinaryUri, commands...)

	err := hoverflyCmd.Start()

	binaryErrorCheck(err, hoverflyBinaryUri)
	healthcheck(adminPort)

	return hoverflyCmd
}

func buildBinaryPath() string {
	workingDirectory, _ := os.Getwd()
	return filepath.Join(workingDirectory, "bin/hoverfly")
}

func binaryErrorCheck(err error, binaryPath string) {
	if err != nil {
		fmt.Println("Unable to start Hoverfly")
		fmt.Println(binaryPath)
		fmt.Println("Is the binary there?")
		os.Exit(1)
	}
}

func healthcheck(adminPort int) {
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
