package cmd

import (
	"fmt"
	"strconv"

	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var followLogs bool

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Get the logs from Hoverfly",
	Long: `
Shows the Hoverfly logs.
`,

	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)

		jsonLogs, _ := cmd.Flags().GetBool("json")

		if jsonLogs {
			logfile := wrapper.NewLogFile(hoverflyDirectory, strconv.Itoa(target.AdminPort), strconv.Itoa(target.ProxyPort))

			if followLogs {
				err := logfile.Tail()
				handleIfError(err)
			} else {

				logs, err := wrapper.GetLogs(*target, "json")
				handleIfError(err)
				for _, log := range logs {
					if log != "" {
						fmt.Println(log)
					}
				}
			}
		} else {
			logs, err := wrapper.GetLogs(*target, "plain")
			handleIfError(err)
			for _, log := range logs {
				if log != "" {
					fmt.Println(log)
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(logsCmd)

	logsCmd.Flags().Bool("json", false, "Retrieve the logs in JSON format")
	logsCmd.Flags().BoolVar(&followLogs, "follow-logs", false, "Follows the Hoverfly logs")
	logsCmd.Flag("follow-logs")
}
