package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"io/ioutil"
	"strings"
	"strconv"
	"net/http"
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/dghubble/sling"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
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

	exportCommand = kingpin.Command("export", "Exports data out of Hoverfly")
	exportNameArg = exportCommand.Arg("name", "Name of exported simulation").Required().String()

	pushCommand = kingpin.Command("push", "Pushes the data to Specto Hub")
	pushNameArg = pushCommand.Arg("name", "Name of exported simulation").Required().String()

	wipeCommand = kingpin.Command("wipe", "Wipe Hoverfly database")
)

type SpectoHubSimulation struct {
	Vendor      string `json:"vendor"`
	Api         string `json:"api"`
	Version     string `json:"version"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func main() {
	hoverflyDirectory := createHomeDirectory()
	cacheDirectory := createCacheDirectory(hoverflyDirectory)

	setConfigurationDefaults()
	viper.SetConfigName("config")
	viper.AddConfigPath(hoverflyDirectory)
	err := viper.ReadInConfig()
	if err != nil {
		// Not sure what to do here
	}
	hoverfly := Hoverfly{
		Host: viper.Get("hoverfly.host").(string),
		AdminPort: viper.Get("hoverfly.admin.port").(string),
		ProxyPort: viper.Get("hoverfly.proxy.port").(string),
		httpClient: http.DefaultClient,
	}

	switch kingpin.Parse() {
		case modeCommand.FullCommand():
			mode, _ := hoverfly.GetMode()
			fmt.Println("Hoverfly is set to", mode, "mode")

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
		case exportCommand.FullCommand():
			vendor, name := splitHoverfileName(*exportNameArg)
			exportHandler(vendor, name, cacheDirectory)
		case pushCommand.FullCommand():
			pushHandler(*pushNameArg, cacheDirectory)
		case wipeCommand.FullCommand():
			err := hoverfly.WipeDatabase()
			if err == nil {
				fmt.Println("Hoverfly has been wiped")
			} else {
				fmt.Println("There was an error wiping Hoverfly")
			}




	}
}

func pushHandler(name string, cacheDirectory string) {
	vendor, name := splitHoverfileName(name)
	hoverfileName := buildHoverfileName(vendor, name)
	hoverfileUri := buildHoverfileUri(hoverfileName, cacheDirectory)

	if _, err := os.Stat(hoverfileUri); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Simulation not found")
			return
		}
	}

	hoverfileData, _ := ioutil.ReadFile(hoverfileUri)

	spectoHubSimulation := SpectoHubSimulation{Vendor: vendor, Api: "build-pipeline", Version: "none", Name: name, Description: "test"}
	getStatusCode := checkIfSimulationExists(spectoHubSimulation)
	if getStatusCode == 200 {
		fmt.Println("Updating Specto Hub")

		putStatusCode := uploadSimulation(spectoHubSimulation, string(hoverfileData))
		if putStatusCode == 200 {
			fmt.Println(name, "has been pushed to the Specto Hub")
		}

	} else {
		fmt.Println("Creating a new simulation on the Specto Hub")

		postStatusCode := createSimulation(spectoHubSimulation)
		if postStatusCode == 201 {
			putStatusCode := uploadSimulation(spectoHubSimulation, string(hoverfileData))
			if putStatusCode == 200 {
				fmt.Println(name, "has been pushed to the Specto Hub")
			}
		} else {
			fmt.Println("Failed to create a new simulation on the Specto Hub")
		}
	}
}

func checkIfSimulationExists(simulation SpectoHubSimulation) int {
	url := fmt.Sprintf("http://%v:%v/api/v1/users/%v/vendors/%v/apis/%v/versions/%v/%v", viper.Get("specto.hub.host"), viper.Get("specto.hub.port"), simulation.Vendor, simulation.Vendor, simulation.Api, simulation.Version, simulation.Name)
	authHeaderValue := fmt.Sprintf("Bearer %v", viper.Get("specto.hub.api.key"))

	request, _ := sling.New().Get(url).Add("Authorization", authHeaderValue).Request()
	response, _ := http.DefaultClient.Do(request)
	defer response.Body.Close()

	return response.StatusCode
}

func createSimulation(simulation SpectoHubSimulation) int {
	postUrl := fmt.Sprintf("http://%v:%v/api/v1/simulations", viper.Get("specto.hub.host"), viper.Get("specto.hub.port"))
	authHeaderValue := fmt.Sprintf("Bearer %v", viper.Get("specto.hub.api.key"))

	request, _ := sling.New().Post(postUrl).Add("Authorization", authHeaderValue).BodyJSON(simulation).Request()
	response, _ := http.DefaultClient.Do(request)
	defer response.Body.Close()
	return response.StatusCode
}

func uploadSimulation(simulation SpectoHubSimulation, body string) int {
	url := fmt.Sprintf("http://%v:%v/api/v1/users/%v/vendors/%v/apis/%v/versions/%v/%v/data", viper.Get("specto.hub.host"), viper.Get("specto.hub.port"), simulation.Vendor, simulation.Vendor, simulation.Api, simulation.Version, simulation.Name)
	authHeaderValue := fmt.Sprintf("Bearer %v", viper.Get("specto.hub.api.key"))
	request, _ := sling.New().Put(url).Add("Authorization", authHeaderValue).Add("Content-Type", "application/json").Body(strings.NewReader(body)).Request()
	response, _ := http.DefaultClient.Do(request)
	defer response.Body.Close()

	return response.StatusCode
}



func splitHoverfileName(hoverfileKey string) (string, string) {
	s := strings.Split(hoverfileKey, "/", )
	vendor := s[0]
	name := s[1]

	return vendor, name

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

func createCacheDirectory(baseUri string) string {
	cacheDirectory := filepath.Join(baseUri, "cache/")

	if _, err := os.Stat(cacheDirectory); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(cacheDirectory, 0777)
		}
	}

	return cacheDirectory
}

func setConfigurationDefaults() {
	viper.SetDefault("hoverfly.host", "localhost")
	viper.SetDefault("hoverfly.admin.port", "8888")
	viper.SetDefault("hoverfly.proxy.port", "8500")
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
	if _, err := os.Stat(hoverflyPidFile); err != nil {
                if os.IsNotExist(err) {
			fmt.Println("Hoverfly is not running")
		}
	} else {
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
}

func exportHandler(vendor string, name string, cacheDirectory string) {
	url := fmt.Sprintf("http://%v:%v/api/records", viper.Get("hoverfly.host"), viper.Get("hoverfly.admin.port"))
	request, _ := sling.New().Get(url).Request()
	response, _ := http.DefaultClient.Do(request)
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	hoverfileName := buildHoverfileName(vendor, name)

	hoverfileUri := buildHoverfileUri(hoverfileName, cacheDirectory)

	ioutil.WriteFile(hoverfileUri, []byte(body), 0644)
	fmt.Printf("%v/%v exported successfully", vendor, name)
}

func buildHoverfileName(vendor string, api string) string {
	return fmt.Sprintf("%v.%v.hfile", vendor, api)
}

func buildHoverfileUri(fileName string, baseUri string) string {
	return filepath.Join(baseUri, fileName)
}



func setHoverflyMode(mode string) (*http.Response) {
	url := fmt.Sprintf("http://%v:%v/api/state", viper.Get("hoverfly.host"), viper.Get("hoverfly.admin.port"))
	request, _ := sling.New().Post(url).Body(strings.NewReader(`{"mode":"` + mode + `"}`)).Request()
	response, _ := http.DefaultClient.Do(request)
	return response
}
