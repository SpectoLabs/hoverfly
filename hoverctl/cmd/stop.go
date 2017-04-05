package cmd

import (
	"fmt"

	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Hoverfly",
	Long: `
Stops Hoverfly.

The "pid" file created in the ".hoverfly" directory by the
"start" command will be deleted when the instance of Hoverfly
is stopped.
`,

	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target, "Cannot stop an instance of Hoverfly without a target")

		err := wrapper.Stop(target, hoverflyDirectory)
		handleIfError(err)

		config.NewTarget(*target)
		handleIfError(config.WriteToFile(hoverflyDirectory))

		fmt.Println("Hoverfly has been stopped")
	},
}

func init() {
	RootCmd.AddCommand(stopCmd)
}
