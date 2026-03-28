package cmd

import (
	"fmt"
	"os"

	"github.com/aarons22/skylight-tools/skylight/internal/client"
	"github.com/aarons22/skylight-tools/skylight/internal/output"
	"github.com/spf13/cobra"
)

var (
	choresUpdateChoreID          string
	choresUpdateSummary          string
	choresUpdateEmoji            string
	choresUpdateRewardPoints     int
	choresUpdateStart            string
	choresUpdateRecurrenceSet    string
	choresUpdateCategoryID       string
	choresUpdateRoutine          bool
	choresUpdateUpForGrabs       bool
)

var choresUpdateChoreCmd = &cobra.Command{
	Use:   "updateChore",
	Short: "Update a chore (full object PUT)",
	Long: `Update a chore. Sends a full PUT request.
Chore ID formats:
  Non-recurring: numeric ID (e.g. 72279767)
  Recurring instance: composite ID (e.g. 72279767-2026-03-21)`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		baseURL, _ := cmd.Root().PersistentFlags().GetString("base-url")
		configPath, _ := cmd.Root().PersistentFlags().GetString("config")
		frameIDFlag, _ := cmd.Root().PersistentFlags().GetString("frame-id")

		token, err := client.ResolveToken(configPath)
		if err != nil {
			return err
		}
		frameID, err := client.ResolveFrameID(frameIDFlag, configPath)
		if err != nil {
			return err
		}

		body := map[string]interface{}{}
		if cmd.Flags().Changed("summary") {
			body["summary"] = choresUpdateSummary
		}
		if cmd.Flags().Changed("emoji") {
			body["emoji_icon"] = choresUpdateEmoji
		}
		if cmd.Flags().Changed("reward-points") {
			body["reward_points"] = choresUpdateRewardPoints
		}
		if cmd.Flags().Changed("start") {
			body["start"] = choresUpdateStart
		}
		if cmd.Flags().Changed("recurrence-set") {
			body["recurrence_set"] = choresUpdateRecurrenceSet
		}
		if cmd.Flags().Changed("category-id") {
			body["category_id"] = choresUpdateCategoryID
		}
		if cmd.Flags().Changed("routine") {
			body["routine"] = choresUpdateRoutine
		}
		if cmd.Flags().Changed("up-for-grabs") {
			body["up_for_grabs"] = choresUpdateUpForGrabs
		}

		c := client.NewClient(baseURL, token)
		resp, err := c.Do("PUT", "/frames/{frameId}/chores/{choreId}",
			map[string]string{"frameId": frameID, "choreId": choresUpdateChoreID}, nil, body)
		if err != nil {
			return err
		}

		jsonMode, _ := cmd.Root().PersistentFlags().GetBool("json")
		noColor, _ := cmd.Root().PersistentFlags().GetBool("no-color")
		if jsonMode {
			fmt.Fprintf(os.Stdout, "%s\n", string(resp))
			return nil
		}
		if err := output.PrintTable(resp, noColor); err != nil {
			fmt.Println(string(resp))
		}
		return nil
	},
}

func init() {
	choresCmd.AddCommand(choresUpdateChoreCmd)
	choresUpdateChoreCmd.Flags().StringVar(&choresUpdateChoreID, "chore-id", "", "Chore ID")
	choresUpdateChoreCmd.Flags().StringVar(&choresUpdateSummary, "summary", "", "Chore summary")
	choresUpdateChoreCmd.Flags().StringVar(&choresUpdateEmoji, "emoji", "", "Emoji icon")
	choresUpdateChoreCmd.Flags().IntVar(&choresUpdateRewardPoints, "reward-points", 0, "Reward points")
	choresUpdateChoreCmd.Flags().StringVar(&choresUpdateStart, "start", "", "Start date (YYYY-MM-DD)")
	choresUpdateChoreCmd.Flags().StringVar(&choresUpdateRecurrenceSet, "recurrence-set", "", "RRULE string (empty to clear)")
	choresUpdateChoreCmd.Flags().StringVar(&choresUpdateCategoryID, "category-id", "", "Category (profile) ID")
	choresUpdateChoreCmd.Flags().BoolVar(&choresUpdateRoutine, "routine", false, "Routine flag")
	choresUpdateChoreCmd.Flags().BoolVar(&choresUpdateUpForGrabs, "up-for-grabs", false, "Up for grabs flag")
	_ = choresUpdateChoreCmd.MarkFlagRequired("chore-id")
}
