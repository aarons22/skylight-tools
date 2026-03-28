package cmd

import (
	"fmt"
	"os"

	"github.com/aarons22/skylight-tools/skylight/internal/client"
	"github.com/aarons22/skylight-tools/skylight/internal/output"
	"github.com/spf13/cobra"
)

var listsListItemsListID string

var listsListItemsCmd = &cobra.Command{
	Use:   "listItems",
	Short: "List items in a list",
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
		resp, err := c.Do("GET", "/frames/{frameId}/lists/{listId}/items",
			map[string]string{"frameId": frameID, "listId": listsListItemsListID}, nil, nil)
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
	listsCmd.AddCommand(listsListItemsCmd)
	listsListItemsCmd.Flags().StringVar(&listsListItemsListID, "list-id", "", "List ID")
	_ = listsListItemsCmd.MarkFlagRequired("list-id")
}
