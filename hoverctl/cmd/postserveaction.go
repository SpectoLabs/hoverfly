package cmd

import (
	"fmt"

	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var binary, scriptPath, actionNameToBeSet, actionNameToBeDeleted, remote string
var delayInMs int

var postServeActionCommand = &cobra.Command{
	Use:   "post-serve-action",
	Short: "Manage the post-serve-action for Hoverfly",
	Long: `
		This allows you to manage post serve action in Hoverfly. 
	`,
}

var postServeActionGetCommand = &cobra.Command{
	Use:   "get-all",
	Short: "Get all post serve actions for Hoverfly",
	Long:  `Get all post serve actions for Hoverfly`,
	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)

		if len(args) == 0 {
			postServeActions, err := wrapper.GetAllPostServeActions(*target)
			handleIfError(err)
			localPostServeActionData, remotePostServeActionData := getPostServeActionsTabularData(postServeActions)
			drawTable(localPostServeActionData, true)
			drawTable(remotePostServeActionData, true)
		}
	},
}

var postServeActionSetCommand = &cobra.Command{
	Use:   "set",
	Short: "Set postServeAction for Hoverfly",
	Long: `
Hoverfly Local PostServeAction can be set using the following flags: 
	 --name --binary --script --delay
Hoverfly Remote PostServeAction can be set using the following flags:
	 --name --remote --delay
`,
	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)
		if remote != "" && actionNameToBeSet != "" {
			err := wrapper.SetRemotePostServeAction(actionNameToBeSet, remote, delayInMs, *target)
			handleIfError(err)
			fmt.Println("Success")
		} else if binary != "" && scriptPath != "" && actionNameToBeSet != "" {
			script, err := configuration.ReadFile(scriptPath)
			handleIfError(err)
			err = wrapper.SetLocalPostServeAction(actionNameToBeSet, binary, string(script), delayInMs, *target)
			handleIfError(err)
			fmt.Println("Success")
		} else {
			fmt.Println("(Binary and script path/remote) and action name are compulsory to set post serve action")
		}
	},
}

var postServeActionDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete postServeAction for Hoverfly",
	Long: `
Hoverfly PostServeAction can be deleted using the following flags: 
	 --name
`,
	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)
		if actionNameToBeDeleted == "" {
			fmt.Println("action name to be deleted not provided")
		} else {
			err := wrapper.DeletePostServeAction(actionNameToBeDeleted, *target)
			handleIfError(err)
			fmt.Println("Success")
		}
	},
}

func init() {
	RootCmd.AddCommand(postServeActionCommand)
	postServeActionCommand.AddCommand(postServeActionGetCommand)
	postServeActionCommand.AddCommand(postServeActionSetCommand)
	postServeActionCommand.AddCommand(postServeActionDeleteCommand)

	postServeActionSetCommand.PersistentFlags().StringVar(&actionNameToBeSet, "name", "", "Action Name to be set")
	postServeActionSetCommand.PersistentFlags().StringVar(&binary, "binary", "",
		"An absolute or relative path to a binary that Hoverfly will execute as post serve action")
	postServeActionSetCommand.PersistentFlags().StringVar(&scriptPath, "script", "",
		"An absolute or relative path to a script that will be executed by the binary")
	postServeActionSetCommand.PersistentFlags().IntVar(&delayInMs, "delay", 0, "Delay in milli seconds after which action needs to be executed")
	postServeActionSetCommand.PersistentFlags().StringVar(&remote, "remote", "", "Remote host to be set for triggering post serve action")

	postServeActionDeleteCommand.PersistentFlags().StringVar(&actionNameToBeDeleted, "name", "", "Action Name to be deleted")

}

func getPostServeActionsTabularData(postServeActions v2.PostServeActionDetailsView) ([][]string, [][]string) {
	localPostServeActionsData := [][]string{{"Action Name", "Binary", "Script", "Delay(Ms)"}}
	remotePostServeActionData := [][]string{{"Action Name", "Remote", "Delay(Ms)"}}
	for _, action := range postServeActions.Actions {
		if action.Remote == "" {
			actionData := []string{action.ActionName, action.Binary, getContentShorthand(action.ScriptContent), fmt.Sprint(action.DelayInMs)}
			localPostServeActionsData = append(localPostServeActionsData, actionData)
		} else {
			actionData := []string{action.ActionName, action.Remote, fmt.Sprint(action.DelayInMs)}
			remotePostServeActionData = append(remotePostServeActionData, actionData)
		}
	}
	return localPostServeActionsData, remotePostServeActionData
}
