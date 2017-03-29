package cmd

import (
	"errors"
	"strconv"

	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

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
			[]string{"Target name", "Host", "Admin port"},
		}

		for key, target := range config.Targets {
			data = append(data, []string{key, target.Host, strconv.Itoa(target.AdminPort)})
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
		if targetName != "" {
			if !askForConfirmation("Are you sure you want to delete the target " + targetName + "?") {
				return
			}
			config.DeleteTarget(wrapper.TargetHoverfly{
				Name: targetName,
			})

			handleIfError(config.WriteToFile(hoverflyDirectory))
		} else {
			handleIfError(errors.New("Cannot delete a target without a name"))
		}

		targetsCmd.Run(cmd, args)
	},
}

var targetsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create target",
	Long: `
Create target"
`,

	Run: func(cmd *cobra.Command, args []string) {

		if targetName == "" {
			targetName = "default"
		}

		targetAdminPort := 8888

		if adminPort != "" {
			adminPort, err := strconv.Atoi(adminPort)
			handleIfError(err)
			targetAdminPort = adminPort
		}

		if host == "" {
			host = "localhost"
		}

		newTarget := wrapper.TargetHoverfly{
			Name:      targetName,
			AdminPort: targetAdminPort,
			Host:      host,
		}
		config.NewTarget(newTarget)

		handleIfError(config.WriteToFile(hoverflyDirectory))

		targetsCmd.Run(cmd, args)
	},
}

func init() {
	RootCmd.AddCommand(targetsCmd)

	targetsCmd.AddCommand(targetsDeleteCmd)
	targetsCmd.AddCommand(targetsCreateCmd)
}
