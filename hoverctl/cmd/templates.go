package cmd

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Get and set request templates in Hoverfly",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			requestTemplatesData, err := hoverfly.GetRequestTemplates()
			handleIfError(err)
			requestTemplatesJson, err := json.MarshalIndent(requestTemplatesData, "", "    ")
			if err != nil {
				log.Error("Error marshalling JSON for printing request templates: " + err.Error())
			}
			fmt.Println(string(requestTemplatesJson))
		} else {
			requestTemplatesData, err := hoverfly.SetRequestTemplates(args[0])
			handleIfError(err)
			fmt.Println("Request template data set in Hoverfly: ")
			requestTemplatesJson, err := json.MarshalIndent(requestTemplatesData, "", "    ")
			if err != nil {
				log.Error("Error marshalling JSON for printing request templates: " + err.Error())
			}
			fmt.Println(string(requestTemplatesJson))
		}
	},
}

func init() {
	RootCmd.AddCommand(templatesCmd)
}
