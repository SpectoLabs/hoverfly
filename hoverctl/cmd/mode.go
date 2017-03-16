package cmd

import (
	"fmt"
	"strings"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/spf13/cobra"
)

var specficHeaders string
var allHeaders bool

var modeCmd = &cobra.Command{
	Use:   "mode [capture|simulate|modify|synthesize (optional)]",
	Short: "Get and set the Hoverfly mode",
	Long: `
Sets Hoverfly to the mode specified. The mode
determines how Hoverfly will process incoming
requests.

If a mode is not specified, the current Hoverfly 
mode is shown.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			mode, err := hoverfly.GetMode()
			handleIfError(err)

			fmt.Println("Hoverfly is currently set to", mode, "mode")
		} else {
			modeView := v2.ModeView{
				Mode: args[0],
			}

			var headersMessage string
			if allHeaders {
				modeView.Arguments.Headers = append(modeView.Arguments.Headers, "*")

				headersMessage = "and will capture all request headers"
			} else if len(specficHeaders) > 0 {
				modeView.Arguments.Headers = append(modeView.Arguments.Headers, strings.Split(specficHeaders, ",")...)

				if strings.Contains(specficHeaders, ",") {
					headersMessage = "and will capture " + specficHeaders + " request headers"
				} else {
					headersMessage = "and will capture " + specficHeaders + " request header"
				}
			}

			mode, err := hoverfly.SetModeWithArguments(modeView)
			handleIfError(err)

			fmt.Println("Hoverfly has been set to", mode, "mode", headersMessage)
		}
	},
}

func init() {
	RootCmd.AddCommand(modeCmd)
	modeCmd.PersistentFlags().StringVar(&specficHeaders, "headers", "",
		"?")
	modeCmd.PersistentFlags().BoolVar(&allHeaders, "all-headers", false,
		"?")
}
