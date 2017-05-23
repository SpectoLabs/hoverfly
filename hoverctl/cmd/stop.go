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
		checkTargetAndExit(target)

		if !wrapper.IsLocal(target.Host) {
			handleIfError(fmt.Errorf("Unable to stop an instance of Hoverfly on a remote host (%s host: %s)\n\nRun `hoverctl start --new-target <name>` to start it", target.Name, target.Host))
		}

		err := wrapper.CheckIfRunning(*target)
		if err != nil {
			handleIfError(err)
		}

		err = wrapper.Stop(*target)
		handleIfError(err)

		config.NewTarget(*target)
		handleIfError(config.WriteToFile(hoverflyDirectory))

		fmt.Println("Hoverfly has been stopped")
	},
}

func init() {
	RootCmd.AddCommand(stopCmd)
}
