package cmd

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete test data from Hoverfly",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			args = append(args, "")
		}

		switch args[0] {
		case "all":
			err := hoverfly.DeleteSimulations()
			handleIfError(err)
			err = hoverfly.DeleteDelays()
			handleIfError(err)
			err = hoverfly.DeleteRequestTemplates()
			handleIfError(err)
			_, err = hoverfly.SetMiddleware("", "", "")
			handleIfError(err)

			log.Info("Delays, middleware, request templates and simulations have all been deleted from Hoverfly")
		case "simulations":
			err := hoverfly.DeleteSimulations()
			handleIfError(err)

			log.Info("Simulations have been deleted from Hoverfly")

		case "delays":
			err := hoverfly.DeleteDelays()
			handleIfError(err)

			log.Info("Delays have been deleted from Hoverfly")
		case "templates":
			err := hoverfly.DeleteRequestTemplates()
			handleIfError(err)

			log.Info("Request templates have been deleted from Hoverfly")

		case "middleware":
			_, err := hoverfly.SetMiddleware("", "", "")
			handleIfError(err)

			log.Info("Middleware has been deleted from Hoverfly")
		case "":
			err := errors.New("You have not specified a resource to delete from Hoverfly")
			handleIfError(err)
		default:
			err := errors.New("You have not specified a valid resource to delete from Hoverfly")
			handleIfError(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
