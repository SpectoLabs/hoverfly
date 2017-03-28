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
			[]string{"Target name", "Admin port"},
		}

		for key, target := range config.Targets {
			data = append(data, []string{key, strconv.Itoa(target.AdminPort)})
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
			if !force {
				if !askForConfirmation("Are you sure you want to delete the target " + targetName + "?") {
					return
				}
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
		if targetName != "" {
			adminPort, err := strconv.Atoi(config.HoverflyAdminPort)
			handleIfError(err)
			target := wrapper.TargetHoverfly{
				Name:      targetName,
				AdminPort: adminPort,
			}
			config.NewTarget(target)

			handleIfError(config.WriteToFile(hoverflyDirectory))
		} else {
			handleIfError(errors.New("Cannot create a target without a name"))
		}

		targetsCmd.Run(cmd, args)
	},
}

func init() {
	RootCmd.AddCommand(targetsCmd)

	targetsCmd.AddCommand(targetsDeleteCmd)
	targetsDeleteCmd.Flags().BoolVar(&force, "force", false,
		"Delete the target without prompting for confirmation")

	targetsDeleteCmd.Flag("force").Shorthand = "f"

	targetsCmd.AddCommand(targetsCreateCmd)
}
