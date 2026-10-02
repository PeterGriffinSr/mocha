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
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := fmt.Fprintf(
			cmd.OutOrStdout(),
			"mocha %s\ncommit: %s\nbuilt: %s\n",
			version,
			commit,
			buildDate,
		)
		return err
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
