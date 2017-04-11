package cmd

import (
	"fmt"
	"strconv"

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
		checkTargetAndExit(target)

		newTargetFlag, _ := cmd.Flags().GetString("new-target")

		if newTargetFlag != "" {
			target = wrapper.NewTarget(newTargetFlag, hostFlag, adminPortFlag, proxyPortFlag)
		}

		if adminPortFlag != 0 {
			target.AdminPort = adminPortFlag
		}

		if proxyPortFlag != 0 {
			target.ProxyPort = proxyPortFlag
		}

		target.Webserver = len(args) > 0
		target.CachePath, _ = cmd.Flags().GetString("cache")
		target.DisableCache, _ = cmd.Flags().GetBool("disable-cache")

		target.CertificatePath, _ = cmd.Flags().GetString("certificate")
		target.KeyPath, _ = cmd.Flags().GetString("key")
		target.DisableTls, _ = cmd.Flags().GetBool("disable-tls")

		target.UpstreamProxyUrl, _ = cmd.Flags().GetString("upstream-proxy")

		err := wrapper.Start(target, hoverflyDirectory)
		handleIfError(err)

		data := [][]string{
			[]string{"admin-port", strconv.Itoa(target.AdminPort)},
		}

		if target.Webserver {
			fmt.Println("Hoverfly is now running as a webserver")
			data = append(data, []string{"webserver-port", strconv.Itoa(target.ProxyPort)})
		} else {
			fmt.Println("Hoverfly is now running")
			data = append(data, []string{"proxy-port", strconv.Itoa(target.ProxyPort)})
		}

		drawTable(data, false)

		config.NewTarget(*target)
		handleIfError(config.WriteToFile(hoverflyDirectory))
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
	startCmd.Flags().String("new-target", "", "?")

	startCmd.Flags().String("cache", "", "A path to a persisted Hoverfly cache. If the cache doesn't exist, Hoverfly will create it")
	startCmd.Flags().Bool("disable-cache", false, "?")
	startCmd.Flags().String("certificate", "", "A path to a certificate file. Overrides the default Hoverfly certificate")
	startCmd.Flags().String("key", "", "A path to a key file. Overrides the default Hoverfly TLS key")
	startCmd.Flags().Bool("disable-tls", false, "Disables TLS verification")
	startCmd.Flags().String("upstream-proxy", "", "A host for which Hoverfly will proxy its requests to")
}
