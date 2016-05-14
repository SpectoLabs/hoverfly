package main

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"strings"
	"net/http"
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/dghubble/sling"
)

var (
	modeCommand = kingpin.Command("mode", "Get Hoverfly's current mode")
	simulateCommand = kingpin.Command("simulate", "Set Hoverfly to simulate mode")
	captureCommand = kingpin.Command("capture", "Set Hoverfly to capture mode")
	modifyCommand = kingpin.Command("modify", "Set Hoverfly to modify mode")
	synthesizeCommand = kingpin.Command("synthesize", "Set Hoverfly to synthesize mode")
)

type ApiStateResponse struct {
	Mode string        `json:"mode"`
	Destination string `json"destination"`
}

func main() {
	switch kingpin.Parse() {
		case modeCommand.FullCommand():
			modeHandler()
		case simulateCommand.FullCommand():
			simulateHandler()
		case captureCommand.FullCommand():
			captureHandler()
		case modifyCommand.FullCommand():
			modifyHandler()
		case synthesizeCommand.FullCommand():
			synthesizeHandler()
		
	}
}

func modeHandler() {
	response := getHoverflyMode()
	fmt.Println("Hoverfly is currently set to " + response.Mode + " mode")
}

func simulateHandler() {
	response := setHoverflyMode("simulate")
	defer response.Body.Close()
	fmt.Println("Hoverfly set to simulate mode")
}

func captureHandler() {
	response := setHoverflyMode("capture")
	defer response.Body.Close()
	fmt.Println("Hoverfly set to capture mode")
}

func modifyHandler() {
	response := setHoverflyMode("modify")
	defer response.Body.Close()
	fmt.Println("Hoverfly set to modify mode")
}

func synthesizeHandler() {
	response := setHoverflyMode("synthesize")
	defer response.Body.Close()
	fmt.Println("Hoverfly set to synthesize mode")
}

func getHoverflyMode() (ApiStateResponse) {
	request, _ := sling.New().Get("http://localhost:8888/api/state").Request()
	response, _ := http.DefaultClient.Do(request)
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	var jsonResponse ApiStateResponse 
	json.Unmarshal(body, &jsonResponse)
	return jsonResponse
}

func setHoverflyMode(mode string) (*http.Response) {
	request, _ := sling.New().Post("http://localhost:8888/api/state").Body(strings.NewReader(`{"mode":"` + mode + `"}`)).Request()
	response, _ := http.DefaultClient.Do(request)
	return response
}
