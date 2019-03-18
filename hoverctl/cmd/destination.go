package cmd

import (
	"errors"
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var dryRun string

var destinationCmd = &cobra.Command{
	Use:   "destination [host (optional)]",
	Short: "Get and set Hoverfly destination",
	Long: `
The "destination" setting allows you to specify which 
HTTP requests Hoverfly will process by supplying a 
Golang regular expression.

The default "destination" setting is ".", meaning that 
Hoverfly will process all HTTP requests.

If you use "destination" without supplying a value, 
hoverctl will show the current Hoverfly destination 
setting.
`,

	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)

		if len(args) == 0 {
			destination, err := wrapper.GetDestination(*target)
			handleIfError(err)

			fmt.Println("Current Hoverfly destination is set to", destination)
		} else {
			regexPattern, err := regexp.Compile(args[0])
			if err != nil {
				log.Debug(err.Error())
				handleIfError(errors.New("Regex pattern does not compile"))
			}

			if dryRun != "" {
				if regexPattern.MatchString(dryRun) {
					fmt.Println("The regex provided matches the dry-run URL")
				} else {
					handleIfError(errors.New("The regex provided does not match the dry-run URL"))
				}
			} else {
				destination, err := wrapper.SetDestination(*target, args[0])
				handleIfError(err)

				fmt.Println("Hoverfly destination has been set to", destination)
			}

		}
	},
}

func init() {
	RootCmd.AddCommand(destinationCmd)
	destinationCmd.Flags().StringVar(&dryRun, "dry-run", "",
		"The destination regexp will be applied to the URL provided. This allows the regexp to be tested.")
}
