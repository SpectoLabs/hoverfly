package main

import (
	log "github.com/Sirupsen/logrus"
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

	if *verboseFlag {
		log.SetLevel(log.DebugLevel)
	}

	SetConfigurationDefaults()
	SetConfigurationPaths()

	config := GetConfig(*hostFlag, *adminPortFlag, *proxyPortFlag)

	hoverflyDirectory := getHoverflyDirectory(config)

	cacheDirectory, err := createCacheDirectory(hoverflyDirectory)
	if err != nil {
		log.Fatal("Could not create local cache")
	}

	localCache := LocalCache{
		Uri: cacheDirectory,
	}


	hoverfly := NewHoverfly(config)

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
					log.Info("Hoverfly is set to ", mode, " mode")
				} else {
					log.Fatal("Could not get Hoverfly's mode")
				}

			} else {

				mode, err := hoverfly.SetMode(*modeNameArg)
				if err == nil {
					log.Info("Hoverfly has been set to ", mode, " mode")
				} else {
					log.Fatal("Could not set Hoverfly's mode")
				}

			}

		case startCommand.FullCommand():
			err := hoverfly.start(hoverflyDirectory)
			if err != nil {
				log.Fatal("Could not start Hoverfly")
			}

		case stopCommand.FullCommand():
			hoverfly.stop(hoverflyDirectory)

		case exportCommand.FullCommand():
			simulation, err := NewSimulation(*exportNameArg)
			if err != nil {
				log.Fatal("Could not export from Hoverfly")
			}

			simulationData, err := hoverfly.ExportSimulation()

			if err != nil {
				log.Fatal("Could not export from Hoverfly")
			}

			if err = localCache.WriteSimulation(simulation, simulationData); err == nil {
				log.Info(simulation.String(), " exported successfully")
			} else {
				log.Fatal("Could not write simulation to local cache")
			}

		case importCommand.FullCommand():
			simulation, err := NewSimulation(*importNameArg)
			if err != nil {
				log.Fatal("Could not import into Hoverfly")
			}

			simulationData, err := localCache.ReadSimulation(simulation)
			if err != nil {
				log.Fatal("Could not read simulation from local cache")
			}

			if err = hoverfly.ImportSimulation(string(simulationData)); err == nil {
				log.Info(simulation.String(), " imported successfully")
			} else {
				log.Fatal("Could not import into Hoverfly")
			}

		case pushCommand.FullCommand():
			simulation, err := NewSimulation(*pushNameArg)
			if err != nil {
				log.Fatal("Could not push to Specto Labs")
			}

			simulationData, err := localCache.ReadSimulation(simulation)
			if err != nil {
				log.Fatal("Could not read simulation from local cache")
			}


			statusCode, err := spectoLab.UploadSimulation(simulation, simulationData)
			if err != nil {
				log.Fatal("Could not upload simulation to Specto Labs")
			}

			if statusCode {
				log.Info(simulation.String(), " has been pushed to the Specto Lab")
			}

		case pullCommand.FullCommand():
			simulation, err := NewSimulation(*pullNameArg)
			if err != nil {
				log.Fatal("Could not pull from Specto Labs")
			}

			simulationData, err := spectoLab.GetSimulation(simulation, *pullOverrideHostFlag)
			if err != nil {
				log.Fatal("Could not pull simulation from Specto Labs")
			}

			if err := localCache.WriteSimulation(simulation, simulationData); err == nil {
				log.Info(simulation.String(), " has been pulled from the Specto Lab")
			} else {
				log.Fatal("Could not write simulation to local cache")
			}

		case wipeCommand.FullCommand():
			if err := hoverfly.Wipe(); err == nil {
				log.Info("Hoverfly has been wiped")
			} else {
				log.Fatal("Could not wipe Hoverfly")
			}
	}
}

func getHoverflyDirectory(config Config) string {
	if len(config.GetFilepath()) == 0 {
		log.Info("Missing a config file")
		log.Info("Creating a new  a config file")

		hoverflyDir := createHoverflyDirectory(getHomeDirectory())

		err := config.WriteToFile(hoverflyDir)

		if err != nil {
			log.Fatal("Could not write new config to disk")
		}

		return hoverflyDir
	}

	return path.Dir(config.GetFilepath())
}