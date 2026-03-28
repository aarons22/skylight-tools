package cmd

import "github.com/spf13/cobra"

var rewardpointsCmd = &cobra.Command{
	Use:   "reward-points",
	Short: "View reward point balances",
}

func init() {
	rootCmd.AddCommand(rewardpointsCmd)
}
