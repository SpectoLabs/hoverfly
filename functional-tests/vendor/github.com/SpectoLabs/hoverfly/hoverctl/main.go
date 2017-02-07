package main

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	hoverctlVersion string
	hostFlag        = kingpin.Flag("host", "Set the host of Hoverfly").String()
	adminPortFlag   = kingpin.Flag("admin-port", "Set the admin port of Hoverfly").String()
	proxyPortFlag   = kingpin.Flag("proxy-port", "Set the admin port of Hoverfly").String()
	verboseFlag     = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()

	modeCommand = kingpin.Command("mode", "Get Hoverfly's current mode")
	modeNameArg = modeCommand.Arg("name", "Set Hoverfly's mode").String()

	middlewareCommand = kingpin.Command("middleware", "Get Hoverfly's middleware")
	middlewarePathArg = middlewareCommand.Arg("path", "Set Hoverfly's middleware").String()

	startCommand = kingpin.Command("start", "Start a local instance of Hoverfly")
	startArg     = startCommand.Arg("server type", "Choose the configuration of Hoverfly (proxy/webserver)").String()
	stopCommand  = kingpin.Command("stop", "Stop a local instance of Hoverfly")

	exportCommand = kingpin.Command("export", "Exports data out of Hoverfly")
	exportNameArg = exportCommand.Arg("name", "Name of exported simulation").Required().String()

	importCommand = kingpin.Command("import", "Imports data into Hoverfly")
	importNameArg = importCommand.Arg("name", "Name of imported simulation").Required().String()

	pushCommand = kingpin.Command("push", "Pushes the data to SpectoLab")
	pushNameArg = pushCommand.Arg("name", "Name of exported simulation").Required().String()

	pullCommand          = kingpin.Command("pull", "Pushes the data to SpectoLab")
	pullNameArg          = pullCommand.Arg("name", "Name of imported simulation").Required().String()
	pullOverrideHostFlag = pullCommand.Flag("override-host", "Name of the host you want to virtualise").String()

	deleteCommand = kingpin.Command("delete", "Delete test data from Hoverfly")
	deleteArg     = deleteCommand.Arg("resource", "A collection of data that can be deleted").String()

	delaysCommand = kingpin.Command("delays", "Get per-host response delay config currently loaded in Hoverfly")
	delaysPathArg = delaysCommand.Arg("path", "Set per-host response delay config from JSON file").String()

	logsCommand    = kingpin.Command("logs", "Get the logs from Hoverfly")
	followLogsFlag = logsCommand.Flag("follow", "Follow the logs from Hoverfly").Bool()

	templatesCommand = kingpin.Command("templates", "Get set of request templates currently loaded in Hoverfly")
	templatesPathArg = templatesCommand.Arg("path", "Add JSON config to set of request templates in Hoverfly").String()
)

func main() {
	deleteCommand.Alias("wipe")
	kingpin.Version(hoverctlVersion)

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
		Host:   "https://lab.specto.io",
		APIKey: config.SpectoLabAPIKey,
	}

	switch kingpin.Parse() {
	case modeCommand.FullCommand():
		if *modeNameArg == "" || *modeNameArg == "status" {
			mode, err := hoverfly.GetMode()
			handleIfError(err)

			log.Info("Hoverfly is set to ", mode, " mode")
		} else {
			mode, err := hoverfly.SetMode(*modeNameArg)
			handleIfError(err)

			log.Info("Hoverfly has been set to ", mode, " mode")
		}

	case middlewareCommand.FullCommand():
		var middleware string
		if *middlewarePathArg == "" || *modeNameArg == "status" {
			middleware, err = hoverfly.GetMiddleware()
			handleIfError(err)
			log.Info("Hoverfly is currently set to run the following as middleware")
		} else {
			middleware, err = hoverfly.SetMiddleware(*middlewarePathArg)
			handleIfError(err)
			log.Info("Hoverfly is now set to run the following as middleware")
		}

		log.Info(middleware)

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
		switch *deleteArg {
		case "all":
			err := hoverfly.DeleteSimulations()
			handleIfError(err)
			err = hoverfly.DeleteDelays()
			handleIfError(err)
			err = hoverfly.DeleteRequestTemplates()
			handleIfError(err)
			_, err = hoverfly.SetMiddleware("")
			handleIfError(err)

			log.Info("Delays, middleware, request templates and simulations have all been deleted from Hoverfly")
		case "simulations":
			err := hoverfly.DeleteSimulations()
			handleIfError(err)

			log.Info("Simulations have been deleted from Hoverfly")

		case "delays":
			err := hoverfly.DeleteDelays()
			handleIfError(err)

			log.Info("Delays have been deleted from Hoverfly")
		case "templates":
			err := hoverfly.DeleteRequestTemplates()
			handleIfError(err)

			log.Info("Request templates have been deleted from Hoverfly")

		case "middleware":
			_, err := hoverfly.SetMiddleware("")
			handleIfError(err)

			log.Info("Middleware has been deleted from Hoverfly")
		case "":
			err := errors.New("You have not specified a resource to delete from Hoverfly")
			handleIfError(err)
		default:
			err := errors.New("You have not specified a valid resource to delete from Hoverfly")
			handleIfError(err)
		}

	case delaysCommand.FullCommand():
		if *delaysPathArg == "" || *delaysPathArg == "status" {
			delays, err := hoverfly.GetDelays()
			handleIfError(err)
			if len(delays) == 0 {
				log.Info("Hoverfly has no delays configured")
			} else {
				log.Info("Hoverfly has been configured with these delays")
				printResponseDelays(delays)
			}

		} else {
			delays, err := hoverfly.SetDelays(*delaysPathArg)
			handleIfError(err)
			log.Info("Response delays set in Hoverfly: ")
			printResponseDelays(delays)
		}
	case logsCommand.FullCommand():
		logfile := NewLogFile(hoverflyDirectory, hoverfly.AdminPort, hoverfly.ProxyPort)

		if *followLogsFlag {
			err := logfile.Tail()
			handleIfError(err)
		} else {
			err := logfile.Print()
			handleIfError(err)
		}
	case templatesCommand.FullCommand():
		if *templatesPathArg == "" || *templatesPathArg == "status" {
			requestTemplatesData, err := hoverfly.GetRequestTemplates()
			handleIfError(err)
			requestTemplatesJson, err := json.MarshalIndent(requestTemplatesData, "", "    ")
			if err != nil {
				log.Error("Error marshalling JSON for printing request templates: " + err.Error())
			}
			fmt.Println(string(requestTemplatesJson))
		} else {
			requestTemplatesData, err := hoverfly.SetRequestTemplates(*templatesPathArg)
			handleIfError(err)
			fmt.Println("Request template data set in Hoverfly: ")
			requestTemplatesJson, err := json.MarshalIndent(requestTemplatesData, "", "    ")
			if err != nil {
				log.Error("Error marshalling JSON for printing request templates: " + err.Error())
			}
			fmt.Println(string(requestTemplatesJson))
		}
	}
}

func printResponseDelays(delays []ResponseDelaySchema) {
	for _, delay := range delays {
		var delayString string
		if delay.HttpMethod != "" {
			delayString = fmt.Sprintf("%v | %v - %vms", delay.HttpMethod, delay.UrlPattern, delay.Delay)
		} else {
			delayString = fmt.Sprintf("%v - %vms", delay.UrlPattern, delay.Delay)
		}
		log.Info(delayString)
	}
}

func handleIfError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
