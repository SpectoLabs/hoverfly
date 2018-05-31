package cmd

import (
	"fmt"

	"bytes"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Manage the diffs for Hoverfly",
	Long: `
This allows you to get or clean the differences 
between expected and actual responses stored by 
the Diff mode in Hoverfly. The diffs are represented 
as lists of strings grouped by the same requests.
	`,
}

const errorMsgTemplate = "\"%s\"\nthe expected value was [%s], but actual value was [%s]\n\n"

var getAllDiffCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets all diffs stored in Hoverfly",
	Long: `
Returns all differences between expected and actual responses from Hoverfly.
	`,
	Run: func(cmd *cobra.Command, args []string) {

		checkTargetAndExit(target)

		if len(args) == 0 {
			diffs, err := wrapper.GetAllDiffs(*target)
			handleIfError(err)
			var output bytes.Buffer

			for _, diffsWithRequest := range diffs {

				diffString := "diff"
				if len(diffsWithRequest.DiffReport) > 1 {
					diffString = "diffs"
				}
				output.WriteString(
					fmt.Sprintf("For request:\n"+
						"\n Method: %s \n Host: %s \n Path: %s \n Query:  %s \n\n%s %s recorded:\n",
						diffsWithRequest.Request.Method,
						diffsWithRequest.Request.Host,
						diffsWithRequest.Request.Path,
						diffsWithRequest.Request.Query,
						fmt.Sprint(len(diffsWithRequest.DiffReport)),
						diffString,
					))

				for index, diff := range diffsWithRequest.DiffReport {
					output.WriteString(fmt.Sprintf("\n%s. %s\n%s\n",
						fmt.Sprint(index+1), diff.Timestamp, diffReportMessage(diff)))
				}
			}

			if len(output.Bytes()) == 0 {
				fmt.Println("There are no diffs stored in Hoverfly")
			} else {
				fmt.Println(output.String())
			}
		}
	},
}

var deleteDiffsCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes all diffs",
	Long: `
Deletes all differences between expected and actual responses stored in Hoverfly.
	`,
	Run: func(cmd *cobra.Command, args []string) {

		checkTargetAndExit(target)

		err := wrapper.DeleteAllDiffs(*target)
		handleIfError(err)
		fmt.Println("All diffs have been deleted")
	},
}

func diffReportMessage(report v2.DiffReport) string {
	var msg bytes.Buffer
	for index, entry := range report.DiffEntries {
		msg.Write([]byte(fmt.Sprintf("(%d)"+errorMsgTemplate, index+1, entry.Field, entry.Expected, entry.Actual)))
	}
	return msg.String()
}

func init() {
	RootCmd.AddCommand(diffCmd)
	diffCmd.AddCommand(getAllDiffCmd)
	diffCmd.AddCommand(deleteDiffsCmd)
}
