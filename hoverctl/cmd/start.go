package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start an instance of Hoverfly",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			config.SetWebserver(args[0])
			hoverfly = wrapper.NewHoverfly(*config)
		}

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
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
}
