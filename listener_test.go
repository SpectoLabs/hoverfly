package hoverfly

import (
	"fmt"
	"github.com/SpectoLabs/hoverfly/core/testutil"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestHoverflyListener(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	proxyPort := "9777"

	dbClient.Cfg.ProxyPort = proxyPort
	// starting hoverfly
	dbClient.UpdateProxy()
	dbClient.StartProxy()

	// checking whether it's running
	response, err := http.Get(fmt.Sprintf("http://localhost:%s/", proxyPort))
	testutil.Expect(t, err, nil)

	testutil.Expect(t, response.StatusCode, 500)

	body, err := ioutil.ReadAll(response.Body)
	testutil.Expect(t, err, nil)
	testutil.Expect(t, strings.Contains(string(body), "is a proxy server"), true)
}

func TestStopHoverflyListener(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	proxyPort := "9778"

	dbClient.Cfg.ProxyPort = proxyPort
	// starting hoverfly
	dbClient.UpdateProxy()
	dbClient.StartProxy()

	dbClient.StopProxy()

	// checking whether it's stopped
	_, err := http.Get(fmt.Sprintf("http://localhost:%s/", proxyPort))
	// should get error
	testutil.Refute(t, err, nil)
}

func TestRestartHoverflyListener(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	proxyPort := "9779"

	dbClient.Cfg.ProxyPort = proxyPort
	// starting hoverfly
	dbClient.UpdateProxy()
	dbClient.StartProxy()

	// checking whether it's running
	response, err := http.Get(fmt.Sprintf("http://localhost:%s/", proxyPort))
	testutil.Expect(t, err, nil)

	testutil.Expect(t, response.StatusCode, 500)

	// stopping proxy
	dbClient.StopProxy()

	// starting again
	dbClient.StartProxy()

	newResponse, err := http.Get(fmt.Sprintf("http://localhost:%s/", proxyPort))
	testutil.Expect(t, err, nil)
	testutil.Expect(t, newResponse.StatusCode, 500)
}
