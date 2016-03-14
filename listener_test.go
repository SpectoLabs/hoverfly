package hoverfly

import (
	"fmt"
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
	proxy, _ := GetNewHoverfly(dbClient.Cfg, dbClient.Cache)
	StartHoverflyProxy(dbClient.Cfg, proxy)

	// checking whether it's running
	response, err := http.Get(fmt.Sprintf("http://localhost:%s/", proxyPort))
	expect(t, err, nil)

	expect(t, response.StatusCode, 500)

	body, err := ioutil.ReadAll(response.Body)
	expect(t, err, nil)
	expect(t, strings.Contains(string(body), "is a proxy server"), true)
}

func TestStopHoverflyListener(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	proxyPort := "9778"

	dbClient.Cfg.ProxyPort = proxyPort
	// starting hoverfly
	proxy, _ := GetNewHoverfly(dbClient.Cfg, dbClient.Cache)
	StartHoverflyProxy(dbClient.Cfg, proxy)

	dbClient.Cfg.SL.Stop()
	dbClient.Cfg.ProxyControlWG.Wait()

	// checking whether it's stopped
	_, err := http.Get(fmt.Sprintf("http://localhost:%s/", proxyPort))
	// should get error
	refute(t, err, nil)
}

func TestRestartHoverflyListener(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	proxyPort := "9779"

	dbClient.Cfg.ProxyPort = proxyPort
	// starting hoverfly
	proxy, _ := GetNewHoverfly(dbClient.Cfg, dbClient.Cache)
	StartHoverflyProxy(dbClient.Cfg, proxy)

	// checking whether it's running
	response, err := http.Get(fmt.Sprintf("http://localhost:%s/", proxyPort))
	expect(t, err, nil)

	expect(t, response.StatusCode, 500)

	// stopping proxy
	dbClient.Cfg.SL.Stop()
	dbClient.Cfg.ProxyControlWG.Wait()

	// starting again
	StartHoverflyProxy(dbClient.Cfg, proxy)

	newResponse, err := http.Get(fmt.Sprintf("http://localhost:%s/", proxyPort))
	expect(t, err, nil)
	expect(t, newResponse.StatusCode, 500)
}
