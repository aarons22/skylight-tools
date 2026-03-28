package cmd

import (
	"fmt"
	"os"

	"github.com/aarons22/skylight-mcp/skylight/internal/client"
	"github.com/aarons22/skylight-mcp/skylight/internal/output"
	"github.com/spf13/cobra"
)

var rewardpointsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get reward point balances for all profiles",
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
		resp, err := c.Do("GET", "/frames/{frameId}/reward_points",
			map[string]string{"frameId": frameID}, nil, nil)
		if err != nil {
			return err
		}

		jsonMode, _ := cmd.Root().PersistentFlags().GetBool("json")
		noColor, _ := cmd.Root().PersistentFlags().GetBool("no-color")
		if jsonMode {
			fmt.Fprintf(os.Stdout, "%s\n", string(resp))
			return nil
		}
		// reward_points returns a plain array
		if err := output.PrintTable(resp, noColor); err != nil {
			fmt.Println(string(resp))
		}
		return nil
	},
}

func init() {
	rewardpointsCmd.AddCommand(rewardpointsGetCmd)
}
