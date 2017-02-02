package cmd

import (
	"errors"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var dryRun string

var destinationCmd = &cobra.Command{
	Use:   "destination",
	Short: "Get and set Hoverfly's current destination",
	Long:  ``,

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
	destinationCmd.Flags().StringVar(&dryRun, "dry-run", "", "Test a url against a regex pattern")
}
