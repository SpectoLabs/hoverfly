package cmd

import (
	"fmt"

	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export [path to simulation]",
	Short: "Export a simulation from Hoverfly",
	Long: `
Exports a simulation from Hoverfly. The simulation JSON
will be written to the file path provided.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		checkArgAndExit(args, "You have not provided a path to simulation", "export")
		simulationData, err := hoverfly.ExportSimulation()
		handleIfError(err)

		err = wrapper.WriteFile(args[0], simulationData)
		handleIfError(err)

		fmt.Println("Successfully exported simulation to", args[0])
	},
}

func init() {
	RootCmd.AddCommand(exportCmd)
}
