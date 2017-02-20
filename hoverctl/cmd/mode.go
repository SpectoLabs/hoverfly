package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var modeCmd = &cobra.Command{
	Use:   "mode [capture|simulate|modify|synthesize (optional)]",
	Short: "Get and set the Hoverfly mode",
	Long: `
Sets Hoverfly to the mode specified. The mode
determines how Hoverfly will process incoming
requests.

If a mode is not specified, the current Hoverfly 
mode is shown.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			mode, err := hoverfly.GetMode()
			handleIfError(err)

			fmt.Println("Hoverfly is currently set to", mode, "mode")
		} else {
			mode, err := hoverfly.SetMode(args[0])
			handleIfError(err)

			fmt.Println("Hoverfly has been set to", mode, "mode")
		}
	},
}

func init() {
	RootCmd.AddCommand(modeCmd)
}
