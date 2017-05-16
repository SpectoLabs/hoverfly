package cmd

import (
	"fmt"
	"strings"

	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show hoverctl configuration information",
	Long: `
Shows the path to the configuration file being used by 
hoverctl, along with all of the current configuration values.

The configuration YAML file can found in the "$HOME/.hoverfly" 
directory. This directory and the file are created when you use
hoverctl for the first time.

Configuration values can be overridden using global flags. 
For example, setting the Hoverfly admin and proxy ports using 
the "--admin-port" and "--proxy-port" flags when calling
"hoverctl start" will override the corresponding values 
in the configuration file.
`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.GetFilepath())
		configData, _ := configuration.ReadFile(config.GetFilepath())
		configLines := strings.Split(string(configData), "\n")
		for _, line := range configLines {
			if line != "" {
				fmt.Println(line)
			}
		}
	},
}

var configHostCmd = &cobra.Command{
	Use:   "host",
	Short: "Get target host",
	Long: `
Gets the config value for the target host"
`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(target.Host)
	},
}

var configAdminPortCmd = &cobra.Command{
	Use:   "admin-port",
	Short: "Get target host",
	Long: `
Gets the config value for the target admin port"
`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(target.AdminPort)
	},
}

var configProxyPortCmd = &cobra.Command{
	Use:   "proxy-port",
	Short: "Get target host",
	Long: `
Gets the config value for the target proxy port"
`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(target.ProxyPort)
	},
}

var configAuthTokenCmd = &cobra.Command{
	Use:   "auth-token",
	Short: "Get target API token",
	Long: `
Gets the config value for the target API token if hoverctl has been logged in"
`,

	Run: func(cmd *cobra.Command, args []string) {
		if target.AuthToken == "" {
			handleIfError(fmt.Errorf("No auth token"))
		}
		fmt.Println(target.AuthToken)
	},
}

func init() {
	RootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configHostCmd)
	configCmd.AddCommand(configAdminPortCmd)
	configCmd.AddCommand(configProxyPortCmd)
	configCmd.AddCommand(configAuthTokenCmd)
}
