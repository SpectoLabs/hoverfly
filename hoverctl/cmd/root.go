package cmd

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var adminPort, proxyPort, host, certificate, key, database, upstreamProxy string
var disableTls, verbose bool

var hoverfly wrapper.Hoverfly
var hoverflyDirectory wrapper.HoverflyDirectory
var config *wrapper.Config

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
	RootCmd.PersistentFlags().StringVar(&upstreamProxy, "upstream-proxy", "",
		"A host for which Hoverfly will proxy its requests to")

	RootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Verbose logging from hoverctl")
	RootCmd.Flag("verbose").Shorthand = "v"
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

	var err error
	hoverflyDirectory, err = wrapper.NewHoverflyDirectory(*config)
	handleIfError(err)

	hoverfly = wrapper.NewHoverfly(*config)
}

func handleIfError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func checkArgAndExit(args []string, message, command string) {
	if len(args) == 0 {
		fmt.Println(message)
		fmt.Println("\nTry hoverctl " + command + " --help for more information")
		os.Exit(1)
	}
}

func drawTable(data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	// table.SetHeader([]string{"Name", "Sign", "Rating"})

	for _, v := range data {
		table.Append(v)
	}
	fmt.Print("\n")
	table.Render()
}
