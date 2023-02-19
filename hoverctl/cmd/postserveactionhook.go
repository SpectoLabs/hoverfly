package cmd

import (
	"fmt"

	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var binary, scriptPath, hookName string
var delayInMilliSeconds int

var postServeActionHookCommand = &cobra.Command{
	Use:   "postserveactionhook",
	Short: "Set or Delete Hoverfly PostServeAction Hook",
	Long: `
Hoverfly PostServe Action can be set using the following
combinations of flags: 
	--hookName --binary --script --delay

`,
	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)
		if binary == "" || scriptPath == "" || hookName == "" {
			fmt.Println("Binary, script path and hookname are compulsory to set post serve action hook")
		} else {
			script, err := configuration.ReadFile(scriptPath)
			handleIfError(err)
			err = wrapper.SetPostServeActionHook(hookName, binary, string(script), delayInMilliSeconds, *target)
			handleIfError(err)
			fmt.Println("Success")
		}
	},
}

func init() {
	RootCmd.AddCommand(postServeActionHookCommand)
	postServeActionHookCommand.PersistentFlags().StringVar(&hookName, "hookName", "", "Hook Name")
	postServeActionHookCommand.PersistentFlags().StringVar(&binary, "binary", "",
		"An absolute or relative path to a binary that Hoverfly will execute as post serve action hook")
	postServeActionHookCommand.PersistentFlags().StringVar(&scriptPath, "script", "",
		"An absolute or relative path to a script that will be executed by the binary")
	postServeActionHookCommand.PersistentFlags().IntVar(&delayInMilliSeconds, "delay", 0, "Delay after which hook needs to be executed")

}
