package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Hoverfly simulation",
	Long: `
Deletes simulation data from the Hoverfly instance.
`,

	Run: func(cmd *cobra.Command, args []string) {
		if !force {
			if !askForConfirmation("Are you sure you want to delete the current simulation?") {
				return
			}
		}
		err := hoverfly.DeleteSimulations()
		handleIfError(err)

		fmt.Println("Simulation data has been deleted from Hoverfly")
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&force, "force", false,
		"Delete the simulation without prompting for confirmation")
}
