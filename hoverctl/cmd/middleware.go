package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var middlewareBinary, middlewareScript, middlewareRemote string

var middlewareCmd = &cobra.Command{
	Use:   "middleware",
	Short: "Get and set Hoverfly's middleware",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		var middleware v2.MiddlewareView
		var err error

		if middlewareBinary == "" && middlewareScript == "" && middlewareRemote == "" {
			middleware, err = hoverfly.GetMiddleware()
			handleIfError(err)
			log.Info("Hoverfly is currently set to run the following as middleware")
		} else {
			if middlewareRemote != "" {
				middleware, err = hoverfly.SetMiddleware("", "", middlewareRemote)
				handleIfError(err)
				log.Info("Hoverfly is now set to run the following as middleware")
			} else {
				script, err := wrapper.ReadFile(middlewareScript)
				handleIfError(err)

				middleware, err = hoverfly.SetMiddleware(middlewareBinary, string(script), "")
				handleIfError(err)
				log.Info("Hoverfly is now set to run the following as middleware")
			}
		}

		if middleware.Binary != "" {
			log.Info("Binary: " + middleware.Binary)
		}

		if middleware.Script != "" {
			log.Info("Script: " + middleware.Script)
		}

		if middleware.Remote != "" {
			log.Info("Remote: " + middleware.Remote)
		}
	},
}

func init() {
	RootCmd.AddCommand(middlewareCmd)
	middlewareCmd.PersistentFlags().StringVar(&middlewareBinary, "binary", "", "The binary that middleware should execute")
	middlewareCmd.PersistentFlags().StringVar(&middlewareScript, "script", "", "The script that middleware should execute")
	middlewareCmd.PersistentFlags().StringVar(&middlewareRemote, "remote", "", "The remote address that middleware should execute")
}
