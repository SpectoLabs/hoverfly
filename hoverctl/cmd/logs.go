package cmd

import (
	"fmt"
	"time"

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

		format := "plain"

		jsonLogs, _ := cmd.Flags().GetBool("json")
		if jsonLogs {
			format = "json"
		}

		var lastLogRequestTime *time.Time

		for followLogs || lastLogRequestTime == nil {
			logs, err := wrapper.GetLogs(*target, format, lastLogRequestTime)
			currentLogRequestTime := time.Now()
			handleIfError(err)

			for _, log := range logs {
				fmt.Println(log)
			}

			lastLogRequestTime = &currentLogRequestTime
			if followLogs {
				time.Sleep(time.Second * 2)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(logsCmd)

	logsCmd.Flags().Bool("json", false, "Retrieve the logs in JSON format")
	logsCmd.Flags().BoolVar(&followLogs, "follow", false, "Follows the Hoverfly logs")
}
