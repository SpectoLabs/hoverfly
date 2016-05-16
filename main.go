package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"encoding/json"
	"io/ioutil"
	"strings"
	"strconv"
	"net/http"
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/dghubble/sling"
	"github.com/mitchellh/go-homedir"
)

var (
	modeCategory = kingpin.Command("mode", "Get Hoverfly's current mode")
	modeCommand = modeCategory.Command("status", "Get Hoverfly's current mode").Default()
	simulateCommand = modeCategory.Command("simulate", "Set Hoverfly to simulate mode")
	captureCommand = modeCategory.Command("capture", "Set Hoverfly to capture mode")
	modifyCommand = modeCategory.Command("modify", "Set Hoverfly to modify mode")
	synthesizeCommand = modeCategory.Command("synthesize", "Set Hoverfly to synthesize mode")

	startCommand = kingpin.Command("start", "Start a local instance of Hoverfly")
	stopCommand = kingpin.Command("stop", "Stop a local instance of Hoverfly")
)

type ApiStateResponse struct {
	Mode string        `json:"mode"`
	Destination string `json"destination"`
}

func main() {
	hoverflyDirectory := createHomeDirectory()

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
			startHandler(hoverflyDirectory)
		case stopCommand.FullCommand():
			stopHandler(hoverflyDirectory)
		
	}
}

func createHomeDirectory() string {
	homeDirectory, _ := homedir.Dir()
	hoverflyDirectory := filepath.Join(homeDirectory, "/.hoverfly")

	if _, err := os.Stat(hoverflyDirectory); err != nil {
    		if os.IsNotExist(err) {
        		os.Mkdir(hoverflyDirectory, 0777)
    		}
	}
	
	return hoverflyDirectory
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

func startHandler(hoverflyDirectory string) {
	hoverflyPidFile := filepath.Join(hoverflyDirectory, "hoverfly.pid")

	if _, err := os.Stat(hoverflyPidFile); err != nil {
                if os.IsNotExist(err) {
			cmd := exec.Command("/Users/benjih/Downloads/hoverfly/hoverfly_v0.5.17_OSX_amd64")
			cmd.Start()
			ioutil.WriteFile(hoverflyPidFile, []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
			fmt.Println("Hoverfly is now running")
		}
        } else {
		fmt.Println("Hoverfly is already running")
	}
}

func stopHandler(hoverflyDirectory string) {
	hoverflyPidFile := filepath.Join(hoverflyDirectory, "hoverfly.pid")
	pidFileData, _ := ioutil.ReadFile(hoverflyPidFile)
	pid, _ := strconv.Atoi(string(pidFileData))
	hoverflyProcess := os.Process{Pid: pid}
	err := hoverflyProcess.Kill()
	if err == nil {
		fmt.Println("Hoverfly has been killed")
		os.Remove(hoverflyPidFile)
	} else {
		fmt.Println("Failed to kill Hoverfly")
		fmt.Println(err.Error())
		fmt.Printf("Pid: %#v", pid)
	}
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
