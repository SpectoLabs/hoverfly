package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop an instance of Hoverfly",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		err := hoverfly.Stop(hoverflyDirectory)
		handleIfError(err)

		log.Info("Hoverfly has been stopped")
	},
}

func init() {
	RootCmd.AddCommand(stopCmd)
}
