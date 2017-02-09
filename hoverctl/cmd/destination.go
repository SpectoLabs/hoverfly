package cmd

import (
	"errors"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var dryRun string

var destinationCmd = &cobra.Command{
	Use:   "destination [host (optional)]",
	Short: "Get and set Hoverfly's current destination",
	Long: `
Without specifying a host, destination will print the
destination configuration value from Hoverfly.

When a host is specified, that host will be set on
Hoverfly. That host will be used to whitelist which
HTTP requests Hoverfly will process. This host can be
specified as Golang regexp. The default destination is ".".
This will match against all incoming HTTP requests.
`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			destination, err := hoverfly.GetDestination()
			handleIfError(err)

			log.Info("The destination in Hoverfly is set to ", destination)
		} else {
			regexPattern, err := regexp.Compile(args[0])
			if err != nil {
				log.Debug(err.Error())
				handleIfError(errors.New("Regex pattern does not compile"))
			}

			if dryRun != "" {
				if regexPattern.MatchString(dryRun) {
					log.Info("The regex provided matches the dry run URL")
				} else {
					log.Fatal("The regex provided does not match the dry run URL")
				}
			} else {
				destination, err := hoverfly.SetDestination(args[0])
				handleIfError(err)

				log.Info("The destination in Hoverfly has been set to ", destination)
			}

		}
	},
}

func init() {
	RootCmd.AddCommand(destinationCmd)
	destinationCmd.Flags().StringVar(&dryRun, "dry-run", "",
		"Given a URL, the host regexp will be applied to the URL to allow testing of host regexp")
}
