package cmd

import "github.com/spf13/cobra"

var choresCmd = &cobra.Command{
	Use:   "chores",
	Short: "Manage scheduled chores",
}

func init() {
	rootCmd.AddCommand(choresCmd)
}
