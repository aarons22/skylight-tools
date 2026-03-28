package cmd

import (
	"fmt"
	"os"

	"github.com/aarons22/skylight-mcp/skylight/internal/client"
	"github.com/aarons22/skylight-mcp/skylight/internal/output"
	"github.com/spf13/cobra"
)

var choresCompleteChoreID string

var choresCompleteChoreCmd = &cobra.Command{
	Use:   "completeChore",
	Short: "Mark a chore as complete",
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
			"status": "complete",
		}

		c := client.NewClient(baseURL, token)
		resp, err := c.Do("PUT", "/frames/{frameId}/chores/{choreId}",
			map[string]string{"frameId": frameID, "choreId": choresCompleteChoreID}, nil, body)
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
	choresCmd.AddCommand(choresCompleteChoreCmd)
	choresCompleteChoreCmd.Flags().StringVar(&choresCompleteChoreID, "chore-id", "", "Chore ID")
	_ = choresCompleteChoreCmd.MarkFlagRequired("chore-id")
}
