package cmd

import (
	"fmt"
	"strings"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/modes"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var specficHeaders string
var allHeaders bool
var matchingStrategy string

var modeCmd = &cobra.Command{
	Use:   "mode [capture|simulate|spy|modify|synthesize (optional)]",
	Short: "Get and set the Hoverfly mode",
	Long: `
Sets Hoverfly to the mode specified. The mode
determines how Hoverfly will process incoming
requests.

If a mode is not specified, the current Hoverfly 
mode is shown.
`,
	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)

		if len(args) == 0 {
			mode, err := wrapper.GetMode(*target)
			handleIfError(err)

			var extraInformation string

			if mode.Mode == modes.Simulate {
				extraInformation = fmt.Sprintf("with a matching strategy of '%s'", *mode.Arguments.MatchingStrategy)
			}

			fmt.Println("Hoverfly is currently set to", mode.Mode, "mode", extraInformation)

		} else {
			modeView := v2.ModeView{
				Mode: args[0],
			}

			var extraInformation string

			//TODO: For @benji, convert this whole thing to a switch case for each mode, only allowing the correct functionality for each one
			if modeView.Mode == modes.Simulate && len(matchingStrategy) > 0 {
				extraInformation = fmt.Sprintf("with a matching strategy of '%s'", matchingStrategy)
				modeView.Arguments.MatchingStrategy = &matchingStrategy
			} else if allHeaders {
				modeView.Arguments.Headers = append(modeView.Arguments.Headers, "*")
				extraInformation = "and will capture all request headers"
			} else if len(specficHeaders) > 0 {
				splitHeaders := strings.Split(specficHeaders, ",")
				modeView.Arguments.Headers = append(modeView.Arguments.Headers, splitHeaders...)

				extraInformation = fmt.Sprintln("and will capture the following request headers:", splitHeaders)
			}

			mode, err := wrapper.SetModeWithArguments(*target, modeView)
			handleIfError(err)

			fmt.Println("Hoverfly has been set to", mode, "mode", extraInformation)
		}
	},
}

func init() {

	RootCmd.AddCommand(modeCmd)
	modeCmd.PersistentFlags().StringVar(&specficHeaders, "headers", "",
		"A comma separated list of headers to record in capture mode `Content-Type,Authorization`")
	modeCmd.PersistentFlags().BoolVar(&allHeaders, "all-headers", false,
		"Record all headers in capture mode")
	modeCmd.PersistentFlags().StringVar(&matchingStrategy, "matching-strategy", "strongest",
		"Sets the matching strategy - 'strongest | first'")
}
