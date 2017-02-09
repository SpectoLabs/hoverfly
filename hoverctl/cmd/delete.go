package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Hoverfly simulation",
	Long: `
Deletes simulation data from the Hoverfly instance.
`,

	Run: func(cmd *cobra.Command, args []string) {
		err := hoverfly.DeleteSimulations()
		handleIfError(err)

		log.Info("Simulation data has been deleted from Hoverfly")
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
