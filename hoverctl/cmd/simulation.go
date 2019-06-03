package cmd

import (
	"fmt"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
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
	Use:   "add",
	Short: "Add simulations into Hoverfly",
	Long: `Appends simulations to existing Hoverfly 
simulation. You may provide absolute or relative paths 
to multiple Hoverfly simulation JSON files. Any pairs that 
have identical requests to those in the existing data 
will be ignored.
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
