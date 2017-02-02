package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var modeCmd = &cobra.Command{
	Use:   "mode",
	Short: "Get and set the mode of Hoverfly",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			mode, err := hoverfly.GetMode()
			handleIfError(err)

			log.Info("Hoverfly is set to ", mode, " mode")
		} else {
			mode, err := hoverfly.SetMode(args[0])
			handleIfError(err)

			log.Info("Hoverfly has been set to ", mode, " mode")
		}
	},
}

func init() {
	RootCmd.AddCommand(modeCmd)
}
