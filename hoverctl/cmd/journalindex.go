package cmd

import (
	"fmt"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var indexName string

var journalIndexCommand = &cobra.Command{
	Use:   "journal-index",
	Short: "Manage the journal index for Hoverfly",
	Long: `
		This allows you to manage journal index in Hoverfly. 
	`,
}

var journalIndexGetCommand = &cobra.Command{
	Use:   "get-all",
	Short: "Get all journal index for Hoverfly",
	Long:  `Get all journal index for Hoverfly`,
	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)

		if len(args) == 0 {
			journalIndexes, err := wrapper.GetAllJournalIndexes(*target)
			handleIfError(err)
			drawTable(getJournalIndexesTabularData(journalIndexes), true)
		}
	},
}

var journalIndexSetCommand = &cobra.Command{
	Use:   "set",
	Short: "Set journal index for Hoverfly",
	Long: `
Hoverfly journal index can be set using the following flags: 
	 --key
`,
	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)
		if indexName == "" {
			fmt.Println("Key is compulsory to set journal index")
		} else {
			err := wrapper.SetJournalIndex(indexName, *target)
			handleIfError(err)
			fmt.Println("Success")
		}
	},
}

var journalIndexDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete journal index for Hoverfly",
	Long: `
Hoverfly journal index can be deleted using the following flags: 
	 --key
`,
	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)
		if indexName == "" {
			fmt.Println("Key is compulsory to delete journal index")
		} else {
			err := wrapper.DeleteJournalIndex(indexName, *target)
			handleIfError(err)
			fmt.Println("Success")
		}
	},
}

func init() {
	RootCmd.AddCommand(journalIndexCommand)
	journalIndexCommand.AddCommand(journalIndexGetCommand)
	journalIndexCommand.AddCommand(journalIndexSetCommand)
	journalIndexCommand.AddCommand(journalIndexDeleteCommand)
	journalIndexSetCommand.PersistentFlags().StringVar(&indexName, "key", "", "Index Key to be set")
	journalIndexDeleteCommand.PersistentFlags().StringVar(&indexName, "key", "", "Index Key to be deleted")
}

func getJournalIndexesTabularData(indexes []v2.JournalIndexView) [][]string {

	journalIndexesData := [][]string{{"Index Name"}}
	for _, index := range indexes {
		indexData := []string{index.Name}
		journalIndexesData = append(journalIndexesData, indexData)
	}
	return journalIndexesData
}
