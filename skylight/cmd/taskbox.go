package cmd

import "github.com/spf13/cobra"

var taskboxCmd = &cobra.Command{
	Use:   "task-box",
	Short: "Manage task bank items",
}

func init() {
	rootCmd.AddCommand(taskboxCmd)
}
