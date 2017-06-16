package cmd

import (
	"fmt"
	"strconv"

	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Hoverfly",
	Long: `
Starts an instance of Hoverfly using the current hoverctl
target configuration.

To start an instance of Hoverfly as a webserver, add the 
argument "webserver" to the start command.

The Hoverfly process ID is stored against the target in the
hoverctl configuration file.
`,
	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)

		if !wrapper.IsLocal(target.Host) {
			handleIfError(fmt.Errorf("Unable to start an instance of Hoverfly on a remote host (%s host: %s)\n\nRun `hoverctl start --new-target <name>`", target.Name, target.Host))
		}

		if wrapper.CheckIfRunning(*target) == nil {
			if _, err := wrapper.GetMode(*target); err == nil {
				handleIfError(fmt.Errorf("Target Hoverfly is already running \n\nRun `hoverctl stop -t %s` to stop it", target.Name))
			}
		}

		newTargetFlag, _ := cmd.Flags().GetString("new-target")

		if newTargetFlag != "" {
			if config.GetTarget(newTargetFlag) != nil {
				handleIfError(fmt.Errorf("Target %s already exists\n\nUse a different target name or run `hoverctl targets update %[1]s`", newTargetFlag))
			}
			target = configuration.NewTarget(newTargetFlag, hostFlag, adminPortFlag, proxyPortFlag)
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
		target.HttpsOnly, _ = cmd.Flags().GetBool("https-only")

		if enableAuth, _ := cmd.Flags().GetBool("auth"); enableAuth {
			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")

			if username == "" {
				username = askForInput("Username", false)
			}
			if password == "" {
				password = askForInput("Password", true)
				fmt.Println("")
			}

			target.AuthEnabled = true
			target.Username = username
			target.Password = password
		}

		err := wrapper.Start(target)
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
	startCmd.Flags().String("new-target", "", "A name for a new target that hoverctl will create and associate the Hoverfly instance to")

	startCmd.Flags().String("cache", "", "A path to a persisted Hoverfly cache. If the cache doesn't exist, Hoverfly will create it")
	startCmd.Flags().Bool("disable-cache", false, "Disables the request response cache on Hoverfly")
	startCmd.Flags().String("certificate", "", "A path to a certificate file. Overrides the default Hoverfly certificate")
	startCmd.Flags().String("key", "", "A path to a key file. Overrides the default Hoverfly TLS key")
	startCmd.Flags().Bool("disable-tls", false, "Disables TLS verification")
	startCmd.Flags().String("upstream-proxy", "", "A host for which Hoverfly will proxy its requests to")
	startCmd.Flags().Bool("https-only", false, "Disables insecure HTTP traffic in Hoverfly")

	startCmd.Flags().Bool("auth", false, "Enable authenticiation on Hoverfly")
	startCmd.Flags().String("username", "", "Username to authenticate Hoverfly")
	startCmd.Flags().String("password", "", "Password to authenticate Hoverfly")
}
