package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export [path to simulation]",
	Short: "export a simulation from Hoverfly",
	Long: `
Exports a simulation from Hoverfly. The simulation JSON
will be written to the file path provided.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		simulationData, err := hoverfly.ExportSimulation()
		handleIfError(err)

		err = wrapper.WriteFile(args[0], simulationData)
		handleIfError(err)

		log.Info("Successfully exported simulation to ", args[0])
	},
}

func init() {
	RootCmd.AddCommand(exportCmd)
}
