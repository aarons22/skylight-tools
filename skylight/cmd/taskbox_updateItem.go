package cmd

import (
	"fmt"
	"os"

	"github.com/aarons22/skylight-tools/skylight/internal/client"
	"github.com/aarons22/skylight-tools/skylight/internal/output"
	"github.com/spf13/cobra"
)

var (
	taskboxUpdateItemID           string
	taskboxUpdateItemSummary      string
	taskboxUpdateItemRoutine      bool
	taskboxUpdateItemEmoji        string
	taskboxUpdateItemRewardPoints int
)

var taskboxUpdateItemCmd = &cobra.Command{
	Use:   "updateItem",
	Short: "Update a task box item",
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

		body := map[string]interface{}{}
		if cmd.Flags().Changed("summary") {
			body["summary"] = taskboxUpdateItemSummary
		}
		if cmd.Flags().Changed("routine") {
			body["routine"] = taskboxUpdateItemRoutine
		}
		if cmd.Flags().Changed("emoji") {
			body["emoji_icon"] = taskboxUpdateItemEmoji
		}
		if cmd.Flags().Changed("reward-points") {
			body["reward_points"] = taskboxUpdateItemRewardPoints
		}

		c := client.NewClient(baseURL, token)
		resp, err := c.Do("PATCH", "/frames/{frameId}/task_box/items/{itemId}",
			map[string]string{"frameId": frameID, "itemId": taskboxUpdateItemID}, nil, body)
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
	taskboxCmd.AddCommand(taskboxUpdateItemCmd)
	taskboxUpdateItemCmd.Flags().StringVar(&taskboxUpdateItemID, "item-id", "", "Task box item ID")
	taskboxUpdateItemCmd.Flags().StringVar(&taskboxUpdateItemSummary, "summary", "", "New task summary")
	taskboxUpdateItemCmd.Flags().BoolVar(&taskboxUpdateItemRoutine, "routine", false, "Routine flag")
	taskboxUpdateItemCmd.Flags().StringVar(&taskboxUpdateItemEmoji, "emoji", "", "Emoji icon")
	taskboxUpdateItemCmd.Flags().IntVar(&taskboxUpdateItemRewardPoints, "reward-points", 0, "Reward points")
	_ = taskboxUpdateItemCmd.MarkFlagRequired("item-id")
}
