package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aarons22/skylight-tools/skylight/internal/client"
	"github.com/aarons22/skylight-tools/skylight/internal/output"
	"github.com/spf13/cobra"
)

var framesListFramesCmd = &cobra.Command{
	Use:   "listFrames",
	Short: "List all Skylight frames",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		baseURL, _ := cmd.Root().PersistentFlags().GetString("base-url")
		configPath, _ := cmd.Root().PersistentFlags().GetString("config")

		token, err := client.ResolveToken(configPath)
		if err != nil {
			return err
		}

		c := client.NewClient(baseURL, token)
		resp, err := c.Do("GET", "/frames", nil, nil, nil)
		if err != nil {
			return err
		}

		jsonMode, _ := cmd.Root().PersistentFlags().GetBool("json")
		noColor, _ := cmd.Root().PersistentFlags().GetBool("no-color")
		if jsonMode {
			fmt.Fprintf(os.Stdout, "%s\n", string(resp))
			return nil
		}

		// /frames returns {"frames": [...]} not JSON:API format
		var fr struct {
			Frames []map[string]interface{} `json:"frames"`
		}
		if err := json.Unmarshal(resp, &fr); err != nil || len(fr.Frames) == 0 {
			fmt.Println(string(resp))
			return nil
		}

		plain, err := json.Marshal(fr.Frames)
		if err != nil {
			fmt.Println(string(resp))
			return nil
		}
		if err := output.PrintTable(plain, noColor); err != nil {
			fmt.Println(string(resp))
		}
		return nil
	},
}

func init() {
	framesCmd.AddCommand(framesListFramesCmd)
}
