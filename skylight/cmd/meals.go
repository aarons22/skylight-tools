package cmd

import "github.com/spf13/cobra"

var mealsCmd = &cobra.Command{
	Use:   "meals",
	Short: "Manage meal planning",
}

func init() {
	rootCmd.AddCommand(mealsCmd)
}
