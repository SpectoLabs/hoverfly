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

var targetName, adminPort, proxyPort, host, certificate, key, database, upstreamProxy string
var disableTls, verbose bool

var force bool

var cacheDisable bool

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
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolVar(&force, "force", false,
		"Bypass any confirmation when using hoverctl")
	RootCmd.Flag("force").Shorthand = "f"

	RootCmd.PersistentFlags().StringVar(&targetName, "target", "",
		"A name for an instance of Hoverfly you are trying to communicate with. Overrides the default target (default)")
	RootCmd.PersistentFlags().StringVar(&adminPort, "admin-port", "",
		"A port number for the Hoverfly API/GUI. Overrides the default Hoverfly admin port (8888)")
	RootCmd.PersistentFlags().StringVar(&proxyPort, "proxy-port", "",
		"A port number for the Hoverfly proxy. Overrides the default Hoverfly proxy port (8500)")
	RootCmd.PersistentFlags().StringVar(&host, "host", "",
		"A host on which a Hoverfly instance is running. Overrides the default Hoverfly host (localhost)")
	RootCmd.PersistentFlags().StringVar(&certificate, "certificate", "",
		"A path to a certificate file. Overrides the default Hoverfly certificate")
	RootCmd.PersistentFlags().StringVar(&key, "key", "",
		"A path to a key file. Overrides the default Hoverfly TLS key")
	RootCmd.PersistentFlags().BoolVar(&disableTls, "disable-tls", false,
		"Disables TLS verification")
	RootCmd.PersistentFlags().StringVar(&database, "database", "",
		"A database type [memory|boltdb]. Overrides the default Hoverfly database type (memory)")
	RootCmd.PersistentFlags().BoolVar(&cacheDisable, "disable-cache", false,
		"?")
	RootCmd.PersistentFlags().StringVar(&upstreamProxy, "upstream-proxy", "",
		"A host for which Hoverfly will proxy its requests to")

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
	config = config.SetHost(host)
	config = config.SetAdminPort(adminPort)
	config = config.SetProxyPort(proxyPort)
	config = config.SetUsername("")
	config = config.SetPassword("")
	config = config.SetCertificate(certificate)
	config = config.SetKey(key)
	config = config.DisableTls(disableTls)
	config = config.SetDbType(database)
	config = config.SetUpstreamProxy(upstreamProxy)
	config = config.DisableCache(cacheDisable)

	target = config.GetTarget(targetName)
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

func checkTargetAndExit(target *wrapper.Target, error string) {
	if target == nil {
		fmt.Fprintln(os.Stderr, error)
		os.Exit(1)
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
