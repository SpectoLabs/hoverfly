package cmd

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Get the configuration being used by hoverctl",
	Long: `
Will print the path to the configuration file currently
being used by hoverctl along with all of the configuration
values set.

A configuration file is a YAML file that be can found in
the $HOME/.hoverfly directory. This directory and the
configuration file are created when you first use hoverctl.

Configuration in hoverctl can be overridden by using the 
global flags for each of the configuration values. An
example of this is providing both the --admin-port and
--proxy-port flags when calling the start command. This
would create several instances of Hoverfly without errors
due to ports already being used.
`,

	Run: func(cmd *cobra.Command, args []string) {
		log.Info(config.GetFilepath())
		configData, _ := wrapper.ReadFile(config.GetFilepath())
		configLines := strings.Split(string(configData), "\n")
		for _, line := range configLines {
			if line != "" {
				log.Info(line)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(configCmd)
}
