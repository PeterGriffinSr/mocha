package mocha

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	commit    = "unknown"
	buildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the Mocha version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(
			cmd.OutOrStdout(),
			"mocha %s\ncommit: %s\nbuilt: %s\n",
			version,
			commit,
			buildDate,
		)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
