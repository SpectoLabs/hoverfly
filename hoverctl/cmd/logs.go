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

		logs, err := wrapper.GetLogs(*target, format)
		handleIfError(err)

		logsPrinted := map[string]string{
			"": "x",
		}

		for i := 0; i < len(logs); i++ {

			if logs[i] != "" && logsPrinted[logs[i]] != "x" {
				fmt.Println(logs[i])
				logsPrinted[logs[i]] = "x"
			}

			if i == len(logs)-1 && followLogs {
				logs, err = wrapper.GetLogs(*target, format)
				handleIfError(err)

				i = 0
				time.Sleep(time.Second * 5)
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
