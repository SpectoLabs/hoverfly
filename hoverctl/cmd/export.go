package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var urlPattern string
var exportCmd = &cobra.Command{
	Use:   "export [path to simulation]",
	Short: "Export a simulation from Hoverfly",
	Long: `
Exports a simulation from Hoverfly. The simulation JSON
will be written to the file path provided.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)

		checkArgAndExit(args, "You have not provided a path to simulation", "export")

		simulationView, err := wrapper.ExportSimulation(*target, urlPattern)
		handleIfError(err)

		for i, pair := range simulationView.DataViewV6.RequestResponsePairs {
			bodyFile := pair.Response.GetBodyFile()
			if len(bodyFile) == 0 {
				continue
			}

			if err := configuration.WriteFile(bodyFile, []byte(pair.Response.GetBody())); err != nil {
				handleIfError(err)
			}

			simulationView.DataViewV6.RequestResponsePairs[i].Response.Body = ""
		}

		simulationData, err := json.MarshalIndent(simulationView, "", "\t")
		handleIfError(err)

		err = configuration.WriteFile(args[0], simulationData)
		handleIfError(err)

		fmt.Println("Successfully exported simulation to", args[0])
	},
}

func init() {
	RootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVar(&urlPattern, "url-pattern", "", "Export simulation for the urls that matches a pattern, eg. foo.com/api/v(.+)")
}
