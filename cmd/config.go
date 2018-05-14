package cmd

import (
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Modify elastic config",
	Long:  `Manage elastic configuration details`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	configCmd.AddCommand(viewCmd)
	configCmd.AddCommand(setCmd)
}
