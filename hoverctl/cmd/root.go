package cmd

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
	"github.com/spf13/cobra"
)

var targetNameFlag, hostFlag string
var adminPortFlag, proxyPortFlag int

var force, verbose, setDefaultTargetFlag bool

var hoverflyDirectory configuration.HoverflyDirectory
var config *configuration.Config
var target *configuration.Target

var version string

var RootCmd = &cobra.Command{
	Use:   "hoverctl",
	Short: "hoverctl is the command line tool for Hoverfly",
	Long:  ``,
}

func Execute(hoverctlVersion string) {
	version = hoverctlVersion

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if setDefaultTargetFlag && target != nil {
		config.DefaultTarget = target.Name
	}
	handleIfError(config.WriteToFile(hoverflyDirectory))
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolVar(&force, "force", false,
		"Bypass any confirmation when using hoverctl")
	RootCmd.Flag("force").Shorthand = "f"

	RootCmd.PersistentFlags().StringVar(&targetNameFlag, "target", "",
		"A name for an instance of Hoverfly you are trying to communicate with. Overrides the default target (default)")
	RootCmd.PersistentFlags().BoolVar(&setDefaultTargetFlag, "set-default", false,
		"Sets the current target as the default target for hoverctl")

	RootCmd.PersistentFlags().IntVar(&adminPortFlag, "admin-port", 0,
		"A port number for the Hoverfly API/GUI. Overrides the default Hoverfly admin port (8888)")
	RootCmd.PersistentFlags().IntVar(&proxyPortFlag, "proxy-port", 0,
		"A port number for the Hoverfly proxy. Overrides the default Hoverfly proxy port (8500)")
	RootCmd.PersistentFlags().StringVar(&hostFlag, "host", "",
		"A host on which a Hoverfly instance is running. Overrides the default Hoverfly host (localhost)")

	RootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Verbose logging from hoverctl")

	RootCmd.Flag("verbose").Shorthand = "v"
	RootCmd.Flag("target").Shorthand = "t"
}

func initConfig() {

	log.SetOutput(os.Stdout)
	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	configuration.SetConfigurationDefaults()
	configuration.SetConfigurationPaths()

	config = configuration.GetConfig()

	if config.GetTarget(config.DefaultTarget) == nil {
		fmt.Printf("Default target `%v` not found, changing default target to `local`", config.DefaultTarget)
		config.DefaultTarget = "local"

	}

	target = config.GetTarget(targetNameFlag)
	if targetNameFlag == "" && target == nil {
		target = configuration.NewDefaultTarget()
	}

	if verbose && target != nil {
		fmt.Println("Current target: " + target.Name + "\n")
	}

	var err error
	hoverflyDirectory, err = configuration.NewHoverflyDirectory(*config)
	handleIfError(err)
}
