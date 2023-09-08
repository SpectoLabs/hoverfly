package cmd

import (
	"fmt"

	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var name, filePath string

var templateCsvDataSourceCommand = &cobra.Command{
	Use:   "templating-data-source",
	Short: "Manage the templating data source for Hoverfly",
	Long: `
		This allows you to manage templating data source for Hoverfly. Only CSV datasource is supported as of now  
	`,
}

var templateCsvDataSourceGetCommand = &cobra.Command{
	Use:   "get-all",
	Short: "Get all templating data source for Hoverfly",
	Long:  `Get all templating data source for Hoverfly`,
	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)

		if len(args) == 0 {
			templatingDataSourceView, err := wrapper.GetAllTemplateDataSources(*target)
			handleIfError(err)
			drawTable(getTemplatingDataSourceTabularData(templatingDataSourceView), true)
		}
	},
}

var templateCsvDataSourceSetCommand = &cobra.Command{
	Use:   "set",
	Short: "Set csv templating data source for Hoverfly",
	Long: `
Hoverfly Templating CSV DataSource can be set using the following flags: 
	 --name --filePath
`,
	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)
		if name == "" || filePath == "" {
			fmt.Println("data source name and file path are compulsory to set csv templating data source")
		} else {
			data, err := configuration.ReadFile(filePath)
			handleIfError(err)
			err = wrapper.SetCsvTemplateDataSource(name, string(data), *target)
			handleIfError(err)
			fmt.Println("Success")
		}
	},
}

var templateCsvDataSourceDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete csv templating data source for Hoverfly",
	Long: `
Hoverfly CSV templating datasource can be deleted using the following flags: 
	 --name
`,
	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)
		if name == "" {
			fmt.Println("csv datasource name to be deleted not provided")
		} else {
			err := wrapper.DeleteCsvDataSource(name, *target)
			handleIfError(err)
			fmt.Println("Success")
		}
	},
}

func init() {
	RootCmd.AddCommand(templateCsvDataSourceCommand)
	templateCsvDataSourceCommand.AddCommand(templateCsvDataSourceGetCommand)
	templateCsvDataSourceCommand.AddCommand(templateCsvDataSourceSetCommand)
	templateCsvDataSourceCommand.AddCommand(templateCsvDataSourceDeleteCommand)

	templateCsvDataSourceSetCommand.PersistentFlags().StringVar(&name, "name", "", "Datasource Name to be set")
	templateCsvDataSourceSetCommand.PersistentFlags().StringVar(&filePath, "filePath", "",
		"An absolute or relative path to a csv file that Hoverfly will use for templating")
	templateCsvDataSourceDeleteCommand.PersistentFlags().StringVar(&name, "name", "", "Datasource Name to be set")

}

func getTemplatingDataSourceTabularData(templatingDataSourceView v2.TemplateDataSourceView) [][]string {

	templateDataSourceDetails := [][]string{{"CSV DataSource Name", "Content"}}
	for _, csvDataSourceView := range templatingDataSourceView.DataSources {
		csvDataSource := []string{csvDataSourceView.Name, getContentShorthand(csvDataSourceView.Data)}
		templateDataSourceDetails = append(templateDataSourceDetails, csvDataSource)
	}
	return templateDataSourceDetails
}
