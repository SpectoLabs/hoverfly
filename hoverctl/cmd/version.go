package cmd

import (
	"os/exec"

	"github.com/kardianos/osext"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the version of hoverctl",
	Long: `
Shows the hoverctl version.
`,

	Run: func(cmd *cobra.Command, args []string) {

		binaryLocation, err := osext.ExecutableFolder()
		handleIfError(err)

		hoverflyCmd := exec.Command(binaryLocation+"/hoverfly", "-version")

		hoverflyVersion, _ := hoverflyCmd.CombinedOutput()

		data := [][]string{
			{"hoverctl", version},
			{"hoverfly", string(hoverflyVersion)},
		}

		drawTable(data, false)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
