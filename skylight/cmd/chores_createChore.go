package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/aarons22/skylight-tools/skylight/internal/client"
	"github.com/aarons22/skylight-tools/skylight/internal/output"
	"github.com/spf13/cobra"
)

var (
	choresCreateSummary        string
	choresCreateStart          string
	choresCreateRoutine        bool
	choresCreateCategoryIDs    string
	choresCreateEmoji          string
	choresCreateRewardPoints   int
	choresCreateStartTime      string
	choresCreateRecurrenceSet  string
	choresCreateRecurringUntil string
)

var choresCreateChoreCmd = &cobra.Command{
	Use:   "createChore",
	Short: "Create one or more chores (one per category ID)",
	Args:  cobra.NoArgs,
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

		catIDs := strings.Split(choresCreateCategoryIDs, ",")
		for i := range catIDs {
			catIDs[i] = strings.TrimSpace(catIDs[i])
		}

		body := map[string]interface{}{
			"summary":      choresCreateSummary,
			"start":        choresCreateStart,
			"routine":      choresCreateRoutine,
			"category_ids": catIDs,
		}
		if choresCreateEmoji != "" {
			body["emoji_icon"] = choresCreateEmoji
		}
		if cmd.Flags().Changed("reward-points") {
			body["reward_points"] = choresCreateRewardPoints
		}
		if choresCreateStartTime != "" {
			body["start_time"] = choresCreateStartTime
		}
		if choresCreateRecurrenceSet != "" {
			body["recurrence_set"] = choresCreateRecurrenceSet
		}
		if choresCreateRecurringUntil != "" {
			body["recurring_until"] = choresCreateRecurringUntil
		}

		c := client.NewClient(baseURL, token)
		resp, err := c.Do("POST", "/frames/{frameId}/chores/create_multiple",
			map[string]string{"frameId": frameID}, nil, body)
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
	choresCmd.AddCommand(choresCreateChoreCmd)
	choresCreateChoreCmd.Flags().StringVar(&choresCreateSummary, "summary", "", "Chore summary/title")
	choresCreateChoreCmd.Flags().StringVar(&choresCreateStart, "start", "", "Start date (YYYY-MM-DD)")
	choresCreateChoreCmd.Flags().BoolVar(&choresCreateRoutine, "routine", false, "Mark as a routine")
	choresCreateChoreCmd.Flags().StringVar(&choresCreateCategoryIDs, "category-ids", "", "Comma-separated category (profile) IDs")
	choresCreateChoreCmd.Flags().StringVar(&choresCreateEmoji, "emoji", "", "Emoji icon")
	choresCreateChoreCmd.Flags().IntVar(&choresCreateRewardPoints, "reward-points", 0, "Reward points")
	choresCreateChoreCmd.Flags().StringVar(&choresCreateStartTime, "start-time", "", "Optional start time (HH:MM)")
	choresCreateChoreCmd.Flags().StringVar(&choresCreateRecurrenceSet, "recurrence-set", "", "RRULE string (e.g. RRULE:FREQ=WEEKLY;INTERVAL=1;BYDAY=SA)")
	choresCreateChoreCmd.Flags().StringVar(&choresCreateRecurringUntil, "recurring-until", "", "End date for recurring chores (YYYY-MM-DD)")
	_ = choresCreateChoreCmd.MarkFlagRequired("summary")
	_ = choresCreateChoreCmd.MarkFlagRequired("start")
	_ = choresCreateChoreCmd.MarkFlagRequired("category-ids")
}
