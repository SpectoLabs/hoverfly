package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var force bool

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Hoverfly simulation",
	Long: `
Deletes simulation data from the Hoverfly instance.
`,

	Run: func(cmd *cobra.Command, args []string) {
		if !force {
			if !askForConfirmation() {
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

func askForConfirmation() bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("Are you sure you want to delete the current simulation? [y/n]: ")

		response, err := reader.ReadString('\n')
		handleIfError(err)

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
