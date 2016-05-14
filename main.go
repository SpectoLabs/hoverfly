package main

import (
	"fmt"
	"strings"
	"net/http"
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/dghubble/sling"
)

var (
	captureFlag = kingpin.Command("capture", "Set hoverfly to capture mode")
)


func main() {
	switch kingpin.Parse() {
		case "capture":
			captureHandler()
	}
}

func captureHandler() {
	response := setHoverflyMode("capture")
	defer response.Body.Close()
	fmt.Println("Hoverfly set to capture mode")
}

func setHoverflyMode(mode string) (*http.Response) {
	request, _ := sling.New().Post("http://localhost:8888/api/state").Body(strings.NewReader(`{"mode":"` + mode + `"}`)).Request()
	response, _ := http.DefaultClient.Do(request)
	return response
}
