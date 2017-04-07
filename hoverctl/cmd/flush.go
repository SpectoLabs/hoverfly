package cmd

import (
	"fmt"

	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

// flushCmd represents the flush command
var flushCmd = &cobra.Command{
	Use:   "flush cache",
	Short: "Flush the internal cache in Hoverfly",
	Long: `
Hoverfly has a cache that is used to store incoming 
requests against matching requests and responses. This cache is flushed
when changing mode.

When changing the mode to simulate, the cache will be
flushed and rebuilt, pre-caching cachable matching requests.

This command will flush this cache regardless of mode.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)

		if !askForConfirmation("Are you sure you want to flush the cache?") {
			return
		}

		err := wrapper.FlushCache(*target)
		handleIfError(err)

		fmt.Println("Successfully flushed cache")
	},
}

func init() {
	RootCmd.AddCommand(flushCmd)
}
