package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Hoverfly",
	Long: `
Stops Hoverfly.

The "pid" file created in the ".hoverfly" directory by the
"start" command will be deleted when the instance of Hoverfly
is stopped.
`,

	Run: func(cmd *cobra.Command, args []string) {
		err := hoverfly.Stop(hoverflyDirectory)
		handleIfError(err)

		log.Info("Hoverfly has been stopped")
	},
}

func init() {
	RootCmd.AddCommand(stopCmd)
}
