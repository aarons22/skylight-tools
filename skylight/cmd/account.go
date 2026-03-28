package cmd

import "github.com/spf13/cobra"

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Account management",
}

func init() {
	rootCmd.AddCommand(accountCmd)
}
