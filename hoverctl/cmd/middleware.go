package cmd

import (
	"fmt"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var middlewareBinary, middlewareScript, middlewareRemote string

var middlewareCmd = &cobra.Command{
	Use:   "middleware",
	Short: "Get and set Hoverfly middleware",
	Long: `
Hoverfly middleware can be set using the following
combinations of flags: 

	--binary
	--binary --script
	--remote

If flags are not used, the current Hoverfly middleware
configuration will be shown.

`,

	Run: func(cmd *cobra.Command, args []string) {
		var middleware v2.MiddlewareView
		var err error

		if middlewareBinary == "" && middlewareScript == "" && middlewareRemote == "" {
			middleware, err = hoverfly.GetMiddleware()
			handleIfError(err)
			fmt.Println("Hoverfly middleware configuration is currently set to")
		} else {
			if middlewareRemote != "" {
				middleware, err = hoverfly.SetMiddleware("", "", middlewareRemote)
				handleIfError(err)
				fmt.Println("Hoverfly middleware configuration has been set to")
			} else {
				script, err := wrapper.ReadFile(middlewareScript)
				handleIfError(err)

				middleware, err = hoverfly.SetMiddleware(middlewareBinary, string(script), "")
				handleIfError(err)
				fmt.Println("Hoverfly middleware configuration has been set to")
			}
		}

		if middleware.Binary != "" {
			fmt.Println("Binary: " + middleware.Binary)
		}

		if middleware.Script != "" {
			fmt.Println("Script: " + middleware.Script)
		}

		if middleware.Remote != "" {
			fmt.Println("Remote: " + middleware.Remote)
		}
	},
}

func init() {
	RootCmd.AddCommand(middlewareCmd)
	middlewareCmd.PersistentFlags().StringVar(&middlewareBinary, "binary", "",
		"An absolute or relative path to a binary that Hoverfly will execute as middleware")
	middlewareCmd.PersistentFlags().StringVar(&middlewareScript, "script", "",
		"An absolute or relative path to a script that will be executed by the middleware binary")
	middlewareCmd.PersistentFlags().StringVar(&middlewareRemote, "remote", "",
		"A URL to a remote address that will be called by Hoverfly as middleware")
}
