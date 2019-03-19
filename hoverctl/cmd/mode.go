package cmd

import (
	"fmt"
	"strings"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/modes"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var specificHeaders string
var allHeaders bool
var stateful bool
var matchingStrategy string

var modeCmd = &cobra.Command{
	Use:   "mode [capture|diff|simulate|spy|modify|synthesize (optional)]",
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

			fmt.Println("Hoverfly is currently set to", mode.Mode, "mode", getExtraInfo(mode))

		} else {
			modeView := &v2.ModeView{
				Mode: args[0],
			}

			switch modeView.Mode {
			case modes.Simulate:
				if len(matchingStrategy) > 0 {
					modeView.Arguments.MatchingStrategy = &matchingStrategy
				}
				break
			case modes.Capture:
				modeView.Arguments.Stateful = stateful
				setHeaderArgument(modeView)
				break
			case modes.Diff:
				setHeaderArgument(modeView)
				break
			}

			mode, err := wrapper.SetModeWithArguments(*target, modeView)
			handleIfError(err)

			fmt.Println("Hoverfly has been set to", mode, "mode", getExtraInfo(modeView))
		}
	},
}

func setHeaderArgument(mode *v2.ModeView) {
	if allHeaders {
		mode.Arguments.Headers = append(mode.Arguments.Headers, "*")
	} else if len(specificHeaders) > 0 {
		splitHeaders := strings.Split(specificHeaders, ",")
		mode.Arguments.Headers = append(mode.Arguments.Headers, splitHeaders...)
	}
}

func getExtraInfo(mode *v2.ModeView) string {
	var extraInfo string
	switch mode.Mode {
	case modes.Simulate:
		if len(*mode.Arguments.MatchingStrategy) > 0 {
			extraInfo = fmt.Sprintf("with a matching strategy of '%s'", *mode.Arguments.MatchingStrategy)
		}
		break
	case modes.Capture:
		if len(mode.Arguments.Headers) > 0 {
			if len(mode.Arguments.Headers) == 1 && mode.Arguments.Headers[0] == "*" {
				extraInfo = "and will capture all request headers"
			} else {
				extraInfo = fmt.Sprintf("and will capture the following request headers: %s", mode.Arguments.Headers)
			}
		}
		break
	case modes.Diff:
		if len(mode.Arguments.Headers) > 0 {
			if len(mode.Arguments.Headers) == 1 && mode.Arguments.Headers[0] == "*" {
				extraInfo = "and will exclude all response headers from diffing"
			} else {
				extraInfo = fmt.Sprintf("and will exclude the following response headers from diffing: %s", mode.Arguments.Headers)
			}
		}
		break
	}

	return extraInfo
}

func init() {

	RootCmd.AddCommand(modeCmd)
	modeCmd.PersistentFlags().StringVar(&specificHeaders, "headers", "",
		"A comma separated list of request headers to record (for capture mode) or response headers to ignore (for diff mode) `Content-Type,Authorization`")
	modeCmd.PersistentFlags().BoolVar(&allHeaders, "all-headers", false,
		"Record all request headers (for capture mode) or ignore all response headers (for diff mode)")
	modeCmd.PersistentFlags().StringVar(&matchingStrategy, "matching-strategy", "strongest",
		"Sets the matching strategy - 'strongest | first'")
	modeCmd.PersistentFlags().BoolVar(&stateful, "stateful", false,
		"Record stateful responses as a sequence in capture mode")
}
