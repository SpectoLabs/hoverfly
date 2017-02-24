package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// flushCmd represents the flush command
var flushCmd = &cobra.Command{
	Use:   "flush cache",
	Short: "TBC",
	Long: `
TBC
	`,

	Run: func(cmd *cobra.Command, args []string) {
		if !force {
			if !askForConfirmation("Are you sure you want to flush the cache?") {
				return
			}
		}

		err := hoverfly.FlushCache()
		handleIfError(err)

		fmt.Println("Successfully flushed cache")
	},
}

func init() {
	RootCmd.AddCommand(flushCmd)
	flushCmd.Flags().BoolVar(&force, "force", false,
		"Flush the cache without prompting for confirmation")
}
