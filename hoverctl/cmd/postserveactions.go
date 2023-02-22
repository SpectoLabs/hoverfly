package cmd

import (
	"fmt"

	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var binary, scriptPath, actionNameToBeSet, actionNameToBeDeleted string
var delayInMs int

var postServeActionCommand = &cobra.Command{
	Use:   "post-serve-actions",
	Short: "Get, Set & Delete Hoverfly PostServeAction",
	Long: `
Hoverfly PostServeAction can be set using the following flags: 
	--set --binary --script --delay

Hoverfly PostServeAction can be deleted by passing action name to be deleted:
	--delete
`,
	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)
		if isNoOptionsPassedForPostServeAction() {
			postServeActions, err := wrapper.GetAllPostServeActions(*target)
			handleIfError(err)
			drawTable(getPostServeActionsTabularData(postServeActions), true)

		} else if actionNameToBeDeleted != "" {
			err := wrapper.DeletePostServeAction(actionNameToBeDeleted, *target)
			handleIfError(err)
			fmt.Println("Success")
		} else {
			if binary == "" || scriptPath == "" || actionNameToBeSet == "" {
				fmt.Println("Binary, script path and action name are compulsory to set post serve action")
			} else {
				script, err := configuration.ReadFile(scriptPath)
				handleIfError(err)
				err = wrapper.SetPostServeAction(actionNameToBeSet, binary, string(script), delayInMs, *target)
				handleIfError(err)
				fmt.Println("Success")
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(postServeActionCommand)
	postServeActionCommand.PersistentFlags().StringVar(&actionNameToBeSet, "set", "", "Action Name to be set")
	postServeActionCommand.PersistentFlags().StringVar(&actionNameToBeDeleted, "delete", "", "Action Name to be deleted")
	postServeActionCommand.PersistentFlags().StringVar(&binary, "binary", "",
		"An absolute or relative path to a binary that Hoverfly will execute as post serve action")
	postServeActionCommand.PersistentFlags().StringVar(&scriptPath, "script", "",
		"An absolute or relative path to a script that will be executed by the binary")
	postServeActionCommand.PersistentFlags().IntVar(&delayInMs, "delay", 0, "Delay in milli seconds after which action needs to be executed")

}

func isNoOptionsPassedForPostServeAction() bool {

	return actionNameToBeSet == "" && actionNameToBeDeleted == "" && binary == "" && scriptPath == "" && delayInMs == 0
}

func getPostServeActionsTabularData(postServeActions v2.PostServeActionDetailsView) [][]string {

	postServeActionsData := [][]string{{"Action Name", "Binary", "Script", "Delay(Ms)"}}
	for _, action := range postServeActions.Actions {
		actionData := []string{action.ActionName, action.Binary, getScriptShorthand(action.ScriptContent), fmt.Sprint(action.DelayInMs)}
		postServeActionsData = append(postServeActionsData, actionData)
	}
	return postServeActionsData
}
