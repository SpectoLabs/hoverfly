package cmd

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Get the config being used by hoverctl",
	Long:  ``,

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
