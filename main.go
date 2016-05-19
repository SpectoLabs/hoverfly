package main

import (
	"fmt"
	"os"
	"net/http"
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/spf13/viper"
)

var (
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

	wipeCommand = kingpin.Command("wipe", "Wipe Hoverfly database")
)

func main() {
	hoverflyDirectory, err := createHomeDirectory()
	if err != nil {
		failAndExit(err)
	}

	cacheDirectory, err := createCacheDirectory(hoverflyDirectory)
	if err != nil {
		failAndExit(err)
	}

	localCache := LocalCache{
		uri: cacheDirectory,
	}


	setConfigurationDefaults(hoverflyDirectory)

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println("You are missing a config file")
	}

	hoverfly := Hoverfly {
		Host: viper.GetString("hoverfly.host"),
		AdminPort: viper.GetString("hoverfly.admin.port"),
		ProxyPort: viper.GetString("hoverfly.proxy.port"),
		httpClient: http.DefaultClient,
	}

	spectoHub := SpectoHub {
		Host: viper.GetString("specto.hub.host"),
		Port: viper.GetString("specto.hub.port"),
		ApiKey: viper.GetString("specto.hub.api.key"),
	}

	switch kingpin.Parse() {
		case modeCommand.FullCommand():
			if *modeNameArg == "" || *modeNameArg == "status"{

				mode, err := hoverfly.GetMode()
				if err == nil {
					fmt.Println("Hoverfly is set to", mode, "mode")
				} else {
					fmt.Println(err.Error())
				}

			} else {

				mode, err := hoverfly.SetMode(*modeNameArg)
				if err == nil {
					fmt.Println("Hoverfly has been set to", mode, "mode")
				} else {
					fmt.Println(err.Error())
				}

			}

		case startCommand.FullCommand():
			startHandler(hoverflyDirectory)

		case stopCommand.FullCommand():
			stopHandler(hoverflyDirectory)

		case exportCommand.FullCommand():

			exportedData, err := hoverfly.ExportSimulation()

			if err != nil {
				failAndExit(err)
			}

			if err = localCache.WriteSimulation(*exportNameArg, exportedData); err == nil {
				fmt.Println(*exportNameArg, "exported successfully")
			} else {
				failAndExit(err)
			}

		case importCommand.FullCommand():

			data, err := localCache.ReadSimulation(*importNameArg)
			if err != nil {
				failAndExit(err)
			}

			if err = hoverfly.ImportSimulation(string(data)); err == nil {
				fmt.Println(*importNameArg, "imported successfully")
			} else {
				failAndExit(err)
			}

		case pushCommand.FullCommand():
			data, err := localCache.ReadSimulation(*pushNameArg)
			if err != nil {
				failAndExit(err)
			}

			_, name := splitHoverfileName(*pushNameArg)
			getStatusCode := spectoHub.CheckSimulation(*pushNameArg)
			if getStatusCode == 200 {
				fmt.Println("Updating Specto Hub")

				putStatusCode := spectoHub.UploadSimulation(*pushNameArg, data)
				if putStatusCode == 200 {
					fmt.Println(name, "has been pushed to the Specto Hub")
				}

			} else {
				fmt.Println("Creating a new simulation on the Specto Hub")

				postStatusCode := spectoHub.CreateSimulation(*pushNameArg)
				if postStatusCode == 201 {
					putStatusCode := spectoHub.UploadSimulation(*pushNameArg, data)
					if putStatusCode == 200 {
						fmt.Println(name, "has been pushed to the Specto Hub")
					}
				} else {
					fmt.Println("Failed to create a new simulation on the Specto Hub")
				}
			}

		case pullCommand.FullCommand():
			data := spectoHub.GetSimulation(*pullNameArg)

			if err := localCache.WriteSimulation(*pullNameArg, data); err == nil {
				fmt.Println(*pullNameArg, "has been pulled from the Specto Hub")
			} else {
				failAndExit(err)
			}

		case wipeCommand.FullCommand():
			if err := hoverfly.Wipe(); err == nil {
				fmt.Println("Hoverfly has been wiped")
			} else {
				failAndExit(err)
			}
	}
}

func failAndExit(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}


func setConfigurationDefaults(hoverflyDirectory string) {
	viper.AddConfigPath(hoverflyDirectory)
	viper.SetDefault("hoverfly.host", "localhost")
	viper.SetDefault("hoverfly.admin.port", "8888")
	viper.SetDefault("hoverfly.proxy.port", "8500")
	viper.SetDefault("specto.hub.host", "localhost")
	viper.SetDefault("specto.hub.port", "81")
}