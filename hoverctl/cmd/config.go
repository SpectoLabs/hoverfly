package cmd

import (
	"fmt"
	"strings"

	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
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
		configData, _ := wrapper.ReadFile(config.GetFilepath())
		configLines := strings.Split(string(configData), "\n")
		for _, line := range configLines {
			if line != "" {
				fmt.Println(line)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(configCmd)
}
