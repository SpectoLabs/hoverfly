package cmd

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var adminPort, proxyPort, host, certificate, key, database string
var disableTls, verbose bool

var hoverfly wrapper.Hoverfly
var hoverflyDirectory wrapper.HoverflyDirectory
var config *wrapper.Config

var RootCmd = &cobra.Command{
	Use:   "hoverctl",
	Short: "hoverctl is the command line tool for Hoverfly",
	Long:  ``,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&adminPort, "admin-port", "",
		"Given a port number, the port number is used to override the default Hoverfly admin port (8888)")
	RootCmd.PersistentFlags().StringVar(&proxyPort, "proxy-port", "",
		"Given a port number, the port number is used to override the default Hoverfly proxy port (8500)")
	RootCmd.PersistentFlags().StringVar(&host, "host", "",
		"Given a host, the host is used to override the default Hoverfly host (localhost)")
	RootCmd.PersistentFlags().StringVar(&certificate, "certificate", "",
		"Given a path, the certificate is used to override the default Hoverfly certificate")
	RootCmd.PersistentFlags().StringVar(&key, "key", "",
		"Given a path, the key is used to override the default Hoverfly TLS key")
	RootCmd.PersistentFlags().BoolVar(&disableTls, "disable-tls", false,
		"Disable TLS verification")
	RootCmd.PersistentFlags().StringVar(&database, "database", "",
		"Given a database type [memory|boltdb], the database type is used to override the default Hoverfly database type (memory)")

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

	var err error
	hoverflyDirectory, err = wrapper.NewHoverflyDirectory(*config)
	handleIfError(err)

	hoverfly = wrapper.NewHoverfly(*config)
}

func handleIfError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
