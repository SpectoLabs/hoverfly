package cmd

import (
	"fmt"
	"github.com/SpectoLabs/hoverfly/v2/hoverctl/configuration"
	"github.com/SpectoLabs/hoverfly/v2/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var simulationCmd = &cobra.Command{
	Use:   "simulation",
	Short: "Manage the simulation for Hoverfly",
	Long: `
This allows you to manage simulation data in Hoverfly. 
	`,
}

var addSimulationCmd = &cobra.Command{
	Use:   "add [path to simulations]",
	Short: "Add one or more simulations into Hoverfly",
	Long: `
Adds one or more simulation files to Hoverfly to the 
existing simulation data. 

Any request/response pairs that have an identical request 
to those in the existing data will be discarded with a 
warning message. 

You may provide an absolute or relative path to each 
simulation file.
	`,
	Run: func(cmd *cobra.Command, args []string) {

		checkTargetAndExit(target)

		checkArgAndExit(args, "You have not provided a path to simulation", "simulation add")

		for _, arg := range args {

			simulationData, err := configuration.ReadFile(arg)
			handleIfError(err)

			err = wrapper.AddSimulation(*target, string(simulationData))
			handleIfError(err)
			fmt.Println("Successfully added simulation from", arg)
		}

	},
}

func init() {
	RootCmd.AddCommand(simulationCmd)
	simulationCmd.AddCommand(addSimulationCmd)
}
