package cmd

import (
	"fmt"

	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
	"os"
)

// currentStateCmd represents the flush command
var currentStateCmd = &cobra.Command{
	Use:   "current-state",
	Short: "Manage the current state for Hoverfly",
	Long: `
This allows you to inspect and modify the current
state of Hoverfly. By current state, we mean the set
of keys and values which are set and matched against
during matching.
	`,
}

var getAllCurrentStateCommand = &cobra.Command{
	Use:   "get-all",
	Short: "Gets the current state",
	Long: `
Returns all of the current state keys and values of
Hoverfly. By current state, we mean the set of keys
and values which are set and matched against during matching.
	`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			currentState, err := wrapper.GetCurrentState(*target)
			handleIfError(err)

			output := ""
			for k, v := range currentState {
				output = output + k + " " + v
			}

			if len(output) == 0 {
				fmt.Println("The current state for Hoverfly is empty")
			} else {
				fmt.Println("Current state of Hoverfly:\n", output)
			}
		} else {
			fmt.Fprintln(os.Stderr, "This command should not take an argument")
			fmt.Fprintln(os.Stderr, "\nTry hoverctl current-state get-all --help for more information")
			os.Exit(1)
		}
	},
}

var getCurrentStateCommand = &cobra.Command{
	Use:   "get",
	Short: "Gets the current state of a single key",
	Long: `
Returns the current state of Hoverfly by key. By
current state, we mean the set of keys and values which
are set and matched against during matching.

Provide a single argument, the state key.
	`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "You must provide a state key as an argument")
			fmt.Fprintln(os.Stderr, "\nTry hoverctl current-state get --help for more information")
			os.Exit(1)
		}

		key := args[0]
		currentState, err := wrapper.GetCurrentState(*target)
		handleIfError(err)
		state := currentState[key]

		if len(state) == 0 {
			fmt.Println("State is not set for the key:", key)
		} else {
			fmt.Printf("Current state of %s:\n%s", key, state)
		}
	},
}

var setCurrentStateCommand = &cobra.Command{
	Use:   "set",
	Short: "Sets the current state",
	Long: `
Sets the current state of Hoverfly. By current state, we
mean the set of keys and values which are set and matched
against during matching.

Provide two arguments, the state key and the state value,
separated by a space.
	`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "You must only provide a key and a value, separated by a space")
			fmt.Fprintln(os.Stderr, "\nTry hoverctl current-state set --help for more information")
			os.Exit(1)
		}

		err := wrapper.PatchCurrentState(*target, args[0], args[1])
		handleIfError(err)
		fmt.Println("Successfully set current-state key and value:\n", args[0], args[1])
	},
}

var deleteCurrentStateCommand = &cobra.Command{
	Use:   "delete",
	Short: "Deletes current state",
	Long: `
Deletes the current state of Hoverfly. By current state, we
mean the set of keys and values which are set and matched
against during matching.

Provide two arguments, the state key and the state value.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		err := wrapper.DeleteCurrentState(*target)
		handleIfError(err)
		fmt.Println("Current state has been deleted")
	},
}

func init() {
	RootCmd.AddCommand(currentStateCmd)
	currentStateCmd.AddCommand(getCurrentStateCommand)
	currentStateCmd.AddCommand(getAllCurrentStateCommand)
	currentStateCmd.AddCommand(setCurrentStateCommand)
	currentStateCmd.AddCommand(deleteCurrentStateCommand)
}
