package hoverfly

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/onsi/gomega"
)

func TestHoverflyListener(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	proxyPort := "9777"

	unit.Cfg.ProxyPort = proxyPort
	// starting hoverfly
	unit.Proxy = NewProxy(unit)
	unit.StartProxy()

	// checking whether it's running
	response, err := http.Get(fmt.Sprintf("http://localhost:%s/", proxyPort))
	Expect(err).To(BeNil())

	Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))

	body, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())
	Expect(string(body)).To(ContainSubstring("is a proxy server"))
}

func TestStopHoverflyListener(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	proxyPort := "9778"

	unit.Cfg.ProxyPort = proxyPort
	// starting hoverfly
	unit.Proxy = NewProxy(unit)
	unit.StartProxy()

	unit.StopProxy()

	// checking whether it's stopped
	_, err := http.Get(fmt.Sprintf("http://localhost:%s/", proxyPort))
	// should get error
	Expect(err).ToNot(BeNil())
}

func TestRestartHoverflyListener(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	proxyPort := "9779"

	dbClient.Cfg.ProxyPort = proxyPort
	// starting hoverfly
	dbClient.Proxy = NewProxy(dbClient)
	dbClient.StartProxy()

	// checking whether it's running
	response, err := http.Get(fmt.Sprintf("http://localhost:%s/", proxyPort))
	Expect(err).To(BeNil())

	Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))

	// stopping proxy
	dbClient.StopProxy()

	// starting again
	dbClient.StartProxy()

	newResponse, err := http.Get(fmt.Sprintf("http://localhost:%s/", proxyPort))
	Expect(err).To(BeNil())
	Expect(newResponse.StatusCode).To(Equal(http.StatusInternalServerError))
}
