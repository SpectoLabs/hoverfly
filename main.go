package main

import (
	"fmt"
	"os"
	"net/http"
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/spf13/viper"
)

var (
	hostFlag = kingpin.Flag("host", "Set the host of Hoverfly").Short('h').String()
	adminPortFlag = kingpin.Flag("port", "Set the admin port of Hoverfly").Short('p').String()

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
		Uri: cacheDirectory,
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

	kingpin.Parse()

	if len(*hostFlag) > 0 {
		hoverfly.Host = *hostFlag
	}

	if len(*adminPortFlag) > 0 {
		hoverfly.AdminPort = *adminPortFlag
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
					failAndExit(err)
				}

			} else {

				mode, err := hoverfly.SetMode(*modeNameArg)
				if err == nil {
					fmt.Println("Hoverfly has been set to", mode, "mode")
				} else {
					failAndExit(err)
				}

			}

		case startCommand.FullCommand():
			err := startHandler(hoverflyDirectory)
			if err != nil {
				failAndExit(err)
			}

		case stopCommand.FullCommand():
			stopHandler(hoverflyDirectory)

		case exportCommand.FullCommand():
			hoverfile, err := NewHoverfile(*exportNameArg)
			if err != nil {
				failAndExit(err)
			}

			exportedData, err := hoverfly.ExportSimulation()

			if err != nil {
				failAndExit(err)
			}

			if err = localCache.WriteSimulation(hoverfile, exportedData); err == nil {
				fmt.Println(*exportNameArg, "exported successfully")
			} else {
				failAndExit(err)
			}

		case importCommand.FullCommand():
			hoverfile, err := NewHoverfile(*importNameArg)
			if err != nil {
				failAndExit(err)
			}

			data, err := localCache.ReadSimulation(hoverfile)
			if err != nil {
				failAndExit(err)
			}

			if err = hoverfly.ImportSimulation(string(data)); err == nil {
				fmt.Println(hoverfile.String(), "imported successfully")
			} else {
				failAndExit(err)
			}

		case pushCommand.FullCommand():
			hoverfile, err := NewHoverfile(*pushNameArg)
			if err != nil {
				failAndExit(err)
			}

			data, err := localCache.ReadSimulation(hoverfile)
			if err != nil {
				failAndExit(err)
			}


			statusCode, err := spectoHub.UploadSimulation(hoverfile, data)
			if err != nil {
				failAndExit(err)
			}

			if statusCode == 200 {
				fmt.Println(hoverfile.String(), "has been pushed to the Specto Hub")
			}

		case pullCommand.FullCommand():
			hoverfile, err := NewHoverfile(*pullNameArg)
			if err != nil {
				failAndExit(err)
			}

			data := spectoHub.GetSimulation(hoverfile)

			if err := localCache.WriteSimulation(hoverfile, data); err == nil {
				fmt.Println(hoverfile.String(), "has been pulled from the Specto Hub")
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