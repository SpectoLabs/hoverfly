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
