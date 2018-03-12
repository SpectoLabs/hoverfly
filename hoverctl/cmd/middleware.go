package cmd

import (
	"fmt"
	"strings"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
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
		checkTargetAndExit(target)

		var middleware v2.MiddlewareView
		var err error
		if middlewareBinary == "" && middlewareScript == "" && middlewareRemote == "" {
			middleware, err = wrapper.GetMiddleware(*target)
			handleIfError(err)
			fmt.Println("Hoverfly middleware configuration is currently set to")
		} else {
			if middlewareRemote != "" {
				fmt.Println("Testing middleware against Hoverfly...")
				middleware, err = wrapper.SetMiddleware(*target, "", "", middlewareRemote)
				handleIfError(err)
				fmt.Println("Hoverfly middleware configuration has been set to")
			} else {
				var script []byte
				if middlewareScript != "" {
					script, err = configuration.ReadFile(middlewareScript)
					handleIfError(err)
				}

				fmt.Println("Testing middleware against Hoverfly...")
				middleware, err = wrapper.SetMiddleware(*target, middlewareBinary, string(script), "")
				handleIfError(err)

				fmt.Println("Hoverfly middleware configuration has been set to")
			}
		}

		if middleware.Binary != "" {
			fmt.Println("Binary: " + middleware.Binary)
		}

		if middleware.Script != "" {
			middlewareScript := strings.Split(middleware.Script, "\n")
			if verbose || len(middlewareScript) < 5 {
				fmt.Println("Script: " + middleware.Script)
			} else {
				fmt.Println("Script: " + middlewareScript[0] + "\n" +
					middlewareScript[1] + "\n" +
					middlewareScript[2] + "\n" +
					middlewareScript[3] + "\n" +
					middlewareScript[4] + "\n" +
					"...")
			}
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
