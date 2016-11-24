package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"regexp"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	hoverctlVersion string
	hostFlag        = kingpin.Flag("host", "Set the host of Hoverfly").String()
	adminPortFlag   = kingpin.Flag("admin-port", "Set the admin port of Hoverfly").String()
	proxyPortFlag   = kingpin.Flag("proxy-port", "Set the proxy port of Hoverfly").String()
	verboseFlag     = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()

	modeCommand = kingpin.Command("mode", "Get Hoverfly's current mode")
	modeNameArg = modeCommand.Arg("name", "Set Hoverfly's mode").String()

	destinationCommand = kingpin.Command("destination", "Get Hoverfly's current destination")
	destinationNameArg = destinationCommand.Arg("name", "Set Hoverfly's destination").String()
	destinationDryRun  = destinationCommand.Flag("dry-run", "Test a url against a regex pattern").String()

	middlewareCommand = kingpin.Command("middleware", "Get Hoverfly's middleware")
	middlewarePathArg = middlewareCommand.Arg("path", "Set Hoverfly's middleware").String()

	startCommand         = kingpin.Command("start", "Start a local instance of Hoverfly")
	startArg             = startCommand.Arg("server type", "Choose the configuration of Hoverfly (proxy/webserver)").String()
	startCertificateFlag = startCommand.Flag("certificate", "Supply path for custom certificate").String()
	startKeyFlag         = startCommand.Flag("key", "Supply path for custom key").String()
	startTlsFlag         = startCommand.Flag("disable-tls", "Disable TLS verification").Bool()

	stopCommand = kingpin.Command("stop", "Stop a local instance of Hoverfly")

	exportCommand = kingpin.Command("export", "Exports data out of Hoverfly")
	exportNameArg = exportCommand.Arg("name", "Name of exported simulation").Required().String()

	importCommand = kingpin.Command("import", "Imports data into Hoverfly")
	importV1Flag  = importCommand.Flag("v1", "Imports v1 formatted data into Hoverfly").Bool()
	importNameArg = importCommand.Arg("name", "Name of imported simulation").Required().String()

	deleteCommand = kingpin.Command("delete", "Delete test data from Hoverfly")
	deleteArg     = deleteCommand.Arg("resource", "A collection of data that can be deleted").String()

	delaysCommand = kingpin.Command("delays", "Get per-host response delay config currently loaded in Hoverfly")
	delaysPathArg = delaysCommand.Arg("path", "Set per-host response delay config from JSON file").String()

	logsCommand    = kingpin.Command("logs", "Get the logs from Hoverfly")
	followLogsFlag = logsCommand.Flag("follow", "Follow the logs from Hoverfly").Bool()

	templatesCommand = kingpin.Command("templates", "Get set of request templates currently loaded in Hoverfly")
	templatesPathArg = templatesCommand.Arg("path", "Add JSON config to set of request templates in Hoverfly").String()

	configCommand = kingpin.Command("config", "Get the config being used by hoverctl and Hoverfly")
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

	config := GetConfig()
	config = config.SetHost(*hostFlag)
	config = config.SetAdminPort(*adminPortFlag)
	config = config.SetProxyPort(*proxyPortFlag)
	config = config.SetUsername("")
	config = config.SetPassword("")
	config = config.SetWebserver(*startArg)
	config = config.SetCertificate(*startCertificateFlag)
	config = config.SetKey(*startKeyFlag)
	config = config.DisableTls(*startTlsFlag)

	hoverflyDirectory, err := NewHoverflyDirectory(*config)
	handleIfError(err)

	hoverfly := NewHoverfly(*config)

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
	case destinationCommand.FullCommand():
		if *destinationNameArg == "" || *destinationNameArg == "status" {
			destination, err := hoverfly.GetDestination()
			handleIfError(err)

			log.Info("The destination in Hoverfly is set to ", destination)
		} else {
			regexPattern, err := regexp.Compile(*destinationNameArg)
			if err != nil {
				log.Debug(err.Error())
				handleIfError(errors.New("Regex pattern does not compile"))
			}

			if *destinationDryRun != "" {
				if regexPattern.MatchString(*destinationDryRun) {
					log.Info("The regex provided matches the dry run URL")
				} else {
					log.Fatal("The regex provided does not match the dry run URL")
				}
			} else {
				destination, err := hoverfly.SetDestination(*destinationNameArg)
				handleIfError(err)

				log.Info("The destination in Hoverfly has been set to ", destination)
			}

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
		err := hoverfly.start(hoverflyDirectory)
		handleIfError(err)
		if config.HoverflyWebserver {
			log.Info("Hoverfly is now running as a webserver")
		} else {
			log.Info("Hoverfly is now running")
		}

	case stopCommand.FullCommand():
		err := hoverfly.stop(hoverflyDirectory)
		handleIfError(err)

		log.Info("Hoverfly has been stopped")

	case exportCommand.FullCommand():
		simulationData, err := hoverfly.ExportSimulation()
		handleIfError(err)

		err = WriteFile(*exportNameArg, simulationData)
		handleIfError(err)

		log.Info("Successfully exported to ", *exportNameArg)

	case importCommand.FullCommand():
		simulationData, err := ReadFile(*importNameArg)
		handleIfError(err)

		err = hoverfly.ImportSimulation(string(simulationData), *importV1Flag)
		handleIfError(err)

		log.Info("Successfully imported from ", *importNameArg)

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
	case configCommand.FullCommand():
		log.Info(config.GetFilepath())
		configData, _ := ReadFile(config.GetFilepath())
		configLines := strings.Split(string(configData), "\n")
		for _, line := range configLines {
			if line != "" {
				log.Info(line)
			}
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
