package cmd

import (
	"fmt"
	"os"

	"github.com/aarons22/skylight-mcp/skylight/internal/client"
	"github.com/aarons22/skylight-mcp/skylight/internal/output"
	"github.com/spf13/cobra"
)

var (
	taskboxCreateItemSummary      string
	taskboxCreateItemRoutine      bool
	taskboxCreateItemEmoji        string
	taskboxCreateItemRewardPoints int
)

var taskboxCreateItemCmd = &cobra.Command{
	Use:   "createItem",
	Short: "Create a task box item",
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

		body := map[string]interface{}{
			"summary": taskboxCreateItemSummary,
			"routine": taskboxCreateItemRoutine,
		}
		if taskboxCreateItemEmoji != "" {
			body["emoji_icon"] = taskboxCreateItemEmoji
		}
		if cmd.Flags().Changed("reward-points") {
			body["reward_points"] = taskboxCreateItemRewardPoints
		}

		c := client.NewClient(baseURL, token)
		resp, err := c.Do("POST", "/frames/{frameId}/task_box/items",
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
	taskboxCmd.AddCommand(taskboxCreateItemCmd)
	taskboxCreateItemCmd.Flags().StringVar(&taskboxCreateItemSummary, "summary", "", "Task summary")
	taskboxCreateItemCmd.Flags().BoolVar(&taskboxCreateItemRoutine, "routine", false, "Mark as a routine")
	taskboxCreateItemCmd.Flags().StringVar(&taskboxCreateItemEmoji, "emoji", "", "Emoji icon")
	taskboxCreateItemCmd.Flags().IntVar(&taskboxCreateItemRewardPoints, "reward-points", 0, "Reward points for completing the task")
	_ = taskboxCreateItemCmd.MarkFlagRequired("summary")
}
