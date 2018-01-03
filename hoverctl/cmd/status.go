package cmd

import (
	"fmt"
	"strconv"

	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the current status of Hoverfly",
	Long: `
If Hoverfly is running, this command will show an overview
of the instance of Hoverfly. This includes reporting the
mode and middleware set.
`,

	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)

		mode, err := wrapper.GetMode(*target)
		handleIfError(err)

		middleware, err := wrapper.GetMiddleware(*target)
		handleIfError(err)

		var proxyType string
		if target.Webserver {
			proxyType = "reverse (webserver)"
		} else {
			proxyType = "forward"
		}

		var middlewareStatus string

		if middleware.Binary == "" && middleware.Script == "" && middleware.Remote == "" {
			middlewareStatus = "disabled"
		} else {
			middlewareStatus = "enabled"
		}

		data := [][]string{
			{"Hoverfly", "running"},
			{"Admin port", strconv.Itoa(target.AdminPort)},
			{"Proxy port", strconv.Itoa(target.ProxyPort)},
			{"Proxy type", proxyType},
			{"Mode", mode.Mode},
			{"Middleware", middlewareStatus},
		}

		drawTable(data, false)
		if middlewareStatus == "enabled" {
			fmt.Println("")
			if middleware.Remote != "" {
				fmt.Println("Hoverfly is using remote middleware:\n" + middleware.Remote)
			} else {
				fmt.Println("Hoverfly is using local middleware with the command " + middleware.Binary + " and the script:\n" + middleware.Script)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
