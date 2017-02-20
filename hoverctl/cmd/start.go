package cmd

import (
	"fmt"

	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Hoverfly",
	Long: `
Starts an instance of Hoverfly using the current hoverctl
configuration.

The Hoverfly process ID will be written to a "pid" file in the 
".hoverfly" directory.

The "pid" file name is composed of the Hoverfly admin
port and proxy port.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			config.SetWebserver(args[0])
			hoverfly = wrapper.NewHoverfly(*config)
		}

		err := hoverfly.Start(hoverflyDirectory)
		handleIfError(err)

		data := [][]string{
			[]string{"admin-port", config.HoverflyAdminPort},
		}

		if config.HoverflyWebserver {
			fmt.Println("Hoverfly is now running as a webserver")
			data = append(data, []string{"webserver-port", config.HoverflyProxyPort})
		} else {
			fmt.Println("Hoverfly is now running")
			data = append(data, []string{"proxy-port", config.HoverflyProxyPort})
		}

		drawTable(data)

	},
}

func init() {
	RootCmd.AddCommand(startCmd)
}
