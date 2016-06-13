package main

import (
	"fmt"
	"os"
	"net/http"
	"gopkg.in/alecthomas/kingpin.v2"
	"path"
)

var (
	hostFlag = kingpin.Flag("host", "Set the host of Hoverfly").String()
	adminPortFlag = kingpin.Flag("admin-port", "Set the admin port of Hoverfly").String()
	proxyPortFlag = kingpin.Flag("proxy-port", "Set the admin port of Hoverfly").String()
	verboseFlag = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()

	modeCommand = kingpin.Command("mode", "Get Hoverfly's current mode")
	modeNameArg = modeCommand.Arg("name", "Set Hoverfly's mode").String()

	startCommand = kingpin.Command("start", "Start a local instance of Hoverfly")
	stopCommand = kingpin.Command("stop", "Stop a local instance of Hoverfly")

	exportCommand = kingpin.Command("export", "Exports data out of Hoverfly")
	exportNameArg = exportCommand.Arg("name", "Name of exported simulation").Required().String()

	importCommand = kingpin.Command("import", "Imports data into Hoverfly")
	importNameArg = importCommand.Arg("name", "Name of imported simulation").Required().String()

	pushCommand = kingpin.Command("push", "Pushes the data to Specto Hub")
	pushNameArg = pushCommand.Arg("name", "Name of exported simulation").Required().String()

	pullCommand = kingpin.Command("pull", "Pushes the data to Specto Hub")
	pullNameArg = pullCommand.Arg("name", "Name of imported simulation").Required().String()
	pullOverrideHostFlag = pullCommand.Flag("override-host", "Name of the host you want to virtualise").String()

	wipeCommand = kingpin.Command("wipe", "Wipe Hoverfly database")
)

func main() {
	kingpin.Parse()

	SetConfigurationDefaults()
	SetConfigurationPaths()

	config := GetConfig(*hostFlag, *adminPortFlag, *proxyPortFlag)

	hoverflyDirectory := getHoverflyDirectory(config)

	cacheDirectory, err := createCacheDirectory(hoverflyDirectory)
	if err != nil {
		failAndExit("Could not create local cache", err, *verboseFlag)
	}

	localCache := LocalCache{
		Uri: cacheDirectory,
	}


	hoverfly := Hoverfly {
		Host: config.HoverflyHost,
		AdminPort: config.HoverflyAdminPort,
		ProxyPort: config.HoverflyProxyPort,
		httpClient: http.DefaultClient,
	}

	spectoLab := SpectoLab{
		Host: config.SpectoLabHost,
		Port: config.SpectoLabPort,
		ApiKey: config.SpectoLabApiKey,
	}

	switch kingpin.Parse() {
		case modeCommand.FullCommand():
			if *modeNameArg == "" || *modeNameArg == "status"{

				mode, err := hoverfly.GetMode()
				if err == nil {
					fmt.Println("Hoverfly is set to", mode, "mode")
				} else {
					failAndExit("Could not get Hoverfly's mode", err, *verboseFlag)
				}

			} else {

				mode, err := hoverfly.SetMode(*modeNameArg)
				if err == nil {
					fmt.Println("Hoverfly has been set to", mode, "mode")
				} else {
					failAndExit("Could not set Hoverfly's mode", err, *verboseFlag)
				}

			}

		case startCommand.FullCommand():
			err := startHandler(hoverflyDirectory, hoverfly)
			if err != nil {
				failAndExit("Could not start Hoverfly", err, *verboseFlag)
			}

		case stopCommand.FullCommand():
			stopHandler(hoverflyDirectory, hoverfly)

		case exportCommand.FullCommand():
			simulation, err := NewSimulation(*exportNameArg)
			if err != nil {
				failAndExit("Could not export from Hoverfly", err, *verboseFlag)
			}

			simulationData, err := hoverfly.ExportSimulation()

			if err != nil {
				failAndExit("Could not export from Hoverfly", err, *verboseFlag)
			}

			if err = localCache.WriteSimulation(simulation, simulationData); err == nil {
				fmt.Println(*exportNameArg, "exported successfully")
			} else {
				failAndExit("Could not write simulation to local cache", err, *verboseFlag)
			}

		case importCommand.FullCommand():
			simulation, err := NewSimulation(*importNameArg)
			if err != nil {
				failAndExit("Could not import into Hoverfly", err, *verboseFlag)
			}

			simulationData, err := localCache.ReadSimulation(simulation)
			if err != nil {
				failAndExit("Could not read simulation from local cache", err, *verboseFlag)
			}

			if err = hoverfly.ImportSimulation(string(simulationData)); err == nil {
				fmt.Println(simulation.String(), "imported successfully")
			} else {
				failAndExit("Could not import into Hoverfly", err, *verboseFlag)
			}

		case pushCommand.FullCommand():
			simulation, err := NewSimulation(*pushNameArg)
			if err != nil {
				failAndExit("Could not push to Specto Labs", err, *verboseFlag)
			}

			simulationData, err := localCache.ReadSimulation(simulation)
			if err != nil {
				failAndExit("Could not read simulation from local cache", err, *verboseFlag)
			}


			statusCode, err := spectoLab.UploadSimulation(simulation, simulationData)
			if err != nil {
				failAndExit("Could not upload simulation to Specto Labs", err, *verboseFlag)
			}

			if statusCode {
				fmt.Println(simulation.String(), "has been pushed to the Specto Lab")
			}

		case pullCommand.FullCommand():
			simulation, err := NewSimulation(*pullNameArg)
			if err != nil {
				failAndExit("Could not pull from Specto Labs", err, *verboseFlag)
			}

			simulationData := spectoLab.GetSimulation(simulation, *pullOverrideHostFlag)

			if err := localCache.WriteSimulation(simulation, simulationData); err == nil {
				fmt.Println(simulation.String(), "has been pulled from the Specto Lab")
			} else {
				failAndExit("Could not write simulation to local cache", err, *verboseFlag)
			}

		case wipeCommand.FullCommand():
			if err := hoverfly.Wipe(); err == nil {
				fmt.Println("Hoverfly has been wiped")
			} else {
				failAndExit("Could not wipe Hoverfly", err, *verboseFlag)
			}
	}
}

func failAndExit(message string, err error, verbose bool) {
	fmt.Println(message)
	if verbose {
		fmt.Println(err.Error())
	}
	os.Exit(1)
}

func getHoverflyDirectory(config Config) string {
	if len(config.GetFilepath()) == 0 {
		fmt.Println("Missing a config file")
		fmt.Println("Creating a new  a config file")

		hoverflyDir, err := createHomeDirectory()

		if err != nil {
			failAndExit("Could not get .hoverfly directory", err, *verboseFlag)
		}

		err = config.WriteToFile(hoverflyDir)

		if err != nil {
			failAndExit("Could not write new config to disk", err, *verboseFlag)
		}

		return hoverflyDir
	}

	return path.Dir(config.GetFilepath())
}