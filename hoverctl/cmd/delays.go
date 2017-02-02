package cmd

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var delaysCmd = &cobra.Command{
	Use:   "delays",
	Short: "Get and set response delay config currently loaded in Hoverfly",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			delays, err := hoverfly.GetDelays()
			handleIfError(err)
			if len(delays) == 0 {
				log.Info("Hoverfly has no delays configured")
			} else {
				log.Info("Hoverfly has been configured with these delays")
				printResponseDelays(delays)
			}

		} else {
			delays, err := hoverfly.SetDelays(args[0])
			handleIfError(err)
			log.Info("Response delays set in Hoverfly: ")
			printResponseDelays(delays)
		}
	},
}

func init() {
	RootCmd.AddCommand(delaysCmd)
}

func printResponseDelays(delays []wrapper.ResponseDelaySchema) {
	for _, delay := range delays {
		var delayString string
		if delay.HttpMethod != "" {
			delayString = fmt.Sprintf("%v | %v - %vms", delay.HttpMethod, delay.UrlPattern, delay.Delay)
		} else {
			delayString = fmt.Sprintf("%v - %vms", delay.UrlPattern, delay.Delay)
		}
		log.Info(delayString)
	}
}
