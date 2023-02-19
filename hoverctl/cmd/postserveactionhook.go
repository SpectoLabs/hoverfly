package cmd

import (
	"fmt"

	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var binary, scriptPath, hookNameToBeAdded, hookNameToBeDeleted string
var delayInMilliSeconds int

var postServeActionHookCommand = &cobra.Command{
	Use:   "postserveactionhook",
	Short: "Set or Delete Hoverfly PostServeAction Hook",
	Long: `
Hoverfly PostServeAction hook can be set using the following flags: 
	--add --binary --script --delay

Hoverfly PostServeAction hook can be deleted by passing hook name to be deleted:
	--delete
`,
	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)
		if hookNameToBeDeleted != "" {
			err := wrapper.DeletePostServeActionHook(hookNameToBeDeleted, *target)
			handleIfError(err)
			fmt.Println("Success")
		} else {
			if binary == "" || scriptPath == "" || hookNameToBeAdded == "" {
				fmt.Println("Binary, script path and hookname are compulsory to set post serve action hook")
			} else {
				script, err := configuration.ReadFile(scriptPath)
				handleIfError(err)
				err = wrapper.SetPostServeActionHook(hookNameToBeAdded, binary, string(script), delayInMilliSeconds, *target)
				handleIfError(err)
				fmt.Println("Success")
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(postServeActionHookCommand)
	postServeActionHookCommand.PersistentFlags().StringVar(&hookNameToBeAdded, "add", "", "Hook Name to be added")
	postServeActionHookCommand.PersistentFlags().StringVar(&hookNameToBeDeleted, "delete", "", "Hook Name to be deleted")
	postServeActionHookCommand.PersistentFlags().StringVar(&binary, "binary", "",
		"An absolute or relative path to a binary that Hoverfly will execute as post serve action hook")
	postServeActionHookCommand.PersistentFlags().StringVar(&scriptPath, "script", "",
		"An absolute or relative path to a script that will be executed by the binary")
	postServeActionHookCommand.PersistentFlags().IntVar(&delayInMilliSeconds, "delay", 0, "Delay in milliseconds after which hook needs to be executed")

}
