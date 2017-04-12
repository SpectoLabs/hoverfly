package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var targetNameFlag, hostFlag string
var adminPortFlag, proxyPortFlag int

var force, verbose, setDefaultTargetFlag bool

var hoverflyDirectory wrapper.HoverflyDirectory
var config *wrapper.Config
var target *wrapper.Target

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

	if setDefaultTargetFlag && targetNameFlag != "" {
		config.DefaultTarget = targetNameFlag
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

	wrapper.SetConfigurationDefaults()
	wrapper.SetConfigurationPaths()

	config = wrapper.GetConfig()

	target = config.GetTarget(targetNameFlag)
	if targetNameFlag == "" && target == nil {
		target = wrapper.NewDefaultTarget()
	}

	if verbose && target != nil {
		fmt.Println("Current target: " + target.Name + "\n")
	}

	var err error
	hoverflyDirectory, err = wrapper.NewHoverflyDirectory(*config)
	handleIfError(err)
}

func handleIfError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func checkArgAndExit(args []string, message, command string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, message)
		fmt.Fprintln(os.Stderr, "\nTry hoverctl "+command+" --help for more information")
		os.Exit(1)
	}
}

func checkTargetAndExit(target *wrapper.Target) {
	if target == nil {
		handleIfError(fmt.Errorf("%[1]s is not a target\n\nRun `hoverctl targets new %[1]s`", targetNameFlag))
	}
}

func askForConfirmation(message string) bool {
	if force {
		return true
	}

	for {
		response := askForInput(message+" [y/n]", false)

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func askForInput(value string, sensitive bool) string {
	if force {
		return ""
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf(value + ": ")
		if sensitive {
			responseBytes, err := terminal.ReadPassword(0)
			handleIfError(err)
			fmt.Println("")

			return strings.TrimSpace(string(responseBytes))
		} else {
			response, err := reader.ReadString('\n')
			handleIfError(err)

			return strings.TrimSpace(response)
		}
	}

	return ""
}

func drawTable(data [][]string, header bool) {
	table := tablewriter.NewWriter(os.Stdout)
	if header {
		table.SetHeader(data[0])
		data = data[1:]
	}

	for _, v := range data {
		table.Append(v)
	}
	fmt.Print("\n")
	table.Render()
}
