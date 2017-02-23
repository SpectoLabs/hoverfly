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
			[]string{"hoverctl", version},
			[]string{"hoverfly", string(hoverflyVersion)},
		}

		drawTable(data)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
