package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var importV1 bool

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a simulation into Hoverfly",
	Long: `Will import a simulation into Hoverfly. The simulation
provided should be an absolute or relative path to a 
JSON file.
	
	hoverctl import simulation.json

	hoverctl import /home/user/simulation.json
	`,

	Run: func(cmd *cobra.Command, args []string) {
		simulationData, err := wrapper.ReadFile(args[0])
		handleIfError(err)

		err = hoverfly.ImportSimulation(string(simulationData), importV1)
		handleIfError(err)

		log.Info("Successfully imported from ", args[0])
	},
}

func init() {
	RootCmd.AddCommand(importCmd)
	importCmd.Flags().BoolVar(&importV1, "v1", false, "This flag can be used to import old simulations from old v1 style schema used pre Hoverfly v0.9.0")
}
