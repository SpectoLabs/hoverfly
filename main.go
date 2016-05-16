package main

import (
	"fmt"
	"os/exec"
	"encoding/json"
	"io/ioutil"
	"strings"
	"net/http"
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/dghubble/sling"
)

var (
	modeCategory = kingpin.Command("mode", "Get Hoverfly's current mode")
	modeCommand = modeCategory.Command("status", "Get Hoverfly's current mode").Default()
	simulateCommand = modeCategory.Command("simulate", "Set Hoverfly to simulate mode")
	captureCommand = modeCategory.Command("capture", "Set Hoverfly to capture mode")
	modifyCommand = modeCategory.Command("modify", "Set Hoverfly to modify mode")
	synthesizeCommand = modeCategory.Command("synthesize", "Set Hoverfly to synthesize mode")

	startCommand = kingpin.Command("start", "Start a local instance of Hoverfly")
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
		case startCommand.FullCommand():
			startHandler()
		
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

func startHandler() {
	//os.Open("/Users/benjih/Downloads/hoverfly/hoverfly_v0.5.17_OSX_amd64")
	cmd := exec.Command("/Users/benjih/Downloads/hoverfly/hoverfly_v0.5.17_OSX_amd64")
	cmd.Start()
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
