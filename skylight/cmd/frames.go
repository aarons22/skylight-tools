package cmd

import "github.com/spf13/cobra"

var framesCmd = &cobra.Command{
	Use:   "frames",
	Short: "Manage Skylight frames",
}

func init() {
	rootCmd.AddCommand(framesCmd)
}
