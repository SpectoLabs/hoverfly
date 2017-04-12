package cmd

import (
	"errors"
	"strconv"

	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var pidFlag int

var targetsCmd = &cobra.Command{
	Use:   "targets",
	Short: "Get the current targets registered with hoverctl",
	Long: `
Get the current targets registered with hoverctl"
`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(config.Targets) == 0 {
			handleIfError(errors.New("No targets registered"))
		}

		data := [][]string{
			[]string{"Target name", "Pid", "Host", "Admin port", "Proxy port"},
		}

		for key, target := range config.Targets {
			data = append(data, []string{key, strconv.Itoa(target.Pid), target.Host, strconv.Itoa(target.AdminPort), strconv.Itoa(target.ProxyPort)})
		}

		drawTable(data, true)
	},
}

var targetsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete target",
	Long: `
Delete target"
`,

	Run: func(cmd *cobra.Command, args []string) {
		checkArgAndExit(args, "Cannot delete a target without a name", "targets delete")

		if !askForConfirmation("Are you sure you want to delete the target " + args[0] + "?") {
			return
		}
		config.DeleteTarget(wrapper.Target{
			Name: args[0],
		})

		handleIfError(config.WriteToFile(hoverflyDirectory))

		targetsCmd.Run(cmd, args)
	},
}

var targetsNewCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new target",
	Long: `
Create target"
`,

	Run: func(cmd *cobra.Command, args []string) {
		checkArgAndExit(args, "Cannot create a target without a name", "targets new")

		newTarget := wrapper.NewTarget(args[0], hostFlag, adminPortFlag, proxyPortFlag)
		newTarget.Pid = pidFlag

		config.NewTarget(*newTarget)

		handleIfError(config.WriteToFile(hoverflyDirectory))

		targetsCmd.Run(cmd, args)
	},
}

func init() {
	RootCmd.AddCommand(targetsCmd)

	targetsCmd.AddCommand(targetsDeleteCmd)
	targetsCmd.AddCommand(targetsNewCmd)

	targetsNewCmd.Flags().IntVar(&pidFlag, "pid", 0, "Process id for a running instance of Hoverfly")
}
