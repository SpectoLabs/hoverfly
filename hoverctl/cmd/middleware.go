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
	Long: `
Without flags, middleware will print the middleware
configuration values from Hoverfly.

When flags are provided, those flags will be set on
Hoverfly. Hoverfly can be set with the following
combinations of flags: 

	--binary
	--binary --script
	--remote
`,

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
	middlewareCmd.PersistentFlags().StringVar(&middlewareBinary, "binary", "",
		"Given a binary, the binary will be executed by Hoverfly as middleware")
	middlewareCmd.PersistentFlags().StringVar(&middlewareScript, "script", "",
		"Given a path, the contents will be read and will be executed with a binary by Hoverfly as middleware")
	middlewareCmd.PersistentFlags().StringVar(&middlewareRemote, "remote", "",
		"Given a URL, the URL will The remote address will be called by Hoverfly as middleware")
}
