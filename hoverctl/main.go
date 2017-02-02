package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	hoverctlVersion string

	verboseFlag = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()

	hostFlag        = kingpin.Flag("host", "Set the host of Hoverfly").String()
	adminPortFlag   = kingpin.Flag("admin-port", "Set the admin port of Hoverfly").String()
	proxyPortFlag   = kingpin.Flag("proxy-port", "Set the proxy port of Hoverfly").String()
	certificateFlag = kingpin.Flag("certificate", "Supply path for custom certificate").String()
	keyFlag         = kingpin.Flag("key", "Supply path for custom key").String()
	disableTlsFlag  = kingpin.Flag("disable-tls", "Disable TLS verification").Bool()
	databaseFlag    = kingpin.Flag("database", "Set persistance storage to use - default to in memory DB").String()

	modeCommand = kingpin.Command("mode", "Get Hoverfly's current mode")
	modeNameArg = modeCommand.Arg("name", "Set Hoverfly's mode").String()

	destinationCommand = kingpin.Command("destination", "Get Hoverfly's current destination")
	destinationNameArg = destinationCommand.Arg("name", "Set Hoverfly's destination").String()
	destinationDryRun  = destinationCommand.Flag("dry-run", "Test a url against a regex pattern").String()

	middlewareCommand    = kingpin.Command("middleware", "Get Hoverfly's middleware")
	middlewareBinaryFlag = middlewareCommand.Flag("binary", "The binary that middleware should execute").String()
	middlewareScriptFlag = middlewareCommand.Flag("script", "The script that middleware should execute").String()
	middlewareRemoteFlag = middlewareCommand.Flag("remote", "The remote address that middleware should execute").String()

	startCommand = kingpin.Command("start", "Start a local instance of Hoverfly")
	startArg     = startCommand.Arg("server type", "Choose the configuration of Hoverfly (proxy/webserver)").String()

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

	templatesCommand = kingpin.Command("templates", "Get request templates currently loaded in Hoverfly")
	templatesPathArg = templatesCommand.Arg("path", "Add JSON config to set of request templates in Hoverfly").String()

	configCommand = kingpin.Command("config", "Get the hoverctl config location and contents")
)

func main() {
	deleteCommand.Alias("wipe")
	kingpin.Version(hoverctlVersion)

	kingpin.Parse()

	log.SetOutput(os.Stdout)
	if *verboseFlag {
		log.SetLevel(log.DebugLevel)
	}

	wrapper.SetConfigurationDefaults()
	wrapper.SetConfigurationPaths()

	config := wrapper.GetConfig()
	config = config.SetHost(*hostFlag)
	config = config.SetAdminPort(*adminPortFlag)
	config = config.SetProxyPort(*proxyPortFlag)
	config = config.SetUsername("")
	config = config.SetPassword("")
	config = config.SetWebserver(*startArg)
	config = config.SetCertificate(*certificateFlag)
	config = config.SetKey(*keyFlag)
	config = config.DisableTls(*disableTlsFlag)
	config = config.SetDbType(*databaseFlag)

	hoverflyDirectory, err := wrapper.NewHoverflyDirectory(*config)
	handleIfError(err)

	hoverfly := wrapper.NewHoverfly(*config)

	switch kingpin.Parse() {
	case startCommand.FullCommand():
		err := hoverfly.Start(hoverflyDirectory)
		handleIfError(err)
		if config.HoverflyWebserver {
			log.WithFields(log.Fields{
				"admin-port":     config.HoverflyAdminPort,
				"webserver-port": config.HoverflyProxyPort,
			}).Info("Hoverfly is now running as a webserver")
		} else {
			log.WithFields(log.Fields{
				"admin-port": config.HoverflyAdminPort,
				"proxy-port": config.HoverflyProxyPort,
			}).Info("Hoverfly is now running")
		}

	case stopCommand.FullCommand():
		err := hoverfly.Stop(hoverflyDirectory)
		handleIfError(err)

		log.Info("Hoverfly has been stopped")

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
		var middleware v2.MiddlewareView
		if *middlewareBinaryFlag == "" && *middlewareScriptFlag == "" && *middlewareRemoteFlag == "" {
			middleware, err = hoverfly.GetMiddleware()
			handleIfError(err)
			log.Info("Hoverfly is currently set to run the following as middleware")
		} else {
			if *middlewareRemoteFlag != "" {
				middleware, err = hoverfly.SetMiddleware("", "", *middlewareRemoteFlag)
				handleIfError(err)
				log.Info("Hoverfly is now set to run the following as middleware")
			} else {
				script, err := wrapper.ReadFile(*middlewareScriptFlag)
				handleIfError(err)

				middleware, err = hoverfly.SetMiddleware(*middlewareBinaryFlag, string(script), "")
				handleIfError(err)
				log.Info("Hoverfly is now set to run the following as middleware")
			}
		}

		if middleware.Binary != "" {
			log.Info("Binary: " + middleware.Binary)
		}

		if middleware.Script != "" {
			log.Info("Script: " + middleware.Script)
		}

		if middleware.Remote != "" {
			log.Info("Remote: " + middleware.Remote)
		}

	case exportCommand.FullCommand():
		simulationData, err := hoverfly.ExportSimulation()
		handleIfError(err)

		err = wrapper.WriteFile(*exportNameArg, simulationData)
		handleIfError(err)

		log.Info("Successfully exported to ", *exportNameArg)

	case importCommand.FullCommand():
		simulationData, err := wrapper.ReadFile(*importNameArg)
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
			_, err = hoverfly.SetMiddleware("", "", "")
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
			_, err := hoverfly.SetMiddleware("", "", "")
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
		logfile := wrapper.NewLogFile(hoverflyDirectory, hoverfly.AdminPort, hoverfly.ProxyPort)

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
		configData, _ := wrapper.ReadFile(config.GetFilepath())
		configLines := strings.Split(string(configData), "\n")
		for _, line := range configLines {
			if line != "" {
				log.Info(line)
			}
		}
	}
}

func printResponseDelays(delays []wrapper.ResponseDelaySchema) {
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
