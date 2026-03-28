package cmd

import (
	"fmt"
	"os"

	"github.com/aarons22/skylight-mcp/skylight/internal/client"
	"github.com/spf13/cobra"
)

var (
	choresDeleteChoreID  string
	choresDeleteApplyTo  string
)

var choresDeleteChoreCmd = &cobra.Command{
	Use:   "deleteChore",
	Short: "Delete a chore or chore series",
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

		queryParams := map[string]string{}
		if choresDeleteApplyTo != "" {
			queryParams["apply_to"] = choresDeleteApplyTo
		}

		c := client.NewClient(baseURL, token)
		resp, err := c.Do("DELETE", "/frames/{frameId}/chores/{choreId}",
			map[string]string{"frameId": frameID, "choreId": choresDeleteChoreID},
			queryParams, nil)
		if err != nil {
			return err
		}

		jsonMode, _ := cmd.Root().PersistentFlags().GetBool("json")
		if jsonMode {
			fmt.Fprintf(os.Stdout, "%s\n", string(resp))
			return nil
		}
		fmt.Println("Chore deleted.")
		return nil
	},
}

func init() {
	choresCmd.AddCommand(choresDeleteChoreCmd)
	choresDeleteChoreCmd.Flags().StringVar(&choresDeleteChoreID, "chore-id", "", "Chore ID")
	choresDeleteChoreCmd.Flags().StringVar(&choresDeleteApplyTo, "apply-to", "one", "Delete scope: 'one' (single instance) or 'all' (entire series)")
	_ = choresDeleteChoreCmd.MarkFlagRequired("chore-id")
}
