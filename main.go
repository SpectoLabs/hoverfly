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

	config := GetConfig(*hostFlag, *adminPortFlag, *proxyPortFlag)

	hoverflyDirectory := getHoverflyDirectory(config)

	cacheDirectory, err := createCacheDirectory(hoverflyDirectory)
	if err != nil {
		failAndExitWithVerboseLevel("Could not create local cache", err, *verboseFlag)
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
					failAndExitWithVerboseLevel("Could not get Hoverfly's mode", err, *verboseFlag)
				}

			} else {

				mode, err := hoverfly.SetMode(*modeNameArg)
				if err == nil {
					fmt.Println("Hoverfly has been set to", mode, "mode")
				} else {
					failAndExitWithVerboseLevel("Could not set Hoverfly's mode", err, *verboseFlag)
				}

			}

		case startCommand.FullCommand():
			err := startHandler(hoverflyDirectory, hoverfly)
			if err != nil {
				failAndExitWithVerboseLevel("Could not start Hoverfly", err, *verboseFlag)
			}

		case stopCommand.FullCommand():
			stopHandler(hoverflyDirectory, hoverfly)

		case exportCommand.FullCommand():
			hoverfile, err := NewHoverfile(*exportNameArg)
			if err != nil {
				failAndExitWithVerboseLevel("Could not export from Hoverfly", err, *verboseFlag)
			}

			exportedData, err := hoverfly.ExportSimulation()

			if err != nil {
				failAndExitWithVerboseLevel("Could not export from Hoverfly", err, *verboseFlag)
			}

			if err = localCache.WriteSimulation(hoverfile, exportedData); err == nil {
				fmt.Println(*exportNameArg, "exported successfully")
			} else {
				failAndExitWithVerboseLevel("Could not write simulation to local cache", err, *verboseFlag)
			}

		case importCommand.FullCommand():
			hoverfile, err := NewHoverfile(*importNameArg)
			if err != nil {
				failAndExitWithVerboseLevel("Could not import into Hoverfly", err, *verboseFlag)
			}

			data, err := localCache.ReadSimulation(hoverfile)
			if err != nil {
				failAndExitWithVerboseLevel("Could not read simulation from local cache", err, *verboseFlag)
			}

			if err = hoverfly.ImportSimulation(string(data)); err == nil {
				fmt.Println(hoverfile.String(), "imported successfully")
			} else {
				failAndExitWithVerboseLevel("Could not import into Hoverfly", err, *verboseFlag)
			}

		case pushCommand.FullCommand():
			hoverfile, err := NewHoverfile(*pushNameArg)
			if err != nil {
				failAndExitWithVerboseLevel("Could not push to Specto Labs", err, *verboseFlag)
			}

			data, err := localCache.ReadSimulation(hoverfile)
			if err != nil {
				failAndExitWithVerboseLevel("Could not read simulation from local cache", err, *verboseFlag)
			}


			statusCode, err := spectoLab.UploadSimulation(hoverfile, data)
			if err != nil {
				failAndExitWithVerboseLevel("Could not upload simulation to Specto Labs", err, *verboseFlag)
			}

			if statusCode == 200 {
				fmt.Println(hoverfile.String(), "has been pushed to the Specto Lab")
			}

		case pullCommand.FullCommand():
			hoverfile, err := NewHoverfile(*pullNameArg)
			if err != nil {
				failAndExitWithVerboseLevel("Could not pull from Specto Labs", err, *verboseFlag)
			}

			data := spectoLab.GetSimulation(hoverfile, *pullOverrideHostFlag)

			if err := localCache.WriteSimulation(hoverfile, data); err == nil {
				fmt.Println(hoverfile.String(), "has been pulled from the Specto Lab")
			} else {
				failAndExitWithVerboseLevel("Could not write simulation to local cache", err, *verboseFlag)
			}

		case wipeCommand.FullCommand():
			if err := hoverfly.Wipe(); err == nil {
				fmt.Println("Hoverfly has been wiped")
			} else {
				failAndExitWithVerboseLevel("Could not wipe Hoverfly", err, *verboseFlag)
			}
	}
}

//func failAndExit(err error) {
//	fmt.Println(err.Error())
//	os.Exit(1)
//}

func failAndExitWithVerboseLevel(message string, err error, verbose bool) {
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
			failAndExitWithVerboseLevel("Could not get .hoverfly directory", err, *verboseFlag)
		}

		err = config.WriteToFile(hoverflyDir)

		if err != nil {
			failAndExitWithVerboseLevel("Could not write new config to disk", err, *verboseFlag)
		}

		return hoverflyDir
	}

	return path.Dir(config.GetFilepath())
}