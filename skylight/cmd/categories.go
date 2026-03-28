package cmd

import "github.com/spf13/cobra"

var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "Manage family member profile categories",
}

func init() {
	rootCmd.AddCommand(categoriesCmd)
}
