package mocha

import "github.com/spf13/cobra"

var buildCmd = &cobra.Command{
	Use:   "build <file>",
	Short: "Compile a Mocha program",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// compiler.Build(args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
