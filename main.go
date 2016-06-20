package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
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

	pushCommand = kingpin.Command("push", "Pushes the data to SpectoLab")
	pushNameArg = pushCommand.Arg("name", "Name of exported simulation").Required().String()

	pullCommand = kingpin.Command("pull", "Pushes the data to SpectoLab")
	pullNameArg = pullCommand.Arg("name", "Name of imported simulation").Required().String()
	pullOverrideHostFlag = pullCommand.Flag("override-host", "Name of the host you want to virtualise").String()

	wipeCommand = kingpin.Command("wipe", "Wipe Hoverfly database")
)

func main() {
	kingpin.Parse()

	if *verboseFlag {
		log.SetLevel(log.DebugLevel)
	}

	SetConfigurationDefaults()
	SetConfigurationPaths()

	config := GetConfig(*hostFlag, *adminPortFlag, *proxyPortFlag)

	hoverflyDirectory, err := NewHoverflyDirectory(config)
	handleIfError(err)

	cacheDirectory, err := createCacheDirectory(hoverflyDirectory)
	handleIfError(err)

	localCache := LocalCache{
		URI: cacheDirectory,
	}

	hoverfly := NewHoverfly(config)

	spectoLab := SpectoLab{
		Host: "https://lab.specto.io",
		APIKey: config.SpectoLabAPIKey,
	}

	switch kingpin.Parse() {
		case modeCommand.FullCommand():
			if *modeNameArg == "" || *modeNameArg == "status"{
				mode, err := hoverfly.GetMode()
				handleIfError(err)

				log.Info("Hoverfly is set to ", mode, " mode")
			} else {
				mode, err := hoverfly.SetMode(*modeNameArg)
				handleIfError(err)

				log.Info("Hoverfly has been set to ", mode, " mode")
			}

		case startCommand.FullCommand():
			err := hoverfly.start(hoverflyDirectory)
			handleIfError_old(err, "Could not start Hoverfly")

			log.Info("Hoverfly is now running")

		case stopCommand.FullCommand():
			err := hoverfly.stop(hoverflyDirectory)
			handleIfError(err)

			log.Info("Hoverfly has been stopped")

		case exportCommand.FullCommand():
			simulation, err := NewSimulation(*exportNameArg)
			handleIfError_old(err, "Could not export from Hoverfly with that name")

			simulationData, err := hoverfly.ExportSimulation()
			handleIfError_old(err, "Could not export from Hoverfly")

			err = localCache.WriteSimulation(simulation, simulationData)
			handleIfError_old(err, "Could not write simulation to local cache")

			log.Info(simulation.String(), " exported successfully")

		case importCommand.FullCommand():
			simulation, err := NewSimulation(*importNameArg)
			handleIfError_old(err, "Could not import into Hoverfly")

			simulationData, err := localCache.ReadSimulation(simulation)
			handleIfError_old(err, "Could not read simulation from local cache")

			err = hoverfly.ImportSimulation(string(simulationData))
			handleIfError_old(err, "Could not import into Hoverfly")

			log.Info(simulation.String(), " imported successfully")

		case pushCommand.FullCommand():
			err := spectoLab.CheckAPIKey()
			handleIfError_old(err, "API key not valid")

			simulation, err := NewSimulation(*pushNameArg)
			handleIfError_old(err, "Could not push to SpectoLab")

			simulationData, err := localCache.ReadSimulation(simulation)
			handleIfError_old(err, "Could not read simulation from local cache")

			statusCode, err := spectoLab.UploadSimulation(simulation, simulationData)
			handleIfError_old(err, "Could not upload simulation to SpectoLab")

			if statusCode {
				log.Info(simulation.String(), " has been pushed to the SpectoLab")
			}

		case pullCommand.FullCommand():
			err := spectoLab.CheckAPIKey()
			handleIfError_old(err, "API key not valid")

			simulation, err := NewSimulation(*pullNameArg)
			handleIfError_old(err, "Could not pull from SpectoLab")

			simulationData, err := spectoLab.GetSimulation(simulation, *pullOverrideHostFlag)
			handleIfError_old(err, "Could not pull simulation from SpectoLab")

			err = localCache.WriteSimulation(simulation, simulationData)
			handleIfError_old(err, "Could not write simulation to local cache")

			log.Info(simulation.String(), " has been pulled from the SpectoLab")

		case wipeCommand.FullCommand():
			err := hoverfly.Wipe()
			handleIfError(err)

			log.Info("Hoverfly has been wiped")
	}
}

func handleIfError_old(err error, message string) {
	if err != nil {
		log.Debug(err.Error())
		log.Fatal(message)
	}
}

func handleIfError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}