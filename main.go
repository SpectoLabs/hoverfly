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
	request, _ := sling.New().Post("http://localhost:8888/api/state").Body(strings.NewReader(`{"mode":"capture"}`)).Request()
	response, _ := http.DefaultClient.Do(request)
	defer response.Body.Close()
	fmt.Println(response.Status)
	fmt.Println("I am capturing")
}
