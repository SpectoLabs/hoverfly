package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
	"github.com/spf13/cobra"
)

var targetNameFlag string

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
