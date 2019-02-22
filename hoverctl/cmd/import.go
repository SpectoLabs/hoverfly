package cmd

import (
	"fmt"

	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import [path to simulation]",
	Short: "Import a simulation into Hoverfly",
	Long: `
Imports a simulation into Hoverfly. An absolute or
relative path to a Hoverfly simulation JSON file
must be provided.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)

		checkArgAndExit(args, "You have not provided a path to simulation", "import")
		simulationData, err := configuration.ReadFile(args[0])
		handleIfError(err)

		err = wrapper.ImportSimulation(*target, string(simulationData))
		handleIfError(err)

		fmt.Println("Successfully imported simulation from", args[0])
	},
}

func init() {
	RootCmd.AddCommand(importCmd)
}
