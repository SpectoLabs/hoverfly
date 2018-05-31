package cmd

import (
	"fmt"

	"os"

	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var stateCmd = &cobra.Command{
	Use:     "state",
	Aliases: []string{"state-store"},
	Short:   "Manage the state for Hoverfly",
	Long: `
This allows you to inspect and modify the
state stored in Hoverfly. The state is a map
of string keys and values which can be used
for matching.
	`,
}

var getAllStateCmd = &cobra.Command{
	Use:   "get-all",
	Short: "Gets all of the the state",
	Long: `
Returns all of the state keys and their values from Hoverfly.
	`,
	Run: func(cmd *cobra.Command, args []string) {

		checkTargetAndExit(target)

		if len(args) == 0 {
			currentState, err := wrapper.GetCurrentState(*target)
			handleIfError(err)

			output := ""
			for k, v := range currentState {
				output = output + "\n\"" + k + "\"=\"" + v + "\""
			}

			if len(output) < 3 {
				fmt.Println("The state for Hoverfly is empty")
			} else {
				fmt.Println("State of Hoverfly:" + output)
			}
		}
	},
}

var getStateCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets the state of a single key",
	Long: `
Returns the state of Hoverfly by key. 

Provide a single argument, the state key.
	`,
	Run: func(cmd *cobra.Command, args []string) {

		checkTargetAndExit(target)

		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "You must provide a state key as an argument")
			fmt.Fprintln(os.Stderr, "\nTry hoverctl state-store get --help for more information")
			os.Exit(1)
		}

		key := args[0]
		currentState, err := wrapper.GetCurrentState(*target)
		handleIfError(err)
		state := currentState[key]

		if len(state) == 0 {
			fmt.Println("State is not set for the key:", key)
		} else {
			fmt.Printf("State of \"%s\":\n%s", key, state)
		}
	},
}

var setStateCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets the state",
	Long: `
Sets the state of Hoverfly by key.

Provide two arguments, the state key and the state value,
separated by a space.
	`,
	Run: func(cmd *cobra.Command, args []string) {

		checkTargetAndExit(target)

		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "You must only provide a key and a value, separated by a space")
			fmt.Fprintln(os.Stderr, "\nTry hoverctl state-store set --help for more information")
			os.Exit(1)
		}

		err := wrapper.PatchCurrentState(*target, args[0], args[1])
		handleIfError(err)
		fmt.Println("Successfully set state key and value:\n" + "\"" + args[0] + "\"=\"" + args[1] + "\"")
	},
}

var deleteStateCmd = &cobra.Command{
	Use:   "delete-all",
	Short: "Deletes all state",
	Long: `
Deletes the  state of Hoverfly. 

Provide two arguments, the state key and the state value.
	`,
	Run: func(cmd *cobra.Command, args []string) {

		checkTargetAndExit(target)

		err := wrapper.DeleteCurrentState(*target)
		handleIfError(err)
		fmt.Println("State has been deleted")
	},
}

func init() {
	RootCmd.AddCommand(stateCmd)
	stateCmd.AddCommand(getStateCmd)
	stateCmd.AddCommand(getAllStateCmd)
	stateCmd.AddCommand(setStateCmd)
	stateCmd.AddCommand(deleteStateCmd)
}
