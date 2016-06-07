package main

import (
	"fmt"
	"os"
	"net/http"
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/spf13/viper"
	"path"
)

var (
	hostFlag = kingpin.Flag("host", "Set the host of Hoverfly").String()
	adminPortFlag = kingpin.Flag("admin-port", "Set the admin port of Hoverfly").String()
	proxyPortFlag = kingpin.Flag("proxy-port", "Set the admin port of Hoverfly").String()

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
	setConfigurationDefaults()

	viper.ReadInConfig()
	configUri := viper.ConfigFileUsed()

	hoverflyDirectory := getHoverflyDirectory(configUri)

	cacheDirectory, err := createCacheDirectory(hoverflyDirectory)
	if err != nil {
		failAndExit(err)
	}

	localCache := LocalCache{
		Uri: cacheDirectory,
	}

	kingpin.Parse()

	hoverfly := createHoverfly(*hostFlag, *adminPortFlag, *proxyPortFlag)

	spectoLab := SpectoLab{
		Host: viper.GetString("specto.lab.host"),
		Port: viper.GetString("specto.lab.port"),
		ApiKey: viper.GetString("specto.lab.api.key"),
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
			err := startHandler(hoverflyDirectory, hoverfly)
			if err != nil {
				failAndExit(err)
			}

		case stopCommand.FullCommand():
			stopHandler(hoverflyDirectory, hoverfly)

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


			statusCode, err := spectoLab.UploadSimulation(hoverfile, data)
			if err != nil {
				failAndExit(err)
			}

			if statusCode == 200 {
				fmt.Println(hoverfile.String(), "has been pushed to the Specto Lab")
			}

		case pullCommand.FullCommand():
			hoverfile, err := NewHoverfile(*pullNameArg)
			if err != nil {
				failAndExit(err)
			}

			data := spectoLab.GetSimulation(hoverfile, *pullOverrideHostFlag)

			if err := localCache.WriteSimulation(hoverfile, data); err == nil {
				fmt.Println(hoverfile.String(), "has been pulled from the Specto Lab")
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

func createHoverfly(hostOverride, adminPortOverride, proxyPortOverride string) Hoverfly {
	hoverfly := Hoverfly {
		Host: viper.GetString("hoverfly.host"),
		AdminPort: viper.GetString("hoverfly.admin.port"),
		ProxyPort: viper.GetString("hoverfly.proxy.port"),
		httpClient: http.DefaultClient,
	}

	if len(*hostFlag) > 0 {
		hoverfly.Host = *hostFlag
	}

	if len(*adminPortFlag) > 0 {
		hoverfly.AdminPort = *adminPortFlag
	}

	if len(*proxyPortFlag) > 0 {
		hoverfly.ProxyPort = *proxyPortFlag
	}

	return hoverfly
}


func setConfigurationDefaults() {
	viper.AddConfigPath("./.hoverfly")
	viper.AddConfigPath("$HOME/.hoverfly")
	viper.SetDefault("hoverfly.host", "localhost")
	viper.SetDefault("hoverfly.admin.port", "8888")
	viper.SetDefault("hoverfly.proxy.port", "8500")
	viper.SetDefault("specto.lab.host", "localhost")
	viper.SetDefault("specto.lab.port", "81")
}

func getHoverflyDirectory(configUri string) string {
	if len(configUri) == 0 {
		fmt.Println("Missing a config file")
		fmt.Println("Creating a new  a config file")

		hoverflyDir, err := createHomeDirectory()

		if err != nil {
			failAndExit(err)
		}

		return hoverflyDir
	}

	return path.Dir(configUri)
}