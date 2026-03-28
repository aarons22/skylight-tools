package cmd

import (
	"fmt"
	"os"

	"github.com/aarons22/skylight-tools/skylight/internal/client"
	"github.com/spf13/cobra"
)

var taskboxDeleteItemID string

var taskboxDeleteItemCmd = &cobra.Command{
	Use:   "deleteItem",
	Short: "Delete a task box item",
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

		c := client.NewClient(baseURL, token)
		resp, err := c.Do("DELETE", "/frames/{frameId}/task_box/items/{itemId}",
			map[string]string{"frameId": frameID, "itemId": taskboxDeleteItemID}, nil, nil)
		if err != nil {
			return err
		}

		jsonMode, _ := cmd.Root().PersistentFlags().GetBool("json")
		if jsonMode {
			fmt.Fprintf(os.Stdout, "%s\n", string(resp))
			return nil
		}
		fmt.Println("Task box item deleted.")
		return nil
	},
}

func init() {
	taskboxCmd.AddCommand(taskboxDeleteItemCmd)
	taskboxDeleteItemCmd.Flags().StringVar(&taskboxDeleteItemID, "item-id", "", "Task box item ID")
	_ = taskboxDeleteItemCmd.MarkFlagRequired("item-id")
}
