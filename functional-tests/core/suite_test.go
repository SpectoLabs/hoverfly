package hoverfly_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	"github.com/phayes/freeport"
)

var (
	hoverflyAdminUrl string
	hoverflyProxyUrl string

	hoverflyCmd *exec.Cmd

	adminPort = freeport.GetPort()

	proxyPort = freeport.GetPort()
)

func TestCore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Core Suite")
}

var _ = BeforeSuite(func() {
	hoverflyAdminUrl = fmt.Sprintf("http://localhost:%v", adminPort)
	hoverflyProxyUrl = fmt.Sprintf("http://localhost:%v", proxyPort)
})

var _ = AfterSuite(func() {
	stopHoverfly()
})

func startHoverfly(adminPort, proxyPort int) *exec.Cmd {
	return startHoverflyInternal("-ap", strconv.Itoa(adminPort), "-pp", strconv.Itoa(proxyPort))
}

func startHoverflyInternal(commands ...string) *exec.Cmd {
	hoverflyBinaryUri := functional_tests.BuildBinaryPath()
	hoverflyCmd := exec.Command(hoverflyBinaryUri, commands...)

	err := hoverflyCmd.Start()

	functional_tests.BinaryErrorCheck(err, hoverflyBinaryUri)
	functional_tests.Healthcheck(adminPort)

	return hoverflyCmd
}

func stopHoverfly() {
	hoverflyCmd.Process.Kill()
}

func DoRequestThroughProxy(r *sling.Sling) *http.Response {
	req, err := r.Request()
	Expect(err).To(BeNil())

	proxy, _ := url.Parse(hoverflyProxyUrl)
	proxyHttpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}, CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }}
	response, err := proxyHttpClient.Do(req)

	Expect(err).To(BeNil())

	return response
}

func SetHoverflyMode(mode string) {
	req := sling.New().Put(hoverflyAdminUrl + "/api/v2/hoverfly/mode").Body(strings.NewReader(`{"mode":"` + mode + `"}`))
	res := functional_tests.DoRequest(req)
	Expect(res.StatusCode).To(Equal(200))
}

func EraseHoverflySimulation() {
	req := sling.New().Delete(hoverflyAdminUrl + "/api/v2/simulation")
	res := functional_tests.DoRequest(req)
	Expect(res.StatusCode).To(Equal(200))
}

func ExportHoverflySimulation() io.Reader {
	res := sling.New().Get(hoverflyAdminUrl + "/api/v2/simulation")
	req := functional_tests.DoRequest(res)
	Expect(req.StatusCode).To(Equal(200))
	return req.Body
}

func ImportHoverflySimulation(payload io.Reader) *http.Response {
	req := sling.New().Put(hoverflyAdminUrl + "/api/v2/simulation").Body(payload)
	return functional_tests.DoRequest(req)
}

func CallFakeServerThroughProxy(server *httptest.Server) *http.Response {
	return DoRequestThroughProxy(sling.New().Get(server.URL))
}
