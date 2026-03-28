package cmd

import "github.com/spf13/cobra"

var listsCmd = &cobra.Command{
	Use:   "lists",
	Short: "Manage grocery/shopping lists",
}

func init() {
	rootCmd.AddCommand(listsCmd)
}
