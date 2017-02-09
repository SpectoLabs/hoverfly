package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete simulation in Hoverfly",
	Long: `
Will delete simulation from Hoverfly.
`,

	Run: func(cmd *cobra.Command, args []string) {
		err := hoverfly.DeleteSimulations()
		handleIfError(err)

		log.Info("Simulations have been deleted from Hoverfly")
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
