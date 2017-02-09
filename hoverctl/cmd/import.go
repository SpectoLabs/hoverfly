package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var importV1 bool

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
		simulationData, err := wrapper.ReadFile(args[0])
		handleIfError(err)

		err = hoverfly.ImportSimulation(string(simulationData), importV1)
		handleIfError(err)

		log.Info("Successfully imported simulation from ", args[0])
	},
}

func init() {
	RootCmd.AddCommand(importCmd)
	importCmd.Flags().BoolVar(&importV1, "v1", false, "Tells Hoverfly that the simulation is formatted according to the old v1 simulation JSON schema used in Hoverfly pre v0.9.0")
}
