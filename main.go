package main

import (
	"fmt"
	"strings"
	"net/http"
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/dghubble/sling"
)

var (
	simulateCommand = kingpin.Command("simulate", "Set hoverfly to simulate mode")
	captureCommand = kingpin.Command("capture", "Set hoverfly to capture mode")
	modifyCommand = kingpin.Command("modify", "Set hoverfly to modify mode")
	synthesizeCommand = kingpin.Command("synthesize", "Set hoverfly to synthesize mode")
)


func main() {
	switch kingpin.Parse() {
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

func setHoverflyMode(mode string) (*http.Response) {
	request, _ := sling.New().Post("http://localhost:8888/api/state").Body(strings.NewReader(`{"mode":"` + mode + `"}`)).Request()
	response, _ := http.DefaultClient.Do(request)
	return response
}
