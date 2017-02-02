package cmd

import (
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var followLogs bool

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Get the logs from Hoverfly",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		logfile := wrapper.NewLogFile(hoverflyDirectory, hoverfly.AdminPort, hoverfly.ProxyPort)

		if followLogs {
			err := logfile.Tail()
			handleIfError(err)
		} else {
			err := logfile.Print()
			handleIfError(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(logsCmd)

	logsCmd.Flags().BoolVar(&followLogs, "follow-logs", false, "Follow the logs from Hoverfly")
	logsCmd.Flag("follow-logs").Shorthand("f")
}
