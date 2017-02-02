package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Import a simulation into Hoverfly",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		simulationData, err := hoverfly.ExportSimulation()
		handleIfError(err)

		err = wrapper.WriteFile(args[0], simulationData)
		handleIfError(err)

		log.Info("Successfully exported to ", args[0])
	},
}

func init() {
	RootCmd.AddCommand(exportCmd)
}
