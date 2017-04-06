package cmd

import (
	"fmt"

	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var cachePathFlag string

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

		if target == nil {
			var err error
			target, err = wrapper.NewTarget(targetName, host, adminPort, proxyPort)
			handleIfError(err)
		}

		target.Webserver = len(args) > 0
		target.CachePath = cachePathFlag
		target.DisableCache = cacheDisable

		target.CertificatePath = certificate
		target.KeyPath = key
		target.DisableTls = disableTls

		target.UpstreamProxyUrl = upstreamProxy

		err := hoverfly.Start(target, hoverflyDirectory)
		handleIfError(err)

		data := [][]string{
			[]string{"admin-port", config.HoverflyAdminPort},
		}

		if target.Webserver {
			fmt.Println("Hoverfly is now running as a webserver")
			data = append(data, []string{"webserver-port", config.HoverflyProxyPort})
		} else {
			fmt.Println("Hoverfly is now running")
			data = append(data, []string{"proxy-port", config.HoverflyProxyPort})
		}

		drawTable(data, false)

		config.NewTarget(*target)
		handleIfError(config.WriteToFile(hoverflyDirectory))
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVar(&cachePathFlag, "cache-path", "", "Something")
}
