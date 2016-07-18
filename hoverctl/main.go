package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"fmt"
	"errors"
	"github.com/hpcloud/tail"
)

var (
	hostFlag = kingpin.Flag("host", "Set the host of Hoverfly").String()
	adminPortFlag = kingpin.Flag("admin-port", "Set the admin port of Hoverfly").String()
	proxyPortFlag = kingpin.Flag("proxy-port", "Set the admin port of Hoverfly").String()
	verboseFlag = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()

	modeCommand = kingpin.Command("mode", "Get Hoverfly's current mode")
	modeNameArg = modeCommand.Arg("name", "Set Hoverfly's mode").String()

	startCommand = kingpin.Command("start", "Start a local instance of Hoverfly")
	startArg = startCommand.Arg("server type", "Choose the configuration of Hoverfly (proxy/webserver)").String()
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

	deleteCommand = kingpin.Command("delete", "Delete test data from Hoverfly")
	deleteArg = deleteCommand.Arg("resource", "A collection of data that can be deleted").String()

	delaysCommand = kingpin.Command("delays", "Get per-host response delay config currently loaded into Hoverfly")
	delaysPathArg = delaysCommand.Arg("path", "Set per-host response delay config from JSON file").String()

	logsCommand = kingpin.Command("logs", "Get the logs from Hoverfly")
	followLogsFlag = logsCommand.Flag("follow", "Follow the logs from Hoverfly").Bool()
)

func main() {
	deleteCommand.Alias("wipe")

	kingpin.Parse()

	log.SetOutput(os.Stdout)
	if *verboseFlag {
		log.SetLevel(log.DebugLevel)
	}

	SetConfigurationDefaults()
	SetConfigurationPaths()

	config := GetConfig(*hostFlag, *adminPortFlag, *proxyPortFlag, "", "")

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
			if *startArg == "webserver" {
				err := hoverfly.startWithFlags(hoverflyDirectory, "-webserver")
				handleIfError(err)

				log.Info("Hoverfly is now running as a webserver")
			} else {
				err := hoverfly.start(hoverflyDirectory)
				handleIfError(err)

				log.Info("Hoverfly is now running")
			}

		case stopCommand.FullCommand():
			err := hoverfly.stop(hoverflyDirectory)
			handleIfError(err)

			log.Info("Hoverfly has been stopped")

		case exportCommand.FullCommand():
			simulation, err := NewSimulation(*exportNameArg)
			handleIfError(err)

			simulationData, err := hoverfly.ExportSimulation()
			handleIfError(err)

			err = localCache.WriteSimulation(simulation, simulationData)
			handleIfError(err)

			log.Info(simulation.String(), " exported successfully")

		case importCommand.FullCommand():
			simulation, err := NewSimulation(*importNameArg)
			handleIfError(err)

			simulationData, err := localCache.ReadSimulation(simulation)
			handleIfError(err)

			err = hoverfly.ImportSimulation(string(simulationData))
			handleIfError(err)

			log.Info(simulation.String(), " imported successfully")

		case pushCommand.FullCommand():
			err := spectoLab.CheckAPIKey()
			handleIfError(err)

			simulation, err := NewSimulation(*pushNameArg)
			handleIfError(err)

			simulationData, err := localCache.ReadSimulation(simulation)
			handleIfError(err)

			statusCode, err := spectoLab.UploadSimulation(simulation, simulationData)
			handleIfError(err)

			if statusCode {
				log.Info(simulation.String(), " has been pushed to the SpectoLab")
			}

		case pullCommand.FullCommand():
			err := spectoLab.CheckAPIKey()
			handleIfError(err)

			simulation, err := NewSimulation(*pullNameArg)
			handleIfError(err)

			simulationData, err := spectoLab.GetSimulation(simulation, *pullOverrideHostFlag)
			handleIfError(err)

			err = localCache.WriteSimulation(simulation, simulationData)
			handleIfError(err)

			log.Info(simulation.String(), " has been pulled from the SpectoLab")

		case deleteCommand.FullCommand():
			if *deleteArg == "all" {
				err := hoverfly.DeleteSimulations()
				handleIfError(err)
				err = hoverfly.DeleteDelays()
				handleIfError(err)

				log.Info("Delays and simulations have been deleted from Hoverfly")
			}
			if *deleteArg == "simulations" {
				err := hoverfly.DeleteSimulations()
				handleIfError(err)

				log.Info("Simulations have been deleted from Hoverfly")
			}
			if *deleteArg == "delays" {
				err := hoverfly.DeleteDelays()
				handleIfError(err)

				log.Info("Delays have been deleted from Hoverfly")
			}

			if *deleteArg == "" {
				err := errors.New("You have not specified what to delete from Hoverfly")
				handleIfError(err)
			}

		case delaysCommand.FullCommand():
			if *delaysPathArg == "" || *delaysPathArg == "status"{
				delays, err := hoverfly.GetDelays()
				handleIfError(err)
				for _, delay := range delays {
					fmt.Printf("%+v\n", delay)
				}
			} else {
				delays, err := hoverfly.SetDelays(*delaysPathArg)
				handleIfError(err)
				log.Info("Response delays set in Hoverfly: ")
				for _, delay := range delays {
					fmt.Printf("%+v\n", delay)
				}
			}
		case logsCommand.FullCommand():
			logfile := NewLogFile(hoverflyDirectory, hoverfly.AdminPort, hoverfly.ProxyPort)

			logs, err := logfile.GetLogs()
			handleIfError(err)

			fmt.Print(logs)

			if *followLogsFlag {
				tail, err := tail.TailFile(logfile.Path, tail.Config{Follow: true})
				if err != nil {
					log.Debug(err.Error())
					handleIfError(errors.New("Could not follow Hoverfly log file"))
				}

				for line := range tail.Lines {
					fmt.Println(line.Text)
				}
			}
	}
}

func handleIfError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}